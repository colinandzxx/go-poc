package poc

import (
	"github.com/colinandzxx/go-consensus/types"
	"testing"
)

func Test_calculateGenerationSignature(t *testing.T) {
	lastGenSig := types.Byte32{}
	var lastGenId uint64 = 0xFFFFFFFFFFFFFFFF
	sig := calculateGenerationSignature(lastGenSig, lastGenId)
	t.Logf("%x", sig)
}

func Test_calculateScoop(t *testing.T) {
	lastGenSig := types.Byte32{}
	var lastGenId uint64 = 0xFFFFFFFFFFFFFFFF
	sig := calculateGenerationSignature(lastGenSig, lastGenId)
	scoop := calculateScoop(sig, 1)
	t.Logf("%v", scoop)
}
