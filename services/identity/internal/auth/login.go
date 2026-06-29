package auth

import (
	"context"
	"fmt"

	identitydomain "github.com/stratflow-labs/stratflow/services/identity/internal/domain"
)

func (s *Service) Login(ctx context.Context, input LoginInput) (LoginOutput, error) {
	identity, password, err := identitydomain.NormalizeLoginInput(input.Login, input.Password)
	if err != nil {
		return LoginOutput{}, err
	}

	creds, err := s.user.FindByIdentity(ctx, identity)
	if err != nil {
		return LoginOutput{}, fmt.Errorf("find user credentials: %w", err)
	}

	match, err := s.password.Compare(ctx, password, creds.PasswordHash)
	if err != nil {
		return LoginOutput{}, fmt.Errorf("compare password hash: %w", err)
	}
	if !match {
		return LoginOutput{}, identitydomain.ErrPasswordMismatch
	}

	var out LoginOutput
	err = s.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		access, err := s.token.IssueAccessToken(txCtx, identitydomain.TokenClaims{
			UserID: creds.UserID,
			Role:   creds.Role,
		})
		if err != nil {
			return fmt.Errorf("issue access token: %w", err)
		}

		out = LoginOutput{AccessToken: access.Value}
		return nil
	})
	if err != nil {
		if identitydomain.IsLoginValidation(err) {
			return LoginOutput{}, err
		}
		return LoginOutput{}, fmt.Errorf("login transaction: %w", err)
	}

	return out, nil
}
