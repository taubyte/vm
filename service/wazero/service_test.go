package service

import (
	"errors"
	"testing"
	"time"

	"github.com/taubyte/go-interfaces/vm/mocks"
	"gotest.tools/v3/assert"
)

func TestService(t *testing.T) {
	_, service, err := newService()
	assert.NilError(t, err)

	err = service.Close()
	assert.NilError(t, err)
}

func TestModuleFunctionFailure(t *testing.T) {
	module, err := newModuleInstance()
	assert.NilError(t, err)

	_, err = module.Function("hello_world")
	assertError(t, err)
}

func TestInstance(t *testing.T) {
	instance, err := newBasicInstance()
	assert.NilError(t, err)
	if instance.Stderr() == nil {
		t.Error("stderr is nil")
	}

	if instance.Filesystem() == nil {
		t.Error("stderr is nil")
	}

	if instance.Stdout() == nil {
		t.Error("stdout is nil")
	}

	if instance.Context() == nil {
		t.Error("context is nil")
	}

	err = instance.Close()
	assert.NilError(t, err)
}

func TestRuntime(t *testing.T) {
	instance, err := newLoadedInstance()
	assert.NilError(t, err)

	if instance.Stderr() == nil {
		t.Error("stderr is nil")
	}

	if instance.Stdout() == nil {
		t.Error("stdout is nil")
	}

	err = instance.Close()
	assert.NilError(t, err)
}

func TestRuntimeCall(t *testing.T) {
	err := callFuncs([]string{"tou32", "tof32", "toi32", "tof64"})
	assert.NilError(t, err)

	mi, err := newModuleInstance()
	assert.NilError(t, err)

	fi, err := mi.Function("tou32")
	assert.NilError(t, err)

	// Coverage

	ret := fi.Call(float64(42))
	assert.NilError(t, ret.Error())

	ret = fi.Call(float32(42))
	assert.NilError(t, ret.Error())

	ret = fi.Call(int(42))
	assert.NilError(t, ret.Error())

	fi.Timeout(10 * time.Second)

	// Failures

	// Type Error: String is not supported
	ret = fi.Call("string")
	assertError(t, ret.Error())
}

func TestReflectFailures(t *testing.T) {
	functions, err := newFuncs([]string{"tou32", "tof64"})
	assert.NilError(t, err)

	retu32 := functions["tou32"].Call(theAnswer)
	assert.NilError(t, retu32.Error())

	err = retu32.Reflect(&f32RetVal)
	assertError(t, err)

	retu32.Reflect(&f64RetVal)
	assertError(t, err)

	retf64 := functions["tof64"].Call(theAnswer)
	assert.NilError(t, retf64.Error())

	err = retf64.Reflect(&u32RetVal)
	assertError(t, err)

	err = retf64.Reflect(&i32RetVal)
	assertError(t, err)

	err = retf64.Reflect(&f64RetVal, &f64RetVal)
	assertError(t, err)

	err = retf64.Reflect(&controlRetVal)
	assertError(t, err)

	err = retf64.Reflect(f64RetVal)
	assertError(t, err)

	ret := retf64.(*wasmReturn)
	ret.err = errors.New("mock error")

	err = ret.Reflect(&f64RetVal)
	assertError(t, err)
}

func TestRuntimePlugin(t *testing.T) {
	instance, err := newLoadedInstance()
	assert.NilError(t, err)

	plugin := mocks.NewPlugin(false)
	_, _, err = instance.Attach(plugin)
	assert.NilError(t, err)

	// mock New error
	plugin = mocks.NewPlugin(true)
	_, _, err = instance.Attach(plugin)
	assertError(t, err)
}
