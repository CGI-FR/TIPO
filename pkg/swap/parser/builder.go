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

var emptySlice = []any{} //nolint:gochecknoglobals

// Builder is a basic type builder. It uses a stack model to build where maps
// (objects) and slices (arrays) add pushed on the stack and closed with a
// pop.
type Builder struct {
	stack  []any
	starts []int
}

func NewBuilder() *Builder {
	return &Builder{
		stack:  []any{},
		starts: []int{},
	}
}

// Reset the builder.
func (b *Builder) Reset() {
	if 0 < cap(b.stack) && 0 < len(b.stack) {
		b.stack = b.stack[:0]
		b.starts = b.starts[:0]
	} else {
		b.stack = make([]any, 0, 64)  //nolint:gomnd
		b.starts = make([]int, 0, 16) //nolint:gomnd
	}
}

// Object pushs an object onto the stack. A key must be
// provided if the top of the stack is an object and must not be
// provided if the op of the stack is an array or slice.
func (b *Builder) Object(key ...string) error {
	newObj := NewObject()

	if 0 < len(key) {
		if len(b.starts) == 0 || 0 <= b.starts[len(b.starts)-1] {
			return ErrCannotUseKeyWhenPushingToArray
		}

		if obj, _ := b.stack[len(b.stack)-1].(*Object); obj != nil {
			obj.m[key[0]] = newObj
			obj.keys = append(obj.keys, key[0])
		}
	} else if 0 < len(b.starts) && b.starts[len(b.starts)-1] < 0 {
		return ErrMustHaveKeyWhenPushingToObject
	}

	b.starts = append(b.starts, -1)
	b.stack = append(b.stack, newObj)

	return nil
}

// Array pushs a []any onto the stack. A key must be provided if the
// top of the stack is an object (map) and must not be provided if the op of
// the stack is an array or slice.
func (b *Builder) Array(key ...string) error {
	if 0 < len(key) {
		if len(b.starts) == 0 || 0 <= b.starts[len(b.starts)-1] {
			return ErrCannotUseKeyWhenPushingToArray
		}

		b.stack = append(b.stack, key[0])
	} else if 0 < len(b.starts) && b.starts[len(b.starts)-1] < 0 {
		return ErrMustHaveKeyWhenPushingToObject
	}

	b.starts = append(b.starts, len(b.stack))
	b.stack = append(b.stack, NewArray(emptySlice))

	return nil
}

// Value pushs a value onto the stack. A key must be provided if the top of
// the stack is an object (map) and must not be provided if the op of the
// stack is an array or slice.
func (b *Builder) Value(value any, key ...string) error {
	switch {
	case 0 < len(key):
		if len(b.starts) == 0 || 0 <= b.starts[len(b.starts)-1] {
			return ErrCannotUseKeyWhenPushingToArray
		}

		if obj, _ := b.stack[len(b.stack)-1].(*Object); obj != nil {
			obj.m[key[0]] = NewValue(value)
			obj.keys = append(obj.keys, key[0])
		}
	case 0 < len(b.starts) && b.starts[len(b.starts)-1] < 0:
		return ErrMustHaveKeyWhenPushingToObject
	default:
		b.stack = append(b.stack, NewValue(value))
	}

	return nil
}

// Pop the stack, closing an array or object.
func (b *Builder) Pop() {
	if 0 < len(b.starts) { //nolint:nestif
		start := b.starts[len(b.starts)-1]
		if 0 <= start { // array
			start++
			size := len(b.stack) - start
			array := make([]any, size)
			copy(array, b.stack[start:len(b.stack)])
			b.stack = b.stack[:start]
			b.stack[start-1] = array

			if 2 < len(b.stack) { //nolint:gomnd
				if k, ok := b.stack[len(b.stack)-2].(string); ok {
					if obj, _ := b.stack[len(b.stack)-3].(*Object); obj != nil {
						obj.m[k] = NewArray(array)
						obj.keys = append(obj.keys, k)
						b.stack = b.stack[:len(b.stack)-2]
					}
				}
			}
		} else if 1 < len(b.starts) && b.starts[len(b.starts)-2] < 0 {
			b.stack = b.stack[:len(b.stack)-1]
		}

		b.starts = b.starts[:len(b.starts)-1]
	}
}

// PopAll repeats Pop until all open arrays or objects are closed.
func (b *Builder) PopAll() {
	for 0 < len(b.starts) {
		b.Pop()
	}
}

// Result of the builder is returned.
// This is the first item pushed on to the stack.
func (b *Builder) Result() Node { //nolint:ireturn
	var (
		result Node
		isNode bool
	)

	if 0 < len(b.stack) {
		if result, isNode = b.stack[0].(Node); !isNode {
			panic("cannot build node")
		}
	}

	return result
}
