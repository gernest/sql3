package dax

import (
	"fmt"

	"errors"
)

var (
	ErrOrganizationIDDoesNotExist = errors.New("OrganizationIDDoesNotExist")

	ErrDatabaseIDExists         = errors.New("DatabaseIDExists")
	ErrDatabaseIDDoesNotExist   = errors.New("DatabaseIDDoesNotExist")
	ErrDatabaseNameDoesNotExist = errors.New("DatabaseNameDoesNotExist")
	ErrDatabaseNameExists       = errors.New("DatabaseNameExists")

	ErrTableIDExists         = errors.New("TableIDExists")
	ErrTableKeyExists        = errors.New("TableKeyExists")
	ErrTableNameExists       = errors.New("TableNameExists")
	ErrTableIDDoesNotExist   = errors.New("TableIDDoesNotExist")
	ErrTableKeyDoesNotExist  = errors.New("TableKeyDoesNotExist")
	ErrTableNameDoesNotExist = errors.New("TableNameDoesNotExist")

	ErrFieldExists       = errors.New("FieldExists")
	ErrFieldDoesNotExist = errors.New("FieldDoesNotExist")

	ErrInvalidTransaction = errors.New("InvalidTransaction")

	ErrUnimplemented = errors.New("Unimplemented")
)

// The following are helper functions for constructing coded errors containing
// relevant information about the specific error.

func NewErrOrganizationIDDoesNotExist(orgID OrganizationID) error {
	return newError(
		ErrOrganizationIDDoesNotExist,
		fmt.Sprintf("Organization ID '%s' does not exist", orgID),
	)
}

func NewErrDatabaseIDExists(qdbid QualifiedDatabaseID) error {
	return newError(
		ErrDatabaseIDExists,
		fmt.Sprintf("database ID '%s' already exists", qdbid),
	)
}

func NewErrDatabaseIDDoesNotExist(qdbid QualifiedDatabaseID) error {
	return newError(
		ErrDatabaseIDDoesNotExist,
		fmt.Sprintf("database ID '%s' does not exist", qdbid),
	)
}

func NewErrDatabaseNameDoesNotExist(dbName DatabaseName) error {
	return newError(
		ErrDatabaseNameDoesNotExist,
		fmt.Sprintf("database name '%s' does not exist", dbName),
	)
}

func NewErrDatabaseNameExists(dbName DatabaseName) error {
	return newError(
		ErrDatabaseNameExists,
		fmt.Sprintf("database name %s already exists", dbName),
	)
}

func NewErrTableIDDoesNotExist(qtid QualifiedTableID) error {
	return newError(
		ErrTableIDDoesNotExist,
		fmt.Sprintf("table ID '%s' does not exist", qtid),
	)
}

func NewErrTableKeyDoesNotExist(tkey TableKey) error {
	return newError(
		ErrTableKeyDoesNotExist,
		fmt.Sprintf("table key '%s' does not exist", tkey),
	)
}

func NewErrTableNameDoesNotExist(tableName TableName) error {
	return newError(
		ErrTableNameDoesNotExist,
		fmt.Sprintf("table name '%s' does not exist", tableName),
	)
}

func NewErrTableIDExists(qtid QualifiedTableID) error {
	return newError(
		ErrTableIDExists,
		fmt.Sprintf("table ID '%s' already exists", qtid),
	)
}

func NewErrTableKeyExists(tkey TableKey) error {
	return newError(
		ErrTableKeyExists,
		fmt.Sprintf("table key '%s' already exists", tkey),
	)
}

func NewErrTableNameExists(tableName TableName) error {
	return newError(
		ErrTableNameExists,
		fmt.Sprintf("table name '%s' already exists", tableName),
	)
}

func NewErrFieldDoesNotExist(fieldName FieldName) error {
	return newError(
		ErrFieldDoesNotExist,
		fmt.Sprintf("field '%s' does not exist", fieldName),
	)
}

func NewErrFieldExists(fieldName FieldName) error {
	return newError(
		ErrFieldExists,
		fmt.Sprintf("field '%s' already exists", fieldName),
	)
}

func NewErrInvalidTransaction(txType string) error {
	return newError(
		ErrInvalidTransaction,
		fmt.Sprintf("tx is not expected type: '%s'", txType),
	)
}

func newError(err error, msg string) error {
	return fmt.Errorf("%s %w", msg, err)
}
