package services

import (
	"github.com/ssamsara98/go-clean-arch/src/api/dto"
	"github.com/ssamsara98/go-clean-arch/src/infrastructure"
	"github.com/ssamsara98/go-clean-arch/src/lib"
	"github.com/ssamsara98/go-clean-arch/src/models"
	"github.com/ssamsara98/go-clean-arch/src/utils"
)

type UsersService struct {
	logger *lib.Logger
	db     *infrastructure.Database
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

func (u UsersService) GetUserList(limit, page *int64) (*[]models.User, *int64, error) {
	items := new([]models.User)
	count := new(int64)

	u.db = u.db.SetHandle(u.db.Scopes(utils.Paginate(limit, page)))
	u.db.Order("id DESC")

	err := u.db.Find(items).Offset(-1).Limit(-1).Count(count).Error
	if err != nil {
		return nil, nil, err
	}

	return items, count, nil
}
func (u UsersService) GetUserListCursor(limit, cursor *int64) (*[]models.User, error) {
	items := new([]models.User)

	u.db = u.db.SetHandle(u.db.Scopes(utils.PaginateCursor(limit)))
	u.db.Order("id DESC")
	if cursor != nil && *cursor > 1 {
		u.db.Where("id < ?", *cursor)
	}

	err := u.db.Find(items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (u UsersService) GetUserByID(uri *dto.GetUserByIDParams) (*models.User, error) {
	user := new(models.User)
	return user, u.db.First(user, "id = ?", uri.ID).Error
}
