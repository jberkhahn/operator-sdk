// Code generated by counterfeiter. DO NOT EDIT.
package generatefakes

import (
	"sync"

	"github.com/operator-framework/operator-sdk/internal/generate"
)

type FakePackageManifestGenerator struct {
	GenerateStub        func(*generate.PkgOptions) error
	generateMutex       sync.RWMutex
	generateArgsForCall []struct {
		arg1 *generate.PkgOptions
	}
	generateReturns struct {
		result1 error
	}
	generateReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakePackageManifestGenerator) Generate(arg1 *generate.PkgOptions) error {
	fake.generateMutex.Lock()
	ret, specificReturn := fake.generateReturnsOnCall[len(fake.generateArgsForCall)]
	fake.generateArgsForCall = append(fake.generateArgsForCall, struct {
		arg1 *generate.PkgOptions
	}{arg1})
	fake.recordInvocation("Generate", []interface{}{arg1})
	fake.generateMutex.Unlock()
	if fake.GenerateStub != nil {
		return fake.GenerateStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.generateReturns
	return fakeReturns.result1
}

func (fake *FakePackageManifestGenerator) GenerateCallCount() int {
	fake.generateMutex.RLock()
	defer fake.generateMutex.RUnlock()
	return len(fake.generateArgsForCall)
}

func (fake *FakePackageManifestGenerator) GenerateCalls(stub func(*generate.PkgOptions) error) {
	fake.generateMutex.Lock()
	defer fake.generateMutex.Unlock()
	fake.GenerateStub = stub
}

func (fake *FakePackageManifestGenerator) GenerateArgsForCall(i int) *generate.PkgOptions {
	fake.generateMutex.RLock()
	defer fake.generateMutex.RUnlock()
	argsForCall := fake.generateArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakePackageManifestGenerator) GenerateReturns(result1 error) {
	fake.generateMutex.Lock()
	defer fake.generateMutex.Unlock()
	fake.GenerateStub = nil
	fake.generateReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakePackageManifestGenerator) GenerateReturnsOnCall(i int, result1 error) {
	fake.generateMutex.Lock()
	defer fake.generateMutex.Unlock()
	fake.GenerateStub = nil
	if fake.generateReturnsOnCall == nil {
		fake.generateReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.generateReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakePackageManifestGenerator) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.generateMutex.RLock()
	defer fake.generateMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakePackageManifestGenerator) recordInvocation(key string, args []interface{}) {
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

var _ generate.PackageManifestGenerator = new(FakePackageManifestGenerator)
