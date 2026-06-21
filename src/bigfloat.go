package BigFloat

import (
	"fmt"
	"math/big"

	Error "github.com/go-composites/error/src"
	MethodNotImplementedError "github.com/go-composites/error/src/method_not_implemented"
	Result "github.com/go-composites/result/src"
)

// prec is the fixed mantissa precision (in bits) of every constructed
// BigFloat. Pinning it makes results deterministic across operations and
// platforms.
const prec = 256

/*
BigFloat is an arbitrary-precision floating-point composite over a
math/big.Float.

It mirrors Ruby's BigDecimal-style precision: every value carries 256 bits of
mantissa, so sums such as 0.1 + 0.2 are exact rather than subject to binary
float64 rounding. Its fallible operations (notably Div) return a
Result.Interface so that failures — such as a division by zero — are values
rather than panics, and they never return a bare nil.
*/
type Interface interface {
	ToGoString() string
	ToFloat64() float64
	ToInt64() int64
	IsNull() bool
	IsZero() bool
	Add(Interface) Result.Interface
	Sub(Interface) Result.Interface
	Mul(Interface) Result.Interface
	Div(Interface) Result.Interface
	Abs() Result.Interface
	Neg() Result.Interface
	Floor() Result.Interface
	Ceil() Result.Interface
	Round() Result.Interface
	Power(int) Result.Interface
	Sqrt() Result.Interface
	Equal(Interface) bool
	LessThan(Interface) bool
	GreaterThan(Interface) bool
	Inspect() String
}

// String is the lightweight inspection representation of a BigFloat.
type String = string

type data struct {
	value *big.Float
}

// newFloat returns a fresh big.Float pinned to the package precision.
func newFloat() *big.Float {
	return new(big.Float).SetPrec(prec)
}

/*
FromFloat64 is the BigFloat constructor from a Go float64.

	x := BigFloat.FromFloat64(0.5) // 0.5
*/
func FromFloat64(f float64) Interface {
	return &data{value: newFloat().SetFloat64(f)}
}

/*
FromString parses a decimal string into a BigFloat at 256 bits of precision.

It returns a Result whose payload is the parsed BigFloat. When the input is not
a valid floating-point literal the Result carries an Error instead of a
payload — the parse never panics and never returns nil. This is how values with
more significant digits than a float64 can represent are constructed.

	r := BigFloat.FromString("0.12345678901234567890123456789012345678901234567890")
	if !r.HasError() {
	    x := r.Payload().(BigFloat.Interface)
	}
*/
func FromString(s string) Result.Interface {
	value, _, err := big.ParseFloat(s, 10, prec, big.ToNearestEven)
	if err != nil {
		return Result.New(
			Result.WithError(
				Error.New("invalid float: " + s),
			),
		)
	}
	return Result.New(
		Result.WithPayload(
			&data{value: value},
		),
	)
}

/*
Null returns the Null-Object variant of BigFloat.

It is defined in src/null; this thin re-export keeps a Null next to the
concrete constructors. The returned value satisfies Interface and reports
IsNull() == true.
*/
func Null() Interface {
	return newNull()
}

/*
ToGoString returns the shortest decimal representation that round-trips the
value, using math/big.Float.Text('g', -1).
*/
func (d data) ToGoString() string {
	return d.value.Text('g', -1)
}

/*
ToFloat64 returns the value as a Go float64.

When the value does not fit in a float64 the result follows
math/big.Float.Float64 (the nearest float64, possibly an infinity), so callers
handling arbitrary precision should prefer ToGoString.
*/
func (d data) ToFloat64() float64 {
	f, _ := d.value.Float64()
	return f
}

/*
ToInt64 returns the value truncated toward zero as a Go int64.

It follows math/big.Float.Int64: when the value does not fit in an int64 the
result saturates at the int64 bounds.
*/
func (d data) ToInt64() int64 {
	i, _ := d.value.Int64()
	return i
}

/*
IsNull reports whether the BigFloat is the Null-Object variant.

A concrete BigFloat is never null.
*/
func (d data) IsNull() bool {
	return false
}

/*
IsZero reports whether the value is exactly zero.
*/
func (d data) IsZero() bool {
	return d.value.Sign() == 0
}

