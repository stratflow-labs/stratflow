package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
	strategyregistrydbsqlc "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/adapters/postgres/sqlc/gen"
	manualqueries "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/adapters/postgres/sqlc/queries"
	attributevalue "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/attributevalue"
	attributeValuedomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"
	attributedomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"
	strategydomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

// Manual queries для batch/сложных операций
const (
	pgForeignKeyViolationCode           = "23503"
	attributeValueSlugUniqueConstraint  = "uq_strategy_param_value_param_slug"
	attributeValueAttributeFKConstraint = "strategy_attribute_value_strategy_attribute_id_fkey"
)

type AttributeValueRepository struct {
	db *sql.DB
}

func NewAttributeValueRepository(db *sql.DB) *AttributeValueRepository {
	return &AttributeValueRepository{db: db}
}

var _ attributevalue.AttributeValueRepository = (*AttributeValueRepository)(nil)

func (r *AttributeValueRepository) Create(ctx context.Context, attributeValue *attributeValuedomain.AttributeValue) (attributeValuedomain.AttributeValue, error) {
	if attributeValue == nil {
		return attributeValuedomain.AttributeValue{}, fmt.Errorf("create attribute value: nil attribute value")
	}
	if attributeValue.ID == uuid.Nil {
		attributeValue.ID = uuid.New()
	}

	strategyID, err := r.lookupStrategyIDByAttributeID(ctx, attributeValue.AttributeID)
	if err != nil {
		return attributeValuedomain.AttributeValue{}, fmt.Errorf("create attribute value: %w", err)
	}
	if attributeValue.StrategyID != uuid.Nil && strategyID != attributeValue.StrategyID {
		return attributeValuedomain.AttributeValue{}, apperr.NotFound[attributedomain.Attribute]("lookup", attributeValue.AttributeID)
	}

	q := strategyregistrydbsqlc.NewQuerier(ctx, r.db)
	row, err := q.CreateStrategyAttributeValue(ctx, attributeValueToCreateParams(attributeValue))
	if err != nil {
		return attributeValuedomain.AttributeValue{}, fmt.Errorf("create attribute value: %w", mapAttributeValuePersistenceError(err))
	}

	return attributeValueToDomain(&row, strategyID), nil
}

func (r *AttributeValueRepository) GetByID(ctx context.Context, id uuid.UUID) (attributeValuedomain.AttributeValue, error) {
	q := strategyregistrydbsqlc.NewQuerier(ctx, r.db)
	row, err := q.GetStrategyAttributeValueByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return attributeValuedomain.AttributeValue{}, apperr.NotFound[attributeValuedomain.AttributeValue]("lookup", id)
	}
	if err != nil {
		return attributeValuedomain.AttributeValue{}, fmt.Errorf("get attribute value by id: %w", err)
	}

	strategyID, err := r.lookupStrategyIDByAttributeID(ctx, row.StrategyAttributeID)
	if err != nil {
		return attributeValuedomain.AttributeValue{}, fmt.Errorf("get attribute value by id: %w", err)
	}

	return attributeValueToDomain(&row, strategyID), nil
}

func (r *AttributeValueRepository) GetBySlug(ctx context.Context, attributeID uuid.UUID, slug string) (attributeValuedomain.AttributeValue, error) {
	q := strategyregistrydbsqlc.NewQuerier(ctx, r.db)
	row, err := q.GetStrategyAttributeValueBySlug(ctx, strategyregistrydbsqlc.GetStrategyAttributeValueBySlugParams{
		StrategyAttributeID: attributeID,
		Slug:                strings.TrimSpace(slug),
	})
	if errors.Is(err, sql.ErrNoRows) {
		return attributeValuedomain.AttributeValue{}, apperr.NotFound[attributeValuedomain.AttributeValue]("lookup", slug)
	}
	if err != nil {
		return attributeValuedomain.AttributeValue{}, fmt.Errorf("get attribute value by slug: %w", err)
	}

	strategyID, err := r.lookupStrategyIDByAttributeID(ctx, row.StrategyAttributeID)
	if err != nil {
		return attributeValuedomain.AttributeValue{}, fmt.Errorf("get attribute value by slug: %w", err)
	}

	return attributeValueToDomain(&row, strategyID), nil
}

