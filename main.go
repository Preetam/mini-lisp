package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
)

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
	return `"` + strings.ReplaceAll(string(s), `"`, `\"`) + `"`
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

type Procedure struct {
	args []Symbol
	body Expression
	env  *Environment
	f    func(args []Expression) Expression
}

func (p *Procedure) ExprToStr() string {
	if p.f != nil {
		return "#built-in-function"
	}
	args := List{}
	for _, x := range p.args {
		args = append(args, x)
	}
	return "(lambda " + args.ExprToStr() + " " + p.body.ExprToStr() + ")"
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

func (env *Environment) Get(key string) (Expression, bool) {
	if v, ok := env.values[key]; ok {
		return v, ok
	}
	if env.outer == nil {
		return Nil{}, false
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

func pop(a *[]string) string {
	v := (*a)[0]
	*a = (*a)[1:]
	return v
}

func tokenize(str string) *[]string {
	tokens := []string{}
	re := regexp.MustCompile(`[\s,]*(~@|[\[\]{}()'` + "`" +
		`~^@]|"(?:\\.|[^\\"])*"|;.*|[^\s\[\]{}('"` + "`" +
		`,;)]*)`)
	for _, match := range re.FindAllStringSubmatch(str, -1) {
		if (match[1] == "") ||
			// comment
			(match[1][0] == ';') {
			continue
		}
		tokens = append(tokens, match[1])
	}
	return &tokens
}

func atom(token string) Expression {
	switch token {
	case "#t":
		return Bool(true)
	case "#f":
		return Bool(false)
	case "nil":
		return Nil{}
	}
	if token[0] == '"' {
		return String(strings.ReplaceAll(strings.Trim(token, `"`), `\"`, `"`))
	}
	f, err := strconv.ParseFloat(token, 64)
	if err == nil {
		return Number(f)
	}
	return Symbol(token)
}

func readFromTokens(tokens *[]string) (Expression, error) {
	if len(*tokens) == 0 {
		return nil, errors.New("unexpected EOF")
	}
	token := pop(tokens)
	switch token {
	case "'":
		// '... => (quote ...)
		quoted, err := readFromTokens(tokens)
		if err != nil {
			return nil, err
		}
		return List{atom("quote"), quoted}, nil
	case "(":
		if len(*tokens) == 0 {
			return nil, errors.New("unexpected EOF")
		}
		list := List{}
		for (*tokens)[0] != ")" {
			expr, err := readFromTokens(tokens)
			if err != nil {
				return nil, err
			}
			list = append(list, expr)
			if len(*tokens) == 0 {
				return nil, errors.New("unexpected EOF")
			}
		}
		pop(tokens)

		if len(list) > 0 && list[0] == Symbol("define") {
			// (define (f ...) (...)) => (define f (lambda (...) (...)))
			if argsList, ok := list[1].(List); ok {
				return List{atom("define"), argsList[0], List{atom("lambda"), argsList[1:], list[2]}}, nil
			}
		}

		return list, nil
	case ")":
		return nil, errors.New("unexpected ')'")
	default:
		return atom(token), nil
	}
}

func eval(exp Expression, env *Environment) Expression {
	for {
		switch exp.(type) {
		case Symbol:
			v, _ := env.Get(string(exp.(Symbol)))
			return v
		case Number, Bool, String, Nil:
			return exp
		case List:
			listExp := exp.(List)
			if len(listExp) == 0 {
				return listExp
			}
			switch listExp[0] {
			case Symbol("begin"):
				var ret Expression
				for _, x := range listExp[1:] {
					ret = eval(x, env)
				}
				return ret
			case Symbol("quote"):
				return listExp[1]
			case Symbol("define"):
				val := eval(listExp[2], env)
				env.Set(string(listExp[1].(Symbol)), val)
				return Nil{}
			case Symbol("set!"):
				val := eval(listExp[2], env)
				env.SetOuter(string(listExp[1].(Symbol)), val)
				return Nil{}
			case Symbol("let"):
				newEnv := NewEnvironment(env)
				bindingsList := listExp[1].(List)
				for _, binding := range bindingsList {
					bindingVal := eval(binding.(List)[1], env)
					newEnv.Set(string(binding.(List)[0].(Symbol)), bindingVal)
				}
				var ret Expression
				for _, x := range listExp[2:] {
					ret = eval(x, newEnv)
				}
				return ret
			case Symbol("if"):
				test := eval(listExp[1], env)
				if b, ok := test.(Bool); (ok && !bool(b)) || IsNil(test) {
					if len(listExp) < 4 {
						return Nil{}
					}
					exp = listExp[3]
					continue
				}
				exp = listExp[2]
				continue
			case Symbol("lambda"):
				args := []Symbol{}
				for _, x := range listExp[1].(List) {
					args = append(args, x.(Symbol))
				}
				return &Procedure{
					args: args,
					body: listExp[2],
					env:  env,
				}
			default:
				procExp := eval(listExp[0], env)
				proc, ok := procExp.(*Procedure)
				if !ok {
					return Error("argument not a function")
				}
				args := []Expression{}
				for _, argExp := range listExp[1:] {
					evalArgExp := eval(argExp, env)
					args = append(args, evalArgExp)
				}
				if proc.f != nil {
					return proc.f(args)
				} else {
					env = NewEnvironment(proc.env)
					for i, x := range proc.args {
						env.Set(string(x), args[i])
					}
					exp = proc.body
				}
			}
		}
	}
}

func main() {
	env := DefaultEnv()
	if len(os.Args) > 1 {
		filename := os.Args[1]
		content, err := ioutil.ReadFile(string(filename))
		if err != nil {
			return
		}
		tokens := tokenize(string(content))
		for len(*tokens) > 0 {
			expr, err := readFromTokens(tokens)
			if err != nil {
				return
			}
			eval(expr, env)
		}
		return
	}
	rl, err := readline.New("mini-lisp> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	buf := &bytes.Buffer{}
	for {
		line, err := rl.Readline()
		if err != nil {
			return
		}
		if strings.TrimSpace(line) == "" {
			continue
		}

		buf.WriteString(line)
		expression, err := readFromTokens(tokenize(buf.String()))
		if err != nil {
			if err.Error() == "unexpected EOF" {
				rl.SetPrompt("| ")
				continue
			}
			fmt.Println("error:", err)
			return
		}
		buf.Reset()
		rl.SetPrompt("mini-lisp> ")
		result := eval(expression, env)
		fmt.Println(result.ExprToStr())
	}
}
