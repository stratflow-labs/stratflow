package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/stratflow-labs/stratflow/internal/foundation/apperr"
	strategyregistrydbsqlc "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/adapters/postgres/sqlc/gen"
	attribute "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/attribute"
	attributedomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"
	strategydomain "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	attributeSlugUniqueConstraint = "uq_strategy_param_strategy_slug"
	attributeStrategyFKConstraint = "strategy_attribute_strategy_id_fkey"
)

type AttributeRepository struct {
	db *sql.DB
}

func NewAttributeRepository(db *sql.DB) *AttributeRepository {
	return &AttributeRepository{db: db}
}

var _ attribute.AttributeRepository = (*AttributeRepository)(nil)

func (r *AttributeRepository) Create(ctx context.Context, item *attributedomain.Attribute) (attributedomain.Attribute, error) {
	if item == nil {
		return attributedomain.Attribute{}, fmt.Errorf("create attribute: nil attribute")
	}
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}

	q := strategyregistrydbsqlc.NewQuerier(ctx, r.db)
	row, err := q.CreateStrategyAttribute(ctx, attributeToCreateParams(item))
	if err != nil {
		return attributedomain.Attribute{}, fmt.Errorf("create attribute: %w", mapAttributePersistenceError(err))
	}

	return attributeToDomain(&row), nil
}

func (r *AttributeRepository) GetByID(ctx context.Context, id uuid.UUID) (attributedomain.Attribute, error) {
	q := strategyregistrydbsqlc.NewQuerier(ctx, r.db)
	row, err := q.GetStrategyAttributeByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return attributedomain.Attribute{}, apperr.NotFound[attributedomain.Attribute]("lookup", id)
	}
	if err != nil {
		return attributedomain.Attribute{}, fmt.Errorf("get attribute by id: %w", err)
	}

	return attributeToDomain(&row), nil
}

func (r *AttributeRepository) GetBySlug(ctx context.Context, strategyID uuid.UUID, slug string) (attributedomain.Attribute, error) {
	q := strategyregistrydbsqlc.NewQuerier(ctx, r.db)
	row, err := q.GetStrategyAttributeBySlug(ctx, strategyregistrydbsqlc.GetStrategyAttributeBySlugParams{
		StrategyID: strategyID,
		Slug:       strings.TrimSpace(slug),
	})
	if errors.Is(err, sql.ErrNoRows) {
		return attributedomain.Attribute{}, apperr.NotFound[attributedomain.Attribute]("lookup", slug)
	}
	if err != nil {
		return attributedomain.Attribute{}, fmt.Errorf("get attribute by slug: %w", err)
	}

	return attributeToDomain(&row), nil
}

func (r *AttributeRepository) List(ctx context.Context, filter attribute.ListFilter) ([]attributedomain.Attribute, int64, error) {
	search := strings.TrimSpace(filter.Search)

	limit, err := intToInt32(filter.PageSize)
	if err != nil {
		return nil, 0, apperr.New(apperr.KindInvalidArgument, "attribute.pageTooLarge", "page size out of range")
	}
	offset, err := intToInt32((filter.Page - 1) * filter.PageSize)
	if err != nil {
		return nil, 0, apperr.New(apperr.KindInvalidArgument, "attribute.pageTooLarge", "page out of range")
	}

	sortAsc := filter.Sort == attribute.AttributeSortCreatedAtAsc
	searchNull := sql.NullString{String: search, Valid: search != ""}

	q := strategyregistrydbsqlc.NewQuerier(ctx, r.db)

	items, err := q.ListStrategyAttributesWithSearch(ctx, strategyregistrydbsqlc.ListStrategyAttributesWithSearchParams{
		StrategyID: filter.StrategyID,
		Search:     searchNull,
		SortAsc:    sortAsc,
		PageLimit:  limit,
		PageOffset: offset,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("list attributes: %w", err)
	}

	total, err := q.CountStrategyAttributesWithSearch(ctx, strategyregistrydbsqlc.CountStrategyAttributesWithSearchParams{
		StrategyID: filter.StrategyID,
		Search:     searchNull,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("count attributes: %w", err)
	}

	out := make([]attributedomain.Attribute, len(items))
	for i := range items {
		out[i] = attributeToDomain(&items[i])
	}
	return out, total, nil
}

func (r *AttributeRepository) Update(ctx context.Context, item *attributedomain.Attribute) (attributedomain.Attribute, error) {
	if item == nil {
		return attributedomain.Attribute{}, fmt.Errorf("update attribute: nil attribute")
	}

	q := strategyregistrydbsqlc.NewQuerier(ctx, r.db)
	row, err := q.UpdateStrategyAttribute(ctx, attributeToUpdateParams(item))
	if errors.Is(err, sql.ErrNoRows) {
		return attributedomain.Attribute{}, apperr.NotFound[attributedomain.Attribute]("lookup", item.ID)
	}
	if err != nil {
		return attributedomain.Attribute{}, fmt.Errorf("update attribute: %w", mapAttributePersistenceError(err))
	}

	return attributeToDomain(&row), nil
}

func (r *AttributeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	q := strategyregistrydbsqlc.NewQuerier(ctx, r.db)
	affected, err := q.DeleteStrategyAttribute(ctx, id)
	if err != nil {
		return fmt.Errorf("delete attribute: %w", err)
	}
	if affected == 0 {
		return apperr.NotFound[attributedomain.Attribute]("lookup", id)
	}
	return nil
}

func mapAttributePersistenceError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}

	if pgErr.Code == pgUniqueViolationCode && isAttributeSlugUniqueViolation(pgErr) {
		return apperr.AlreadyExists[attributedomain.Attribute]("create", pgErr.Detail)
	}
	if pgErr.Code == pgForeignKeyViolationCode && isAttributeStrategyForeignKeyViolation(pgErr) {
		return apperr.NotFound[strategydomain.Strategy]("lookup", pgErr.Detail)
	}

	return err
}

func isAttributeSlugUniqueViolation(pgErr *pgconn.PgError) bool {
	if pgErr == nil {
		return false
	}
	if pgErr.ConstraintName == attributeSlugUniqueConstraint {
		return true
	}

	details := strings.ToLower(pgErr.ConstraintName + " " + pgErr.Detail + " " + pgErr.Message)
	return strings.Contains(details, "slug")
}

func isAttributeStrategyForeignKeyViolation(pgErr *pgconn.PgError) bool {
	if pgErr == nil {
		return false
	}
	if pgErr.ConstraintName == attributeStrategyFKConstraint {
		return true
	}

	details := strings.ToLower(pgErr.ConstraintName + " " + pgErr.Detail + " " + pgErr.Message)
	return strings.Contains(details, "strategy_id")
}
