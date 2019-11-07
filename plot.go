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
 * @File: plot.go
 * @LastModified: 2019-10-08 17:57:55
 */

package poc

import (
	"encoding/binary"
	"github.com/colinandzxx/go-consensus/types"
	"github.com/moonfruit/go-shabal"
)

const hashSize = 32
const hashesPerScoop = 2
const scoopSize = hashesPerScoop * hashSize
const scoopsPerPlot = 4096 // original 1MB/plot = 16384
const plotSize = scoopsPerPlot * scoopSize
const hashCap = 4096

const baseLen = 16

type SimplePlotter struct {
	data [plotSize]byte
}

func (self *SimplePlotter) PlotPoC1(addr uint64, nonce uint64) {
	var base [baseLen]byte
	// use BigEndian in burst code !
	binary.BigEndian.PutUint64(base[:], addr)
	binary.BigEndian.PutUint64(base[8:], nonce)
	genData := make([]byte, plotSize)
	genData = append(genData, base[:]...)

	s256 := shabal.NewShabal256()
	for i := plotSize; i > 0; i -= hashSize {
		s256.Reset()

		len := plotSize + baseLen - i
		if len > hashCap {
			len = hashCap
		}
		s256.Write(genData[i:(i + len)])
		copy(genData[i - hashSize:], s256.Sum(nil))
	}

	s256.Reset()
	s256.Write(genData)
	finalHash := s256.Sum(nil)
	for i := 0; i < plotSize; i++ {
		self.data[i] = genData[i] ^ finalHash[i % hashSize]
	}
}

func (self *SimplePlotter) PlotPoC2(addr uint64, nonce uint64) {
	self.PlotPoC1(addr, nonce)

	//PoC2 Rearrangement
	var hashBuffer [hashSize]byte
	revPos := plotSize - hashSize //Start at second hash in last scoop
	for pos := hashSize; pos < plotSize / 2; pos += scoopSize { //Start at second hash in first scoop
		copy(hashBuffer[:], self.data[pos:(pos + hashSize)]) 		//Copy low scoop second hash to buffer
		copy(self.data[pos:], self.data[revPos:(revPos + hashSize)]) 	//Copy high scoop second hash to low scoop second hash
		copy(self.data[revPos:], hashBuffer[:hashSize]) 	//Copy buffer to high scoop second hash
		revPos -= scoopSize //move backwards
	}
}

func (self SimplePlotter) GetScoop(pos int32) types.Byte64 {
	var data types.Byte64
	copy(data[:], self.data[pos * scoopSize : (pos * scoopSize + scoopSize)])
	return data
}
