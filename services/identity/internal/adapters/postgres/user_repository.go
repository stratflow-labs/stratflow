package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	identitydbsqlc "github.com/stratflow-labs/stratflow/services/identity/internal/adapters/postgres/sqlc/gen"
	identitydomain "github.com/stratflow-labs/stratflow/services/identity/internal/domain"
	usecases "github.com/stratflow-labs/stratflow/services/identity/internal/user"

	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

var _ usecases.UserRepository = (*UserRepository)(nil)

func (r *UserRepository) Create(ctx context.Context, user *identitydomain.User) (identitydomain.User, error) {
	if user == nil {
		return identitydomain.User{}, fmt.Errorf("create user: nil user")
	}
	now := time.Now().UTC()
	q := identitydbsqlc.NewQuerier(ctx, r.db)
	params := userToCreateParams(user, now)
	row, err := q.CreateUser(ctx, params)
	if err != nil {
		return identitydomain.User{}, err
	}
	return userToDomain(&row), nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (identitydomain.User, error) {
	q := identitydbsqlc.NewQuerier(ctx, r.db)
	row, err := q.GetUserByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return identitydomain.User{}, identitydomain.ErrUserNotFound
	}
	if err != nil {
		return identitydomain.User{}, err
	}
	return userToDomain(&row), nil
}

func (r *UserRepository) List(ctx context.Context, filter usecases.ListFilter) ([]identitydomain.User, int64, error) {
	search := strings.TrimSpace(filter.Search)
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * pageSize
	order, ok := sanitizeSort(filter.Sort)
	if !ok {
		order = "created_at DESC"
	}

	limitInt32, err := intToInt32(pageSize)
	if err != nil {
		return nil, 0, err
	}
	offsetInt32, err := intToInt32(offset)
	if err != nil {
		return nil, 0, err
	}

	rows, total, err := r.listWithTotal(ctx, search, limitInt32, offsetInt32, order)
	if err != nil {
		return nil, 0, err
	}

	out := make([]identitydomain.User, len(rows))
	for i := range rows {
		out[i] = userToDomain(&rows[i])
	}
	return out, total, nil
}

