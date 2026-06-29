package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
	strategyregistrydbsqlc "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/adapters/postgres/sqlc/gen"
	strategydomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"
	strategy "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/strategy"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

// Manual queries для batch INSERT (пока нет sqlc support)
const (
	insertCloneAttributeSQL = `
INSERT INTO strategy_attribute (
	id, strategy_id, slug, name, description, created_at, updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	insertCloneAttributeValueSQL = `
INSERT INTO strategy_attribute_value (
	id, strategy_attribute_id, slug, value, created_at, updated_at
) VALUES ($1, $2, $3, $4, $5, $6)`

	insertCloneRelationSQL = `
INSERT INTO strategy_attribute_value_relation (
	from_attribute_id, from_value_id, to_attribute_id, to_value_id, created_at, updated_at
) VALUES ($1, $2, $3, $4, $5, $6)`
)

const (
	strategySlugUniqueConstraint = "strategy_slug_key"
)

type StrategyRepository struct {
	db *sql.DB
}

func NewStrategyRepository(db *sql.DB) *StrategyRepository {
	return &StrategyRepository{db: db}
}

var _ strategy.StrategyRepository = (*StrategyRepository)(nil)
var _ strategy.StrategyCloneRepository = (*StrategyRepository)(nil)

func (r *StrategyRepository) Create(ctx context.Context, item *strategydomain.Strategy) (strategydomain.Strategy, error) {
	if item == nil {
		return strategydomain.Strategy{}, fmt.Errorf("create strategy: nil strategy")
	}
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}

	q := strategyregistrydbsqlc.NewQuerier(ctx, r.db)
	row, err := q.CreateStrategy(ctx, strategyToCreateParams(item))
	if err != nil {
		return strategydomain.Strategy{}, fmt.Errorf("create strategy: %w", mapStrategyPersistenceError(err, item.Slug))
	}

	return strategyToDomain(&row), nil
}

func (r *StrategyRepository) GetByID(ctx context.Context, id uuid.UUID) (strategydomain.Strategy, error) {
	q := strategyregistrydbsqlc.NewQuerier(ctx, r.db)
	row, err := q.GetStrategyByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return strategydomain.Strategy{}, apperr.NotFound[strategydomain.Strategy]("lookup", id)
	}
	if err != nil {
		return strategydomain.Strategy{}, fmt.Errorf("get strategy by id: %w", err)
	}
	return strategyToDomain(&row), nil
}

func (r *StrategyRepository) GetBySlug(ctx context.Context, slug string) (strategydomain.Strategy, error) {
	q := strategyregistrydbsqlc.NewQuerier(ctx, r.db)
	row, err := q.GetStrategyBySlug(ctx, strings.TrimSpace(slug))
	if errors.Is(err, sql.ErrNoRows) {
		return strategydomain.Strategy{}, apperr.NotFound[strategydomain.Strategy]("lookup", slug)
	}
	if err != nil {
		return strategydomain.Strategy{}, fmt.Errorf("get strategy by slug: %w", err)
	}
	return strategyToDomain(&row), nil
}

func (r *StrategyRepository) List(ctx context.Context, filter strategy.ListFilter) ([]strategydomain.Strategy, int64, error) {
	search := strings.TrimSpace(filter.Search)

	limit, err := intToInt32(filter.PageSize)
	if err != nil {
		return nil, 0, apperr.New(apperr.KindInvalidArgument, "strategy.pageTooLarge", "page size out of range")
	}
	offset, err := intToInt32((filter.Page - 1) * filter.PageSize)
	if err != nil {
		return nil, 0, apperr.New(apperr.KindInvalidArgument, "strategy.pageTooLarge", "page out of range")
	}

	sortAsc := filter.Sort == strategydomain.StrategySortCreatedAtAsc
	searchNull := sql.NullString{String: search, Valid: search != ""}

	q := strategyregistrydbsqlc.NewQuerier(ctx, r.db)

	items, err := q.ListStrategiesWithSearch(ctx, strategyregistrydbsqlc.ListStrategiesWithSearchParams{
		Search:     searchNull,
		SortAsc:    sortAsc,
		PageLimit:  limit,
		PageOffset: offset,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("list strategies: %w", err)
	}

	total, err := q.CountStrategiesWithSearch(ctx, searchNull)
	if err != nil {
		return nil, 0, fmt.Errorf("count strategies: %w", err)
	}

	out := make([]strategydomain.Strategy, len(items))
	for i := range items {
		out[i] = strategyToDomain(&items[i])
	}
	return out, total, nil
}

func (r *StrategyRepository) Update(ctx context.Context, item *strategydomain.Strategy) (strategydomain.Strategy, error) {
	if item == nil {
		return strategydomain.Strategy{}, fmt.Errorf("update strategy: nil strategy")
	}

	q := strategyregistrydbsqlc.NewQuerier(ctx, r.db)
	row, err := q.UpdateStrategy(ctx, strategyToUpdateParams(item))
	if errors.Is(err, sql.ErrNoRows) {
		return strategydomain.Strategy{}, apperr.NotFound[strategydomain.Strategy]("lookup", item.ID)
	}
	if err != nil {
		return strategydomain.Strategy{}, fmt.Errorf("update strategy: %w", mapStrategyPersistenceError(err, item.Slug))
	}

	return strategyToDomain(&row), nil
}

func (r *StrategyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	q := strategyregistrydbsqlc.NewQuerier(ctx, r.db)
	affected, err := q.DeleteStrategy(ctx, id)
	if err != nil {
		return fmt.Errorf("delete strategy: %w", err)
	}
	if affected == 0 {
		return apperr.NotFound[strategydomain.Strategy]("lookup", id)
	}
	return nil
}

func (r *StrategyRepository) CloneBatch(ctx context.Context, items []strategy.CloneStrategySpec) ([]strategydomain.Strategy, error) {
	if len(items) == 0 {
		return []strategydomain.Strategy{}, nil
	}

	cloned := make([]strategydomain.Strategy, 0, len(items))
	for i := range items {
		clonedStrategy, err := r.cloneOne(ctx, items[i])
		if err != nil {
			return nil, fmt.Errorf("clone strategy[%d]: %w", i, err)
		}
		cloned = append(cloned, clonedStrategy)
	}

	return cloned, nil
}

// =============================================================================
// Clone internals
// =============================================================================

type cloneSourceAttribute struct {
	ID          uuid.UUID
	Slug        string
	Name        string
	Description string
}

type cloneSourceValue struct {
	ID          uuid.UUID
	AttributeID uuid.UUID
	Slug        string
	Value       string
}

type cloneSourceRelation struct {
	FromAttributeID uuid.UUID
	FromValueID     uuid.UUID
	ToAttributeID   uuid.UUID
	ToValueID       uuid.UUID
}

func (r *StrategyRepository) cloneOne(ctx context.Context, spec strategy.CloneStrategySpec) (strategydomain.Strategy, error) {
	source, err := r.GetByID(ctx, spec.SourceStrategyID)
	if err != nil {
		return strategydomain.Strategy{}, err
	}

	now := timeNowUTC()
	clonedID := uuid.New()

	q := strategyregistrydbsqlc.NewQuerier(ctx, r.db)
	inserted, err := q.CreateStrategy(ctx, strategyregistrydbsqlc.CreateStrategyParams{
		ID:          clonedID,
		Slug:        strings.TrimSpace(spec.Slug),
		Name:        source.Name,
		Description: source.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		return strategydomain.Strategy{}, mapStrategyPersistenceError(err, spec.Slug)
	}

	// Шаг 1: клонируем атрибуты
	attrIDMap, err := r.cloneAttributes(ctx, source.ID, clonedID, now)
	if err != nil {
		return strategydomain.Strategy{}, err
	}
	if len(attrIDMap) == 0 {
		return strategyToDomain(&inserted), nil
	}

	// Шаг 2: клонируем значения
	valIDMap, err := r.cloneValues(ctx, attrIDMap, now)
	if err != nil {
		return strategydomain.Strategy{}, err
	}
	if len(valIDMap) == 0 {
		return strategyToDomain(&inserted), nil
	}

	// Шаг 3: клонируем связи
	if err := r.cloneRelations(ctx, valIDMap, attrIDMap, now); err != nil {
		return strategydomain.Strategy{}, err
	}

	return strategyToDomain(&inserted), nil
}

// cloneAttributes копирует атрибуты из source в target и возвращает map oldID → newID.
func (r *StrategyRepository) cloneAttributes(
	ctx context.Context,
	sourceStrategyID, targetStrategyID uuid.UUID,
	now time.Time,
) (map[uuid.UUID]uuid.UUID, error) {
	sources, err := r.listCloneSourceAttributes(ctx, sourceStrategyID)
	if err != nil {
		return nil, err
	}
	if len(sources) == 0 {
		return nil, nil
	}

	idMap := make(map[uuid.UUID]uuid.UUID, len(sources))
	for i := range sources {
		newID := uuid.New()
		idMap[sources[i].ID] = newID

		if _, err := execContext(ctx, r.db, insertCloneAttributeSQL,
			newID, targetStrategyID,
			sources[i].Slug, sources[i].Name, sources[i].Description,
			now, now,
		); err != nil {
			return nil, fmt.Errorf("clone attribute[%d]: %w", i, err)
		}
	}
	return idMap, nil
}

// cloneValues копирует значения атрибутов и возвращает map oldID → newID.
func (r *StrategyRepository) cloneValues(
	ctx context.Context,
	attrIDMap map[uuid.UUID]uuid.UUID,
	now time.Time,
) (map[uuid.UUID]uuid.UUID, error) {
	sourceAttrIDs := make([]uuid.UUID, 0, len(attrIDMap))
	for oldID := range attrIDMap {
		sourceAttrIDs = append(sourceAttrIDs, oldID)
	}

	sources, err := r.listCloneSourceValues(ctx, sourceAttrIDs)
	if err != nil {
		return nil, err
	}
	if len(sources) == 0 {
		return nil, nil
	}

	idMap := make(map[uuid.UUID]uuid.UUID, len(sources))
	for i := range sources {
		newAttrID, ok := attrIDMap[sources[i].AttributeID]
		if !ok {
			return nil, fmt.Errorf("missing cloned attribute id for source attribute %s", sources[i].AttributeID)
		}

		newID := uuid.New()
		idMap[sources[i].ID] = newID

		if _, err := execContext(ctx, r.db, insertCloneAttributeValueSQL,
			newID, newAttrID,
			sources[i].Slug, sources[i].Value,
			now, now,
		); err != nil {
			return nil, fmt.Errorf("clone attribute value[%d]: %w", i, err)
		}
	}
	return idMap, nil
}

// cloneRelations копирует связи между значениями атрибутов.
func (r *StrategyRepository) cloneRelations(
	ctx context.Context,
	valIDMap, attrIDMap map[uuid.UUID]uuid.UUID,
	now time.Time,
) error {
	sourceValueIDs := make([]uuid.UUID, 0, len(valIDMap))
	for oldID := range valIDMap {
		sourceValueIDs = append(sourceValueIDs, oldID)
	}

	sources, err := r.listCloneSourceRelations(ctx, sourceValueIDs)
	if err != nil {
		return err
	}

	for i := range sources {
		rel := sources[i]

		newFromAttrID, ok := attrIDMap[rel.FromAttributeID]
		if !ok {
			return fmt.Errorf("missing cloned from_attribute id for source attribute %s", rel.FromAttributeID)
		}
		newToAttrID, ok := attrIDMap[rel.ToAttributeID]
		if !ok {
			return fmt.Errorf("missing cloned to_attribute id for source attribute %s", rel.ToAttributeID)
		}
		newFromValID, ok := valIDMap[rel.FromValueID]
		if !ok {
			return fmt.Errorf("missing cloned from_value id for source value %s", rel.FromValueID)
		}
		newToValID, ok := valIDMap[rel.ToValueID]
		if !ok {
			return fmt.Errorf("missing cloned to_value id for source value %s", rel.ToValueID)
		}

		if _, err := execContext(ctx, r.db, insertCloneRelationSQL,
			newFromAttrID, newFromValID,
			newToAttrID, newToValID,
			now, now,
		); err != nil {
			return fmt.Errorf("clone relation[%d]: %w", i, err)
		}
	}
	return nil
}

// =============================================================================
// Clone source readers (sqlc)
// =============================================================================

func (r *StrategyRepository) listCloneSourceAttributes(ctx context.Context, sourceStrategyID uuid.UUID) ([]cloneSourceAttribute, error) {
	q := strategyregistrydbsqlc.NewQuerier(ctx, r.db)

	items, err := q.ListStrategyAttributesByStrategyIDAsc(ctx, strategyregistrydbsqlc.ListStrategyAttributesByStrategyIDAscParams{
		StrategyID: sourceStrategyID,
		Limit:      maxCloneBatchSize,
		Offset:     0,
	})
	if err != nil {
		return nil, fmt.Errorf("list source attributes: %w", err)
	}

	out := make([]cloneSourceAttribute, len(items))
	for i := range items {
		out[i] = cloneSourceAttribute{
			ID:          items[i].ID,
			Slug:        items[i].Slug,
			Name:        items[i].Name,
			Description: items[i].Description,
		}
	}
	return out, nil
}

func (r *StrategyRepository) listCloneSourceValues(ctx context.Context, sourceAttributeIDs []uuid.UUID) ([]cloneSourceValue, error) {
	q := strategyregistrydbsqlc.NewQuerier(ctx, r.db)

	items, err := q.ListStrategyAttributeValuesByStrategyAttributeIDArray(ctx, strategyregistrydbsqlc.ListStrategyAttributeValuesByStrategyAttributeIDArrayParams{
		Column1: sourceAttributeIDs,
		Limit:   maxCloneBatchSize,
		Offset:  0,
	})
	if err != nil {
		return nil, fmt.Errorf("list source attribute values: %w", err)
	}

	out := make([]cloneSourceValue, len(items))
	for i := range items {
		out[i] = cloneSourceValue{
			ID:          items[i].ID,
			AttributeID: items[i].StrategyAttributeID,
			Slug:        items[i].Slug,
			Value:       items[i].Value,
		}
	}
	return out, nil
}

func (r *StrategyRepository) listCloneSourceRelations(ctx context.Context, sourceValueIDs []uuid.UUID) ([]cloneSourceRelation, error) {
	q := strategyregistrydbsqlc.NewQuerier(ctx, r.db)

	items, err := q.ListStrategyAttributeValueRelationsByFromValueIDArray(ctx, strategyregistrydbsqlc.ListStrategyAttributeValueRelationsByFromValueIDArrayParams{
		Column1: sourceValueIDs,
		Limit:   maxCloneBatchSize,
		Offset:  0,
	})
	if err != nil {
		return nil, fmt.Errorf("list source attribute value relations: %w", err)
	}

	out := make([]cloneSourceRelation, len(items))
	for i := range items {
		out[i] = cloneSourceRelation{
			FromAttributeID: items[i].FromAttributeID,
			FromValueID:     items[i].FromValueID,
			ToAttributeID:   items[i].ToAttributeID,
			ToValueID:       items[i].ToValueID,
		}
	}
	return out, nil
}

// =============================================================================
// Error mapping
// =============================================================================

func mapStrategyPersistenceError(err error, slug string) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}

	if pgErr.Code == pgUniqueViolationCode && isStrategySlugUniqueViolation(pgErr) {
		return apperr.AlreadyExists[strategydomain.Strategy]("create", slug)
	}

	return err
}

func isStrategySlugUniqueViolation(pgErr *pgconn.PgError) bool {
	if pgErr == nil {
		return false
	}
	if pgErr.ConstraintName == strategySlugUniqueConstraint {
		return true
	}

	details := strings.ToLower(pgErr.ConstraintName + " " + pgErr.Detail + " " + pgErr.Message)
	return strings.Contains(details, "slug")
}
