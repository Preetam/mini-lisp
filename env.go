package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func DefaultEnv() *Environment {
	env := NewEnvironment(nil)

	// Numeric
	ensureNumeric := func(f func(args []Expression) Expression) func(args []Expression) Expression {
		return func(args []Expression) Expression {
			if len(args) != 2 {
				return Error("invalid arguments")
			}
			if _, ok := args[0].(Number); !ok {
				return Error("argument not a number")
			}
			if _, ok := args[1].(Number); !ok {
				return Error("argument not a number")
			}
			return f(args)
		}
	}
	env.Set("+", &Procedure{f: ensureNumeric(func(args []Expression) Expression { return Number(args[0].(Number) + args[1].(Number)) })})
	env.Set("-", &Procedure{f: ensureNumeric(func(args []Expression) Expression { return Number(args[0].(Number) - args[1].(Number)) })})
	env.Set("*", &Procedure{f: ensureNumeric(func(args []Expression) Expression { return Number(args[0].(Number) * args[1].(Number)) })})
	env.Set("/", &Procedure{f: ensureNumeric(func(args []Expression) Expression { return Number(args[0].(Number) / args[1].(Number)) })})
	env.Set("<", &Procedure{f: ensureNumeric(func(args []Expression) Expression { return Bool(args[0].(Number) < args[1].(Number)) })})
	env.Set("<=", &Procedure{f: ensureNumeric(func(args []Expression) Expression { return Bool(args[0].(Number) <= args[1].(Number)) })})
	env.Set(">", &Procedure{f: ensureNumeric(func(args []Expression) Expression { return Bool(args[0].(Number) > args[1].(Number)) })})
	env.Set(">=", &Procedure{f: ensureNumeric(func(args []Expression) Expression { return Bool(args[0].(Number) >= args[1].(Number)) })})
	env.Set("=", &Procedure{f: ensureNumeric(func(args []Expression) Expression { return Bool(args[0].(Number) == args[1].(Number)) })})

	// List
	env.Set("list", &Procedure{f: func(args []Expression) Expression { return List(args) }})
	env.Set("car", &Procedure{f: func(args []Expression) Expression { return args[0].(List)[0] }})
	env.Set("cdr", &Procedure{f: func(args []Expression) Expression { return args[0].(List)[1:] }})
	env.Set("list?", &Procedure{f: func(args []Expression) Expression { _, ok := args[0].(List); return Bool(ok) }})
	env.Set("nil?", &Procedure{f: func(args []Expression) Expression {
		if IsNil(args[0]) {
			return Bool(true)
		}
		if list, ok := args[0].(List); ok {
			return Bool(len(list) == 0)
		}
		return Bool(false)
	}})
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

	// String
	stringsToList := func(strings []string) List {
		list := List{}
		for _, s := range strings {
			list = append(list, String(s))
		}
		return list
	}
	ensureStrings := func(f func(args []Expression) Expression, numArgs int) func(args []Expression) Expression {
		return func(args []Expression) Expression {
			if len(args) != numArgs {
				return Error("invalid arguments")
			}
			for _, arg := range args {
				if _, ok := arg.(String); !ok {
					return Error("argument is not a string")
				}
			}
			return f(args)
		}
	}
	env.Set("strings/split", &Procedure{f: ensureStrings(func(args []Expression) Expression {
		return stringsToList(strings.Split(string(args[0].(String)), string(args[1].(String))))
	}, 2)})
	env.Set("strings/concat", &Procedure{f: ensureStrings(func(args []Expression) Expression {
		return args[0].(String) + args[1].(String)
	}, 2)})

	// Print
	env.Set("print", &Procedure{f: func(args []Expression) Expression { fmt.Println(args[0]); return Nil{} }})
	env.Set("str", &Procedure{f: func(args []Expression) Expression { return String(args[0].ExprToStr()) }})

	// Save and load
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
