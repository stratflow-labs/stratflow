package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/stratflow-labs/stratflow/internal/foundation/logger"

	"github.com/pressly/goose/v3"
)

const (
	gooseSchemaEnvVar  = "MIGRATIONS_GOOSE_SCHEMA"
	defaultGooseSchema = "infra_meta"
)

// RunServiceMigrations runs migrations for a service if the directory exists.
func RunServiceMigrations(sqlDB *sql.DB, serviceName string) error {
	return runIfExists(
		sqlDB,
		serviceName,
		servicePath(serviceName, "migrations"),
		"migrations",
		tableName(serviceName),
		legacyTableName(serviceName),
	)
}

// RunServiceSeeds runs seed migrations for a service if the directory exists.
func RunServiceSeeds(sqlDB *sql.DB, serviceName string) error {
	return runIfExists(
		sqlDB,
		serviceName,
		servicePath(serviceName, "seeds"),
		"seeds",
		seedsTableName(serviceName),
		legacySeedsTableName(serviceName),
	)
}

// ServiceDirExists reports whether the service's db/<kind> directory exists.
func ServiceDirExists(serviceName, kind string) (bool, error) {
	return dirExists(servicePath(serviceName, kind))
}

func runIfExists(sqlDB *sql.DB, serviceName, dir string, label string, table string, legacyTable string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	goose.SetLogger(newServiceGooseLogger(serviceName, label))
	if err := ensureGooseVersionTable(sqlDB, table, legacyTable); err != nil {
		return err
	}
	goose.SetTableName(table)
	exists, err := dirExists(dir)
	if err != nil {
		return err
	}
	if !exists {
		logger.Warn(label+" directory not found, skip", "path", dir)
		return nil
	}

	hasFiles, err := hasGooseFiles(dir)
	if err != nil {
		return err
	}
	if !hasFiles {
		logger.Info("no "+label+" files found, skip", "path", dir)
		return nil
	}

	if err := goose.Up(sqlDB, dir); err != nil {
		if isNoMigrationFilesError(err) {
			logger.Info("no "+label+" files found, skip", "path", dir)
			return nil
		}
		return err
	}
	return nil
}

func dirExists(dir string) (bool, error) {
	checkPath := dir
	if !filepath.IsAbs(dir) {
		if wd, err := os.Getwd(); err == nil {
			checkPath = filepath.Join(wd, dir)
		}
	}
	if _, err := os.Stat(checkPath); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func hasGooseFiles(dir string) (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := strings.ToLower(entry.Name())
		if strings.HasSuffix(name, ".sql") || strings.HasSuffix(name, ".go") {
			return true, nil
		}
	}
	return false, nil
}

func isNoMigrationFilesError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(strings.ToLower(err.Error()), "no migration files found")
}

func servicePath(serviceName, kind string) string {
	return filepath.Join("services", serviceName, "db", kind)
}

func tableName(serviceName string) string {
	return qualifyGooseTable(legacyTableName(serviceName))
}

func legacyTableName(serviceName string) string {
	safe := strings.ReplaceAll(serviceName, "-", "_")
	return "goose_db_version_" + safe
}

func seedsTableName(serviceName string) string {
	return qualifyGooseTable(legacySeedsTableName(serviceName))
}

func legacySeedsTableName(serviceName string) string {
	safe := strings.ReplaceAll(serviceName, "-", "_")
	return "goose_db_version_" + safe + "_seeds"
}

func qualifyGooseTable(baseTable string) string {
	schema := strings.TrimSpace(os.Getenv(gooseSchemaEnvVar))
	if schema == "" {
		schema = defaultGooseSchema
	}
	return schema + "." + baseTable
}

func ensureGooseVersionTable(sqlDB *sql.DB, targetTable, legacyTable string) error {
	schema, table, ok := splitQualifiedName(targetTable)
	if !ok {
		return nil
	}

	if !isSimpleIdentifier(schema) || !isSimpleIdentifier(table) {
		return fmt.Errorf("invalid goose table name: %q", targetTable)
	}
	if !isSimpleIdentifier(legacyTable) {
		return fmt.Errorf("invalid legacy goose table name: %q", legacyTable)
	}

	if _, err := sqlDB.ExecContext(context.Background(), `CREATE SCHEMA IF NOT EXISTS `+quoteIdentifier(schema)); err != nil {
		return fmt.Errorf("create goose schema %q: %w", schema, err)
	}

	targetExists, err := dbTableExists(sqlDB, schema, table)
	if err != nil {
		return err
	}
	if targetExists {
		return nil
	}

	legacyExists, err := dbTableExists(sqlDB, "public", legacyTable)
	if err != nil {
		return err
	}
	if !legacyExists {
		return nil
	}

	alterStmt := `ALTER TABLE public.` + quoteIdentifier(legacyTable) + ` SET SCHEMA ` + quoteIdentifier(schema)
	if _, err := sqlDB.ExecContext(context.Background(), alterStmt); err != nil {
		return fmt.Errorf("move goose table public.%s to %s.%s: %w", legacyTable, schema, table, err)
	}
	logger.Info("moved goose version table to dedicated schema", "from", "public."+legacyTable, "to", schema+"."+table)
	return nil
}

func dbTableExists(sqlDB *sql.DB, schema, table string) (bool, error) {
	const query = `
SELECT EXISTS (
  SELECT 1
  FROM information_schema.tables
  WHERE table_schema = $1 AND table_name = $2
)`

	var exists bool
	if err := sqlDB.QueryRowContext(context.Background(), query, schema, table).Scan(&exists); err != nil {
		return false, fmt.Errorf("check table existence %s.%s: %w", schema, table, err)
	}
	return exists, nil
}

func splitQualifiedName(name string) (schema string, table string, ok bool) {
	parts := strings.Split(name, ".")
	if len(parts) != 2 {
		return "", "", false
	}
	schema = strings.TrimSpace(parts[0])
	table = strings.TrimSpace(parts[1])
	if schema == "" || table == "" {
		return "", "", false
	}
	return schema, table, true
}

func isSimpleIdentifier(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_' || (i > 0 && r >= '0' && r <= '9') {
			continue
		}
		return false
	}
	return true
}

func quoteIdentifier(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}
