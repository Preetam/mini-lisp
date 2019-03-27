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
	env.Set("car", &Procedure{f: func(args []Expression) (Expression, int) { return args[0].(List)[0], 0 }})
	env.Set("cdr", &Procedure{f: func(args []Expression) (Expression, int) { return args[0].(List)[1:], 0 }})
	env.Set("cons", &Procedure{f: func(args []Expression) (Expression, int) {
		if len(args) == 1 {
			return List{args[0]}, 0
		}
		if args[1] == nil || IsNil(args[1]) {
			return List{args[0]}, 0
		}
		if _, ok := args[1].(List); ok {
			return append(List{args[0]}, args[1].(List)...), 0
		}
		return List{args[0], args[1]}, 0
	}})
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
