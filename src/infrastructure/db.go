package infrastructure

import (
	"fmt"
	"go-clean-arch/src/lib"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database modal
type Database struct {
	*gorm.DB
	logger *lib.Logger
}

// NewDatabase creates a new database instance
func NewDatabase(
	env *lib.Env,
	log *lib.Logger,
) *Database {
	url := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s",
		env.DBUsername,
		env.DBPassword,
		env.DBHost,
		env.DBPort,
		env.DBName,
	)

	log.Info("opening db connection")
	var db *gorm.DB
	var err error
	if env.Environment == "production" {
		db, err = gorm.Open(postgres.Open(url), &gorm.Config{
			Logger: log.GetGormLogger(),
			// Logger: logger.Default.LogMode(logger.Warn),
		})
	} else {
		db, err = gorm.Open(postgres.Open(url), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	}
	if err != nil {
		log.Info("Url: ", url)
		log.Panic(err)
	}

	// 	log.Info("creating database if it does't exist")
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
	// 		log.Info("couldn't create database")
	// 		log.Panic(err)
	// 	}

	// log.Info("using given database")
	// if err := db.Exec(fmt.Sprintf("USE %s", env.DBName)).Error; err != nil {
	// 	log.Info("cannot use the given database")
	// 	log.Panic(err)
	// }
	log.Info("database connection established")

	database := &Database{
		db,
		log,
	}
	log.Info("currentDatabase: ", db.Migrator().CurrentDatabase())

	if err := RunMigration(env, log, database); err != nil {
		log.Info("migration failed.")
		log.Panic(err)
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
