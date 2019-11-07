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
 * @File: consensusdata.go
 * @LastModified: 2019-10-08 17:30:36
 */

package poc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/colinandzxx/go-consensus"
	"github.com/colinandzxx/go-consensus/types"
	pocError "github.com/colinandzxx/go-poc/error"
	"github.com/tinylib/msgp/msgp"
)

//go:generate msgp

type ConsensusData struct {
	GenerationSignature types.Byte32  `json: "generationSignature"`
	Nonce               uint64	      `json: "nonce"`

	BaseTarget          *types.BigInt `json: "baseTarget"`
	Deadline            *types.BigInt `json: "deadline"`

	//Timestamp uint64
	//Generator uint64
}

type WrapConsensusData struct {
	GenerationSignature types.Byte32  `json: "generationSignature"`
	Nonce               uint64	      `json: "nonce"`
}

func GetConsensusDataFromHeader(header consensus.Header) (*ConsensusData, error) {
	if header == nil {
		return nil, pocError.ErrNilHeader
	}

	if header.GetConsensusData() == nil {
		return nil, pocError.ErrGetConsensusData
	}

	consensusData, ok := header.GetConsensusData().(*ConsensusData)
	if !ok {
		return nil, pocError.ErrTypeConver
	}

	return consensusData, nil
}

func (self ConsensusData) String() string {
	return fmt.Sprintf("generationSignature: %x, baseTarget: %v, deadline: %v",
		self.GenerationSignature, self.BaseTarget, self.Deadline)
}

func (self *ConsensusData) Wrap(chain consensus.ChainReader, unconsensus consensus.Header, engine consensus.Engine) ([]byte, error) {
	poc, ok := engine.(Poc)
	if !ok {
		return nil, pocError.ErrInvalidEngine
	}

	header := unconsensus
	preHeader := chain.GetHeader(header.GetParentHash(), header.GetHeight() - 1)
	if preHeader == nil {
		return nil, pocError.GetHeaderError{
			Height: header.GetHeight() - 1,
			Hash: header.GetParentHash(),
			Method: pocError.GetHeaderMethod,
		}
	}

	consensusData, err := GetConsensusDataFromHeader(header)
	if err != nil {
		return nil, err
	}

	preConsensusData, err := GetConsensusDataFromHeader(preHeader)
	if err != nil {
		return nil, err
	}

	generator := binary.LittleEndian.Uint64(header.GetGenerator())

	// GenerationSignature
	self.GenerationSignature = CalculateGenerationSignature(preConsensusData.GenerationSignature, generator)

	// BaseTarget
	bt := CalculateBaseTarget(chain, preHeader, &poc)
	if bt == nil {
		return []byte{}, pocError.ErrCalculateBaseTarget
	}
	self.BaseTarget.Put(*bt)

	// Deadline
	var plotter SimplePlotter
	plotter.PlotPoC2(generator, consensusData.Nonce)
	scoopIndex := CalculateScoop(self.GenerationSignature, header.GetHeight())
	dl := CalculateDeadline(self.GenerationSignature, plotter.GetScoop(scoopIndex), preConsensusData.BaseTarget.ToInt().Uint64())
	self.Deadline.Put(*dl)

	// encode as WrapConsensusData
	wrapData := WrapConsensusData {
		self.GenerationSignature,
		self.Nonce,
	}
	var buf bytes.Buffer
	err = msgp.Encode(&buf, &wrapData)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

func (self *ConsensusData) UnWrap(chain consensus.ChainReader, header consensus.Header, engine consensus.Engine) (consensus.Data, error) {
	poc, ok := engine.(Poc)
	if !ok {
		return nil, pocError.ErrInvalidEngine
	}

	preHeader := chain.GetHeader(header.GetParentHash(), header.GetHeight() - 1)
	if preHeader == nil {
		return nil, pocError.GetHeaderError{
			Height: header.GetHeight() - 1,
			Hash: header.GetParentHash(),
			Method: pocError.GetHeaderMethod,
		}
	}

	preConsensusData, err := GetConsensusDataFromHeader(preHeader)
	if err != nil {
		return nil, err
	}

	// decode as WrapConsensusData
	if header.GetOriConsensusData() == nil {
		return nil, pocError.ErrNilOriData
	}
	var wrapData WrapConsensusData
	buf := bytes.NewBuffer(header.GetOriConsensusData())
	err = msgp.Decode(buf, &wrapData)
	if err != nil {
		return nil, err
	}

	// fix
	self.GenerationSignature = wrapData.GenerationSignature
	self.Nonce = wrapData.Nonce

	// BaseTarget
	bt := CalculateBaseTarget(chain, preHeader, &poc)
	if bt == nil {
		return nil, pocError.ErrCalculateBaseTarget
	}
	self.BaseTarget.Put(*bt)

	// Deadline
	generator := binary.LittleEndian.Uint64(header.GetGenerator())
	var plotter SimplePlotter
	plotter.PlotPoC2(generator, self.Nonce)
	scoopIndex := CalculateScoop(self.GenerationSignature, header.GetHeight())
	dl := CalculateDeadline(self.GenerationSignature, plotter.GetScoop(scoopIndex), preConsensusData.BaseTarget.ToInt().Uint64())
	self.Deadline.Put(*dl)

	return self, nil
}
