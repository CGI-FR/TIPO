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

type Driver struct {
	rand   *rand.Rand
	frame  *Frame
	caches *CacheList
}

func NewDriver(seed int64, frameSize int) *Driver {
	return &Driver{
		rand:   rand.New(rand.NewSource(seed)), //nolint:gosec
		frame:  NewFrame(frameSize),
		caches: NewCacheList(),
	}
}

func (d *Driver) AddSelector(selector *Selector) {
	d.caches.Add(NewCache(selector, d.rand.Int63()))
}

func (d *Driver) Run(reader RowReader, collector RowCollector) {
	for d.frame.FillUp(reader, d.caches.Post) {
		nextRow := d.frame.PullFirst()

		for _, cache := range d.caches.All() {
			selector := cache.GetSelector()
			selector.Update(nextRow, cache.PullRandom)
		}

		_ = collector.Collect(nextRow)
	}
}
