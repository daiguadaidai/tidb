// Copyright 2018 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package planner

import (
	"github.com/daiguadaidai/parser/ast"
	"github.com/daiguadaidai/tidb/infoschema"
	"github.com/daiguadaidai/tidb/planner/cascades"
	plannercore "github.com/daiguadaidai/tidb/planner/core"
	"github.com/daiguadaidai/tidb/privilege"
	"github.com/daiguadaidai/tidb/sessionctx"
)

// Optimize does optimization and creates a Plan.
// The node must be prepared first.
func Optimize(ctx sessionctx.Context, node ast.Node, is infoschema.InfoSchema) (plannercore.Plan, error) {
	fp := plannercore.TryFastPlan(ctx, node)
	if fp != nil {
		return fp, nil
	}

	// build logical plan
	ctx.GetSessionVars().PlanID = 0
	ctx.GetSessionVars().PlanColumnID = 0
	builder := plannercore.NewPlanBuilder(ctx, is)
	p, err := builder.Build(node)
	if err != nil {
		return nil, err
	}

	ctx.GetSessionVars().StmtCtx.Tables = builder.GetDBTableInfo()
	activeRoles := ctx.GetSessionVars().ActiveRoles
	// Check privilege. Maybe it's better to move this to the Preprocess, but
	// we need the table information to check privilege, which is collected
	// into the visitInfo in the logical plan builder.
	if pm := privilege.GetPrivilegeManager(ctx); pm != nil {
		if err := plannercore.CheckPrivilege(activeRoles, pm, builder.GetVisitInfo()); err != nil {
			return nil, err
		}
	}

	if err := plannercore.CheckTableLock(ctx, is, builder.GetVisitInfo()); err != nil {
		return nil, err
	}

	// Handle the execute statement.
	if execPlan, ok := p.(*plannercore.Execute); ok {
		err := execPlan.OptimizePreparedPlan(ctx, is)
		return p, err
	}

	// Handle the non-logical plan statement.
	logic, isLogicalPlan := p.(plannercore.LogicalPlan)
	if !isLogicalPlan {
		return p, nil
	}

	// Handle the logical plan statement, use cascades planner if enabled.
	if ctx.GetSessionVars().EnableCascadesPlanner {
		return cascades.FindBestPlan(ctx, logic)
	}
	return plannercore.DoOptimize(builder.GetOptFlag(), logic)
}

func init() {
	plannercore.OptimizeAstNode = Optimize
}
