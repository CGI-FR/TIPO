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

// Object implements Node and Keyed interfaces and represents an odered list of key/value pairs.
type Object struct {
	keys []string
	m    map[string]any
}

func NewObject() *Object {
	return &Object{
		keys: []string{},
		m:    map[string]any{},
	}
}

func (o *Object) MustUpdate(n Node) {
	if other, ok := n.(*Object); ok {
		o.Reset()

		for _, k := range other.keys {
			o.SetValueForKey(k, other.MustValueForKey(k))
		}
	} else {
		panic("cannot update object with non object type")
	}
}

func (o *Object) Reset() {
	o.keys = []string{}
	o.m = map[string]any{}
}

func (o *Object) Duplicate() Node { //nolint:ireturn
	dup := NewObject()

	for _, key := range o.keys {
		value := o.MustValueForKey(key)
		if node, ok := value.(Node); ok {
			dup.SetValueForKey(key, node.Duplicate())
		} else {
			dup.SetValueForKey(key, value)
		}
	}

	return dup
}

// ValueForKey should return the value associated with the key or nil if
// no entry exists for the key.
func (o *Object) ValueForKey(key string) (any, bool) {
	value, has := o.m[key]

	return value, has
}

// MustValueForKey should return the value associated with the key.
func (o *Object) MustValueForKey(key string) any {
	value, has := o.m[key]
	if !has {
		panic("key does not exist")
	}

	return value
}

// SetValueForKey sets the value for a key in the collection.
func (o *Object) SetValueForKey(key string, value any) {
	o.m[key] = value

	for _, k := range o.keys {
		if key == k {
			return
		}
	}

	o.keys = append(o.keys, key)
}

// RemoveValueForKey removes the value for a key in the collection.
func (o *Object) RemoveValueForKey(key string) {
	index := -1

	for i, k := range o.keys {
		if key == k {
			index = i

			break
		}
	}

	delete(o.m, key)

	if index != -1 {
		o.keys = append(o.keys[:index], o.keys[index+1:]...)
	}
}

// Keys should return an list of the keys for all the entries in the
// collection.
func (o *Object) Keys() []string {
	return o.keys
}

func (o *Object) MarshalJSON() ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.WriteRune('{')

	for index, key := range o.keys {
		buf.WriteRune('"')
		buf.WriteString(key)
		buf.WriteString(`":`)

		if b, err := json.Marshal(o.m[key]); err != nil {
			return nil, err //nolint:wrapcheck
		} else { //nolint:golint,revive
			buf.Write(b)
		}

		if index+1 < len(o.keys) {
			buf.WriteRune(',')
		}
	}

	buf.WriteRune('}')

	return buf.Bytes(), nil
}

func (o *Object) MustMarshalJSON() []byte {
	b, err := o.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return b
}

func (o *Object) String() string {
	return string(o.MustMarshalJSON())
}
