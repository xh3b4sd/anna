package isgreater

import (
	"reflect"
	"testing"

	"golang.org/x/net/context"

	"github.com/xh3b4sd/anna/spec"
)

func Test_CLG_New_Error(t *testing.T) {
	newConfig := DefaultConfig()
	newConfig.Log = nil
	_, err := New(newConfig)
	if !IsInvalidConfig(err) {
		t.Fatal("expected", true, "got", false)
	}
}

func Test_CLG_GetName(t *testing.T) {
	newCLG := MustNew()
	clgName := newCLG.GetName()
	if clgName != "is-greater" {
		t.Fatal("expected", "is-greater", "got", clgName)
	}
}

func Test_CLG_GetInputChannel(t *testing.T) {
	newCLG := MustNew()
	inputChannel := newCLG.GetInputChannel()
	if inputChannel == nil {
		t.Fatal("expected", make(chan spec.NetworkPayload, 1000), "got", nil)
	}
}

func Test_CLG_GetInputTypes(t *testing.T) {
	newCLG := MustNew()
	inputTypes := newCLG.GetInputTypes()
	if len(inputTypes) < 1 {
		t.Fatal("expected", 1, "got", 0)
	}

	var ctx context.Context
	if !inputTypes[0].Implements(reflect.TypeOf(&ctx).Elem()) {
		t.Fatal("expected", true, "got", false)
	}
}

func Test_CLG_SetLog(t *testing.T) {
	newCLG := MustNew()
	var rawCLG *clg

	switch c := newCLG.(type) {
	case *clg:
		// all good
		rawCLG = newCLG.(*clg)
	default:
		t.Fatal("expected", "*clg", "got", c)
	}

	if rawCLG.Log == nil {
		t.Fatal("expected", "spec.Log", "got", nil)
	}

	newCLG.SetLog(nil)

	if rawCLG.Log != nil {
		t.Fatal("expected", nil, "got", "spec.Log")
	}
}

func Test_CLG_SetStorage(t *testing.T) {
	newCLG := MustNew()
	var rawCLG *clg

	switch c := newCLG.(type) {
	case *clg:
		// all good
		rawCLG = newCLG.(*clg)
	default:
		t.Fatal("expected", "*clg", "got", c)
	}

	if rawCLG.Storage == nil {
		t.Fatal("expected", "spec.Storage", "got", nil)
	}

	newCLG.SetStorage(nil)

	if rawCLG.Storage != nil {
		t.Fatal("expected", nil, "got", "spec.Storage")
	}
}
