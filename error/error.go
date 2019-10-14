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
 * @File: error.go
 * @LastModified: 2019-10-09 10:36:12
 */

package error

import (
	"fmt"
	"github.com/colinandzxx/go-consensus/types"
)

type GetBlockError struct { Height uint64 }

func (err GetBlockError) Error() string {
	return fmt.Sprintf("can't get block(%v)", err.Height)
}

type GetHeaderEnum int
const (
	GetHeaderMethod = iota
	GetHeaderByHeightMethod
	GetHeaderByHashMethod
)

type GetHeaderError struct {
	Height uint64
	Hash types.Byte32
	Method GetHeaderEnum
}

func (err GetHeaderError) Error() string {
	var ret string
	switch err.Method {
	case GetHeaderMethod:
		ret = fmt.Sprintf("can't get header(%v, %x)", err.Height, err.Hash)
	case GetHeaderByHeightMethod:
		ret = fmt.Sprintf("can't get header(%v)", err.Height)
	case GetHeaderByHashMethod:
		ret = fmt.Sprintf("can't get header(%x)", err.Hash)
	default:
		panic("unkonwn value")
	}
	return ret
}

type decError struct { msg string }

func (err decError) Error() string { return err.msg }

var (
	ErrGetConsensusData   	= &decError{"can't get consensus data"}
	ErrCalculateBaseTarget  = &decError{"can't calculate basetarget"}
	ErrGetDifficulty = &decError{"can't get difficulty"}
	ErrInvalidDifficulty = &decError{"non-positive difficulty"}
	ErrSickTimestamp = &decError{"sick timestamp"}
	ErrTypeConver = &decError{"type conversion error"}
	ErrNilHeader = &decError{"nil header"}
	ErrNilBlock = &decError{"nil block"}
	ErrNilOriData = &decError{"nil origin consensus data"}
)
