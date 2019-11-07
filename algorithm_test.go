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
 * @File: algorithm_test.go
 * @LastModified: 2019-11-07 14:52:23
 */

package poc

import (
	"bytes"
	"encoding/hex"
	"github.com/colinandzxx/go-consensus/types"
	"math/big"
	"testing"
)

func reserve(buf *[]byte) {
	for i := 0; i < len(*buf)/2; i++ {
		(*buf)[i], (*buf)[len(*buf)-i-1] = (*buf)[len(*buf)-i-1], (*buf)[i]
	}
}

func Test_calculateGenerationSignature(t *testing.T) {
	except, _ := hex.DecodeString("24c64309d302086ecc7b03d6cb3f7287fabd45e24c476a0c2cb61a73a920a4ad")

	gs, _ := hex.DecodeString("8f1c281852952d203ade668f2f3114ac6c01f1260ab9567f24a4d75b8efbae5c")
	reserve(&gs)
	var lastGenSig types.Byte32
	copy(lastGenSig[:], gs[:32])
	t.Logf("lastGenSig: %x\n", lastGenSig)
	var lastGenId uint64 = 58970560028650869
	sig := CalculateGenerationSignature(lastGenSig, lastGenId)
	var bufSig []byte = make([]byte, 32)
	copy(bufSig, sig[:32])
	reserve(&bufSig)
	t.Logf("bufSig: %x", bufSig)

	if bytes.Compare(except, bufSig) != 0 {
		t.Errorf("fail to calculateGenerationSignature, ret: %x, want: %x", bufSig, except)
	}
}

func Test_calculateScoop(t *testing.T) {
	lastGenSig := types.Byte32{}
	var lastGenId uint64 = 0xFFFFFFFFFFFFFFFF
	sig := CalculateGenerationSignature(lastGenSig, lastGenId)
	scoop := CalculateScoop(sig, 1)
	t.Logf("%v", scoop)
}

func Test_calculateDeadline(t *testing.T) {
	//lastGenSig := types.Byte32{
	//	0x12, 0xbe, 0xf4, 0x68, 0x72, 0x84, 0x72, 0xd7,
	//	0x82, 0xc5, 0x61, 0x42, 0xb8, 0x5a, 0x9c, 0x7c,
	//	0x6b, 0x19, 0x9d, 0x2c, 0xe4, 0x1a, 0x19, 0xad,
	//	0xf7, 0x75, 0xb3, 0xd0, 0x08, 0x48, 0x17, 0x6f,
	//}
	//var lastGenId uint64 = 10725173944100240
	//var addr uint64 = 17023786578764300
	//var nonce uint64 = 2803155511487343564
	////var scoopIndex int32 = 3320
	//var baseTarget uint64 = 43707
	//var height uint64 = 22532

	lgs, _ := big.NewInt(0).SetString("cbe922a9bcb308e271badc7b30f541131b841c163434899728a3fba68e6d4699", 16)
	blgs := lgs.Bytes()
	for i := 0; i < len(blgs)/2; i++ {
		blgs[i], blgs[len(blgs)-i-1] = blgs[len(blgs)-i-1], blgs[i]
	}

	var lastGenSig types.Byte32
	copy(lastGenSig[:], blgs[:32])
	t.Logf("lastGenSig: %v\n", lastGenSig)

	var lastGenId uint64 = 6986083158196820480
	var addr uint64 = 12645014192979329200
	var nonce uint64 = 15032170525642997731
	//var scoopIndex int32 = 3320
	var baseTarget uint64 = 37848
	var height uint64 = 1

	sig := CalculateGenerationSignature(lastGenSig, lastGenId)
	t.Logf("sig: %x\n", sig)

	scoopIndex := CalculateScoop(sig, height)
	t.Logf("scoopIndex: %v\n", scoopIndex)

	var plotter SimplePlotter
	plotter.PlotPoC2(addr, nonce)

	dl := CalculateDeadline(sig, plotter.GetScoop(scoopIndex), baseTarget)
	t.Logf("deadline: %v\n", dl.Uint64())
}

func Test_calculateHit(t *testing.T) {
	//lastGenSig := types.Byte32{
	//	0xd1, 0xa9, 0x39, 0x96, 0x89, 0x2e, 0x56, 0xa3,
	//	0xf1, 0x9f, 0x33, 0xc7, 0xd1, 0x1f, 0x67, 0x74,
	//	0x40, 0x95, 0x7f, 0xa1, 0xf2, 0x59, 0x0d, 0x44,
	//	0xd3, 0xc6, 0xa6, 0x1e, 0xd8, 0x0b, 0x15, 0xb9,
	//}
	//var lastGenId uint64 = 4640397412688285442
	//var addr uint64 = 17023786578764300
	//var nonce uint64 = 1910725663170410144
	////var scoopIndex int32 = 3320
	//var height uint64 = 22550

	lgs, _ := big.NewInt(0).SetString("cbe922a9bcb308e271badc7b30f541131b841c163434899728a3fba68e6d4699", 16)
	blgs := lgs.Bytes()
	for i := 0; i < len(blgs)/2; i++ {
		blgs[i], blgs[len(blgs)-i-1] = blgs[len(blgs)-i-1], blgs[i]
	}

	var lastGenSig types.Byte32
	copy(lastGenSig[:], blgs[:32])
	t.Logf("lastGenSig: %v\n", lastGenSig)

	var lastGenId uint64 = 6986083158196820480
	var addr uint64 = 12645014192979329200
	var nonce uint64 = 15032170525642997731
	//var scoopIndex int32 = 3320
	var height uint64 = 1

	sig := CalculateGenerationSignature(lastGenSig, lastGenId)
	copy(sig[:], lgs.Bytes()[:32])
	t.Logf("sig: %x\n", sig)

	scoopIndex := CalculateScoop(sig, height)
	t.Logf("scoopIndex: %v\n", scoopIndex)

	var plotter SimplePlotter
	plotter.PlotPoC2(addr, nonce)

	ht := CalculateHit(sig, plotter.GetScoop(scoopIndex))
	t.Logf("hit: %v\n", ht.Uint64())
}

func Test_calculateHit2(t *testing.T) {
	gs, _ := hex.DecodeString("24c64309d302086ecc7b03d6cb3f7287fabd45e24c476a0c2cb61a73a920a4ad")
	for i := 0; i < len(gs)/2; i++ {
		gs[i], gs[len(gs)-i-1] = gs[len(gs)-i-1], gs[i]
	}
	var sig types.Byte32
	copy(sig[:], gs[:32])
	t.Logf("sig: %x\n", sig)

	var addr uint64 = 17023786578764300
	var nonce uint64 = 16186784844367304518
	var height uint64 = 23206

	scoopIndex := CalculateScoop(sig, height)
	t.Logf("scoopIndex: %v\n", scoopIndex)

	var plotter SimplePlotter
	plotter.PlotPoC2(addr, nonce)

	ht := CalculateHit(sig, plotter.GetScoop(scoopIndex))
	t.Logf("hit: %v\n", ht.Uint64())
}
