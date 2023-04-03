// Copyright (C) 2023 CGI France
//
// This file is part of TIPO.
//
// TIPO is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// TIPO is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with TIPO.  If not, see <http://www.gnu.org/licenses/>.

package jsonline

import (
	"fmt"
	"io"

	"github.com/cgi-fr/tipo/pkg/swap"
)

type Collector struct {
	out io.Writer
}

// NewCollector creates a new row collector that writes rows as  format.
func NewCollector(out io.Writer) *Collector {
	return &Collector{out: out}
}

// Collect the finalized row.
func (jc *Collector) Collect(row swap.Row) error {
	fmt.Fprintln(jc.out, row)

	return nil
}
