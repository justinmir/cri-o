package metrics_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	dto "github.com/prometheus/client_model/go"

	"github.com/cri-o/cri-o/internal/storage/references"
	"github.com/cri-o/cri-o/server/metrics"
	. "github.com/cri-o/cri-o/test/framework"
)

// TestMetrics runs the created specs.
func TestMetrics(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Metrics")
}

//nolint:gochecknoglobals // test framework requires global state
var t *TestFramework

var _ = BeforeSuite(func() {
	t = NewTestFramework(NilFunc, NilFunc)
	t.Setup()
})

var _ = AfterSuite(func() {
	t.Teardown()
})

// The actual test suite.
var _ = t.Describe("Metrics", func() {
	t.Describe("SinceInMicroseconds", func() {
		It("should succeed", func() {
			// Given
			// When
			res := metrics.SinceInMicroseconds(
				time.Now().Add(-time.Millisecond))

			// Then
			Expect(res).NotTo(BeZero())
		})

		It("should be zero at time.Now()", func() {
			// Given
			// When
			res := metrics.SinceInMicroseconds(time.Now())

			// Then
			Expect(res).To(BeZero())
		})
	})

	t.Describe("MetricImagePullsSuccessesInc", func() {
		It("should increment the by-registry counter with a resolved registry", func() {
			m := metrics.Instance()
			ref, err := references.ParseRegistryImageReferenceFromOutOfProcessData("docker.io/library/busybox:latest")
			Expect(err).NotTo(HaveOccurred())

			m.MetricImagePullsSuccessesInc(ref, "harbor.internal:5000")

			var metric dto.Metric
			c, err := m.ImagePullsSuccessByRegistryVec().GetMetricWithLabelValues("harbor.internal:5000")
			Expect(err).NotTo(HaveOccurred())
			Expect(c.Write(&metric)).To(Succeed())
			Expect(metric.GetCounter().GetValue()).To(Equal(1.0))
		})

		It("should not panic when registry is empty", func() {
			m := metrics.Instance()
			ref, err := references.ParseRegistryImageReferenceFromOutOfProcessData("docker.io/library/busybox:latest")
			Expect(err).NotTo(HaveOccurred())

			Expect(func() { m.MetricImagePullsSuccessesInc(ref, "") }).NotTo(Panic())
		})
	})

	t.Describe("MetricImagePullsFailuresInc", func() {
		It("should increment the by-registry counter with a resolved registry", func() {
			m := metrics.Instance()
			ref, err := references.ParseRegistryImageReferenceFromOutOfProcessData("docker.io/library/busybox:latest")
			Expect(err).NotTo(HaveOccurred())

			m.MetricImagePullsFailuresInc(ref, "UNKNOWN", "harbor.internal:5000")

			var metric dto.Metric
			rc, err := m.ImagePullsFailureByRegistryVec().GetMetricWithLabelValues("harbor.internal:5000", "UNKNOWN")
			Expect(err).NotTo(HaveOccurred())
			Expect(rc.Write(&metric)).To(Succeed())
			Expect(metric.GetCounter().GetValue()).To(Equal(1.0))
		})

		It("should not increment the by-registry counter when registry is empty", func() {
			m := metrics.Instance()
			ref, err := references.ParseRegistryImageReferenceFromOutOfProcessData("docker.io/library/busybox:latest")
			Expect(err).NotTo(HaveOccurred())

			m.MetricImagePullsFailuresInc(ref, "CONNECTION_REFUSED", "")
			// Should not panic
		})
	})
})
