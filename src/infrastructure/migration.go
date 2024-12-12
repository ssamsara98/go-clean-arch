package infrastructure

import (
	migrate "github.com/rubenv/sql-migrate"
	"github.com/ssamsara98/go-clean-arch/src/constants"
	"github.com/ssamsara98/go-clean-arch/src/lib"
)

/*
Migrations -> Migration Struct
*/

type Migrations struct {
	env    *lib.Env
	logger *lib.Logger
	db     *Database
}

/*
NewMigrations -> return new Migrations struct
*/

// func NewMigrations(
// 	env *lib.Env,
// 	logger *lib.Logger,
// 	db Database,
// ) *Migrations {
// 	return &Migrations{
// 		env:    env,
// 		logger: logger,
// 		db:     db,
// 	}
// }

/*
Migrate migrates all migrations that are defined
*/

func (m Migrations) Migrate() error {
	if m.env.Environment == constants.Production {
		m.logger.Info("no start-up migration on production.")
		return nil
	}

	sqlDB, err := m.db.DB.DB()
	if err != nil {
		return err
	}

	m.logger.Info("running migration.")
	migrations := &migrate.FileMigrationSource{
		Dir: "db/migrations/",
	}
	_, err = migrate.Exec(sqlDB, m.env.DatabaseType, migrations, migrate.Up)
	if err != nil {
		return err
	}
	m.logger.Info("migration completed.")

	return nil
}

/*
RunMigration runs the migration provided logger and database instance
*/

func RunMigration(
	env *lib.Env,
	logger *lib.Logger,
	db *Database,
) error {
	m := &Migrations{
		env:    env,
		logger: logger,
		db:     db,
	}
	return m.Migrate()
}
