package pkg

import (
	"container/list"
	"math/big"
)

// TODO: this should be configure !!!
const consensusInterval = uint64(4 * 60) //s
const maxBaseTarget = uint64(18325193796)

var two64, _ = big.NewInt(0).SetString("18446744073709551616", 10)

type ConsensusData struct {
	BaseTarget *big.Int	`json: baseTarget`
	TimeStamp  uint64
}

type Consensus struct {

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
		avgBaseTarget = avgBaseTarget.Add(avgBaseTarget, prev.BaseTarget)
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
	if front.TimeStamp < back.TimeStamp {
		panic("TimeStamp is sick")
	}
	difTime := front.TimeStamp - back.TimeStamp
	targetTimespan := uint64(listConsensusData.Len()) * consensusInterval

	if difTime < targetTimespan / 2 {
		difTime = targetTimespan / 2
	} else if difTime > targetTimespan * 2 {
		difTime = targetTimespan * 2
	}

	lastBaseTarget := front.BaseTarget
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


