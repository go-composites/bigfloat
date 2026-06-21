package main

import (
	"fmt"

	BigFloat "github.com/go-composites/bigfloat/src"
	Error "github.com/go-composites/error/src"
	Result "github.com/go-composites/result/src"
)

func report(label string, result Result.Interface) {
	if result.HasError() {
		fmt.Printf("%s -> error: %s\n", label, result.Error().Message())
		return
	}
	fmt.Printf("%s -> %s\n", label, result.Payload().(BigFloat.Interface).ToGoString())
}

func mustParse(s string) BigFloat.Interface {
	r := BigFloat.FromString(s)
	return r.Payload().(BigFloat.Interface)
}

func main() {
	six := BigFloat.FromFloat64(6)
	two := BigFloat.FromFloat64(2)
	zero := BigFloat.FromFloat64(0)

	report("6 + 2", six.Add(two))
	report("6 - 2", six.Sub(two))
	report("6 * 2", six.Mul(two))
	report("6 / 2", six.Div(two))

	// The canonical Result use-case: division by zero is a value, not a panic.
	divByZero := six.Div(zero)
	fmt.Println("6 / 0 has error:", divByZero.HasError())
	report("6 / 0", divByZero)

	// Errors are first-class values.
	var _ Error.Interface = divByZero.Error()

	// Precision beyond float64: 0.1 + 0.2 is exact here, not 0.30000000000000004.
	report("0.1 + 0.2", mustParse("0.1").Add(mustParse("0.2")))

	// A 49-significant-digit fraction survives round-trip; a float64 (~15-17
	// significant digits) could not hold it.
	report("49-digit fraction", Result.New(Result.WithPayload(
		mustParse("0.1234567890123456789012345678901234567890123456789"))))

	fmt.Println("6 == 2 :", six.Equal(two))
	fmt.Println("6 < 2  :", six.LessThan(two))
	fmt.Println("6 > 2  :", six.GreaterThan(two))
	fmt.Println(six.Inspect())
}
