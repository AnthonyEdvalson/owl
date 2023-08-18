# Owl

Owl is an interpreted programming language that I made for fun to learn about compiler design and to get some experience with Go.

It shamelessly steals ideas from many different languages, here are a few highlights that I enjoy:

```
// Pattern matching arguments
let factorial = (0) => 1
              : (n) => n * f(n - 1)

// Simple list, set, and map definitions
let l = [1, 2, 3]
let s = {1, 2, 3}
let m = {a: 1, b: 2, c: 3}

// Spread syntax
let swap = (a, b, c) => c, b, a
let [x, ...xs] = swap(...l)

// First class functions
let plus_two = (x) => x + 2
let l_plus_two = l.map(plus_two)

// Null coalesce / access / call
let y = null ?? 4
let z = {a: 5, b: 6}?.c 
let w = {a: v => v + 1}.w?()
let v = l?[1]

// Multiple assignment / return values
x, y = y, x
x, y = (() => return 1, 2)()

// Operator overloading
complex = {real: 2, imag: 3}
complex::mul = (a, b) => {real: a.real * b.real - a.imag * b.imag, imag: a.real * b.imag + a.imag * b.real}
complex * complex == {real: -5, imag: 12}
```

# Building

Bulding the code is as simple as running `build.sh` or  `build.bat` depending on your operating system of choice. They both just run `go build` and copy the standard library into the bin folder. You can optionally add the bin folder to $PATH so you can execute the `owl` command more easily.