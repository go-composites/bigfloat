package BigFloat_test

import (
	BigFloat "github.com/go-composites/bigfloat/src"
	Result "github.com/go-composites/result/src"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// payloadOf unwraps a success Result into a BigFloat.Interface.
func payloadOf(r interface {
	HasError() bool
	Payload() interface{}
}) BigFloat.Interface {
	gomega.ExpectWithOffset(1, r.HasError()).To(gomega.BeFalse())
	return r.Payload().(BigFloat.Interface)
}

// foreign is a BigFloat.Interface implementation that is NOT the package's own
// concrete type. It is used to exercise the string-bridging path of
// fromInterface with a value that DOES parse as a float (the success branch).
type foreign struct{ s string }

func (f foreign) ToGoString() string                    { return f.s }
func (foreign) ToFloat64() float64                      { return 0 }
func (foreign) IsNull() bool                            { return false }
func (foreign) Add(BigFloat.Interface) Result.Interface { return nil }
func (foreign) Sub(BigFloat.Interface) Result.Interface { return nil }
func (foreign) Mul(BigFloat.Interface) Result.Interface { return nil }
func (foreign) Div(BigFloat.Interface) Result.Interface { return nil }
func (foreign) Abs() Result.Interface                   { return nil }
func (foreign) Neg() Result.Interface                   { return nil }
func (foreign) Equal(BigFloat.Interface) bool           { return false }
func (foreign) LessThan(BigFloat.Interface) bool        { return false }
func (foreign) GreaterThan(BigFloat.Interface) bool     { return false }
func (foreign) Inspect() BigFloat.String                { return `` }

var _ = ginkgo.Describe("BigFloat", func() {

	ginkgo.Describe("constructors", func() {
		ginkgo.It("builds from a Go float64", func() {
			x := BigFloat.FromFloat64(0.5)
			gomega.Expect(x.ToFloat64()).To(gomega.Equal(0.5))
			gomega.Expect(x.ToGoString()).To(gomega.Equal("0.5"))
			gomega.Expect(x.IsNull()).To(gomega.BeFalse())
		})
		ginkgo.It("parses a valid decimal string", func() {
			r := BigFloat.FromString("3.25")
			gomega.Expect(r.HasError()).To(gomega.BeFalse())
			gomega.Expect(r.Payload().(BigFloat.Interface).ToFloat64()).
				To(gomega.Equal(3.25))
		})
		ginkgo.It("returns an error Result on a bad string", func() {
			r := BigFloat.FromString("not-a-number")
			gomega.Expect(r.HasError()).To(gomega.BeTrue())
			gomega.Expect(r.Error().Message()).To(gomega.ContainSubstring("invalid float"))
		})
		ginkgo.It("exposes a Null-Object", func() {
			x := BigFloat.Null()
			gomega.Expect(x.IsNull()).To(gomega.BeTrue())
		})
	})

	ginkgo.Describe("arbitrary precision", func() {
		ginkgo.It("adds 0.1 and 0.2 exactly, not as a binary float64", func() {
			a := payloadOf(BigFloat.FromString("0.1"))
			b := payloadOf(BigFloat.FromString("0.2"))
			sum := payloadOf(a.Add(b))
			gomega.Expect(sum.ToGoString()).To(gomega.Equal("0.3"))
		})
		ginkgo.It("preserves a 49-digit fraction beyond float64 range", func() {
			// 49 significant digits — far more than a float64's ~15-17. The
			// trailing zero of a 50-digit literal is not significant, so the
			// shortest round-tripping decimal carries 49 digits.
			s := "0.1234567890123456789012345678901234567890123456789"
			x := payloadOf(BigFloat.FromString(s))
			gomega.Expect(x.ToGoString()).To(gomega.Equal(s))
		})
	})

	ginkgo.Describe("arithmetic", func() {
		var six = BigFloat.FromFloat64(6)
		var two = BigFloat.FromFloat64(2)

		ginkgo.It("adds", func() {
			gomega.Expect(payloadOf(six.Add(two)).ToFloat64()).To(gomega.Equal(8.0))
		})
		ginkgo.It("subtracts", func() {
			gomega.Expect(payloadOf(six.Sub(two)).ToFloat64()).To(gomega.Equal(4.0))
		})
		ginkgo.It("multiplies", func() {
			gomega.Expect(payloadOf(six.Mul(two)).ToFloat64()).To(gomega.Equal(12.0))
		})
		ginkgo.It("divides", func() {
			gomega.Expect(payloadOf(six.Div(two)).ToFloat64()).To(gomega.Equal(3.0))
		})
		ginkgo.It("does not mutate its operands", func() {
			_ = six.Add(two)
			gomega.Expect(six.ToFloat64()).To(gomega.Equal(6.0))
			gomega.Expect(two.ToFloat64()).To(gomega.Equal(2.0))
		})

		ginkgo.Describe("division by zero", func() {
			ginkgo.It("returns a Result carrying an error instead of panicking", func() {
				r := six.Div(BigFloat.FromFloat64(0))
				gomega.Expect(r.HasError()).To(gomega.BeTrue())
				gomega.Expect(r.Error().Message()).To(gomega.Equal("division by zero"))
			})
		})

		ginkgo.Describe("absolute value", func() {
			ginkgo.It("makes a negative number positive", func() {
				gomega.Expect(payloadOf(BigFloat.FromFloat64(-7.5).Abs()).ToFloat64()).
					To(gomega.Equal(7.5))
			})
			ginkgo.It("leaves a positive number unchanged", func() {
				gomega.Expect(payloadOf(BigFloat.FromFloat64(7.5).Abs()).ToFloat64()).
					To(gomega.Equal(7.5))
			})
		})

		ginkgo.Describe("negation", func() {
			ginkgo.It("negates a positive number", func() {
				gomega.Expect(payloadOf(BigFloat.FromFloat64(7.5).Neg()).ToFloat64()).
					To(gomega.Equal(-7.5))
			})
			ginkgo.It("negates a negative number", func() {
				gomega.Expect(payloadOf(BigFloat.FromFloat64(-7.5).Neg()).ToFloat64()).
					To(gomega.Equal(7.5))
			})
		})
	})

	ginkgo.Describe("operations against a Null operand", func() {
		var six = BigFloat.FromFloat64(6)
		var null = BigFloat.Null()

		ginkgo.It("treats a Null operand as zero in addition", func() {
			gomega.Expect(payloadOf(six.Add(null)).ToFloat64()).To(gomega.Equal(6.0))
		})
		ginkgo.It("guards division by a Null operand (zero)", func() {
			gomega.Expect(six.Div(null).HasError()).To(gomega.BeTrue())
		})
		ginkgo.It("bridges a foreign Interface through its decimal string", func() {
			gomega.Expect(payloadOf(six.Add(foreign{s: "4"})).ToFloat64()).
				To(gomega.Equal(10.0))
		})
		ginkgo.It("treats an unparsable foreign operand as zero", func() {
			gomega.Expect(payloadOf(six.Add(foreign{s: "xx"})).ToFloat64()).
				To(gomega.Equal(6.0))
		})
	})

	ginkgo.Describe("comparisons", func() {
		var six = BigFloat.FromFloat64(6)
		var two = BigFloat.FromFloat64(2)

		ginkgo.It("reports equality both ways", func() {
			gomega.Expect(six.Equal(six)).To(gomega.BeTrue())
			gomega.Expect(six.Equal(two)).To(gomega.BeFalse())
		})
		ginkgo.It("reports less-than both ways", func() {
			gomega.Expect(two.LessThan(six)).To(gomega.BeTrue())
			gomega.Expect(six.LessThan(two)).To(gomega.BeFalse())
		})
		ginkgo.It("reports greater-than both ways", func() {
			gomega.Expect(six.GreaterThan(two)).To(gomega.BeTrue())
			gomega.Expect(two.GreaterThan(six)).To(gomega.BeFalse())
		})
	})

	ginkgo.Describe("inspection", func() {
		ginkgo.It("renders a BigFloat", func() {
			gomega.Expect(BigFloat.FromFloat64(6).Inspect()).
				To(gomega.ContainSubstring("value=6"))
		})
	})

	ginkgo.Describe("the package-local Null-Object", func() {
		var n = BigFloat.Null()

		ginkgo.It("converts to zero values", func() {
			gomega.Expect(n.ToFloat64()).To(gomega.Equal(0.0))
			gomega.Expect(n.ToGoString()).To(gomega.Equal(``))
		})
		ginkgo.It("returns error Results for every arithmetic method", func() {
			gomega.Expect(n.Add(n).HasError()).To(gomega.BeTrue())
			gomega.Expect(n.Sub(n).HasError()).To(gomega.BeTrue())
			gomega.Expect(n.Mul(n).HasError()).To(gomega.BeTrue())
			gomega.Expect(n.Div(n).HasError()).To(gomega.BeTrue())
			gomega.Expect(n.Abs().HasError()).To(gomega.BeTrue())
			gomega.Expect(n.Neg().HasError()).To(gomega.BeTrue())
		})
		ginkgo.It("compares as a Null-Object", func() {
			gomega.Expect(n.Equal(BigFloat.Null())).To(gomega.BeTrue())
			gomega.Expect(n.Equal(BigFloat.FromFloat64(0))).To(gomega.BeFalse())
			gomega.Expect(n.LessThan(BigFloat.FromFloat64(1))).To(gomega.BeFalse())
			gomega.Expect(n.GreaterThan(BigFloat.FromFloat64(-1))).To(gomega.BeFalse())
		})
		ginkgo.It("inspects as the null marker", func() {
			gomega.Expect(n.Inspect()).To(gomega.Equal(`<NullBigFloat>`))
		})
	})
})
