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
	"powerdns-auth-proxy/domain/shared/basics"
	"powerdns-auth-proxy/domain/shared/config"
)

func NewDatabaseClient(container *basics.InjectionContainer, config *config.Config) (*gorm.DB, error) {
	if container == nil {
		return nil, basics.NewMissingDependencyError("could not provide database client: passed injection container is nil")
	}

	var db *gorm.DB
	var err error

	gormLogger := newGormLogger(container.Logger)

	switch config.Database {
	case "mariadb", "mysql":
		db, err = newMariadbClient(config.MySQL, gormLogger)
	case "postgres", "postgresql":
		db, err = newPostgresClient(config.PostgreSQL, gormLogger)
	case "sqlite":
		db, err = newSqliteClient(config.Sqlite, gormLogger)
	default:
		return nil, basics.NewMissingDependencyError(fmt.Sprintf("could not provide database client: unsupported database type: %s", config.Database))
	}

	if err != nil {
		return nil, basics.NewMissingDependencyError(fmt.Sprintf("could not provide database client: %v", err))
	}

	return db, nil
}

func newSqliteClient(config *config.SqliteConfig, logger logger.Interface) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(config.Path), &gorm.Config{Logger: logger})

	if err != nil {
		return nil, err
	}

	if db == nil {
		return nil, fmt.Errorf("could not open local database file")
	}

	return db, nil
}

func newMariadbClient(config *config.MySQLDBConfig, logger logger.Interface) (*gorm.DB, error) {
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

func newPostgresClient(sql *config.PostgreSQLConfig, logger logger.Interface) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", sql.Host, sql.User, sql.Password, sql.Database, sql.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger})
	if err != nil {
		return nil, err
	}

	if db == nil {
		return nil, fmt.Errorf("could not connect to remote database")
	}

	return db, nil
}

func NewTestDatabaseClient(container *basics.InjectionContainer) (*gorm.DB, error) {
	if container == nil {
		return nil, basics.NewMissingDependencyError("could not provide test database client: passed injection container is nil")
	}

	return gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
}
