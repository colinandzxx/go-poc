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
 * @File: config.go
 * @LastModified: 2019-11-07 15:00:48
 */

package poc

import (
	"bytes"
	"github.com/colinandzxx/go-consensus/types"
	"github.com/tinylib/msgp/msgp"
)

var genesisGenerationSignature = types.Byte32{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
}

type Config struct {
	AvgBaseTargetNum uint32
	ConsensusInterval uint32
	MaxBaseTarget uint64
	//Loglevel

	Gen Genenis
}

type Genenis struct {
	GenerationSignature types.Byte32
	Nonce uint64
}

func (self *Config) Default() {
	self.AvgBaseTargetNum = 24
	self.ConsensusInterval = 240 //s
	self.MaxBaseTarget = 0x444444444

	self.Gen.GenerationSignature = genesisGenerationSignature
	self.Gen.Nonce = 0
}

func (self *Config) GetGenesisConsensusData() ([]byte, error) {
	genesisData := WrapConsensusData{
		GenerationSignature: self.Gen.GenerationSignature,
		Nonce:               self.Gen.Nonce,
	}
	var buf bytes.Buffer
	err := msgp.Encode(&buf, &genesisData)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