func (r *AttributeValueRepository) GetAttributeByID(ctx context.Context, id uuid.UUID) (attributevalue.AttributeRef, error) {
	q := strategyregistrydbsqlc.NewQuerier(ctx, r.db)
	row, err := q.GetStrategyAttributeByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return attributevalue.AttributeRef{}, apperr.NotFound[attributedomain.Attribute]("lookup", id)
	}
	if err != nil {
		return attributevalue.AttributeRef{}, fmt.Errorf("get attribute by id: %w", err)
	}

	return attributevalue.AttributeRef{
		ID:         row.ID,
		StrategyID: row.StrategyID,
		Slug:       row.Slug,
	}, nil
}

func (r *AttributeValueRepository) GetAttributeBySlug(ctx context.Context, strategyID uuid.UUID, slug string) (attributevalue.AttributeRef, error) {
	q := strategyregistrydbsqlc.NewQuerier(ctx, r.db)
	row, err := q.GetStrategyAttributeBySlug(ctx, strategyregistrydbsqlc.GetStrategyAttributeBySlugParams{
		StrategyID: strategyID,
		Slug:       strings.TrimSpace(slug),
	})
	if errors.Is(err, sql.ErrNoRows) {
		return attributevalue.AttributeRef{}, apperr.NotFound[attributedomain.Attribute]("lookup", slug)
	}
	if err != nil {
		return attributevalue.AttributeRef{}, fmt.Errorf("get attribute by slug: %w", err)
	}

	return attributevalue.AttributeRef{
		ID:         row.ID,
		StrategyID: row.StrategyID,
		Slug:       row.Slug,
	}, nil
}