/*
Add returns a Result whose payload is the sum of the receiver and other.

A fresh big.Float at 256 bits backs the payload; the operands are never
mutated.
*/
func (d data) Add(other Interface) Result.Interface {
	return payload(
		newFloat().Add(d.value, fromInterface(other)),
	)
}

/*
Sub returns a Result whose payload is the difference of the receiver and other.
*/
func (d data) Sub(other Interface) Result.Interface {
	return payload(
		newFloat().Sub(d.value, fromInterface(other)),
	)
}

/*
Mul returns a Result whose payload is the product of the receiver and other.
*/
func (d data) Mul(other Interface) Result.Interface {
	return payload(
		newFloat().Mul(d.value, fromInterface(other)),
	)
}

/*
Div returns a Result whose payload is the quotient of the receiver and other.

When other is zero the Result carries an Error ("division by zero") instead of
a payload — the division never panics and never returns nil.
*/
func (d data) Div(other Interface) Result.Interface {
	rhs := fromInterface(other)
	if rhs.Sign() == 0 {
		return Result.New(
			Result.WithError(
				Error.New("division by zero"),
			),
		)
	}
	return payload(
		newFloat().Quo(d.value, rhs),
	)
}

/*
Abs returns a Result whose payload is the absolute value of the receiver.
*/
func (d data) Abs() Result.Interface {
	return payload(
		newFloat().Abs(d.value),
	)
}

/*
Neg returns a Result whose payload is the negation of the receiver.
*/
func (d data) Neg() Result.Interface {
	return payload(
		newFloat().Neg(d.value),
	)
}

/*
Floor returns a Result whose payload is the largest integer less than or equal
to the receiver.
*/
func (d data) Floor() Result.Interface {
	return payload(floorFloat(d.value))
}

/*
Ceil returns a Result whose payload is the smallest integer greater than or
equal to the receiver.
*/
func (d data) Ceil() Result.Interface {
	return payload(ceilFloat(d.value))
}

/*
Round returns a Result whose payload is the receiver rounded to the nearest
integer, with halves rounded away from zero.
*/
func (d data) Round() Result.Interface {
	return payload(roundFloat(d.value))
}

/*
Power returns a Result whose payload is the receiver raised to the integer
power n.

A negative exponent yields the reciprocal; 0**0 is 1. When the base is zero and
the exponent is negative the Result carries an Error ("division by zero")
instead of a payload — it never panics and never returns nil.
*/
func (d data) Power(n int) Result.Interface {
	if n < 0 && d.value.Sign() == 0 {
		return Result.New(
			Result.WithError(
				Error.New("division by zero"),
			),
		)
	}
	return payload(powFloat(d.value, n))
}

/*
Sqrt returns a Result whose payload is the square root of the receiver.

A negative receiver yields a Result carrying an Error ("square root of a
negative number") instead of a payload — it never panics and never returns nil.
*/
func (d data) Sqrt() Result.Interface {
	if d.value.Sign() < 0 {
		return Result.New(
			Result.WithError(
				Error.New("square root of a negative number"),
			),
		)
	}
	return payload(newFloat().Sqrt(d.value))
}

/*
Equal reports whether the receiver and other hold the same value.
*/
func (d data) Equal(other Interface) bool {
	return d.value.Cmp(fromInterface(other)) == 0
}

/*
LessThan reports whether the receiver is strictly less than other.
*/
func (d data) LessThan(other Interface) bool {
	return d.value.Cmp(fromInterface(other)) < 0
}

/*
GreaterThan reports whether the receiver is strictly greater than other.
*/
func (d data) GreaterThan(other Interface) bool {
	return d.value.Cmp(fromInterface(other)) > 0
}

/*
Inspect returns a one-line representation of the BigFloat with its address and
value — mirroring the style of the other composites.
*/
func (d data) Inspect() String {
	return fmt.Sprintf(
		"<BigFloat:%p value=%s>",
		&d, d.value.Text('g', -1),
	)
}

// nullData is the Null-Object variant returned by Null(). The importable
// NullBigFloat package in src/null mirrors it; this copy keeps a Null next to
// the concrete constructors without creating an import cycle.
type nullData struct{}

func newNull() Interface {
	return &nullData{}
}

func nullNotImplemented(methodName string) Result.Interface {
	return Result.New(
		Result.WithError(
			MethodNotImplementedError.New(methodName),
		),
	)
}

