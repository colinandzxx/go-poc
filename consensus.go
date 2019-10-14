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
 * @File: consensus.go
 * @LastModified: 2019-10-10 14:43:20
 */

package poc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/colinandzxx/go-consensus"
	pocError "github.com/colinandzxx/go-poc/error"
)

// Currently only support poc2 !!!
type Poc struct {

}

//func (self Poc) CalculateDifficulty(chain consensus.ChainReader, header consensus.Header) *big.Int {
//	var consensusData ConsensusData
//	_, err := consensusData.UnWrap(header.GetConsensusData())
//	if err != nil {
//		return nil
//	}
//	return calculateDifficulty(consensusData.BaseTarget.ToInt())
//}

func (self Poc) VerifyHeader(chain consensus.ChainReader, header consensus.Header) error {
	// these fall under the vaildator of chain, not in consensus, so comment out
	//if chain.GetHeader(header.GetHash(), header.GetHeight()) != nil {
	//	// the header is known, return nil
	//	return nil
	//}
	//if chain.GetHeader(header.GetParentHash(), header.GetHeight() - 1) == nil {
	//	return pocError.GetHeaderError{
	//		Height: header.GetHeight() - 1,
	//		Hash: header.GetParentHash(),
	//		Method: pocError.GetHeaderMethod,
	//	}
	//}

	err := self.VerifyHeaderWithoutForge(chain, header)
	if err != nil {
		return err
	}

	//TODO: verify forge info

	return nil
}

func (self Poc) VerifyHeaderWithoutForge(chain consensus.ChainReader, header consensus.Header) error {
	var consensusData ConsensusData
	_, err := consensusData.UnWrap(chain, header)
	if err != nil {
		return pocError.ErrGetConsensusData
	}

	// Verify the block's difficulty
	if header.GetDifficulty() == nil {
		return pocError.ErrGetDifficulty
	}
	expected := CalculateDifficulty(consensusData.BaseTarget.ToInt())
	if expected.Cmp(header.GetDifficulty()) != 0 {
		return fmt.Errorf("invalid difficulty: have %v, want %v", header.GetDifficulty(), expected)
	}

	return nil
}

func (self Poc) VerifyForge(chain consensus.ChainReader, header consensus.Header) error {
	// Ensure that we have a valid difficulty for the block
	if header.GetDifficulty() == nil {
		return pocError.ErrGetDifficulty
	}
	if header.GetDifficulty().Sign() <= 0 {
		return pocError.ErrInvalidDifficulty
	}

	err := self.verifyGenerationSignature(chain, header)
	if err != nil {
		return err
	}

	err = self.verifyDeadline(chain, header)
	if err != nil {
		return err
	}

	return nil
}

func (self Poc) verifyGenerationSignature(chain consensus.ChainReader, header consensus.Header) error {
	preHeader := chain.GetHeader(header.GetHash(), header.GetHeight() - 1)
	if preHeader == nil {
		return pocError.GetHeaderError{
			Height: header.GetHeight() - 1,
			Hash: header.GetHash(),
			Method: pocError.GetHeaderMethod,
		}
	}

	var consensusData ConsensusData
	_, err := consensusData.UnWrap(chain, header)
	if err != nil {
		return err
	}

	preConsensusData, err := GetConsensusDataFromHeader(preHeader)
	if err != nil {
		return err
	}

	// GenerationSignature
	generator := binary.LittleEndian.Uint64(header.GetGenerator())
	generationSignature := calculateGenerationSignature(preConsensusData.GenerationSignature, generator)
	if bytes.Compare(consensusData.GenerationSignature[:], generationSignature[:]) != 0 {
		return fmt.Errorf("invalid generationSignature: have %x, want %x",
			consensusData.GenerationSignature, generationSignature)
	}

	return nil
}

// don't verify basetarget, local calculate
//func (self Poc) verifyBaseTarget(chain consensus.ChainReader, header consensus.Header) error {
//	return nil
//}

func (self Poc) verifyDeadline(chain consensus.ChainReader, header consensus.Header) error {
	preHeader := chain.GetHeader(header.GetHash(), header.GetHeight() - 1)
	if preHeader == nil {
		return pocError.GetHeaderError{
			Height: header.GetHeight() - 1,
			Hash: header.GetHash(),
			Method: pocError.GetHeaderMethod,
		}
	}

	var consensusData ConsensusData
	_, err := consensusData.UnWrap(chain, header)
	if err != nil {
		return err
	}

	preConsensusData, err := GetConsensusDataFromHeader(preHeader)
	if err != nil {
		return err
	}

	var plot simplePlot
	generator := binary.LittleEndian.Uint64(header.GetGenerator())
	plot.plotPoC2(generator, consensusData.Nonce)
	scoopNum := calculateScoop(consensusData.GenerationSignature, header.GetHeight())
	scoopData := plot.getScoop(scoopNum)
	deadline := calculateDeadline(consensusData.GenerationSignature, scoopData, preConsensusData.BaseTarget.ToInt().Uint64())

	if header.GetTimestamp() < preHeader.GetTimestamp() {
		return pocError.ErrSickTimestamp
	}
	elapsedTime := header.GetTimestamp() - preHeader.GetTimestamp()
	if elapsedTime <= deadline.Uint64() {
		return fmt.Errorf("deadline does not match the block timestamp: %v, %v, %v", header.GetHeight(), elapsedTime, deadline)
	}

	return nil
}

func (self Poc) Forge(chain consensus.ChainReader, header consensus.Header) (consensus.Data, error) {
	return nil, nil
}

