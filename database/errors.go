package database

import "errors"

var (
	Err_NotFound            = errors.New("not found")
	Err_DuplicatePrimaryKey = errors.New("duplicate primary key")
)
