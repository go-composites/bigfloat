<p align="center"><img src="https://raw.githubusercontent.com/go-composites/brand/main/social/go-composites.png" alt="go-composites/bigfloat" width="720"></p>

# bigfloat

[![ci](https://github.com/go-composites/bigfloat/actions/workflows/ci.yml/badge.svg)](https://github.com/go-composites/bigfloat/actions/workflows/ci.yml)

An **arbitrary-precision floating-point** composite for **Composition-Oriented
Programming**. A `BigFloat` wraps Go's `math/big.Float` and carries a fixed
256-bit mantissa for every constructed value, so results are deterministic and
sums such as `0.1 + 0.2` are *exact* rather than the binary-float64
`0.30000000000000004`. Its arithmetic is exposed as **fallible operations that
return a `Result`** — so failures (the canonical example being a division by
zero) are *values*, never panics and never `nil`.

```golang
quotient := numerator.Div(denominator)
if quotient.HasError() {
    fmt.Println(quotient.Error().Message()) // "division by zero"
} else {
    fmt.Println(quotient.Payload().(BigFloat.Interface).ToGoString())
}
```

`BigFloat` follows the org's Null-Object / never-nil invariant (enforced by the
`nonnil` CI analyzer): the `NullBigFloat` variant in `src/null` satisfies the
same `Interface` and reports `IsNull() == true`.

## Install

```bash
export GOPRIVATE=github.com/go-composites GOPROXY=direct GOSUMDB=off
go get github.com/go-composites/bigfloat@main
```

## Usage

> [!NOTE] main.go

```golang
package main

import (
    "fmt"

    BigFloat "github.com/go-composites/bigfloat/src"
)

func main() {
    six := BigFloat.FromFloat64(6)
    two := BigFloat.FromFloat64(2)
    zero := BigFloat.FromFloat64(0)

    // Arithmetic returns a Result.
    sum := six.Add(two)
    fmt.Println(sum.Payload().(BigFloat.Interface).ToGoString()) // 8

    // Division by zero is a value, not a panic.
    div := six.Div(zero)
    fmt.Println("has error:", div.HasError())      // true
    fmt.Println(div.Error().Message())             // division by zero

    // Precision beyond float64: 0.1 + 0.2 is exact here.
    a := BigFloat.FromString("0.1").Payload().(BigFloat.Interface)
    b := BigFloat.FromString("0.2").Payload().(BigFloat.Interface)
    fmt.Println(a.Add(b).Payload().(BigFloat.Interface).ToGoString()) // 0.3

    fmt.Println(six.GreaterThan(two)) // true
    fmt.Println(six.Inspect())        // <BigFloat:0x... value=6>
}
```

```bash
$ go run .
```

## API

Constructors

- `FromFloat64(f float64) Interface` — build from a Go float64.
- `FromString(s string) Result.Interface` — parse a decimal string at 256 bits
  of precision; a `Result` carrying `Error.New(...)` when the input is not a
  valid float. This is how values with more significant digits than a float64
  can hold are constructed.
- `Null() Interface` — the `NullBigFloat` Null-Object (`IsNull() == true`).
- `null.New() Interface` — the importable `NullBigFloat` Null-Object.

Conversions

- `ToGoString() string` (shortest round-tripping decimal), `ToFloat64() float64`,
  `IsNull() bool`.

Arithmetic (each returns `Result.Interface`)

- `Add(other)` / `Sub(other)` / `Mul(other)` — sum, difference, product.
- `Div(other)` — quotient; a `Result` carrying `Error.New("division by zero")`
  when `other` is zero.
- `Abs()` / `Neg()` — absolute value and negation.

Every operation works on a fresh 256-bit `big.Float`, so operands are never
mutated.

Comparisons (each returns `bool`)

- `Equal(other)` / `LessThan(other)` / `GreaterThan(other)`.

Inspection

- `Inspect() string` — `<BigFloat:0x... value=...>`.

## License

BSD-3-Clause — see [LICENSE](./LICENSE).
