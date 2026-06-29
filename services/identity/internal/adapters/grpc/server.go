package identitygrpc

import (
	identityv1 "github.com/stratflow-labs/stratflow/services/identity/gen/go/proto/v1"
	"github.com/stratflow-labs/stratflow/services/identity/internal/auth"
	user "github.com/stratflow-labs/stratflow/services/identity/internal/user"
)

type HandlerDependencies struct {
	Auth  *auth.Service
	Users *user.Service
}

type Handler struct {
	identityv1.UnimplementedIdentityServiceServer

	auth  *auth.Service
	users *user.Service
}

func NewHandler(deps HandlerDependencies) *Handler {
	return &Handler{
		auth:  deps.Auth,
		users: deps.Users,
	}
}