func (nullData) ToGoString() string             { return `` }
func (nullData) ToFloat64() float64             { return 0 }
func (nullData) ToInt64() int64                 { return 0 }
func (nullData) IsNull() bool                   { return true }
func (nullData) IsZero() bool                   { return false }
func (nullData) Add(Interface) Result.Interface { return nullNotImplemented(`Add`) }
func (nullData) Sub(Interface) Result.Interface { return nullNotImplemented(`Sub`) }
func (nullData) Mul(Interface) Result.Interface { return nullNotImplemented(`Mul`) }
func (nullData) Div(Interface) Result.Interface { return nullNotImplemented(`Div`) }
func (nullData) Abs() Result.Interface          { return nullNotImplemented(`Abs`) }
func (nullData) Neg() Result.Interface          { return nullNotImplemented(`Neg`) }
func (nullData) Floor() Result.Interface        { return nullNotImplemented(`Floor`) }
func (nullData) Ceil() Result.Interface         { return nullNotImplemented(`Ceil`) }
func (nullData) Round() Result.Interface        { return nullNotImplemented(`Round`) }
func (nullData) Power(int) Result.Interface     { return nullNotImplemented(`Power`) }
func (nullData) Sqrt() Result.Interface         { return nullNotImplemented(`Sqrt`) }
func (nullData) Equal(other Interface) bool     { return other.IsNull() }
func (nullData) LessThan(Interface) bool        { return false }
func (nullData) GreaterThan(Interface) bool     { return false }
func (nullData) Inspect() String                { return `<NullBigFloat>` }

// floorFloat returns a fresh big.Float holding the largest integer less than
// or equal to value.
func floorFloat(value *big.Float) *big.Float {
	truncated, frac := truncate(value)
	if frac && value.Sign() < 0 {
		truncated.Sub(truncated, newFloat().SetInt64(1))
	}
	return truncated
}

// ceilFloat returns a fresh big.Float holding the smallest integer greater
// than or equal to value.
func ceilFloat(value *big.Float) *big.Float {
	truncated, frac := truncate(value)
	if frac && value.Sign() > 0 {
		truncated.Add(truncated, newFloat().SetInt64(1))
	}
	return truncated
}

// roundFloat returns a fresh big.Float holding value rounded to the nearest
// integer, halves away from zero.
func roundFloat(value *big.Float) *big.Float {
	half := newFloat().SetFloat64(0.5)
	if value.Sign() < 0 {
		shifted := newFloat().Sub(value, half)
		result, _ := truncate(shifted)
		return result
	}
	shifted := newFloat().Add(value, half)
	result, _ := truncate(shifted)
	return result
}

// truncate returns a fresh big.Float holding value truncated toward zero, and
// whether value had a fractional part.
func truncate(value *big.Float) (*big.Float, bool) {
	intPart, _ := value.Int(nil)
	truncated := newFloat().SetInt(intPart)
	frac := truncated.Cmp(value) != 0
	return truncated, frac
}

// powFloat returns a fresh big.Float holding base raised to the integer power
// n via repeated multiplication. A negative n yields the reciprocal; n == 0
// yields 1.
func powFloat(base *big.Float, n int) *big.Float {
	exp := n
	if exp < 0 {
		exp = -exp
	}
	result := newFloat().SetInt64(1)
	factor := newFloat().Set(base)
	for exp > 0 {
		if exp&1 == 1 {
			result.Mul(result, factor)
		}
		exp >>= 1
		if exp > 0 {
			factor.Mul(factor, factor)
		}
	}
	if n < 0 {
		return newFloat().Quo(newFloat().SetInt64(1), result)
	}
	return result
}

// payload wraps a fresh big.Float in a success Result.
func payload(value *big.Float) Result.Interface {
	return Result.New(
		Result.WithPayload(
			&data{value: value},
		),
	)
}

// fromInterface extracts a *big.Float from any BigFloat.Interface, parsing its
// decimal string when the concrete type is unknown (e.g. the Null-Object).
// The returned big.Float is always a fresh copy, so operands are never shared.
func fromInterface(other Interface) *big.Float {
	if d, ok := other.(*data); ok {
		return newFloat().Set(d.value)
	}
	value, _, err := big.ParseFloat(other.ToGoString(), 10, prec, big.ToNearestEven)
	if err != nil {
		return newFloat()
	}
	return value
}
