package infrastructure

import (
	"fmt"
	"go-clean-arch/src/lib"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// Database modal
type Database struct {
	*gorm.DB
	logger *lib.Logger
}

// NewDatabase creates a new database instance
func NewDatabase(
	env *lib.Env,
	logger *lib.Logger,
) *Database {
	url := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s",
		env.DBUsername,
		env.DBPassword,
		env.DBHost,
		env.DBPort,
		env.DBName,
	)

	logger.Info("opening db connection")
	var db *gorm.DB
	var err error
	if env.Environment == "production" {
		db, err = gorm.Open(postgres.Open(url), &gorm.Config{
			Logger: logger.GetGormLogger(),
		})
	} else {
		db, err = gorm.Open(postgres.Open(url), &gorm.Config{
			Logger: logger.GetGormLogger().LogMode(gormlogger.Info),
		})
	}
	if err != nil {
		logger.Info("Url: ", url)
		logger.Panic(err)
	}

	// 	logger.Info("creating database if it does't exist")
	// 	pgCreateDb := fmt.Sprintf(`
	// DO
	// $do$
	// DECLARE
	// 	_db TEXT := '%s';
	// 	_user TEXT := '%s';
	// 	_password TEXT := '%s';
	// BEGIN
	// 	CREATE EXTENSION IF NOT EXISTS dblink; -- enable extension
	// 	IF EXISTS (SELECT FROM pg_database WHERE datname = _db) THEN
	// 		RAISE NOTICE 'Database already exists';
	// 	ELSE
	// 		PERFORM dblink_connect('host=localhost user=' || _user || ' password=' || _password || ' dbname=' || current_database());
	// 		PERFORM dblink_exec('CREATE DATABASE ' || _db);
	// 	END IF;
	// END
	// $do$`, env.DBName, env.DBUsername, env.DBPassword)
	// 	if err = db.Exec(pgCreateDb).Error; err != nil {
	// 		logger.Info("couldn't create database")
	// 		logger.Panic(err)
	// 	}

	// logger.Info("using given database")
	// if err := db.Exec(fmt.Sprintf("USE %s", env.DBName)).Error; err != nil {
	// 	logger.Info("cannot use the given database")
	// 	logger.Panic(err)
	// }
	logger.Info("database connection established")

	database := &Database{
		db,
		logger,
	}
	logger.Info("currentDatabase: ", db.Migrator().CurrentDatabase())

	if err := RunMigration(env, logger, database); err != nil {
		logger.Info("migration failed.")
		logger.Panic(err)
	}

	return database
}

// WithTrx delegate transaction from user repository
func (d Database) WithTrx(trxHandle *gorm.DB) *Database {
	if trxHandle != nil {
		d.logger.Debug("using WithTrx as trxHandle is not nil")
		d.DB = trxHandle
	}
	return &d
}
