/*
 * Copyright (c) 2019
 *
 * This project is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 * @File: algorithm.go
 * @LastModified: 2019-11-07 14:52:22
 */

package poc

import (
	"encoding/binary"
	"fmt"
	"github.com/colinandzxx/go-consensus"
	"github.com/colinandzxx/go-consensus/types"
	pocError "github.com/colinandzxx/go-poc/error"
	"github.com/moonfruit/go-shabal"
	"math/big"
)

// TODO: this should be configure !!!
//const consensusInterval = uint64(4 * 60) //s
//const maxBaseTarget = uint64(0x444444444) // 18325193796

var two64, _ = big.NewInt(0).SetString("18446744073709551616", 10) // 0x10000000000000000

func CalculateGenerationSignature(lastGenSig types.Byte32, lastGenId uint64) types.Byte32 {
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

func CalculateScoop(genSig types.Byte32, height uint64) int32 {
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
	return int32(scoopBig.Int64())
}

func CalculateHit(genSig types.Byte32, scoopData types.Byte64) *big.Int {
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
	hitBig.SetBytes([]byte{hitBytes[7], hitBytes[6], hitBytes[5], hitBytes[4], hitBytes[3], hitBytes[2], hitBytes[1], hitBytes[0]})
	return hitBig
}

// baseTarget from prev header !!
func CalculateDeadline(genSig types.Byte32, scoopData types.Byte64, lastBaseTarget uint64) *big.Int {
	hit := CalculateHit(genSig, scoopData)
	return CalculateDeadlineByHit(hit, lastBaseTarget)
}

func CalculateDeadlineByHit(hit *big.Int, lastBaseTarget uint64) *big.Int {
	if hit == nil {
		return nil
	}
	return big.NewInt(0).Div(hit, big.NewInt(0).SetUint64(lastBaseTarget))
}

func CalculateDifficulty(baseTarget *big.Int) *big.Int {
	return big.NewInt(0).Div(two64, baseTarget)
}

func CalculateAvgBaseTarget(chain consensus.ChainReader, from consensus.Header, offset uint32) (*big.Int, consensus.Header)  {
	avgBaseTarget := big.NewInt(0)
	header := from
	var blockCounter int64 = 0
	for ; offset != 0; offset-- {
		prev, err := GetConsensusDataFromHeader(header)
		if err != nil {
			panic(err)
		}

		avgBaseTarget = avgBaseTarget.Mul(avgBaseTarget, big.NewInt(blockCounter))
		avgBaseTarget = avgBaseTarget.Add(avgBaseTarget, prev.BaseTarget.ToInt())
		avgBaseTarget = avgBaseTarget.Div(avgBaseTarget, big.NewInt(blockCounter + 1))

		// return the last header
		if offset != 1 {
			header = chain.GetHeaderByHash(header.GetParentHash())
			if header == nil {
				panic(pocError.GetHeaderError{
					Hash:   header.GetHash(),
					Method: pocError.GetHeaderByHashMethod,
				})
			}
		}

		blockCounter++
	}

	return avgBaseTarget, header
}

func CalculateBaseTarget(chain consensus.ChainReader, prev consensus.Header, poc *Poc) *big.Int  {
	if prev.GetHeight() <= uint64(poc.AvgBaseTargetNum) {
		return big.NewInt(0).SetUint64(poc.MaxBaseTarget)
	}

	avgBaseTarget, back := CalculateAvgBaseTarget(chain, prev, poc.AvgBaseTargetNum)
	front := prev
	if front.GetTimestamp() < back.GetTimestamp() {
		panic("Timestamp is sick")
	}
	difTime := front.GetTimestamp() - back.GetTimestamp()
	targetTimespan := uint64(poc.AvgBaseTargetNum * poc.ConsensusInterval)

	if difTime < targetTimespan / 2 {
		difTime = targetTimespan / 2
	} else if difTime > targetTimespan * 2 {
		difTime = targetTimespan * 2
	}

	data, err := GetConsensusDataFromHeader(prev)
	if err != nil {
		panic(err)
	}

	lastBaseTarget := data.BaseTarget.ToInt()
	newBaseTarget := avgBaseTarget.Mul(avgBaseTarget, big.NewInt(0).SetUint64(difTime))
	newBaseTarget = newBaseTarget.Div(newBaseTarget, big.NewInt(0).SetUint64(targetTimespan))

	if newBaseTarget.Cmp(big.NewInt(0).SetUint64(poc.MaxBaseTarget)) > 0 {
		newBaseTarget = big.NewInt(0).SetUint64(poc.MaxBaseTarget)
	}

	if newBaseTarget.Cmp(big.NewInt(0)) == 0 {
		newBaseTarget.SetUint64(1)
	}

	{
		tmpBaseTarget := big.NewInt(0).Mul(lastBaseTarget, big.NewInt(8))
		tmpBaseTarget = tmpBaseTarget.Div(tmpBaseTarget, big.NewInt(10))
		if newBaseTarget.Cmp(tmpBaseTarget) < 0 {
			newBaseTarget = tmpBaseTarget
		}
	}

	{
		tmpBaseTarget := big.NewInt(0).Mul(lastBaseTarget, big.NewInt(12))
		tmpBaseTarget = tmpBaseTarget.Div(tmpBaseTarget, big.NewInt(10))
		if newBaseTarget.Cmp(tmpBaseTarget) > 0 {
			newBaseTarget = tmpBaseTarget
		}
	}

	return newBaseTarget
}