func (r *AttributeValueRepository) List(ctx context.Context, filter attributevalue.ListFilter) ([]attributeValuedomain.AttributeValue, int64, error) {
	if filter.AttributeID == uuid.Nil {
		return nil, 0, apperr.NotFound[attributedomain.Attribute]("lookup", filter.AttributeID)
	}

	strategyID, err := r.lookupStrategyIDByAttributeID(ctx, filter.AttributeID)
	if err != nil {
		return nil, 0, fmt.Errorf("list attribute values: %w", err)
	}
	if strategyID != filter.StrategyID {
		return nil, 0, apperr.NotFound[attributedomain.Attribute]("lookup", filter.AttributeID)
	}

	search := strings.TrimSpace(filter.Search)

	limit, err := intToInt32(filter.PageSize)
	if err != nil {
		return nil, 0, apperr.New(apperr.KindInvalidArgument, "attributeValue.pageTooLarge", "page size out of range")
	}
	offset, err := intToInt32((filter.Page - 1) * filter.PageSize)
	if err != nil {
		return nil, 0, apperr.New(apperr.KindInvalidArgument, "attributeValue.pageTooLarge", "page out of range")
	}

	sortAsc := filter.Sort == attributevalue.AttributeValueSortCreatedAtAsc
	searchNull := sql.NullString{String: search, Valid: search != ""}

	q := strategyregistrydbsqlc.NewQuerier(ctx, r.db)

	items, err := q.ListStrategyAttributeValuesWithSearch(ctx, strategyregistrydbsqlc.ListStrategyAttributeValuesWithSearchParams{
		StrategyAttributeID: filter.AttributeID,
		Search:              searchNull,
		SortAsc:             sortAsc,
		PageLimit:           limit,
		PageOffset:          offset,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("list attribute values: %w", err)
	}

	total, err := q.CountStrategyAttributeValuesWithSearch(ctx, strategyregistrydbsqlc.CountStrategyAttributeValuesWithSearchParams{
		StrategyAttributeID: filter.AttributeID,
		Search:              searchNull,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("count attribute values: %w", err)
	}

	out := make([]attributeValuedomain.AttributeValue, len(items))
	for i := range items {
		out[i] = attributeValueToDomain(&items[i], strategyID)
	}
	return out, total, nil
}

func (r *AttributeValueRepository) Update(ctx context.Context, attributeValue *attributeValuedomain.AttributeValue) (attributeValuedomain.AttributeValue, error) {
	if attributeValue == nil {
		return attributeValuedomain.AttributeValue{}, fmt.Errorf("update attribute value: nil attribute value")
	}

	q := strategyregistrydbsqlc.NewQuerier(ctx, r.db)
	row, err := q.UpdateStrategyAttributeValue(ctx, attributeValueToUpdateParams(attributeValue))
	if errors.Is(err, sql.ErrNoRows) {
		return attributeValuedomain.AttributeValue{}, apperr.NotFound[attributeValuedomain.AttributeValue]("lookup", attributeValue.ID)
	}
	if err != nil {
		return attributeValuedomain.AttributeValue{}, fmt.Errorf("update attribute value: %w", mapAttributeValuePersistenceError(err))
	}

	strategyID, err := r.lookupStrategyIDByAttributeID(ctx, row.StrategyAttributeID)
	if err != nil {
		return attributeValuedomain.AttributeValue{}, fmt.Errorf("update attribute value: %w", err)
	}

	return attributeValueToDomain(&row, strategyID), nil
}

func (r *AttributeValueRepository) Delete(ctx context.Context, id uuid.UUID) error {
	q := strategyregistrydbsqlc.NewQuerier(ctx, r.db)
	affected, err := q.DeleteStrategyAttributeValue(ctx, id)
	if err != nil {
		return fmt.Errorf("delete attribute value: %w", err)
	}
	if affected == 0 {
		return apperr.NotFound[attributeValuedomain.AttributeValue]("lookup", id)
	}
	return nil
}

func (r *AttributeValueRepository) ReplaceRelations(ctx context.Context, input attributevalue.ReplaceAttributeValueRelationsInput) error {
	if input.StrategyID == uuid.Nil {
		return apperr.NotFound[strategydomain.Strategy]("lookup", input.StrategyID)
	}
	if input.FromAttributeID == uuid.Nil {
		return apperr.NotFound[attributedomain.Attribute]("lookup", input.FromAttributeID)
	}
	if input.FromValueID == uuid.Nil {
		return apperr.NotFound[attributeValuedomain.AttributeValue]("lookup", input.FromValueID)
	}

	if _, err := execContext(ctx, r.db, manualqueries.DeleteRelationsByFromValueSQL,
		input.StrategyID, input.FromAttributeID, input.FromValueID,
	); err != nil {
		return fmt.Errorf("replace relations: delete current relations: %w", err)
	}

	if len(input.Relations) == 0 {
		return nil
	}

	for i := range input.Relations {
		if _, err := execContext(ctx, r.db, manualqueries.InsertRelationSQL,
			input.FromAttributeID, input.FromValueID,
			input.Relations[i].ToAttributeID, input.Relations[i].ToValueID,
		); err != nil {
			return fmt.Errorf("replace relations: insert relation: %w", mapAttributeValuePersistenceError(err))
		}
	}

	return nil
}

func mapAttributeValuePersistenceError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}

	if pgErr.Code == pgUniqueViolationCode && isAttributeValueSlugUniqueViolation(pgErr) {
		return apperr.AlreadyExists[attributeValuedomain.AttributeValue]("create", pgErr.Detail)
	}
	if pgErr.Code == pgForeignKeyViolationCode && isAttributeValueRelationViolation(pgErr) {
		return apperr.New(apperr.KindNotFound, "relation.combinationNotFound", "relation combination not found")
	}
	if pgErr.Code == pgForeignKeyViolationCode && isAttributeForeignKeyViolation(pgErr) {
		return apperr.NotFound[attributedomain.Attribute]("lookup", pgErr.Detail)
	}

	return err
}

func isAttributeValueSlugUniqueViolation(pgErr *pgconn.PgError) bool {
	if pgErr == nil {
		return false
	}
	if pgErr.ConstraintName == attributeValueSlugUniqueConstraint {
		return true
	}
	details := strings.ToLower(pgErr.ConstraintName + " " + pgErr.Detail + " " + pgErr.Message)
	return strings.Contains(details, "slug")
}

