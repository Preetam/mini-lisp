package main

import (
	"fmt"
	"io/ioutil"
)

func DefaultEnv() *Environment {
	env := NewEnvironment(nil)
	env.Set("+", &Procedure{f: func(args []Expression) Expression { return Number(args[0].(Number) + args[1].(Number)) }})
	env.Set("-", &Procedure{f: func(args []Expression) Expression { return Number(args[0].(Number) - args[1].(Number)) }})
	env.Set("*", &Procedure{f: func(args []Expression) Expression { return Number(args[0].(Number) * args[1].(Number)) }})
	env.Set("/", &Procedure{f: func(args []Expression) Expression { return Number(args[0].(Number) / args[1].(Number)) }})
	env.Set("<", &Procedure{f: func(args []Expression) Expression { return Bool(args[0].(Number) < args[1].(Number)) }})
	env.Set("<=", &Procedure{f: func(args []Expression) Expression { return Bool(args[0].(Number) <= args[1].(Number)) }})
	env.Set(">", &Procedure{f: func(args []Expression) Expression { return Bool(args[0].(Number) > args[1].(Number)) }})
	env.Set(">=", &Procedure{f: func(args []Expression) Expression { return Bool(args[0].(Number) >= args[1].(Number)) }})
	env.Set("=", &Procedure{f: func(args []Expression) Expression { return Bool(args[0].(Number) == args[1].(Number)) }})
	env.Set("list", &Procedure{f: func(args []Expression) Expression { return List(args) }})
	env.Set("car", &Procedure{f: func(args []Expression) Expression { return args[0].(List)[0] }})
	env.Set("cdr", &Procedure{f: func(args []Expression) Expression { return args[0].(List)[1:] }})
	env.Set("cons", &Procedure{f: func(args []Expression) Expression {
		if len(args) == 1 {
			return List{args[0]}
		}
		if args[1] == nil || IsNil(args[1]) {
			return List{args[0]}
		}
		if _, ok := args[1].(List); ok {
			return append(List{args[0]}, args[1].(List)...)
		}
		return List{args[0], args[1]}
	}})
	env.Set("print", &Procedure{f: func(args []Expression) Expression { fmt.Println(args[0]); return Nil{} }})
	env.Set("str", &Procedure{f: func(args []Expression) Expression { return String(args[0].ExprToStr()) }})
	env.Set("save", &Procedure{f: func(args []Expression) Expression {
		filename := args[0].(String)
		expr := args[1]
		ioutil.WriteFile(string(filename), []byte(expr.ExprToStr()+"\n"), 0644)
		return Nil{}
	}})
	env.Set("load", &Procedure{f: func(args []Expression) Expression {
		filename := args[0].(String)
		content, err := ioutil.ReadFile(string(filename))
		if err != nil {
			return Error(err.Error())
		}
		expr, err := readFromTokens(tokenize(string(content)))
		if err != nil {
			return Error(err.Error())
		}
		return expr
	}})
	return env
}
