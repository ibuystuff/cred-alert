package metrics

import "github.com/pivotal-golang/lager"

//go:generate counterfeiter . Counter

type Counter interface {
	Inc(lager.Logger, ...string)
	IncN(lager.Logger, int, ...string)
}

type counter struct {
	metric Metric
}

func NewCounter(metric Metric) *counter {
	return &counter{
		metric: metric,
	}
}

func (c *counter) Inc(logger lager.Logger, tags ...string) {
	c.IncN(logger, 1, tags...)
}

func (c *counter) IncN(logger lager.Logger, count int, tags ...string) {
	if count <= 0 {
		return
	}
	c.metric.Update(logger, float32(count), tags...)
}

func NewNullCounter(metric Metric) *nullCounter {
	return &nullCounter{
		metric: metric,
	}
}

type nullCounter struct {
	metric Metric
}

func (c *nullCounter) Inc(logger lager.Logger, tags ...string) {
	c.IncN(logger, 1, tags...)
}

func (c *nullCounter) IncN(logger lager.Logger, count int, tags ...string) {
	c.metric.Update(logger, float32(count), tags...)
}
