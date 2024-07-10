package services

import (
	"go-clean-arch/src/api/dto"
	"go-clean-arch/src/infrastructure"
	"go-clean-arch/src/lib"
	"go-clean-arch/src/models"

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
func (u UsersService) SetPaginationScope(scope func(*gorm.DB) *gorm.DB) *UsersService {
	u.paginationScope = u.db.WithTrx(u.db.Scopes(scope)).DB
	return &u
}

func (u UsersService) GetUserList() (*[]models.User, *int64, error) {
	var items []models.User
	var count int64

	err := u.db.WithTrx(u.paginationScope).Find(&items).Offset(-1).Limit(-1).Count(&count).Error
	if err != nil {
		return nil, nil, err
	}

	return &items, &count, nil
}

func (u UsersService) GetUserByID(uri *dto.GetUserByIDParams) (user models.User, err error) {
	return user, u.db.First(&user, "id = ?", uri.ID).Error
}
