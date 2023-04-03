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
	"fmt"
)

// Value implements the Node interface and holds a pointer to any type.
type Value struct {
	value any
}

func NewValue(value any) *Value {
	return &Value{
		value: value,
	}
}

func (v *Value) MustUpdate(n Node) {
	if other, ok := n.(*Value); ok {
		v.value = other.value
	} else {
		panic("cannot update value with non value type")
	}
}

func (v *Value) Reset() {
	v.value = nil
}

func (v *Value) Duplicate() Node { //nolint:ireturn
	return NewValue(v.value)
}

func (v *Value) MarshalJSON() ([]byte, error) {
	buf := &bytes.Buffer{}

	if _, ok := v.value.(Node); ok {
		panic(fmt.Sprintf("value is of type %T", v.value))
	}

	if b, err := json.Marshal(v.value); err != nil {
		return nil, err //nolint:wrapcheck
	} else { //nolint:golint,revive
		buf.Write(b)
	}

	return buf.Bytes(), nil
}

func (v *Value) MustMarshalJSON() []byte {
	b, err := v.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return b
}

func (v *Value) String() string {
	return string(v.MustMarshalJSON())
}
