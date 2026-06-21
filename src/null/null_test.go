package NullBigFloat_test

import (
	BigFloat "github.com/go-composites/bigfloat/src"
	NullBigFloat "github.com/go-composites/bigfloat/src/null"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("NullBigFloat", func() {
	var n NullBigFloat.Interface
	ginkgo.BeforeEach(func() {
		n = NullBigFloat.New()
	})

	ginkgo.It("satisfies the BigFloat interface", func() {
		var _ BigFloat.Interface = n
	})
	ginkgo.It("reports IsNull() true", func() {
		gomega.Expect(n.IsNull()).To(gomega.BeTrue())
	})
	ginkgo.It("converts to zero values", func() {
		gomega.Expect(n.ToFloat64()).To(gomega.Equal(0.0))
		gomega.Expect(n.ToGoString()).To(gomega.Equal(``))
	})

	ginkgo.It("Add returns an error result", func() {
		r := n.Add(BigFloat.FromFloat64(0))
		gomega.Expect(r.HasError()).To(gomega.BeTrue())
		gomega.Expect(r.Error().Message()).To(gomega.ContainSubstring("Add"))
	})
	ginkgo.It("Sub returns an error result", func() {
		r := n.Sub(BigFloat.FromFloat64(0))
		gomega.Expect(r.HasError()).To(gomega.BeTrue())
		gomega.Expect(r.Error().Message()).To(gomega.ContainSubstring("Sub"))
	})
	ginkgo.It("Mul returns an error result", func() {
		r := n.Mul(BigFloat.FromFloat64(0))
		gomega.Expect(r.HasError()).To(gomega.BeTrue())
		gomega.Expect(r.Error().Message()).To(gomega.ContainSubstring("Mul"))
	})
	ginkgo.It("Div returns an error result", func() {
		r := n.Div(BigFloat.FromFloat64(0))
		gomega.Expect(r.HasError()).To(gomega.BeTrue())
		gomega.Expect(r.Error().Message()).To(gomega.ContainSubstring("Div"))
	})
	ginkgo.It("Abs returns an error result", func() {
		r := n.Abs()
		gomega.Expect(r.HasError()).To(gomega.BeTrue())
		gomega.Expect(r.Error().Message()).To(gomega.ContainSubstring("Abs"))
	})
	ginkgo.It("Neg returns an error result", func() {
		r := n.Neg()
		gomega.Expect(r.HasError()).To(gomega.BeTrue())
		gomega.Expect(r.Error().Message()).To(gomega.ContainSubstring("Neg"))
	})
	ginkgo.It("Equal is true only against another null", func() {
		gomega.Expect(n.Equal(NullBigFloat.New())).To(gomega.BeTrue())
		gomega.Expect(n.Equal(BigFloat.FromFloat64(0))).To(gomega.BeFalse())
	})
	ginkgo.It("LessThan is always false", func() {
		gomega.Expect(n.LessThan(BigFloat.FromFloat64(0))).To(gomega.BeFalse())
	})
	ginkgo.It("GreaterThan is always false", func() {
		gomega.Expect(n.GreaterThan(BigFloat.FromFloat64(0))).To(gomega.BeFalse())
	})
	ginkgo.It("Inspect renders the null marker", func() {
		gomega.Expect(n.Inspect()).To(gomega.Equal(`<NullBigFloat>`))
	})
})
