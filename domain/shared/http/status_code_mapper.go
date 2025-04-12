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

package http

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mattn/go-sqlite3"
	nethttp "net/http"
	databaseerrors "powerdns-auth-proxy/domain/shared/database/repository/errors"
	httperrors "powerdns-auth-proxy/domain/shared/http/errors"
)

type StatusCodeMapper struct {
}

func (m *StatusCodeMapper) MapStatusCode(err error) int {
	if err == nil {
		return nethttp.StatusOK
	}

	var httpError httperrors.HTTPError
	if errors.As(err, &httpError) {
		return httpError.GetCode()
	}

	sqliteErr := sqlite3.Error{}
	if errors.As(err, &sqliteErr) {
		return m.mapSqliteError(err.(sqlite3.Error))
	}

	mysqlError := &mysql.MySQLError{}
	if errors.As(err, &mysqlError) {
		return m.mapMySqlError(err.(*mysql.MySQLError))
	}

	pgError := &pgconn.PgError{}
	if errors.As(err, &pgError) {
		return m.mapPgError(err.(*pgconn.PgError))
	}

	if errors.As(err, &databaseerrors.ItemNotFoundError{}) {
		return nethttp.StatusNotFound
	}

	if errors.As(err, &databaseerrors.ItemAlreadyExistsError{}) {
		return nethttp.StatusConflict
	}

	return nethttp.StatusInternalServerError
}

func (m *StatusCodeMapper) mapSqliteError(err sqlite3.Error) int {
	switch err.Code {
	case sqlite3.ErrConstraint:
		return nethttp.StatusConflict
	default:
		return nethttp.StatusInternalServerError
	}
}

func (m *StatusCodeMapper) mapMySqlError(sqlError *mysql.MySQLError) int {
	switch sqlError.Number {
	case 1062:
		return nethttp.StatusConflict
	default:
		return nethttp.StatusInternalServerError
	}
}

func (m *StatusCodeMapper) mapPgError(pgError *pgconn.PgError) int {
	switch pgError.Code {
	case "23505", "42601":
		return nethttp.StatusConflict
	default:
		return nethttp.StatusInternalServerError
	}
}
