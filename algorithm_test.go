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
 * @LastModified: 2019-10-10 14:42:30
 */

package poc

import (
	"github.com/colinandzxx/go-consensus/types"
	"testing"
)

func Test_calculateGenerationSignature(t *testing.T) {
	lastGenSig := types.Byte32{}
	var lastGenId uint64 = 0xFFFFFFFFFFFFFFFF
	sig := calculateGenerationSignature(lastGenSig, lastGenId)
	t.Logf("%x", sig)
}

func Test_calculateScoop(t *testing.T) {
	lastGenSig := types.Byte32{}
	var lastGenId uint64 = 0xFFFFFFFFFFFFFFFF
	sig := calculateGenerationSignature(lastGenSig, lastGenId)
	scoop := calculateScoop(sig, 1)
	t.Logf("%v", scoop)
}
