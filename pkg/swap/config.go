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
	"os"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Version   string `yaml:"version"`
	Seed      int64  `yaml:"seed"`
	FrameSize int    `yaml:"frameSize"`
	Selectors []any  `yaml:"selectors"`
}

// LoadConfigurationFromYAML returns the configuration of the yaml file in a Configuration object.
func LoadConfigurationFromYAML(filename string) (Configuration, error) {
	source, err := os.ReadFile(filename)
	if err != nil {
		return Configuration{}, fmt.Errorf("%w", err)
	}

	var def Configuration
	err = yaml.Unmarshal(source, &def)

	if err != nil {
		return def, fmt.Errorf("%w", err)
	}

	return def, nil
}

func (c Configuration) MustValidate() {
	if c.FrameSize <= 0 {
		panic(ErrInvalidConfigFrameSize)
	}
}

func (c Configuration) BuildDriver() *Driver {
	driver := NewDriver(c.Seed, c.FrameSize)

	for _, item := range c.Selectors {
		switch selectorConfig := item.(type) {
		case string:
			driver.AddSelector(NewSelector(selectorConfig))
		case []any:
			driver.AddSelector(buildSelectorFromAnyArray(selectorConfig))
		case map[any]any:
			for _, item := range selectorConfig {
				driver.AddSelector(buildSelectorFromAny(item))
			}
		default:
			panic(ErrInvalidConfigSelector)
		}
	}

	return driver
}

func buildSelectorFromAny(something any) *Selector {
	switch typed := something.(type) {
	case string:
		return NewSelector(typed)
	case []any:
		return buildSelectorFromAnyArray(typed)
	default:
		panic(ErrInvalidConfigSelector)
	}
}

func buildSelectorFromAnyArray(array []any) *Selector {
	jsonpaths := []string{}

	for _, item := range array {
		if jsonpath, ok := item.(string); !ok {
			panic(ErrInvalidConfigSelector)
		} else {
			jsonpaths = append(jsonpaths, jsonpath)
		}
	}

	return NewSelector(jsonpaths...)
}
