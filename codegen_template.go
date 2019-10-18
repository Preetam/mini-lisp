package main

var tmpl = `package main

import (
	"fmt"
	"strconv"
	"strings"
)

var _ = fmt.Println
var _ = strconv.FormatFloat
var _ = strings.Split

type Expression interface {
	ExprToStr() string
}

type Nil struct{}

func (_ Nil) ExprToStr() string { return "nil" }

func IsNil(e Expression) bool {
	var n Nil
	return e == n
}

type Number float64

func (n Number) ExprToStr() string { return strconv.FormatFloat(float64(n), 'g', -1, 64) }

type Bool bool

func (b Bool) ExprToStr() string {
	if b {
		return "#t"
	}
	return "#f"
}

type String string

func (s String) ExprToStr() string {
	return "\"" + strings.ReplaceAll(string(s), "\"", "\\\"") + "\""
}

type Error string

func (e Error) ExprToStr() string {
	return "!{" + string(e) + "}"
}

type Symbol string

func (s Symbol) ExprToStr() string { return string(s) }

type List []Expression

func (l List) ExprToStr() string {
	elemStrings := []string{}
	for _, e := range l {
		elemStrings = append(elemStrings, e.ExprToStr())
	}
	return "(" + strings.Join(elemStrings, " ") + ")"
}

type Procedure func(args []Expression) Expression

func (p Procedure) ExprToStr() string {
	return "#builtin"
}

type Environment struct {
	outer  *Environment
	values map[string]Expression
}

func NewEnvironment(outer *Environment) *Environment {
	return &Environment{
		outer:  outer,
		values: map[string]Expression{},
	}
}

func (env *Environment) Get(key string) Expression {
	if v, ok := env.values[key]; ok {
		return v
	}
	if env.outer == nil {
		return Nil{}
	}
	return env.outer.Get(key)
}

func (env *Environment) Set(key string, value Expression) {
	env.values[key] = value
}

func (env *Environment) SetOuter(key string, value Expression) {
	if _, ok := env.values[key]; ok {
		env.values[key] = value
		return
	}
	if env.outer != nil {
		env.outer.SetOuter(key, value)
	}
}

func main() {
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
	env.Set("+", Procedure(ensureNumeric(func(args []Expression) Expression { return Number(args[0].(Number) + args[1].(Number)) })))
	env.Set("-", Procedure(ensureNumeric(func(args []Expression) Expression { return Number(args[0].(Number) - args[1].(Number)) })))
	env.Set("*", Procedure(ensureNumeric(func(args []Expression) Expression { return Number(args[0].(Number) * args[1].(Number)) })))
	env.Set("/", Procedure(ensureNumeric(func(args []Expression) Expression { return Number(args[0].(Number) / args[1].(Number)) })))
	env.Set("<", Procedure(ensureNumeric(func(args []Expression) Expression { return Bool(args[0].(Number) < args[1].(Number)) })))
	env.Set("<=", Procedure(ensureNumeric(func(args []Expression) Expression { return Bool(args[0].(Number) <= args[1].(Number)) })))
	env.Set(">", Procedure(ensureNumeric(func(args []Expression) Expression { return Bool(args[0].(Number) > args[1].(Number)) })))
	env.Set(">=", Procedure(ensureNumeric(func(args []Expression) Expression { return Bool(args[0].(Number) >= args[1].(Number)) })))
	env.Set("=", Procedure(ensureNumeric(func(args []Expression) Expression { return Bool(args[0].(Number) == args[1].(Number)) })))

	fmt.Println({{code}})
}


`
