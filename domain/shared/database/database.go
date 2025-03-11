// MIT License
// Copyright (c) 2025 Toni Liesche
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"powerdns-auth-proxy/domain/shared/config"
	"time"
)

func ProvideDatabaseClient(config *config.Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	gormLogger, err := getLogger(config.LogPath)
	if err != nil {
		return nil, err
	}

	switch config.Database {
	case "mariadb", "mysql":
		db, err = provideMariadbClient(config.MySQL, gormLogger)
	case "postgres", "postgresql":
		db, err = providePostgresClient(config.PostgreSQL, gormLogger)
	case "sqlite":
		db, err = provideSqliteClient(config.Sqlite, gormLogger)
	default:
		return nil, fmt.Errorf("could not provide database client: unsupported database type: %s", config.Database)
	}

	if err != nil {
		return nil, fmt.Errorf("could not provide database client: %w", err)
	}

	return db, nil
}

func provideSqliteClient(config *config.SqliteConfig, logger logger.Interface) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(config.Path), &gorm.Config{Logger: logger})

	if err != nil {
		return nil, err
	}

	if db == nil {
		return nil, fmt.Errorf("could not open local database file")
	}

	return db, nil
}

func provideMariadbClient(config *config.MySQLDBConfig, logger logger.Interface) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.User, config.Password, config.Host, config.Port, config.Database)

	db, err := gorm.Open(mysql.New(mysql.Config{DSN: dsn}), &gorm.Config{Logger: logger})
	if err != nil {
		return nil, err
	}

	if db == nil {
		return nil, fmt.Errorf("could not connect to remote database")
	}

	return db, nil
}

func providePostgresClient(sql *config.PostgreSQLConfig, gormLogger logger.Interface) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", sql.Host, sql.User, sql.Password, sql.Database, sql.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormLogger})
	if err != nil {
		return nil, err
	}

	if db == nil {
		return nil, fmt.Errorf("could not connect to remote database")
	}

	return db, nil
}

func getLogger(logPath string) (logger.Interface, error) {
	logWriter, err := os.Create(fmt.Sprintf("%s/database.log", logPath))
	if err != nil {
		return nil, err
	}

	return logger.New(
		log.New(logWriter, "\n", log.LstdFlags), // io logWriter
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	), nil
}

func ProvideTestDatabaseClient() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
}
