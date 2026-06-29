package client

import (
	"context"
	"os"
	"strings"

	identityv1 "github.com/stratflow-labs/stratflow/services/identity/gen/go/proto/v1"

	"github.com/google/uuid"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// TODO: there is already a client in the service root; it may be better to remove it and merge the logic here.
// TokenVerifier verifies access tokens via Identity service.
type TokenVerifier interface {
	Verify(ctx context.Context, token string) (Claims, error)
}

type IdentityTokenVerifier struct {
	target string
}

// NewTokenVerifier constructs a verifier calling Identity gRPC API.
// If target is empty, falls back to IDENTITY_GRPC_URL, then localhost:9090.
func NewTokenVerifier(target string) *IdentityTokenVerifier {
	return &IdentityTokenVerifier{target: grpcTarget(target)}
}

func (v *IdentityTokenVerifier) Verify(ctx context.Context, token string) (claims Claims, err error) {
	if strings.TrimSpace(token) == "" {
		return Claims{}, ErrUnauthorized
	}

	conn, err := grpc.NewClient(v.target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return Claims{}, err
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	resp, err := identityv1.NewIdentityServiceClient(conn).VerifyToken(ctx, &identityv1.VerifyTokenRequest{
		AccessToken: token,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code().String() == "Unauthenticated" {
			return Claims{}, ErrUnauthorized
		}
		return Claims{}, err
	}

	uid, err := uuid.Parse(strings.TrimSpace(resp.GetUserId()))
	if err != nil {
		return Claims{}, err
	}

	return Claims{
		UserID: uid,
		Role:   mapRole(resp.GetRole()),
	}, nil
}

func grpcTarget(target string) string {
	if strings.TrimSpace(target) != "" {
		return strings.TrimSpace(target)
	}

	if envTarget := strings.TrimSpace(os.Getenv("IDENTITY_GRPC_URL")); envTarget != "" {
		return envTarget
	}
	return "localhost:9090"
}

func mapRole(role identityv1.Role) string {
	switch role {
	case identityv1.Role_ROLE_UNSPECIFIED:
		return ""
	case identityv1.Role_ROLE_USER:
		return "user"
	case identityv1.Role_ROLE_MANAGER:
		return "manager"
	case identityv1.Role_ROLE_ADMIN:
		return "admin"
	default:
		return ""
	}
}
