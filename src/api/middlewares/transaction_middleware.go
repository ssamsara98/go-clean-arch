package middlewares

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/ssamsara98/go-clean-arch/src/constants"
	"github.com/ssamsara98/go-clean-arch/src/infrastructure"
	"github.com/ssamsara98/go-clean-arch/src/lib"
	"github.com/ssamsara98/go-clean-arch/src/utils"
)

/*
DBTransactionMiddleware -> struct for transaction
*/

type DBTransactionMiddleware struct {
	logger *lib.Logger
	db     *infrastructure.Database
}

/*
NewDBTransactionMiddleware -> new instance of transaction
*/

func NewDBTransactionMiddleware(
	logger *lib.Logger,
	db *infrastructure.Database,
) *DBTransactionMiddleware {
	return &DBTransactionMiddleware{
		logger: logger,
		db:     db,
	}
}

/*
Handle -> It setup the database transaction middleware
*/

func (m DBTransactionMiddleware) Handle() fiber.Handler {
	m.logger.Debug("setting up database transaction middleware")

	return func(c *fiber.Ctx) (err error) {
		txHandle := m.db.DB.Begin()
		m.logger.Debug("beginning database transaction")

		defer func() {
			if r := recover(); r != nil {
				txHandle.Rollback()
			}
		}()

		c.Locals(constants.DBTransaction, txHandle)
		err = c.Next()

		if utils.StatusInList(c.Response().StatusCode(), []int{http.StatusOK, http.StatusCreated}) {
			m.logger.Debug("committing transactions")
			if err := txHandle.Commit().Error; err != nil {
				m.logger.Error("trx commit error: ", err)
			}
		} else {
			m.logger.Debug("rolling back transaction due to status code: ", c.Response().StatusCode())
			txHandle.Rollback()
		}

		return err
	}
}
