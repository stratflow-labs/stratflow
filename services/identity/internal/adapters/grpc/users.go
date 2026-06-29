package identitygrpc

import (
	"context"
	"net/http"

	identityv1 "github.com/stratflow-labs/stratflow/services/identity/gen/go/proto/v1"
	user "github.com/stratflow-labs/stratflow/services/identity/internal/user"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Handler) CreateUser(ctx context.Context, req *identityv1.CreateUserRequest) (*identityv1.UserResponse, error) {
	out, err := s.users.Create(ctx, user.CreateInput{
		Login:    req.GetLogin(),
		Name:     req.GetName(),
		LastName: req.GetLastName(),
		Email:    req.GetEmail(),
		Role:     req.GetRole(),
		Password: req.GetPassword(),
		Gender:   int32PtrToIntPtr(req.Gender),
	})
	if err != nil {
		return nil, mapError(err)
	}

	setHTTPStatus(ctx, http.StatusCreated)
	return &identityv1.UserResponse{
		Message: "user created",
		Data:    mapUser(out),
	}, nil
}

func (s *Handler) ListUsers(ctx context.Context, req *identityv1.ListUsersRequest) (*identityv1.UsersListResponse, error) {
	out, err := s.users.List(ctx, user.ListInput{
		Search:   req.GetSearch(),
		Page:     int(req.GetPage()),
		PageSize: int(req.GetPageSize()),
		Sort:     req.GetSort(),
	})
	if err != nil {
		return nil, mapError(err)
	}

	items := make([]*identityv1.User, len(out.Users))
	for i := range out.Users {
		items[i] = mapUser(out.Users[i])
	}

	return &identityv1.UsersListResponse{
		Message: "users listed",
		Data: &identityv1.UsersListData{
			Items: items,
			Total: out.Total,
		},
	}, nil
}

func (s *Handler) GetUser(ctx context.Context, req *identityv1.GetUserRequest) (*identityv1.UserResponse, error) {
	userID, err := parseUserID(req.GetUserId())
	if err != nil {
		return nil, mapError(err)
	}

	out, err := s.users.Get(ctx, userID)
	if err != nil {
		return nil, mapError(err)
	}

	return &identityv1.UserResponse{
		Message: "user retrieved",
		Data:    mapUser(out),
	}, nil
}

func (s *Handler) UpdateUser(ctx context.Context, req *identityv1.UpdateUserByIDRequest) (*identityv1.UserResponse, error) {
	userID, err := parseUserID(req.GetUserId())
	if err != nil {
		return nil, mapError(err)
	}

	patch := updateUserPatch(req.GetPatch())

	out, err := s.users.Update(ctx, user.UpdateInput{
		ID:       userID,
		Name:     patch.Name,
		LastName: patch.LastName,
		Email:    patch.Email,
		Gender:   int32PtrToIntPtr(patch.Gender),
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &identityv1.UserResponse{
		Message: "user updated",
		Data:    mapUser(out),
	}, nil
}

func (s *Handler) DeleteUser(ctx context.Context, req *identityv1.DeleteUserRequest) (*emptypb.Empty, error) {
	userID, err := parseUserID(req.GetUserId())
	if err != nil {
		return nil, mapError(err)
	}

	if err := s.users.Delete(ctx, userID); err != nil {
		return nil, mapError(err)
	}

	return noContent(ctx), nil
}

func (s *Handler) GetCurrentUser(ctx context.Context, _ *identityv1.GetCurrentUserRequest) (*identityv1.UserResponse, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, mapError(err)
	}

	out, err := s.users.Get(ctx, userID)
	if err != nil {
		return nil, mapError(err)
	}

	return &identityv1.UserResponse{
		Message: "user retrieved",
		Data:    mapUser(out),
	}, nil
}

func (s *Handler) UpdateCurrentUser(ctx context.Context, req *identityv1.UpdateCurrentUserRequest) (*identityv1.UserResponse, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, mapError(err)
	}

	patch := updateUserPatch(req.GetPatch())

	out, err := s.users.Update(ctx, user.UpdateInput{
		ID:       userID,
		Name:     patch.Name,
		LastName: patch.LastName,
		Email:    patch.Email,
		Gender:   int32PtrToIntPtr(patch.Gender),
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &identityv1.UserResponse{
		Message: "user updated",
		Data:    mapUser(out),
	}, nil
}

func (s *Handler) DeleteCurrentUser(ctx context.Context, _ *identityv1.DeleteCurrentUserRequest) (*emptypb.Empty, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, mapError(err)
	}

	if err := s.users.Delete(ctx, userID); err != nil {
		return nil, mapError(err)
	}

	return noContent(ctx), nil
}

func int32PtrToIntPtr(value *int32) *int {
	if value == nil {
		return nil
	}

	out := int(*value)
	return &out
}

func updateUserPatch(patch *identityv1.UpdateUserRequest) *identityv1.UpdateUserRequest {
	if patch != nil {
		return patch
	}

	return &identityv1.UpdateUserRequest{}
}
