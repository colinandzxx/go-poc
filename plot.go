package poc

import (
	"encoding/binary"
	"github.com/moonfruit/go-shabal"
)

const hashSize = 32
const hashesPerScoop = 2
const scoopSize = hashesPerScoop * hashSize
const scoopsPerPlot = 4096 // original 1MB/plot = 16384
const plotSize = scoopsPerPlot * scoopSize
const hashCap = 4096

const baseLen = 16

type simplePlot struct {
	data [plotSize]byte
}

func (self *simplePlot) plotPoC1(addr uint64, nonce uint64) {
	var base [baseLen]byte
	// use BigEndian in burst code !
	binary.BigEndian.PutUint64(base[:], addr)
	binary.BigEndian.PutUint64(base[:], nonce)
	genData := make([]byte, plotSize)
	genData = append(genData, base[:]...)

	s256 := shabal.NewShabal256()
	for i := plotSize; i > 0; i += hashSize {
		s256.Reset()

		len := plotSize + baseLen - i
		if len > hashCap {
			len = hashCap
		}
		s256.Write(genData[i:len])
		copy(genData[i - hashSize:], s256.Sum(nil))
	}

	s256.Reset()
	s256.Write(genData)
	finalHash := s256.Sum(nil)
	for i := 0; i < plotSize; i++ {
		self.data[i] = genData[i] ^ finalHash[i % hashSize]
	}
}

func (self *simplePlot) plotPoC2(addr uint64, nonce uint64) {
	self.plotPoC1(addr, nonce)

	//PoC2 Rearrangement
	var hashBuffer [hashSize]byte
	revPos := plotSize - hashSize //Start at second hash in last scoop
	for pos := hashSize; pos < plotSize / 2; pos += scoopSize { //Start at second hash in first scoop
		copy(hashBuffer[:], self.data[pos:hashSize]) 		//Copy low scoop second hash to buffer
		copy(self.data[pos:], self.data[revPos:hashSize]) 	//Copy high scoop second hash to low scoop second hash
		copy(self.data[revPos:], hashBuffer[:hashSize]) 	//Copy buffer to high scoop second hash
		revPos -= scoopSize //move backwards
	}
}

func (self simplePlot) getScoop(pos int32) []byte {
	return self.data[pos * scoopSize : scoopSize]
}
