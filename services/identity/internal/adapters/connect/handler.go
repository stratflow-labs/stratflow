package identityconnect

import (
	"context"
	"net/http"

	"github.com/stratflow-labs/stratflow/internal/httpserver"
	identityv1 "github.com/stratflow-labs/stratflow/services/identity/gen/go/proto/v1"
	identityv1connect "github.com/stratflow-labs/stratflow/services/identity/gen/go/proto/v1/identityv1connect"

	connect "connectrpc.com/connect"

	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	pattern string
	handler http.Handler
}

func NewHandler(server identityv1.IdentityServiceServer) *Handler {
	pattern, handler := identityv1connect.NewIdentityServiceHandler(serviceAdapter{server: server})
	return &Handler{
		pattern: "/connect" + pattern,
		handler: http.StripPrefix("/connect", handler),
	}
}

func (h *Handler) RegisterRoutes(r httpserver.Router) {
	r.Mount(h.pattern, h.handler)
}

type serviceAdapter struct {
	server identityv1.IdentityServiceServer
}

func (a serviceAdapter) Login(ctx context.Context, req *connect.Request[identityv1.LoginRequest]) (*connect.Response[identityv1.TokenEnvelope], error) {
	return unary(ctx, req, a.server.Login)
}

func (a serviceAdapter) Logout(ctx context.Context, req *connect.Request[identityv1.LogoutRequest]) (*connect.Response[emptypb.Empty], error) {
	return unary(ctx, req, a.server.Logout)
}

func (a serviceAdapter) VerifyToken(ctx context.Context, req *connect.Request[identityv1.VerifyTokenRequest]) (*connect.Response[identityv1.VerifyTokenResponse], error) {
	return unary(ctx, req, a.server.VerifyToken)
}

func (a serviceAdapter) CreateUser(ctx context.Context, req *connect.Request[identityv1.CreateUserRequest]) (*connect.Response[identityv1.UserResponse], error) {
	return unary(ctx, req, a.server.CreateUser)
}

func (a serviceAdapter) ListUsers(ctx context.Context, req *connect.Request[identityv1.ListUsersRequest]) (*connect.Response[identityv1.UsersListResponse], error) {
	return unary(ctx, req, a.server.ListUsers)
}

func (a serviceAdapter) GetUser(ctx context.Context, req *connect.Request[identityv1.GetUserRequest]) (*connect.Response[identityv1.UserResponse], error) {
	return unary(ctx, req, a.server.GetUser)
}

func (a serviceAdapter) UpdateUser(ctx context.Context, req *connect.Request[identityv1.UpdateUserByIDRequest]) (*connect.Response[identityv1.UserResponse], error) {
	return unary(ctx, req, a.server.UpdateUser)
}

func (a serviceAdapter) DeleteUser(ctx context.Context, req *connect.Request[identityv1.DeleteUserRequest]) (*connect.Response[emptypb.Empty], error) {
	return unary(ctx, req, a.server.DeleteUser)
}

func (a serviceAdapter) GetCurrentUser(ctx context.Context, req *connect.Request[identityv1.GetCurrentUserRequest]) (*connect.Response[identityv1.UserResponse], error) {
	return unary(ctx, req, a.server.GetCurrentUser)
}

func (a serviceAdapter) UpdateCurrentUser(ctx context.Context, req *connect.Request[identityv1.UpdateCurrentUserRequest]) (*connect.Response[identityv1.UserResponse], error) {
	return unary(ctx, req, a.server.UpdateCurrentUser)
}

func (a serviceAdapter) DeleteCurrentUser(ctx context.Context, req *connect.Request[identityv1.DeleteCurrentUserRequest]) (*connect.Response[emptypb.Empty], error) {
	return unary(ctx, req, a.server.DeleteCurrentUser)
}

func unary[Req any, Resp any](ctx context.Context, req *connect.Request[Req], fn func(context.Context, *Req) (*Resp, error)) (*connect.Response[Resp], error) {
	msg, err := fn(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(msg), nil
}
