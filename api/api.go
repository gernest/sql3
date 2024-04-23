package api

import (
	"context"

	"github.com/gernest/sql3/dax"
)

type Executor interface {
}
type ClusterNode interface{}

type SystemAPI interface {
	ClusterName() string
	Version() string
	PlatformDescription() string
	PlatformVersion() string
	ClusterNodeCount() int
	ClusterReplicaCount() int
	ShardWidth() int
	ClusterState() string
	DataDir() string

	NodeID() string
}

type SchemaAPI interface {
	CreateDatabase(context.Context, *dax.Database) error
	DropDatabase(context.Context, dax.DatabaseID) error

	DatabaseByName(ctx context.Context, dbname dax.DatabaseName) (*dax.Database, error)
	DatabaseByID(ctx context.Context, dbid dax.DatabaseID) (*dax.Database, error)
	SetDatabaseOption(ctx context.Context, dbid dax.DatabaseID, option string, value string) error
	Databases(context.Context, ...dax.DatabaseID) ([]*dax.Database, error)

	TableByName(ctx context.Context, tname dax.TableName) (*dax.Table, error)
	TableByID(ctx context.Context, tid dax.TableID) (*dax.Table, error)
	Tables(ctx context.Context) ([]*dax.Table, error)

	CreateTable(ctx context.Context, tbl *dax.Table) error
	CreateField(ctx context.Context, tname dax.TableName, fld *dax.Field) error

	DeleteTable(ctx context.Context, tname dax.TableName) error
	DeleteField(ctx context.Context, tname dax.TableName, fname dax.FieldName) error
}
