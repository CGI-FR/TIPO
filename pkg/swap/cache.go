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

import (
	"math/rand"
)

type Cache struct {
	selector *Selector
	buffer   []Tuple
	rand     *rand.Rand
}

// NewCache creates a new empty cache.
func NewCache(selector *Selector, seed int64) *Cache {
	return &Cache{selector: selector, buffer: []Tuple{}, rand: rand.New(rand.NewSource(seed))} //nolint:gosec
}

// Post the row to feed the cache with new tuples.
func (c *Cache) Post(row Row) {
	results := c.selector.Select(row)
	duplicates := make([]Tuple, len(results))

	for index, tuple := range results {
		duplicates[index] = tuple.Duplicate().(Tuple) //nolint:forcetypeassert
	}

	c.buffer = append(c.buffer, duplicates...)
}

// PullRandom removes a random tuple from the cache and returns it.
func (c *Cache) PullRandom() Tuple {
	index := c.rand.Intn(len(c.buffer))
	result := c.buffer[index]
	c.buffer = append(c.buffer[:index], c.buffer[index+1:]...)

	return result
}

// GetSelector returns the selector used by the cache.
func (c *Cache) GetSelector() *Selector {
	return c.selector
}
