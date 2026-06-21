package NullBigFloat

import (
	BigFloat "github.com/go-composites/bigfloat/src"
	MethodNotImplementedError "github.com/go-composites/error/src/method_not_implemented"
	Result "github.com/go-composites/result/src"
)

/*
NullBigFloat is the Null-Object variant of BigFloat.

It satisfies BigFloat.Interface so callers never have to test for a bare nil:
its value is zero, its arithmetic yields a Result carrying a
"method not implemented" Error, its comparisons are false (except Equal against
another null), and IsNull() returns true.
*/
type Interface interface {
	BigFloat.Interface
}

type data struct{}

/*
New returns a NullBigFloat.
*/
func New() Interface {
	return &data{}
}

func (d data) ToGoString() string {
	return ``
}

func (d data) ToFloat64() float64 {
	return 0
}

func (d data) ToInt64() int64 {
	return 0
}

func (d data) IsNull() bool {
	return true
}

func (d data) IsZero() bool {
	return false
}

func notImplemented(methodName string) Result.Interface {
	return Result.New(
		Result.WithError(
			MethodNotImplementedError.New(methodName),
		),
	)
}

func (d data) Add(BigFloat.Interface) Result.Interface {
	return notImplemented(`Add`)
}

func (d data) Sub(BigFloat.Interface) Result.Interface {
	return notImplemented(`Sub`)
}

func (d data) Mul(BigFloat.Interface) Result.Interface {
	return notImplemented(`Mul`)
}

func (d data) Div(BigFloat.Interface) Result.Interface {
	return notImplemented(`Div`)
}

func (d data) Abs() Result.Interface {
	return notImplemented(`Abs`)
}

func (d data) Neg() Result.Interface {
	return notImplemented(`Neg`)
}

func (d data) Floor() Result.Interface {
	return notImplemented(`Floor`)
}

func (d data) Ceil() Result.Interface {
	return notImplemented(`Ceil`)
}

func (d data) Round() Result.Interface {
	return notImplemented(`Round`)
}

func (d data) Power(int) Result.Interface {
	return notImplemented(`Power`)
}

func (d data) Sqrt() Result.Interface {
	return notImplemented(`Sqrt`)
}

func (d data) Equal(other BigFloat.Interface) bool {
	return other.IsNull()
}

func (d data) LessThan(BigFloat.Interface) bool {
	return false
}

func (d data) GreaterThan(BigFloat.Interface) bool {
	return false
}

func (d data) Inspect() BigFloat.String {
	return `<NullBigFloat>`
}