func isAttributeValueRelationViolation(pgErr *pgconn.PgError) bool {
	if pgErr == nil {
		return false
	}
	details := strings.ToLower(pgErr.ConstraintName + " " + pgErr.Detail + " " + pgErr.Message)
	return strings.Contains(details, "strategy_attribute_value_relation") ||
		strings.Contains(details, "from_attribute_id") ||
		strings.Contains(details, "from_value_id") ||
		strings.Contains(details, "to_attribute_id") ||
		strings.Contains(details, "to_value_id")
}

func isAttributeForeignKeyViolation(pgErr *pgconn.PgError) bool {
	if pgErr == nil {
		return false
	}
	if pgErr.ConstraintName == attributeValueAttributeFKConstraint {
		return true
	}
	details := strings.ToLower(pgErr.ConstraintName + " " + pgErr.Detail + " " + pgErr.Message)
	return strings.Contains(details, "strategy_attribute_id")
}

func (r *AttributeValueRepository) lookupStrategyIDByAttributeID(ctx context.Context, attributeID uuid.UUID) (uuid.UUID, error) {
	q := strategyregistrydbsqlc.NewQuerier(ctx, r.db)
	row, err := q.GetStrategyAttributeByID(ctx, attributeID)
	if errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, apperr.NotFound[attributedomain.Attribute]("lookup", attributeID)
	}
	if err != nil {
		return uuid.Nil, fmt.Errorf("load attribute by id: %w", err)
	}
	return row.StrategyID, nil
}

