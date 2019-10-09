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

var Cfg Config

type Config struct {
	AvgBaseTargetNum uint32
	ConsensusInterval uint32
	MaxBaseTarget uint64
	//Loglevel
}

func init() {
	Cfg.AvgBaseTargetNum = 24
	Cfg.ConsensusInterval = 240 //s
	Cfg.MaxBaseTarget = 0x444444444
}

func (self *Config) Load() error {
	return nil
}
