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


