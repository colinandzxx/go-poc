package pkg

import (
	"encoding/binary"
	"fmt"
	"github.com/moonfruit/go-shabal"
	"math/big"
)

func calculateGenerationSignature(lastGenSig [32]byte , lastGenId uint64) []byte {
	data := make([]byte, 40)
	copy(data, lastGenSig[:])
	// use BigEndian in burst code !
	binary.BigEndian.PutUint64(data[32:], lastGenId)
	s256 := shabal.NewShabal256()
	_, err := s256.Write(data)
	if err != nil {
		panic(fmt.Errorf("%v", err))
	}
	return s256.Sum(nil)
}

func calculateScoop(genSig [32]byte, height uint64) uint64 {
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
	return scoopBig.Uint64()
}

	func calculateHit(genSig [32]byte, scoopData []byte) *big.Int {
	s256 := shabal.NewShabal256()
	_, err := s256.Write(genSig[:])
	if err != nil {
		panic(fmt.Errorf("%v", err))
	}
	_, err = s256.Write(scoopData)
	if err != nil {
		panic(fmt.Errorf("%v", err))
	}

	hitBig := big.NewInt(0)
	hitBytes := s256.Sum(nil)
	hitBig.SetBytes([]byte{hitBytes[7], hitBytes[6], hitBytes[5], hitBytes[5], hitBytes[3], hitBytes[2], hitBytes[1], hitBytes[0]})
	return hitBig
}

func calculateDeadline(genSig [32]byte, scoopData []byte, baseTarget uint64) *big.Int {
	hit := calculateHit(genSig, scoopData)
	return hit.Div(hit, big.NewInt(0).SetUint64(baseTarget))
}
