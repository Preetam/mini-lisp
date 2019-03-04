package main

import (
	"fmt"
	"io/ioutil"
)

func DefaultEnv() *Environment {
	env := NewEnvironment(nil)
	env.Set("+", Procedure{f: func(args []Expression) Expression { return Number(args[0].(Number) + args[1].(Number)) }})
	env.Set("-", Procedure{f: func(args []Expression) Expression { return Number(args[0].(Number) - args[1].(Number)) }})
	env.Set("*", Procedure{f: func(args []Expression) Expression { return Number(args[0].(Number) * args[1].(Number)) }})
	env.Set("/", Procedure{f: func(args []Expression) Expression { return Number(args[0].(Number) / args[1].(Number)) }})
	env.Set("<", Procedure{f: func(args []Expression) Expression { return Bool(args[0].(Number) < args[1].(Number)) }})
	env.Set("<=", Procedure{f: func(args []Expression) Expression { return Bool(args[0].(Number) <= args[1].(Number)) }})
	env.Set(">", Procedure{f: func(args []Expression) Expression { return Bool(args[0].(Number) > args[1].(Number)) }})
	env.Set(">=", Procedure{f: func(args []Expression) Expression { return Bool(args[0].(Number) >= args[1].(Number)) }})
	env.Set("=", Procedure{f: func(args []Expression) Expression { return Bool(args[0].(Number) == args[1].(Number)) }})
	env.Set("list", Procedure{f: func(args []Expression) Expression { return List(args) }})
	env.Set("print", Procedure{f: func(args []Expression) Expression { fmt.Println(args[0].(String)); return Nil{} }})
	env.Set("str", Procedure{f: func(args []Expression) Expression { return String(args[0].ExprToStr()) }})
	env.Set("save", Procedure{f: func(args []Expression) Expression {
		filename := args[0].(String)
		expr := args[1]
		ioutil.WriteFile(string(filename), []byte(expr.ExprToStr()+"\n"), 0644)
		return Nil{}
	}})
	env.Set("load", Procedure{f: func(args []Expression) Expression {
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
