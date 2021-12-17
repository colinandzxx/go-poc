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
 * @File: plot_test.go
 * @LastModified: 2019-10-08 17:31:49
 */

package poc

import "testing"

func Test_plotPoC1(t *testing.T) {
	var plotter SimplePlotter
	plotter.PlotPoC1(1, 0)
	t.Logf("nonce: %x\n", plotter.data)
}

func Test_plotPoC2(t *testing.T) {
	var plotter SimplePlotter
	plotter.PlotPoC2(1, 1)
	t.Logf("nonce: %x\n", plotter.data)
}
