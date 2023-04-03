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
	"bufio"
	"io"

	"github.com/cgi-fr/tipo/pkg/swap"
	"github.com/cgi-fr/tipo/pkg/swap/parser"
)

type Reader struct {
	scanner *bufio.Scanner
}

func NewReader(reader io.Reader) *Reader {
	return &Reader{
		scanner: bufio.NewScanner(reader),
	}
}

func (jr *Reader) HasRow() bool {
	return jr.scanner.Scan()
}

func (jr *Reader) ReadRow() (swap.Row, error) {
	return parser.MustParseObject(jr.scanner.Bytes()), nil
}
