package main

import "testing"

func TestEverything(t *testing.T) {
	type singleCase struct {
		input  string
		output string
	}
	cases := []singleCase{
		// Simple stuff
		{`'()`, `()`},
		{`(str 4)`, `"4"`},
		{`(print (str 4))`, `nil`},
		{`(+ 1 3)`, `4`},
		{`(+ 1)`, `!{invalid arguments}`},
		{`(/ 1 3) ; comment`, `0.3333333333333333`},

		// Lists
		{`(car '(1 2 3))`, `1`},
		{`(cdr '(1 2 3))`, `(2 3)`},
		{`(cons 'a '(1 2 3))`, `(a 1 2 3)`},
		{`(list 1 2 3 4)`, `(1 2 3 4)`},
		{`(nil? a)`, `#t`},
		{`(nil? nil)`, `#t`},
		{`(nil? '())`, `#t`},
		{`(nil? 5)`, `#f`},
		{`(list? '())`, `#t`},

		// Let
		{`(let ((a 5)) (nil? a))`, `#f`},
		{`(let ((a 5)) a)`, `5`},
		{`(let ((a 5) (b 6)) a (* a b))`, `30`},

		// Strings
		{`(strings/split "foo/bar" "/")`, `("foo" "bar")`},
		{`(strings/split "foo/bar")`, `!{invalid arguments}`},
		{`(strings/split "foo/bar" 1)`, `!{argument is not a string}`},
		{`(strings/concat "foo" "bar")`, `"foobar"`},

		// Complicated stuff
		{`(begin (+ 1 3) (* 4 2))`, `8`},
		{`(begin (define fact (lambda (n) (if (<= n 1) 1 (* n (fact (- n 1)))))) (fact 5))`, `120`},
		{`(begin (define sum2 (lambda (n acc) (if (= n 0) acc (sum2 (- n 1) (+ n acc))))) (sum2 1000 0))`, `500500`},
		{`(let ((x 0)) (begin
			(define (not v) (if v #f #t))
			(define (map l f)
				(if (not (nil? l)) (begin
					(f (car l))
					(map (cdr l) f)
				))
			)
			(map '(1 2 3) (lambda (v) (set! x (+ x v))))
			x
		))`, `6`}, // 1 + 2 + 3
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
