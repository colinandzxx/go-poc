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
	BaseTarget          *types.BigInt `json: "baseTarget"`
	Deadline            *types.BigInt `json: "deadline"`

	//Timestamp uint64
	//Generator uint64
}

func (self ConsensusData) String() string {
	return fmt.Sprintf("generationSignature: %x, baseTarget: %v, deadline: %v",
		self.GenerationSignature, self.BaseTarget, self.Deadline)
}

func (self *ConsensusData) Wrap(chain consensus.ChainReader, unconsensus consensus.Block) ([]byte, error) {
	header := unconsensus.GetHeader()
	preHeader := chain.GetHeader(header.GetHash(), header.GetHeight() - 1)
	if preHeader == nil {
		return []byte{}, pocError.GetHeaderError{
			Height: header.GetHeight() - 1,
			Hash: header.GetHash(),
			Method: pocError.GetHeaderMethod,
		}
	}
	data := preHeader.GetConsensusData()
	if data == nil {
		return []byte{}, pocError.ErrGetConsensusData
	}
	var preConsensusData ConsensusData
	_, err := preConsensusData.UnWrap(data)
	if err != nil {
		return []byte{}, err
	}

	// GenerationSignature
	generator := binary.LittleEndian.Uint64(header.GetGenerator())
	self.GenerationSignature = calculateGenerationSignature(preConsensusData.GenerationSignature, generator)

	// BaseTarget
	bt := CalculateBaseTarget(chain, unconsensus)
	if bt == nil {
		return []byte{}, pocError.ErrCalculateBaseTarget
	}
	self.BaseTarget.Put(*bt)

	// Deadline
	var plotter simplePlot
	plotter.plotPoC1(generator, header.GetNonce())
	scoopIndex := calculateScoop(self.GenerationSignature, header.GetHeight())
	dl := calculateDeadline(self.GenerationSignature, plotter.getScoop(scoopIndex), self.BaseTarget.ToInt().Uint64())
	self.Deadline.Put(*dl)

	// encode
	var buf bytes.Buffer
	err = msgp.Encode(&buf, self)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

func (self *ConsensusData) UnWrap(ori []byte) (consensus.Data, error) {
	buf := bytes.NewBuffer(ori)
	err := msgp.Decode(buf, self)
	if err != nil {
		return nil, err
	}
	return self, nil
}
