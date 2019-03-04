package main

import (
	"fmt"
	"io/ioutil"
)

func DefaultEnv() *Environment {
	env := NewEnvironment(nil)
	env.Set("+", &Procedure{f: func(args []Expression) (Expression, int) { return Number(args[0].(Number) + args[1].(Number)), 0 }})
	env.Set("-", &Procedure{f: func(args []Expression) (Expression, int) { return Number(args[0].(Number) - args[1].(Number)), 0 }})
	env.Set("*", &Procedure{f: func(args []Expression) (Expression, int) { return Number(args[0].(Number) * args[1].(Number)), 0 }})
	env.Set("/", &Procedure{f: func(args []Expression) (Expression, int) { return Number(args[0].(Number) / args[1].(Number)), 0 }})
	env.Set("<", &Procedure{f: func(args []Expression) (Expression, int) { return Bool(args[0].(Number) < args[1].(Number)), 0 }})
	env.Set("<=", &Procedure{f: func(args []Expression) (Expression, int) { return Bool(args[0].(Number) <= args[1].(Number)), 0 }})
	env.Set(">", &Procedure{f: func(args []Expression) (Expression, int) { return Bool(args[0].(Number) > args[1].(Number)), 0 }})
	env.Set(">=", &Procedure{f: func(args []Expression) (Expression, int) { return Bool(args[0].(Number) >= args[1].(Number)), 0 }})
	env.Set("=", &Procedure{f: func(args []Expression) (Expression, int) { return Bool(args[0].(Number) == args[1].(Number)), 0 }})
	env.Set("list", &Procedure{f: func(args []Expression) (Expression, int) { return List(args), 0 }})
	env.Set("print", &Procedure{f: func(args []Expression) (Expression, int) { fmt.Println(args[0].(String)); return Nil{}, 0 }})
	env.Set("str", &Procedure{f: func(args []Expression) (Expression, int) { return String(args[0].ExprToStr()), 0 }})
	env.Set("save", &Procedure{f: func(args []Expression) (Expression, int) {
		filename := args[0].(String)
		expr := args[1]
		ioutil.WriteFile(string(filename), []byte(expr.ExprToStr()+"\n"), 0644)
		return Nil{}, 0
	}})
	env.Set("load", &Procedure{f: func(args []Expression) (Expression, int) {
		filename := args[0].(String)
		content, err := ioutil.ReadFile(string(filename))
		if err != nil {
			return Error(err.Error()), 0
		}
		expr, err := readFromTokens(tokenize(string(content)))
		if err != nil {
			return Error(err.Error()), 0
		}
		return expr, 0
	}})
	return env
}
