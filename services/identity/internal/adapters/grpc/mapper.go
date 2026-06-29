package identitygrpc

import (
	"strings"

	identityv1 "github.com/stratflow-labs/stratflow/services/identity/gen/go/proto/v1"
	identitydomain "github.com/stratflow-labs/stratflow/services/identity/internal/domain"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func mapRole(role string) identityv1.Role {
	switch role {
	case "user":
		return identityv1.Role_ROLE_USER
	case "manager":
		return identityv1.Role_ROLE_MANAGER
	case "admin":
		return identityv1.Role_ROLE_ADMIN
	default:
		return identityv1.Role_ROLE_UNSPECIFIED
	}
}

func mapGender(value int) *int32 {
	if value == 0 {
		return nil
	}

	gender := int32(value)
	return &gender
}

func mapUser(user identitydomain.User) *identityv1.User {
	return &identityv1.User{
		Id:              user.ID.String(),
		Login:           user.Login,
		Name:            user.Name,
		LastName:        user.LastName,
		Email:           optionalString(user.Email),
		Role:            user.Role,
		Gender:          mapGender(user.Gender),
		ImageUrl:        optionalString(user.ImageUrl),
		IsEmailVerified: user.IsEmailVerified,
		IsVerified:      user.IsEmailVerified,
		CreatedAt:       timestamppb.New(user.CreatedAt),
		UpdatedAt:       timestamppb.New(user.UpdatedAt),
	}
}

func optionalString(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
