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

import "github.com/ohler55/ojg/oj"

// Node is the common interface for all elements.
type Node interface {
	MustUpdate(n Node)
	Reset()
	Duplicate() Node
	MarshalJSON() ([]byte, error)
	MustMarshalJSON() []byte
	String() string
}

func MustParseObject(data []byte) *Object {
	handler := NewHandler()
	if err := oj.Tokenize(data, handler); err != nil {
		panic(err)
	}

	if obj, ok := handler.Node().(*Object); ok {
		return obj
	}

	panic("can't parse object")
}

type Handler struct {
	builder *Builder
	key     *string
}

func NewHandler() *Handler {
	return &Handler{
		builder: NewBuilder(),
		key:     nil,
	}
}

func (z *Handler) Node() Node { //nolint:ireturn
	return z.builder.Result()
}

func (z *Handler) Value(v any) {
	if z.key != nil {
		z.dieOnError(z.builder.Value(v, *z.key))
	} else {
		z.dieOnError(z.builder.Value(v))
	}

	z.key = nil
}

// Null is called when a JSON null is encountered.
func (z *Handler) Null() {
	z.Value(nil)
}

// Bool is called when a JSON true or false is encountered.
func (z *Handler) Bool(v bool) {
	z.Value(v)
}

// Int is called when a JSON integer is encountered.
func (z *Handler) Int(v int64) {
	z.Value(v)
}

// Float is called when a JSON decimal is encountered that fits into a
// float64.
func (z *Handler) Float(v float64) {
	z.Value(v)
}

// Number is called when a JSON number is encountered that does not fit
// into an int64 or float64.
func (z *Handler) Number(v string) {
	z.Value(v)
}

// String is called when a JSON string is encountered.
func (z *Handler) String(v string) {
	z.Value(v)
}

// ObjectStart is called when a JSON object start '{' is encountered.
func (z *Handler) ObjectStart() {
	if z.key != nil {
		z.dieOnError(z.builder.Object(*z.key))
	} else {
		z.dieOnError(z.builder.Object())
	}

	z.key = nil
}

// ObjectEnd is called when a JSON object end '}' is encountered.
func (z *Handler) ObjectEnd() {
	z.builder.Pop()
}

// Key is called when a JSON object key is encountered.
func (z *Handler) Key(k string) {
	z.key = &k
}

// ArrayStart is called when a JSON array start '[' is encountered.
func (z *Handler) ArrayStart() {
	if z.key != nil {
		z.dieOnError(z.builder.Array(*z.key))
	} else {
		z.dieOnError(z.builder.Array())
	}

	z.key = nil
}

// ArrayEnd is called when a JSON array end ']' is encountered.
func (z *Handler) ArrayEnd() {
	z.builder.Pop()
}

func (z *Handler) dieOnError(err error) {
	if err != nil {
		panic(err)
	}
}
