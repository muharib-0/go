package service

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	sqlc "github.com/muharib-0/ainyx-user-api/db/sqlc"
	"github.com/muharib-0/ainyx-user-api/internal/logger"
	"github.com/muharib-0/ainyx-user-api/internal/models"
	"github.com/muharib-0/ainyx-user-api/internal/repository"
	"go.uber.org/zap"
)

const dateLayout = "2006-01-02"

type UserService interface {
	CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.UserResponse, error)
	GetUserByID(ctx context.Context, id int32) (*models.UserWithAgeResponse, error)
	ListUsers(ctx context.Context, page, pageSize int) (*models.PaginatedUsersResponse, error)
	UpdateUser(ctx context.Context, id int32, req models.UpdateUserRequest) (*models.UserResponse, error)
	DeleteUser(ctx context.Context, id int32) error
}

type userService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.UserResponse, error) {
	dob, err := time.Parse(dateLayout, req.Dob)
	if err != nil {
		logger.Error("Failed to parse date of birth", zap.Error(err))
		return nil, err
	}

	arg := sqlc.CreateUserParams{
		Name: req.Name,
		Dob:  pgtype.Date{Time: dob, Valid: true},
	}

	user, err := s.repo.CreateUser(ctx, arg)
	if err != nil {
		logger.Error("Failed to create user", zap.Error(err))
		return nil, err
	}

	logger.Info("User created successfully", zap.Int32("user_id", user.ID), zap.String("name", user.Name))

	return &models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob.Time.Format(dateLayout),
	}, nil
}

func (s *userService) GetUserByID(ctx context.Context, id int32) (*models.UserWithAgeResponse, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		logger.Error("Failed to get user by ID", zap.Int32("user_id", id), zap.Error(err))
		return nil, err
	}

	age := models.CalculateAge(user.Dob.Time)

	logger.Info("User retrieved successfully", zap.Int32("user_id", user.ID))

	return &models.UserWithAgeResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob.Time.Format(dateLayout),
		Age:  age,
	}, nil
}

func (s *userService) ListUsers(ctx context.Context, page, pageSize int) (*models.PaginatedUsersResponse, error) {
	// Set defaults
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// User repository ListUsers does not assume pagination, so we fetch all and paginate in memory
	allUsers, err := s.repo.ListUsers(ctx)
	if err != nil {
		logger.Error("Failed to list users", zap.Error(err))
		return nil, err
	}

	total, err := s.repo.CountUsers(ctx)
	if err != nil {
		logger.Error("Failed to count users", zap.Error(err))
		return nil, err
	}

	// Pagination implementation
	start := (page - 1) * pageSize
	if start > len(allUsers) {
		start = len(allUsers)
	}
	end := start + pageSize
	if end > len(allUsers) {
		end = len(allUsers)
	}
	users := allUsers[start:end]

	userResponses := make([]models.UserWithAgeResponse, len(users))
	for i, user := range users {
		userResponses[i] = s.toUserWithAgeResponse(user)
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	logger.Info("Users listed successfully", zap.Int("count", len(users)), zap.Int64("total", total))

	return &models.PaginatedUsersResponse{
		Users:      userResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *userService) UpdateUser(ctx context.Context, id int32, req models.UpdateUserRequest) (*models.UserResponse, error) {
	dob, err := time.Parse(dateLayout, req.Dob)
	if err != nil {
		logger.Error("Failed to parse date of birth", zap.Error(err))
		return nil, err
	}

	arg := sqlc.UpdateUserParams{
		ID:   id,
		Name: req.Name,
		Dob:  pgtype.Date{Time: dob, Valid: true},
	}

	user, err := s.repo.UpdateUser(ctx, arg)
	if err != nil {
		logger.Error("Failed to update user", zap.Int32("user_id", id), zap.Error(err))
		return nil, err
	}

	logger.Info("User updated successfully", zap.Int32("user_id", user.ID), zap.String("name", user.Name))

	return &models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob.Time.Format(dateLayout),
	}, nil
}

func (s *userService) DeleteUser(ctx context.Context, id int32) error {
	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		logger.Error("Failed to delete user", zap.Int32("user_id", id), zap.Error(err))
		return err
	}

	logger.Info("User deleted successfully", zap.Int32("user_id", id))
	return nil
}

func (s *userService) toUserWithAgeResponse(user sqlc.User) models.UserWithAgeResponse {
	return models.UserWithAgeResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob.Time.Format(dateLayout),
		Age:  models.CalculateAge(user.Dob.Time),
	}
}
