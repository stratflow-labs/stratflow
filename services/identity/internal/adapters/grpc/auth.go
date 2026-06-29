package identitygrpc

import (
	"context"

	identityv1 "github.com/stratflow-labs/stratflow/services/identity/gen/go/proto/v1"
	"github.com/stratflow-labs/stratflow/services/identity/internal/auth"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Handler) Login(ctx context.Context, req *identityv1.LoginRequest) (*identityv1.TokenEnvelope, error) {
	out, err := s.auth.Login(ctx, auth.LoginInput{
		Login:    req.GetLogin(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &identityv1.TokenEnvelope{
		Message: "login successful",
		Data: &identityv1.TokenPair{
			AccessToken: out.AccessToken,
		},
	}, nil
}

func (s *Handler) VerifyToken(ctx context.Context, req *identityv1.VerifyTokenRequest) (*identityv1.VerifyTokenResponse, error) {
	out, err := s.auth.VerifyToken(ctx, &auth.VerifyTokenInput{
		AccessToken: req.GetAccessToken(),
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &identityv1.VerifyTokenResponse{
		UserId: out.Claims.UserID.String(),
		Role:   mapRole(out.Claims.Role),
	}, nil
}

func (s *Handler) Logout(ctx context.Context, _ *identityv1.LogoutRequest) (*emptypb.Empty, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, mapError(err)
	}

	accessToken, err := bearerToken(ctx)
	if err != nil {
		return nil, mapError(err)
	}

	if err := s.auth.Logout(ctx, auth.LogoutCommand{
		UserID:      userID,
		AccessToken: accessToken,
	}); err != nil {
		return nil, mapError(err)
	}

	return noContent(ctx), nil
}
