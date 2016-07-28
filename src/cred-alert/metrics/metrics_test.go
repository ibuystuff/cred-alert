package metrics_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.cloudfoundry.org/lager/lagertest"

	"cred-alert/datadog"
	"cred-alert/datadog/datadogfakes"
	"cred-alert/metrics"
)

var _ = Describe("Metrics", func() {
	var (
		logger *lagertest.TestLogger

		client             *datadogfakes.FakeClient
		metric             metrics.Metric
		expectedMetricType string
		expectedMetricName string
		expectedEnv        string
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("metrics")

		client = &datadogfakes.FakeClient{}
		expectedEnv = "env"
		emitter := metrics.NewEmitter(client, expectedEnv)
		expectedMetricType = "name"
		expectedMetricName = "type"
		metric = metrics.NewMetric(expectedMetricName, expectedMetricType, emitter)
	})

	It("calls BuildMetricCallCount and PublishSeries", func() {
		expectedValue := float32(0)
		metric.Update(logger, expectedValue)

		expectedMetric := datadog.Metric{}
		client.BuildMetricReturns(expectedMetric)

		Expect(client.BuildMetricCallCount()).Should(Equal(1))
		metricType, name, value, env := client.BuildMetricArgsForCall(0)
		Expect(metricType).To(Equal(expectedMetricType))
		Expect(name).To(Equal(expectedMetricName))
		Expect(value).To(Equal(expectedValue))
		Expect(env[0]).To(Equal(expectedEnv))
		Expect(client.PublishSeriesCallCount()).Should(Equal(1))
		Expect(client.PublishSeriesArgsForCall(0)).To(ContainElement(expectedMetric))
	})

	It("can have tags", func() {
		metric.Update(logger, 3.52, "name:value", "othername:othervalue")

		Expect(client.BuildMetricCallCount()).Should(Equal(1))
		_, _, _, tags := client.BuildMetricArgsForCall(0)
		Expect(tags).To(ConsistOf(expectedEnv, "name:value", "othername:othervalue"))
	})

	It("logs", func() {
		expectedValue := 1
		metric.Update(logger, float32(expectedValue))

		Expect(len(logger.LogMessages())).To(Equal(2))
		Expect(logger.LogMessages()[0]).To(ContainSubstring("metrics.update.starting"))
		Expect(logger.Logs()[0].Data["name"]).To(Equal(expectedMetricName))
		Expect(logger.Logs()[0].Data["type"]).To(Equal(expectedMetricType))
		Expect(logger.Logs()[0].Data["environment"]).To(Equal(expectedEnv))
		Expect(logger.Logs()[0].Data["value"]).To(Equal(float64(expectedValue)))
	})
})