func (r *UserRepository) listWithTotal(ctx context.Context, search string, limit, offset int32, order string) ([]identitydbsqlc.User, int64, error) {
	// For PostgreSQL 12+, windowed count can reduce queries; for clarity using two queries.
	query := userListQuery(order)
	rows, err := r.db.QueryContext(ctx, query, search, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var users []identitydbsqlc.User
	for rows.Next() {
		var u identitydbsqlc.User
		if err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.LastName,
			&u.Login,
			&u.Email,
			&u.Password,
			&u.Role,
			&u.ImageUrl,
			&u.Gender,
			&u.IsEmailVerified,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	total, err := r.countWithSearch(ctx, search)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *UserRepository) countWithSearch(ctx context.Context, search string) (int64, error) {
	row := r.db.QueryRowContext(ctx, userCountQuery, search)
	var total int64
	if err := row.Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

// sanitizeSort validates and maps user input to safe SQL ORDER BY clause.
// Returns the sanitized ORDER BY clause and a boolean indicating success.
// Only whitelisted columns are allowed to prevent SQL injection.
func sanitizeSort(sort string) (string, bool) {
	s := strings.TrimSpace(sort)
	if s == "" {
		return "", false
	}

	// Whitelist of allowed sort fields mapped to actual column names
	allowedFields := map[string]string{
		"id":         "id",
		"name":       "name",
		"lastname":   "last_name",
		"last_name":  "last_name",
		"email":      "email",
		"role":       "role",
		"created_at": "created_at",
		"updated_at": "updated_at",
		"gender":     "gender",
	}

	// Whitelist of allowed sort directions
	allowedDirections := map[string]bool{
		"ASC":  true,
		"DESC": true,
	}

	direction := "ASC"
	field := s

	// Parse field and direction from input
	parts := strings.Fields(s)
	if len(parts) > 0 {
		field = parts[0]
		if len(parts) > 1 {
			dir := strings.ToUpper(parts[1])
			if !allowedDirections[dir] {
				return "", false
			}
			direction = dir
		}
	}

	// Support "-field" syntax for DESC ordering
	if strings.HasPrefix(field, "-") {
		direction = "DESC"
		field = strings.TrimPrefix(field, "-")
	}

	field = strings.ToLower(field)

	// Validate field against whitelist
	columnName, ok := allowedFields[field]
	if !ok {
		return "", false
	}

	// Safe to use in SQL as both field and direction are whitelisted
	return fmt.Sprintf("%s %s", columnName, direction), true
}

func userListQuery(order string) string {
	if q, ok := userListQueries[order]; ok {
		return q
	}
	return userListQueryCreatedAtDesc
}

const userListQueryBase = `
SELECT
    id, name, last_name, login, email, password, role, image_url, gender,
    is_email_verified, created_at, updated_at
FROM users
WHERE (
    $1 = '' OR
    name ILIKE '%' || $1 || '%' OR
    last_name ILIKE '%' || $1 || '%' OR
    email ILIKE '%' || $1 || '%'
)
ORDER BY `

const userListQuerySuffix = `
LIMIT $2 OFFSET $3`

const userCountQuery = `
SELECT COUNT(*) AS total
FROM users
WHERE (
    $1 = '' OR
    name ILIKE '%' || $1 || '%' OR
    last_name ILIKE '%' || $1 || '%' OR
    email ILIKE '%' || $1 || '%'
)`

const userListQueryIDAsc = userListQueryBase + "id ASC" + userListQuerySuffix
const userListQueryIDDesc = userListQueryBase + "id DESC" + userListQuerySuffix
const userListQueryNameAsc = userListQueryBase + "name ASC" + userListQuerySuffix
const userListQueryNameDesc = userListQueryBase + "name DESC" + userListQuerySuffix
const userListQueryLastNameAsc = userListQueryBase + "last_name ASC" + userListQuerySuffix
const userListQueryLastNameDesc = userListQueryBase + "last_name DESC" + userListQuerySuffix
const userListQueryEmailAsc = userListQueryBase + "email ASC" + userListQuerySuffix
const userListQueryEmailDesc = userListQueryBase + "email DESC" + userListQuerySuffix
const userListQueryRoleAsc = userListQueryBase + "role ASC" + userListQuerySuffix
const userListQueryRoleDesc = userListQueryBase + "role DESC" + userListQuerySuffix
const userListQueryCreatedAtAsc = userListQueryBase + "created_at ASC" + userListQuerySuffix
const userListQueryCreatedAtDesc = userListQueryBase + "created_at DESC" + userListQuerySuffix
const userListQueryUpdatedAtAsc = userListQueryBase + "updated_at ASC" + userListQuerySuffix
const userListQueryUpdatedAtDesc = userListQueryBase + "updated_at DESC" + userListQuerySuffix
const userListQueryGenderAsc = userListQueryBase + "gender ASC" + userListQuerySuffix
const userListQueryGenderDesc = userListQueryBase + "gender DESC" + userListQuerySuffix

var userListQueries = map[string]string{
	"id ASC":          userListQueryIDAsc,
	"id DESC":         userListQueryIDDesc,
	"name ASC":        userListQueryNameAsc,
	"name DESC":       userListQueryNameDesc,
	"last_name ASC":   userListQueryLastNameAsc,
	"last_name DESC":  userListQueryLastNameDesc,
	"email ASC":       userListQueryEmailAsc,
	"email DESC":      userListQueryEmailDesc,
	"role ASC":        userListQueryRoleAsc,
	"role DESC":       userListQueryRoleDesc,
	"created_at ASC":  userListQueryCreatedAtAsc,
	"created_at DESC": userListQueryCreatedAtDesc,
	"updated_at ASC":  userListQueryUpdatedAtAsc,
	"updated_at DESC": userListQueryUpdatedAtDesc,
	"gender ASC":      userListQueryGenderAsc,
	"gender DESC":     userListQueryGenderDesc,
}

func intToInt32(v int) (int32, error) {
	if v > math.MaxInt32 || v < math.MinInt32 {
		return 0, fmt.Errorf("value %d is out of int32 range", v)
	}
	return int32(v), nil
}

func (r *UserRepository) Update(ctx context.Context, user *identitydomain.User) (identitydomain.User, error) {
	if user == nil {
		return identitydomain.User{}, fmt.Errorf("update user: nil user")
	}
	q := identitydbsqlc.NewQuerier(ctx, r.db)
	params := userToUpdateParams(user, time.Now().UTC())
	row, err := q.UpdateUser(ctx, params)
	if errors.Is(err, sql.ErrNoRows) {
		return identitydomain.User{}, identitydomain.ErrUserNotFound
	}
	if err != nil {
		return identitydomain.User{}, err
	}
	return userToDomain(&row), nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	rows, err := identitydbsqlc.NewQuerier(ctx, r.db).DeleteUser(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return identitydomain.ErrUserNotFound
	}
	return nil
}
