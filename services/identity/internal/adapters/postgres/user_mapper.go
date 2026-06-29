package db

import (
	"database/sql"
	"strings"
	"time"

	identitydbsqlc "github.com/stratflow-labs/stratflow/services/identity/internal/adapters/postgres/sqlc/gen"
	identitydomain "github.com/stratflow-labs/stratflow/services/identity/internal/domain"
)

func userToDomain(model *identitydbsqlc.User) identitydomain.User {
	if model == nil {
		return identitydomain.User{}
	}
	return identitydomain.User{
		ID:              model.ID,
		Login:           model.Login,
		Name:            model.Name,
		LastName:        model.LastName,
		Email:           model.Email.String,
		PasswordHash:    model.Password,
		Role:            model.Role,
		ImageUrl:        model.ImageUrl,
		Gender:          genderDBToInt(model.Gender),
		IsEmailVerified: model.IsEmailVerified,
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
	}
}

func userToCreateParams(user *identitydomain.User, now time.Time) identitydbsqlc.CreateUserParams {
	email := sql.NullString{Valid: false}
	if v := strings.TrimSpace(user.Email); v != "" {
		email = sql.NullString{String: v, Valid: true}
	}
	return identitydbsqlc.CreateUserParams{
		ID:              user.ID,
		Name:            user.Name,
		LastName:        user.LastName,
		Login:           user.Login,
		Email:           email,
		Password:        user.PasswordHash,
		Role:            user.Role,
		ImageUrl:        user.ImageUrl,
		Gender:          userGenderIntToDB(user.Gender),
		IsEmailVerified: user.IsEmailVerified,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func userToUpdateParams(user *identitydomain.User, updatedAt time.Time) identitydbsqlc.UpdateUserParams {
	email := sql.NullString{Valid: false}
	if v := strings.TrimSpace(user.Email); v != "" {
		email = sql.NullString{String: v, Valid: true}
	}

	return identitydbsqlc.UpdateUserParams{
		ID:              user.ID,
		Name:            user.Name,
		LastName:        user.LastName,
		Login:           user.Login,
		Email:           email,
		Password:        user.PasswordHash,
		Role:            user.Role,
		ImageUrl:        user.ImageUrl,
		Gender:          userGenderIntToDB(user.Gender),
		IsEmailVerified: user.IsEmailVerified,
		UpdatedAt:       updatedAt,
	}
}

func userGenderIntToDB(v int) sql.NullString {
	switch v {
	case 1:
		return sql.NullString{String: "male", Valid: true}
	case 2:
		return sql.NullString{String: "female", Valid: true}
	default:
		return sql.NullString{Valid: false}
	}
}

func genderDBToInt(v sql.NullString) int {
	if !v.Valid {
		return 0
	}
	switch strings.ToLower(strings.TrimSpace(v.String)) {
	case "male":
		return 1
	case "female":
		return 2
	default:
		return 0
	}
}
