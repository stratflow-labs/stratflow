package identitygrpc

import (
	"errors"
	"net/http"

	"github.com/stratflow-labs/stratflow/internal/grpcserver"
	identitydomain "github.com/stratflow-labs/stratflow/services/identity/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapError(err error) error {
	switch {
	case errors.Is(err, identitydomain.ErrAccessTokenNotFound):
		return grpcserver.NewStatusError(codes.Unauthenticated, "auth.accessTokenNotFound", "access token not found")
	case errors.Is(err, identitydomain.ErrAccessTokenInvalid):
		return grpcserver.NewStatusError(codes.Unauthenticated, "auth.accessTokenInvalid", "access token is invalid")
	case errors.Is(err, identitydomain.ErrInvalidCredentials), errors.Is(err, identitydomain.ErrPasswordMismatch):
		return grpcserver.NewStatusError(codes.Unauthenticated, "auth.invalidCredentials", "invalid email or password")
	case errors.Is(err, identitydomain.ErrLoginRequired):
		return grpcserver.NewStatusError(codes.InvalidArgument, "auth.loginRequired", "login is required")
	case errors.Is(err, identitydomain.ErrLoginInvalid):
		return grpcserver.NewStatusError(codes.InvalidArgument, "auth.loginInvalid", "login format is invalid")
	case errors.Is(err, identitydomain.ErrNameRequired):
		return grpcserver.NewStatusError(codes.InvalidArgument, "auth.nameRequired", "name is required")
	case errors.Is(err, identitydomain.ErrLastNameRequired):
		return grpcserver.NewStatusError(codes.InvalidArgument, "auth.lastNameRequired", "last name is required")
	case errors.Is(err, identitydomain.ErrPasswordRequired):
		return grpcserver.NewStatusError(codes.InvalidArgument, "auth.passwordRequired", "password is required")
	case errors.Is(err, identitydomain.ErrPasswordTooShort):
		return grpcserver.NewStatusError(codes.InvalidArgument, "auth.passwordTooShort", "password is too short")
	case errors.Is(err, identitydomain.ErrPasswordTooLong):
		return grpcserver.NewStatusError(codes.InvalidArgument, "auth.passwordTooLong", "password is too long")
	case errors.Is(err, identitydomain.ErrLoginAlreadyUsed):
		return grpcserver.NewStatusError(codes.AlreadyExists, "auth.loginAlreadyUsed", "login is already used")
	case errors.Is(err, identitydomain.ErrEmailAlreadyUsed):
		return grpcserver.NewStatusError(codes.AlreadyExists, "auth.emailAlreadyUsed", "email is already used")
	case errors.Is(err, identitydomain.ErrEmailInvalid):
		return grpcserver.NewStatusError(codes.InvalidArgument, "auth.emailInvalid", "email is invalid")
	case errors.Is(err, identitydomain.ErrGenderInvalid):
		return grpcserver.NewStatusError(codes.InvalidArgument, "auth.genderInvalid", "gender is invalid")
	case errors.Is(err, identitydomain.ErrUserNotFound):
		return grpcserver.NewStatusError(codes.NotFound, "user.notFound", "user not found")
	case errors.Is(err, identitydomain.ErrNameEmpty):
		return grpcserver.NewStatusError(codes.InvalidArgument, "user.nameRequired", "name is required")
	case errors.Is(err, identitydomain.ErrEmailEmpty):
		return grpcserver.NewStatusError(codes.InvalidArgument, "user.emailRequired", "email is required")
	case errors.Is(err, identitydomain.ErrRoleEmpty):
		return grpcserver.NewStatusError(codes.InvalidArgument, "user.roleRequired", "role is required")
	case errors.Is(err, identitydomain.ErrPasswordEmpty):
		return grpcserver.NewStatusError(codes.InvalidArgument, "user.passwordRequired", "password is required")
	case errors.Is(err, identitydomain.ErrPasswordLength):
		return grpcserver.NewStatusError(codes.InvalidArgument, "user.passwordLength", "password must be between 5 and 32 characters")
	case errors.Is(err, identitydomain.ErrPasswordWhitespace):
		return grpcserver.NewStatusError(codes.InvalidArgument, "user.passwordWhitespace", "password must not contain spaces")
	default:
		return grpcserver.NewStatusError(codes.Internal, "identity.internal", http.StatusText(http.StatusInternalServerError))
	}
}

func unimplemented(method string) error {
	return status.Errorf(codes.Unimplemented, "identity grpc: %s is not implemented yet", method)
}
