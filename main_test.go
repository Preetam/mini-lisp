package main

import "testing"

func TestEverything(t *testing.T) {
	type singleCase struct {
		input  string
		output string
	}
	cases := []singleCase{
		{`'()`, `()`},
		{`(str 4)`, `"4"`},
		{`(print (str 4))`, `nil`},
		{`(+ 1 3)`, `4`},
		{`(+ 1)`, `!{invalid arguments}`},
		{`(/ 1 3) ; comment`, `0.3333333333333333`},
		{`(begin (+ 1 3) (* 4 2))`, `8`},
		{`(car '(1 2 3))`, `1`},
		{`(cdr '(1 2 3))`, `(2 3)`},
		{`(cons 'a '(1 2 3))`, `(a 1 2 3)`},
		{`(list 1 2 3 4)`, `(1 2 3 4)`},
		{`(nil? a)`, `#t`},
		{`(nil? '())`, `#t`},
		{`(nil? 5)`, `#f`},
		{`(list? '())`, `#t`},
		{`(let ((a 5)) (nil? a))`, `#f`},
		{`(let ((a 5)) a)`, `5`},
		{`(let ((a 5) (b 6)) a (* a b))`, `30`},
		{`(begin (define fact (lambda (n) (if (<= n 1) 1 (* n (fact (- n 1)))))) (fact 5))`, `120`},
		{`(begin (define sum2 (lambda (n acc) (if (= n 0) acc (sum2 (- n 1) (+ n acc))))) (sum2 1000 0))`, `500500`},
	}
	for _, testcase := range cases {
		expr, err := readFromTokens(tokenize(testcase.input))
		if err != nil {
			t.Error(err)
			continue
		}
		evaluated := eval(expr, DefaultEnv())
		if evaluated.ExprToStr() != testcase.output {
			t.Errorf("%s => %s, not %s", testcase.input, evaluated.ExprToStr(), testcase.output)
		}
	}
}
