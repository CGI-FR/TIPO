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

package swap

type Frame struct {
	size   int
	buffer []Row
}

// NewFrame creates a new empty frame.
func NewFrame(size int) *Frame {
	return &Frame{size: size, buffer: make([]Row, 0, size)}
}

// FillUp fills the frame with Rows read from reader,
// calls the callback function for each row added to the frame,
// and finally returns false if there is no more values in the Frame and it can't read from the reader.
func (f *Frame) FillUp(reader RowReader, callback func(Row)) bool {
	for {
		if len(f.buffer) == f.size {
			return true
		}

		if !reader.HasRow() {
			return len(f.buffer) > 0
		}

		if row, err := reader.ReadRow(); err == nil {
			f.buffer = append(f.buffer, row)
			callback(row)
		} else {
			panic(err)
		}
	}
}

// PullFirst get the first row and removes it from the frame.
func (f *Frame) PullFirst() Row {
	result := f.buffer[0]
	f.buffer = f.buffer[1:]

	return result
}
