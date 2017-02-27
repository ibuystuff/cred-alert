// This file was generated by counterfeiter
package revokfakes

import (
	"cred-alert/revok"
	"sync"

	"code.cloudfoundry.org/lager"
)

type FakeNotificationComposer struct {
	ScanAndNotifyStub        func(lager.Logger, string, string, map[string]struct{}, string, string, string) error
	scanAndNotifyMutex       sync.RWMutex
	scanAndNotifyArgsForCall []struct {
		arg1 lager.Logger
		arg2 string
		arg3 string
		arg4 map[string]struct{}
		arg5 string
		arg6 string
		arg7 string
	}
	scanAndNotifyReturns struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeNotificationComposer) ScanAndNotify(arg1 lager.Logger, arg2 string, arg3 string, arg4 map[string]struct{}, arg5 string, arg6 string, arg7 string) error {
	fake.scanAndNotifyMutex.Lock()
	fake.scanAndNotifyArgsForCall = append(fake.scanAndNotifyArgsForCall, struct {
		arg1 lager.Logger
		arg2 string
		arg3 string
		arg4 map[string]struct{}
		arg5 string
		arg6 string
		arg7 string
	}{arg1, arg2, arg3, arg4, arg5, arg6, arg7})
	fake.recordInvocation("ScanAndNotify", []interface{}{arg1, arg2, arg3, arg4, arg5, arg6, arg7})
	fake.scanAndNotifyMutex.Unlock()
	if fake.ScanAndNotifyStub != nil {
		return fake.ScanAndNotifyStub(arg1, arg2, arg3, arg4, arg5, arg6, arg7)
	}
	return fake.scanAndNotifyReturns.result1
}

func (fake *FakeNotificationComposer) ScanAndNotifyCallCount() int {
	fake.scanAndNotifyMutex.RLock()
	defer fake.scanAndNotifyMutex.RUnlock()
	return len(fake.scanAndNotifyArgsForCall)
}

func (fake *FakeNotificationComposer) ScanAndNotifyArgsForCall(i int) (lager.Logger, string, string, map[string]struct{}, string, string, string) {
	fake.scanAndNotifyMutex.RLock()
	defer fake.scanAndNotifyMutex.RUnlock()
	return fake.scanAndNotifyArgsForCall[i].arg1, fake.scanAndNotifyArgsForCall[i].arg2, fake.scanAndNotifyArgsForCall[i].arg3, fake.scanAndNotifyArgsForCall[i].arg4, fake.scanAndNotifyArgsForCall[i].arg5, fake.scanAndNotifyArgsForCall[i].arg6, fake.scanAndNotifyArgsForCall[i].arg7
}

func (fake *FakeNotificationComposer) ScanAndNotifyReturns(result1 error) {
	fake.ScanAndNotifyStub = nil
	fake.scanAndNotifyReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeNotificationComposer) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.scanAndNotifyMutex.RLock()
	defer fake.scanAndNotifyMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeNotificationComposer) recordInvocation(key string, args []interface{}) {
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

var _ revok.NotificationComposer = new(FakeNotificationComposer)