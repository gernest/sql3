// Copyright 2021 Molecula Corp. All rights reserved.

package planner

import (
	"context"
	"errors"
	"log/slog"

	"github.com/gernest/sql3"
	"github.com/gernest/sql3/api"
	"github.com/gernest/sql3/dax"
	"github.com/gernest/sql3/parser"
	"github.com/gernest/sql3/planner/types"
)

func isTableNotFoundError(err error) bool {
	return errors.Is(err, dax.ErrTableNameDoesNotExist) ||
		errors.Is(errors.Unwrap(err), dax.ErrTableNameDoesNotExist)
}

// ExecutionPlanner compiles SQL text into a query plan
type ExecutionPlanner struct {
	executor       api.Executor
	schemaAPI      api.SchemaAPI
	systemAPI      api.SystemAPI
	systemLayerAPI api.SystemLayerAPI
	importer       api.Importer
	logger         slog.Logger
	sql            string
}

func NewExecutionPlanner(executor api.Executor, schemaAPI api.SchemaAPI, systemAPI api.SystemAPI, systemLayerAPI api.SystemLayerAPI, importer api.Importer, logger slog.Logger, sql string) *ExecutionPlanner {
	return &ExecutionPlanner{
		executor:       executor,
		schemaAPI:      schemaAPI,
		systemAPI:      systemAPI,
		systemLayerAPI: systemLayerAPI,
		importer:       importer,
		logger:         logger,
		sql:            sql,
	}
}

// CompilePlan takes an AST (parser.Statement) and compiles into a query plan returning the root
// PlanOperator
// The act of compiling includes an analysis step that does semantic analysis of the AST, this includes
// type checking, and sometimes AST rewriting. The compile phase uses the type-checked and rewritten AST
// to produce a query plan.
func (p *ExecutionPlanner) CompilePlan(ctx context.Context, stmt parser.Statement) (types.PlanOperator, error) {
	// call analyze first
	err := p.analyzePlan(ctx, stmt)
	if err != nil {
		return nil, err
	}

	var rootOperator types.PlanOperator
	switch stmt := stmt.(type) {
	case *parser.SelectStatement:
		rootOperator, err = p.compileSelectStatement(stmt, false)
	default:
		return nil, sql3.NewErrInternalf("cannot plan statement: %T", stmt)
	}
	// optimize the plan
	if err == nil {
		// rootOperator, err = p.optimizePlan(ctx, rootOperator)
	}
	return rootOperator, err
}

func (p *ExecutionPlanner) analyzePlan(ctx context.Context, stmt parser.Statement) error {
	switch stmt := stmt.(type) {
	case *parser.SelectStatement:
		_, err := p.analyzeSelectStatement(ctx, stmt)
		return err
	default:
		return sql3.NewErrInternalf("cannot analyze statement: %T", stmt)
	}
}

type accessType byte

const (
	accessTypeReadData accessType = iota
	accessTypeWriteData
	accessTypeCreateObject
	accessTypeAlterObject
	accessTypeDropObject
)

func (p *ExecutionPlanner) checkAccess(ctx context.Context, objectName string, _ accessType) error {
	return nil
}

type reduceFunc func(ctx context.Context, prev, v types.Rows) (types.Rows, error)

type mapResponse struct {
	node   api.ClusterNode
	result types.Rows
	err    error
}

func (e *ExecutionPlanner) mapReducePlanOp(ctx context.Context, op types.PlanOperator) (result types.Rows, err error) {
	iter, err := op.Iterator(ctx, nil)
	if err != nil {
		return nil, err
	}
	row, err := iter.Next(ctx)
	if err != nil && err != types.ErrNoMoreRows {
		return nil, err
	}
	if err != types.ErrNoMoreRows {
		for {
			result = append(result, row)
			row, err = iter.Next(ctx)
			if err != nil && err != types.ErrNoMoreRows {
				return nil, err
			}
			if err == types.ErrNoMoreRows {
				break
			}
		}
	}
	return result, nil
}
