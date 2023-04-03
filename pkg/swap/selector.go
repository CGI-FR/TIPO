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
	"fmt"

	"github.com/cgi-fr/tipo/pkg/swap/parser"
	"github.com/ohler55/ojg/jp"
	"github.com/rs/zerolog/log"
)

type Selector struct {
	jpexps []jp.Expr
}

// NewSelector creates a new selector.
func NewSelector(jsonpaths ...string) *Selector {
	log.Debug().Strs("jsonpaths", jsonpaths).Msg("Created selector")

	jpexps := []jp.Expr{}
	for _, jsonpath := range jsonpaths {
		jpexps = append(jpexps, jp.MustParseString(jsonpath))
	}

	return &Selector{jpexps: jpexps}
}

// Select all tuples from the given row that are matched by the selector.
func (s *Selector) Select(row Row) []Tuple {
	results := [][]any{}
	cardinality := -1

	for _, jpexp := range s.jpexps {
		column := jpexp.Get(row)
		if cardinality != -1 && len(column) != cardinality {
			panic("group of selectors matched different cardinality")
		} else {
			cardinality = len(column)
		}

		results = append(results, column)
	}

	tuples := make([]Tuple, cardinality)
	for i := 0; i < cardinality; i++ {
		tuples[i] = NewTuple()
		for _, column := range results {
			tuples[i].Append(column[i])
		}
	}

	return tuples
}

// Update all tuples matched by the selector with values provided by the provider.
func (s *Selector) Update(row Row, provider Provider) {
	for _, tuple := range s.Select(row) {
		newTuple := provider()
		log.Debug().Interface("orig", tuple).Interface("new", newTuple).Msg("Swap values")

		for index := 0; index < tuple.Size(); index++ {
			obj := tuple.ValueAtIndex(index)
			if node, ok := obj.(parser.Node); ok {
				node.MustUpdate(newTuple.ValueAtIndex(index).(parser.Node)) //nolint:forcetypeassert
			} else {
				panic(fmt.Sprintf("not a node %v (%T)\n", obj, obj))
			}
		}
	}
}
