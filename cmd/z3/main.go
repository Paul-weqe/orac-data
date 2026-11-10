package main

import (
	"fmt"

	z3 "github.com/Z3Prover/z3/src/api/go"
)

func main() {
	ctx := z3.NewContext()
	x := ctx.MkIntConst("x")
	y := ctx.MkIntConst("y")

	solver := ctx.NewSolver()
	solver.Assert(ctx.MkEq(ctx.MkAdd(x, y), ctx.MkInt(10, ctx.MkIntSort())))
	solver.Assert(ctx.MkEq(ctx.MkSub(x, y), ctx.MkInt(2, ctx.MkIntSort())))

	if solver.Check() == z3.Satisfiable {
		model := solver.Model()
		xVal, _ := model.Eval(x, true)
		yVal, _ := model.Eval(y, true)
		fmt.Printf("x = %s, y = %s\n", xVal.String(), yVal.String())
	}
}
