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
 * @LastModified: 2019-10-09 14:35:23
 */

package poc

import (
	"bytes"
	"github.com/colinandzxx/go-consensus/types"
	"github.com/tinylib/msgp/msgp"
)

var Cfg Config

var genesisGenerationSignature = types.Byte32{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
}

type Config struct {
	AvgBaseTargetNum uint32
	ConsensusInterval uint32
	MaxBaseTarget uint64
	//Loglevel

	GenesisData types.Bytes
}


func init() {
	Cfg.Default()
}

func (self *Config) Default() {
	self.AvgBaseTargetNum = 24
	self.ConsensusInterval = 240 //s
	self.MaxBaseTarget = 0x444444444

	err := self.SetGenesisData(genesisGenerationSignature, 0)
	if err != nil {
		panic(err)
	}
}

func (self *Config) SetGenesisData(generationSignature types.Byte32, nonce uint64) error {
	genesisData := WrapConsensusData{
		GenerationSignature: genesisGenerationSignature,
		Nonce: 0,
	}
	var buf bytes.Buffer
	err := msgp.Encode(&buf, &genesisData)
	if err != nil {
		return err
	}
	self.GenesisData = buf.Bytes()
	return nil
}
