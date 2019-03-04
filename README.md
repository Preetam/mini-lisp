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

;; call/cc examples
(call/cc (lambda (throw) (+ 5 (* 10 (call/cc (lambda (escape) (* 100 (throw 3))))))))
; 3
(call/cc (lambda (throw) (+ 5 (* 10 (call/cc (lambda (escape) (* 100 (escape 3))))))))
; 35

(define the-continuation #f)
(define test
  (lambda ()
    (begin
      (define i 0)
      (call/cc (lambda (k) (set! the-continuation k)))
      (set! i (+ i 1))
      i
    )
  )
)
(test)
; 1
(the-continuation)
; 2
(the-continuation)
; 3
(define another-continuation the-continuation)
(test)
; 1
(another-continuation)
; 4
```

## Features

* REPL
* Lambdas
* Tail recursion optimization
* call/cc

## References

* [mal - Make a Lisp](https://github.com/kanaka/mal/)
* [(How to Write a (Lisp) Interpreter (in Python))](http://norvig.com/lispy.html)
