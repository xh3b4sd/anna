package lesser

// This file is generated by the CLG generator. Don't edit it manually. The CLG
// generator is invoked by go generate. For more information about the usage of
// the CLG generator check https://github.com/xh3b4sd/clggen or have a look at
// the clg package. There is the go generate statement placed to invoke clggen.

import (
	"testing"
)

func Test_CLG_GetMetadata(t *testing.T) {
	newCLG := MustNew()

	if len(newCLG.GetMetadata()) == 0 {
		t.Fatal("expected", "metadata", "got", "nothing")
	}
}
