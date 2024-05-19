package services

import (
	"go-clean-arch/api/dto"
	"go-clean-arch/infrastructure"
	"go-clean-arch/lib"
	"go-clean-arch/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UsersService struct {
	logger          *lib.Logger
	db              *infrastructure.Database
	paginationScope *gorm.DB
}

func NewUsersService(
	logger *lib.Logger,
	db *infrastructure.Database,
) *UsersService {
	return &UsersService{
		logger: logger,
		db:     db,
	}
}

// PaginationScope
func (s *UsersService) SetPaginationScope(scope func(*gorm.DB) *gorm.DB) *UsersService {
	s.paginationScope = s.db.WithTrx(s.db.Scopes(scope)).DB
	return s
}

func (s *UsersService) GetUserList() (response gin.H, err error) {
	var users []models.User
	var count int64

	err = s.db.WithTrx(s.paginationScope).Find(&users).Offset(-1).Limit(-1).Count(&count).Error
	if err != nil {
		return nil, err
	}

	return gin.H{"result": users, "count": count}, nil
}

func (s *UsersService) GetUserByID(uri *dto.GetUserByIDParams) (user models.User, err error) {
	return user, s.db.First(&user, "id = ?", uri.ID).Error
}
