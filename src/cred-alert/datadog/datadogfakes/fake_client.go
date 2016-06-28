// This file was generated by counterfeiter
package datadogfakes

import (
	"cred-alert/datadog"
	"sync"
)

type FakeClient struct {
	PublishSeriesStub        func(series datadog.Series) error
	publishSeriesMutex       sync.RWMutex
	publishSeriesArgsForCall []struct {
		series datadog.Series
	}
	publishSeriesReturns struct {
		result1 error
	}
	BuildCountMetricStub        func(metricName string, count float32, tags ...string) datadog.Metric
	buildCountMetricMutex       sync.RWMutex
	buildCountMetricArgsForCall []struct {
		metricName string
		count      float32
		tags       []string
	}
	buildCountMetricReturns struct {
		result1 datadog.Metric
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeClient) PublishSeries(series datadog.Series) error {
	fake.publishSeriesMutex.Lock()
	fake.publishSeriesArgsForCall = append(fake.publishSeriesArgsForCall, struct {
		series datadog.Series
	}{series})
	fake.recordInvocation("PublishSeries", []interface{}{series})
	fake.publishSeriesMutex.Unlock()
	if fake.PublishSeriesStub != nil {
		return fake.PublishSeriesStub(series)
	} else {
		return fake.publishSeriesReturns.result1
	}
}

func (fake *FakeClient) PublishSeriesCallCount() int {
	fake.publishSeriesMutex.RLock()
	defer fake.publishSeriesMutex.RUnlock()
	return len(fake.publishSeriesArgsForCall)
}

func (fake *FakeClient) PublishSeriesArgsForCall(i int) datadog.Series {
	fake.publishSeriesMutex.RLock()
	defer fake.publishSeriesMutex.RUnlock()
	return fake.publishSeriesArgsForCall[i].series
}

func (fake *FakeClient) PublishSeriesReturns(result1 error) {
	fake.PublishSeriesStub = nil
	fake.publishSeriesReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) BuildCountMetric(metricName string, count float32, tags ...string) datadog.Metric {
	fake.buildCountMetricMutex.Lock()
	fake.buildCountMetricArgsForCall = append(fake.buildCountMetricArgsForCall, struct {
		metricName string
		count      float32
		tags       []string
	}{metricName, count, tags})
	fake.recordInvocation("BuildCountMetric", []interface{}{metricName, count, tags})
	fake.buildCountMetricMutex.Unlock()
	if fake.BuildCountMetricStub != nil {
		return fake.BuildCountMetricStub(metricName, count, tags...)
	} else {
		return fake.buildCountMetricReturns.result1
	}
}

func (fake *FakeClient) BuildCountMetricCallCount() int {
	fake.buildCountMetricMutex.RLock()
	defer fake.buildCountMetricMutex.RUnlock()
	return len(fake.buildCountMetricArgsForCall)
}

func (fake *FakeClient) BuildCountMetricArgsForCall(i int) (string, float32, []string) {
	fake.buildCountMetricMutex.RLock()
	defer fake.buildCountMetricMutex.RUnlock()
	return fake.buildCountMetricArgsForCall[i].metricName, fake.buildCountMetricArgsForCall[i].count, fake.buildCountMetricArgsForCall[i].tags
}

func (fake *FakeClient) BuildCountMetricReturns(result1 datadog.Metric) {
	fake.BuildCountMetricStub = nil
	fake.buildCountMetricReturns = struct {
		result1 datadog.Metric
	}{result1}
}

func (fake *FakeClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.publishSeriesMutex.RLock()
	defer fake.publishSeriesMutex.RUnlock()
	fake.buildCountMetricMutex.RLock()
	defer fake.buildCountMetricMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeClient) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ datadog.Client = new(FakeClient)
