package handler

import (
	"context"

	"github.com/gieart87/gohexaclean/internal/dto/request"
	"github.com/gieart87/gohexaclean/internal/port/inbound"
	pb "github.com/gieart87/gohexaclean/api/proto/user"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserHandlerGRPC implements the gRPC user service
type UserHandlerGRPC struct {
	pb.UnimplementedUserServiceServer
	userService inbound.UserServicePort
}

// NewUserHandlerGRPC creates a new gRPC user handler
func NewUserHandlerGRPC(userService inbound.UserServicePort) *UserHandlerGRPC {
	return &UserHandlerGRPC{
		userService: userService,
	}
}

// CreateUser creates a new user and returns with token (for gRPC, we return LoginResponse with token)
func (h *UserHandlerGRPC) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.LoginResponse, error) {
	createReq := &request.CreateUserRequest{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	}

	// Validate request
	if err := createReq.Validate(); err != nil {
		return nil, err
	}

	registerResp, err := h.userService.CreateUser(ctx, createReq)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		Token: registerResp.Token,
		User: &pb.UserResponse{
			Id:        registerResp.User.ID.String(),
			Email:     registerResp.User.Email,
			Name:      registerResp.User.Name,
			IsActive:  true, // Active by default for new users
			CreatedAt: timestamppb.New(registerResp.User.CreatedAt),
			UpdatedAt: timestamppb.New(registerResp.User.UpdatedAt),
		},
	}, nil
}

// GetUser gets a user by ID
func (h *UserHandlerGRPC) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	user, err := h.userService.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{
		Id:        user.ID.String(),
		Email:     user.Email,
		Name:      user.Name,
		IsActive:  true, // No soft delete check in response, assume active
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}

// UpdateUser updates a user
func (h *UserHandlerGRPC) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	updateReq := &request.UpdateUserRequest{
		Name: req.Name,
	}

	// Validate request
	if err := updateReq.Validate(); err != nil {
		return nil, err
	}

	user, err := h.userService.UpdateUser(ctx, id, updateReq)
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{
		Id:        user.ID.String(),
		Email:     user.Email,
		Name:      user.Name,
		IsActive:  true, // No soft delete check in response, assume active
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}

// DeleteUser deletes a user
func (h *UserHandlerGRPC) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	if err := h.userService.DeleteUser(ctx, id); err != nil {
		return &pb.DeleteUserResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return &pb.DeleteUserResponse{
		Success: true,
		Message: "User deleted successfully",
	}, nil
}

// ListUsers lists users with pagination
func (h *UserHandlerGRPC) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	page := int(req.Page)
	limit := int(req.Limit)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, total, err := h.userService.ListUsers(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	pbUsers := make([]*pb.UserResponse, len(users))
	for i, user := range users {
		pbUsers[i] = &pb.UserResponse{
			Id:        user.ID.String(),
			Email:     user.Email,
			Name:      user.Name,
			IsActive:  true, // No soft delete check in response, assume active
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		}
	}

	return &pb.ListUsersResponse{
		Users: pbUsers,
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}, nil
}

// Login authenticates a user
func (h *UserHandlerGRPC) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	loginReq := &request.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	// Validate request
	if err := loginReq.Validate(); err != nil {
		return nil, err
	}

	loginResp, err := h.userService.Login(ctx, loginReq)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		Token: loginResp.Token,
		User: &pb.UserResponse{
			Id:        loginResp.User.ID.String(),
			Email:     loginResp.User.Email,
			Name:      loginResp.User.Name,
			IsActive:  true, // No soft delete check in response, assume active
			CreatedAt: timestamppb.New(loginResp.User.CreatedAt),
			UpdatedAt: timestamppb.New(loginResp.User.UpdatedAt),
		},
	}, nil
}
