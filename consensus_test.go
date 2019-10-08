package poc

import (
	"testing"

	"github.com/colinandzxx/go-consensus"
)

func Test_calculateGenerationSignature(t *testing.T) {
	lastGenSig := consensus.Byte32{}
	var lastGenId uint64 = 0xFFFFFFFFFFFFFFFF
	sig := calculateGenerationSignature(lastGenSig, lastGenId)
	t.Logf("%x", sig)
}

func Test_calculateScoop(t *testing.T) {
	lastGenSig := consensus.Byte32{}
	var lastGenId uint64 = 0xFFFFFFFFFFFFFFFF
	sig := calculateGenerationSignature(lastGenSig, lastGenId)
	scoop := calculateScoop(sig, 1)
	t.Logf("%v", scoop)
}
