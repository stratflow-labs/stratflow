package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	identitydomain "github.com/stratflow-labs/stratflow/services/identity/internal/domain"
)

func (s *Service) Logout(ctx context.Context, cmd LogoutCommand) error {
	if cmd.UserID == uuid.Nil {
		return identitydomain.ErrAccessTokenNotFound
	}

	rawToken := strings.TrimSpace(cmd.AccessToken)
	if rawToken == "" {
		return identitydomain.ErrAccessTokenNotFound
	}

	err := s.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		if err := s.accessToken.DeleteAccessToken(txCtx, cmd.UserID, rawToken); err != nil {
			return fmt.Errorf("delete access token: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("logout transaction: %w", err)
	}

	return nil
}