// ListValueGraphByAttributeIDs loads attribute values together with their relations in one query.
func (r *AttributeValueRepository) ListValueGraphByAttributeIDs(ctx context.Context, strategyID uuid.UUID, attributeIDs []uuid.UUID) (map[uuid.UUID][]attributevalue.AttributeValueView, error) {
	result := make(map[uuid.UUID][]attributevalue.AttributeValueView, len(attributeIDs))
	if len(attributeIDs) == 0 {
		return result, nil
	}
	for i := range attributeIDs {
		result[attributeIDs[i]] = []attributevalue.AttributeValueView{}
	}

	rows, err := queryContext(ctx, r.db, manualqueries.ListAttributeValueGraphByAttributeIDsSQL, strategyID, attributeIDs)
	if err != nil {
		return nil, fmt.Errorf("list attribute value graph by attribute ids: %w", err)
	}
	defer rows.Close()

	valuesByAttributeID := make(map[uuid.UUID]map[uuid.UUID]int, len(attributeIDs))

	for rows.Next() {
		var (
			attributeID    uuid.UUID
			valueID        uuid.NullUUID
			valueSlug      sql.NullString
			valueValue     sql.NullString
			valueCreatedAt sql.NullTime
			valueUpdatedAt sql.NullTime
			fromAttrID     uuid.NullUUID
			fromValID      uuid.NullUUID
			toAttrID       uuid.NullUUID
			toValID        uuid.NullUUID
			toAttrSlug     sql.NullString
			toValSlug      sql.NullString
		)
		if err := rows.Scan(
			&attributeID,
			&valueID,
			&valueSlug,
			&valueValue,
			&valueCreatedAt,
			&valueUpdatedAt,
			&fromAttrID,
			&fromValID,
			&toAttrID,
			&toValID,
			&toAttrSlug,
			&toValSlug,
		); err != nil {
			return nil, fmt.Errorf("list attribute value graph by attribute ids: scan row: %w", err)
		}

		if !valueID.Valid {
			continue
		}

		indexByValueID, ok := valuesByAttributeID[attributeID]
		if !ok {
			indexByValueID = make(map[uuid.UUID]int)
			valuesByAttributeID[attributeID] = indexByValueID
		}

		valueIndex, ok := indexByValueID[valueID.UUID]
		if !ok {
			result[attributeID] = append(result[attributeID], attributevalue.AttributeValueView{
				ID:          valueID.UUID,
				StrategyID:  strategyID,
				AttributeID: attributeID,
				Slug:        valueSlug.String,
				Value:       valueValue.String,
				Relations:   []attributevalue.AttributeValueRelationView{},
				CreatedAt:   valueCreatedAt.Time,
				UpdatedAt:   valueUpdatedAt.Time,
			})
			valueIndex = len(result[attributeID]) - 1
			indexByValueID[valueID.UUID] = valueIndex
		}

		if !fromAttrID.Valid || !fromValID.Valid || !toAttrID.Valid || !toValID.Valid {
			continue
		}

		result[attributeID][valueIndex].Relations = append(result[attributeID][valueIndex].Relations, attributevalue.AttributeValueRelationView{
			FromAttributeID: fromAttrID.UUID,
			FromValueID:     fromValID.UUID,
			ToAttributeID:   toAttrID.UUID,
			ToValueID:       toValID.UUID,
			ToAttributeSlug: toAttrSlug.String,
			ToValueSlug:     toValSlug.String,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list attribute value graph by attribute ids: iterate rows: %w", err)
	}

	return result, nil
}

// ListRelationsByFromValueIDs — manual query, нужен JOIN для slug'ов.
func (r *AttributeValueRepository) ListRelationsByFromValueIDs(ctx context.Context, strategyID uuid.UUID, fromValueIDs []uuid.UUID) (map[uuid.UUID][]attributevalue.AttributeValueRelationView, error) {
	result := make(map[uuid.UUID][]attributevalue.AttributeValueRelationView, len(fromValueIDs))
	if len(fromValueIDs) == 0 {
		return result, nil
	}

	for i := range fromValueIDs {
		result[fromValueIDs[i]] = []attributevalue.AttributeValueRelationView{}
	}

	rows, err := queryContext(ctx, r.db, manualqueries.ListRelationsByFromValueIDsSQL, strategyID, fromValueIDs)
	if err != nil {
		return nil, fmt.Errorf("list attribute value relations by from value ids: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			fromAttrID uuid.UUID
			fromValID  uuid.UUID
			toAttrID   uuid.UUID
			toValID    uuid.UUID
			toAttrSlug sql.NullString
			toValSlug  sql.NullString
		)
		if err := rows.Scan(&fromAttrID, &fromValID, &toAttrID, &toValID, &toAttrSlug, &toValSlug); err != nil {
			return nil, fmt.Errorf("list attribute value relations by from value ids: scan row: %w", err)
		}

		result[fromValID] = append(result[fromValID], attributevalue.AttributeValueRelationView{
			FromAttributeID: fromAttrID,
			FromValueID:     fromValID,
			ToAttributeID:   toAttrID,
			ToValueID:       toValID,
			ToAttributeSlug: toAttrSlug.String,
			ToValueSlug:     toValSlug.String,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list attribute value relations by from value ids: iterate rows: %w", err)
	}

	return result, nil
}

// ListRelationsByFromValueIDsForAttributeValues конвертирует результаты ListRelationsByFromValueIDs
// в тип, используемый слоем application.
func (r *AttributeValueRepository) ListRelationsByFromValueIDsForAttributeValues(ctx context.Context, strategyID uuid.UUID, fromValueIDs []uuid.UUID) (map[uuid.UUID][]attributevalue.AttributeValueRelationView, error) {
	raw, err := r.ListRelationsByFromValueIDs(ctx, strategyID, fromValueIDs)
	if err != nil {
		return nil, err
	}

	result := make(map[uuid.UUID][]attributevalue.AttributeValueRelationView, len(raw))
	for valueID, relations := range raw {
		converted := make([]attributevalue.AttributeValueRelationView, len(relations))
		for i := range relations {
			converted[i] = attributevalue.AttributeValueRelationView{
				FromAttributeID: relations[i].FromAttributeID,
				FromValueID:     relations[i].FromValueID,
				ToAttributeID:   relations[i].ToAttributeID,
				ToValueID:       relations[i].ToValueID,
				ToAttributeSlug: relations[i].ToAttributeSlug,
				ToValueSlug:     relations[i].ToValueSlug,
			}
		}
		result[valueID] = converted
	}

	return result, nil
}
