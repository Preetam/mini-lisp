# mini-lisp

A small Lisp implementation in Go.

## Examples

```lisp
;; Factorial
(define fact (lambda (n) (if (<= n 1) 1 (* n (fact (- n 1))))))
(fact 5)
; 120

;; sum2 demonstrating tail recursion optimization
(define sum2 (lambda (n acc) (if (= n 0) acc (sum2 (- n 1) (+ n acc)))))
(sum2 1000 0)
; 500500
```

## Features

* REPL
* Lambdas
* Tail recursion optimization
* call/cc

## References

* [mal - Make a Lisp](https://github.com/kanaka/mal/)
* [(How to Write a (Lisp) Interpreter (in Python))](http://norvig.com/lispy.html)
