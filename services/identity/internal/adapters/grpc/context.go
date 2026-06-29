package identitygrpc

import (
	"context"
	"strings"

	"github.com/stratflow-labs/stratflow/internal/authkit"
	"github.com/stratflow-labs/stratflow/internal/grpcserver"
	identitydomain "github.com/stratflow-labs/stratflow/services/identity/internal/domain"

	"github.com/google/uuid"
)

func parseUserID(raw string) (uuid.UUID, error) {
	userID, err := uuid.Parse(strings.TrimSpace(raw))
	if err != nil || userID == uuid.Nil {
		return uuid.Nil, identitydomain.ErrUserNotFound
	}
	return userID, nil
}

func currentUserID(ctx context.Context) (uuid.UUID, error) {
	rawUserID, ok := authkit.UserIDFromContext(ctx)
	if !ok {
		return uuid.Nil, identitydomain.ErrAccessTokenNotFound
	}

	return parseUserID(rawUserID)
}

func bearerToken(ctx context.Context) (string, error) {
	token, err := grpcserver.BearerTokenFromContext(ctx)
	if err != nil {
		return "", identitydomain.ErrAccessTokenNotFound
	}

	return token, nil
}
