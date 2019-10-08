package poc

import (
	"fmt"
	"github.com/colinandzxx/go-consensus/types"
)

//go:generate msgp

type ConsensusData struct {
	GenerationSignature byte     `json: "generationSignature"`
	BaseTarget          *types.BigInt `json: "baseTarget"`
	Deadline            *types.BigInt `json: "deadline"`

	Timestamp uint64
}

func (self ConsensusData) String() string {
	return fmt.Sprintf("generationSignature: %x, baseTarget: %v, deadline: %v, timestamp: %v",
		self.GenerationSignature, self.BaseTarget, self.Deadline, self.Timestamp)
}


