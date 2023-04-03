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

package parser

import (
	"bytes"
	"encoding/json"
)

// Array implements Node and Indexed interfaces and represents an odered list of values.
type Array struct {
	values []any
}

func NewArray(values []any) *Array {
	return &Array{
		values: values,
	}
}

func (a *Array) MustUpdate(n Node) {
	if other, ok := n.(*Array); ok {
		a.Reset()

		a.values = append(a.values, other.values...)
	} else {
		panic("cannot update array with non array type")
	}
}

func (a *Array) Reset() {
	a.values = []any{}
}

func (a *Array) Duplicate() Node { //nolint:ireturn
	values := make([]any, len(a.values))
	for i, v := range a.values {
		if node, ok := v.(Node); ok {
			values[i] = node.Duplicate()
		} else {
			values[i] = v
		}
	}

	return NewArray(values)
}

func (a *Array) Append(value any) {
	a.values = append(a.values, value)
}

// ValueAtIndex should return the value at the provided index or nil if no
// entry exists at the index.
func (a *Array) ValueAtIndex(index int) any {
	return a.values[index]
}

// SetValueAtIndex should set the value at the provided index.
func (a *Array) SetValueAtIndex(index int, value any) {
	a.values[index] = value
}

// Size should return the size for the collection.
func (a *Array) Size() int {
	return len(a.values)
}

func (a *Array) MarshalJSON() ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.WriteRune('[')

	for index, value := range a.values {
		if b, err := json.Marshal(value); err != nil {
			return nil, err //nolint:wrapcheck
		} else { //nolint:golint,revive
			buf.Write(b)
		}

		if index+1 < len(a.values) {
			buf.WriteRune(',')
		}
	}

	buf.WriteRune(']')

	return buf.Bytes(), nil
}

func (a *Array) MustMarshalJSON() []byte {
	b, err := a.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return b
}

func (a *Array) String() string {
	return string(a.MustMarshalJSON())
}
