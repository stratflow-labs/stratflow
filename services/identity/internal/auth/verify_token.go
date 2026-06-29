package auth

import (
	"context"
	"fmt"

	identitydomain "github.com/stratflow-labs/stratflow/services/identity/internal/domain"
)

type VerifyTokenInput struct {
	AccessToken string
}

type VerifyTokenOutput struct {
	Claims identitydomain.TokenClaims
}

func (s *Service) VerifyToken(ctx context.Context, input *VerifyTokenInput) (VerifyTokenOutput, error) {
	if input == nil {
		return VerifyTokenOutput{}, fmt.Errorf("validation failed: input is required")
	}

	token, err := identitydomain.NormalizeAccessToken(input.AccessToken)
	if err != nil {
		return VerifyTokenOutput{}, err
	}

	payload, err := s.token.VerifyAccessToken(ctx, token)
	if err != nil {
		return VerifyTokenOutput{}, fmt.Errorf("verify access token: %w", err)
	}

	return VerifyTokenOutput{
		Claims: payload.Claims,
	}, nil
}
