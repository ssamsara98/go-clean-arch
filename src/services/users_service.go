package services

import (
	"go-clean-arch/lib"
	"go-clean-arch/models"
	"go-clean-arch/repository"
	"go-clean-arch/src/dto"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UsersService struct {
	logger          lib.Logger
	userRepository  repository.UserRepository
	paginationScope *gorm.DB
}

func NewUsersService(
	logger lib.Logger,
	userRepository repository.UserRepository,
) *UsersService {
	return &UsersService{
		logger:         logger,
		userRepository: userRepository,
	}
}

// WithTrx delegates transaction to repository database
func (s UsersService) WithTrx(trxHandle *gorm.DB) UsersService {
	s.userRepository = s.userRepository.WithTrx(trxHandle)
	return s
}

// PaginationScope
func (s UsersService) SetPaginationScope(scope func(*gorm.DB) *gorm.DB) UsersService {
	s.paginationScope = s.userRepository.WithTrx(s.userRepository.Scopes(scope)).DB
	return s
}

func (s UsersService) GetUserList() (response gin.H, err error) {
	var users []models.User
	var count int64

	err = s.userRepository.WithTrx(s.paginationScope).Find(&users).Offset(-1).Limit(-1).Count(&count).Error
	if err != nil {
		return nil, err
	}

	return gin.H{"result": users, "count": count}, nil
}

func (s UsersService) GetUserByID(uri *dto.GetUserByIDParams) (user models.User, err error) {
	return user, s.userRepository.First(&user, "id = ?", uri.ID).Error
}
