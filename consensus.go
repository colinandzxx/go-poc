package poc

import (
	"container/list"
	"encoding/binary"
	"fmt"
	"github.com/colinandzxx/go-consensus/types"
	"github.com/moonfruit/go-shabal"
	"math/big"
)

// TODO: this should be configure !!!
const consensusInterval = uint64(4 * 60) //s
const maxBaseTarget = uint64(0x444444444) // 18325193796

var two64, _ = big.NewInt(0).SetString("18446744073709551616", 10) // 0x10000000000000000

type Consensus struct {

}

func calculateGenerationSignature(lastGenSig types.Byte32, lastGenId uint64) types.Byte32 {
	data := make([]byte, 40)
	copy(data, lastGenSig[:])
	// use BigEndian in burst code !
	binary.BigEndian.PutUint64(data[32:], lastGenId)
	s256 := shabal.NewShabal256()
	_, err := s256.Write(data)
	if err != nil {
		panic(fmt.Errorf("%v", err))
	}
	var ret types.Byte32
	copy(ret[:], s256.Sum(nil)[:32])
	return ret
}

func calculateScoop(genSig types.Byte32, height uint64) uint64 {
	data := make([]byte, 40)
	copy(data, genSig[:])
	// use BigEndian in burst code !
	binary.BigEndian.PutUint64(data[32:], height)
	s256 := shabal.NewShabal256()
	_, err := s256.Write(data)
	if err != nil {
		panic(fmt.Errorf("%v", err))
	}

	scoopBig := big.Int{}
	scoopBig.SetBytes(s256.Sum(nil))
	scoopBig.Mod(&scoopBig, big.NewInt(int64(scoopsPerPlot)))
	return scoopBig.Uint64()
}

func calculateHit(genSig types.Byte32, scoopData types.Byte64) *big.Int {
	s256 := shabal.NewShabal256()
	_, err := s256.Write(genSig[:])
	if err != nil {
		panic(fmt.Errorf("%v", err))
	}
	_, err = s256.Write(scoopData[:])
	if err != nil {
		panic(fmt.Errorf("%v", err))
	}

	hitBig := big.NewInt(0)
	hitBytes := s256.Sum(nil)
	hitBig.SetBytes([]byte{hitBytes[7], hitBytes[6], hitBytes[5], hitBytes[5], hitBytes[3], hitBytes[2], hitBytes[1], hitBytes[0]})
	return hitBig
}

func calculateDeadline(genSig types.Byte32, scoopData types.Byte64, baseTarget uint64) *big.Int {
	hit := calculateHit(genSig, scoopData)
	return hit.Div(hit, big.NewInt(0).SetUint64(baseTarget))
}

func CalculateDifficulty(baseTarget *big.Int) *big.Int {
	return two64.Div(two64, baseTarget)
}

func calculateAvgBaseTarget(listConsensusData list.List) *big.Int  {
	avgBaseTarget := big.NewInt(0)
	var blockCounter int64 = 0
	for e := listConsensusData.Front(); e != nil; e = e.Prev() {
		blockCounter++
		prev, ok := e.Value.(*ConsensusData)
		if !ok {
			panic("listConsensusData can not convert to *ConsensusData")
		}
		avgBaseTarget = avgBaseTarget.Mul(avgBaseTarget, big.NewInt(blockCounter))
		avgBaseTarget = avgBaseTarget.Add(avgBaseTarget, prev.BaseTarget.ToInt())
		avgBaseTarget = avgBaseTarget.Div(avgBaseTarget, big.NewInt(blockCounter + 1))
	}
	return avgBaseTarget
}

func CalculateBaseTarget(listConsensusData list.List) *big.Int  {
	avgBaseTarget := calculateAvgBaseTarget(listConsensusData)
	front, ok := listConsensusData.Front().Value.(*ConsensusData)
	if !ok {
		panic("listConsensusData can not convert to *ConsensusData")
	}
	back, ok := listConsensusData.Back().Value.(*ConsensusData)
	if !ok {
		panic("listConsensusData can not convert to *ConsensusData")
	}
	if front.Timestamp < back.Timestamp {
		panic("Timestamp is sick")
	}
	difTime := front.Timestamp - back.Timestamp
	targetTimespan := uint64(listConsensusData.Len()) * consensusInterval

	if difTime < targetTimespan / 2 {
		difTime = targetTimespan / 2
	} else if difTime > targetTimespan * 2 {
		difTime = targetTimespan * 2
	}

	lastBaseTarget := front.BaseTarget.ToInt()
	newBaseTarget := avgBaseTarget.Mul(avgBaseTarget, big.NewInt(0).SetUint64(difTime))
	newBaseTarget = newBaseTarget.Div(newBaseTarget, big.NewInt(0).SetUint64(targetTimespan))

	if newBaseTarget.Cmp(big.NewInt(0).SetUint64(maxBaseTarget)) > 0 {
		newBaseTarget = big.NewInt(0).SetUint64(maxBaseTarget)
	}

	if newBaseTarget.Cmp(big.NewInt(0)) == 0 {
		newBaseTarget.SetUint64(1)
	}

	{
		tmpBaseTarget := lastBaseTarget.Mul(lastBaseTarget, big.NewInt(8)).Div(lastBaseTarget, big.NewInt(10))
		if newBaseTarget.Cmp(tmpBaseTarget) < 0 {
			newBaseTarget = tmpBaseTarget
		}
	}

	{
		tmpBaseTarget := lastBaseTarget.Mul(lastBaseTarget, big.NewInt(12)).Div(lastBaseTarget, big.NewInt(10))
		if newBaseTarget.Cmp(tmpBaseTarget) > 0 {
			newBaseTarget = tmpBaseTarget
		}
	}

	return newBaseTarget
}
