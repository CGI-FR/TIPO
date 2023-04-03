![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/CGI-FR/TIPO/ci.yml?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/cgi-fr/tipo)](https://goreportcard.com/report/github.com/cgi-fr/tipo)
![GitHub all releases](https://img.shields.io/github/downloads/CGI-FR/TIPO/total)
![GitHub](https://img.shields.io/github/license/CGI-FR/TIPO)
![GitHub Repo stars](https://img.shields.io/github/stars/CGI-FR/TIPO)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/CGI-FR/TIPO)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/CGI-FR/TIPO)

# TIPO: Tidy Input Permuted Output

This tool provides reproducible data swapping on a JSONLine stream.

## Configuration File

Here is an example YAML configuration file

```YAML
version: 1
seed: 42
frameSize: 1000
selectors:
  - group1:
    - attribute1.*
    - attribute2.*
  - attribute3.*
  - attribute4.*
```

- The `seed` parameter is the starting seed of the pseudo-random process.
- The `frameSize` parameter is an essential element affecting the quality of the permutation, it is the size of the processing window. To ensure good permutation quality, its value must be large, in order to have a greater number of values ready to be permuted and to reduce the chance of having permutations with identical data at the origin.
- The `selectors` parameter can be a group of attributes to be swapped together or attributes to be swapped independently of each other.

## Execution example

Suppose our input stream of type JSONLine is stored in a "stream.jsonl" file

```JSON
{"company":"Ese1","employees":[{"lastname":"Martin","firstname":"Lebaron","age":36,"children":[{"lastname":"agathe" ,"age":14}]},{"surname":"Josselin","firstname":"Jireau","age":57,"children":[{"surname":"Pierre","age" :14},{"name":"Damien","age":9}]}]}
{"company":"Ese2","employees":[{"surname":"Jérémie","firstname":"Namie","age":42,"children":[{"surname":"Patrice" ,"age":25},{"name":"Alex","age":10},{"name":"Lilie","age":2}]}]}
{"company":"Ese3","employees":[{"surname":"Océane","firstname":"Dupont","age":42,"children":[{"surname":"Alice" ,"age":25},{"name":"Maélie","age":10}]}]}
```

and the following configuration file is named "swapConf.yml":

```YAML
version: 1
seed: 42
frameSize: 1000
selectors:
   - employees.children.name.*
```

In this case we do not want to swap a group of attributes, but only the names of the children of the employees.

### 1st possibility of execution

```console
< stream.jsonl | tipo -c swapConf.yml
```

The result will be the following

```JSON
{"company":"Ese1","employees":[{"lastname":"Martin","firstname":"Lebaron","age":36,"children":[{"lastname":"Damien" ,"age":14}]},{"surname":"Josselin","firstname":"Jireau","children":[{"surname":"Alex","age":14},{" name":"Peter","age":9}]}]}
{"company":"Ese2","employees":[{"surname":"Jérémie","firstname":"Namie","age":42,"children":[{"surname":"Patrice" ,"age":25},{"name":"agathe","age":10},{"name":"Maélie","age":2}]}]}
{"company":"Ese3","employees":[{"surname":"Océane","firstname":"Dupont","age":42,"children":[{"surname":"Lilie" ,"age":25},{"name":"Alice","age":10}]}]}
```

N.B.: The tipo command takes the path to the configuration file via the `-c` flag, if no path has been provided it will try to look by default for the swap.yml file which must be in the root of the project .

### 2nd possibility of execution

In the case where the configuration file is named "swap.yml", the execution can be done as follows

```console
< stream.jsonl | type
```

## Permutation of an attribute group

When the need is to permute a group of attributes in a coherent way, for example the name and the first name of the employees and that one wishes to have an independent permutation of the names of the children of the employees then the configuration file named " swap.yml" will have the following content

```yaml
version: 1
seed: 42
frameSize: 1000
selectors:
  - group1:
    - employees.name.*
    - employees.firstname.*
  - employees.children.name.*
```

The way to execute is always the same

```console
< stream.jsonl | type
```

The result will be the following

```json
{"company":"Ese1","employees":[{"lastname":"Josselin","firstname":"Jireau","age":36,"children":[{"lastname":"Patrice" ,"age":14}]},{"surname":"Jérémie","firstname":"Namie","age":57,"children":[{"surname":"Alex","age" :14},{"name":"agathe","age":9}]}]}
{"company":"Ese2","employees":[{"lastname":"Martin","firstname":"Lebaron","age":42,"children":[{"lastname":"Lilie" ,"age":25},{"name":"Damien","age":10},{"name":"Peter","age":2}]}]}
{"company":"Ese3","employees":[{"surname":"Océane","firstname":"Dupont","age":42,"children":[{"surname":"Alice" ,"age":25},{"name":"Maélie","age":10}]}]}
```

## Swap Multiple Attribute Groups

Suppose the following incoming stream is stored in a file named stream.jsonl

```json
{"company":"company1","employees":[{"lastname":"Martin","firstname":"Lebaron","nationality":"italian","age":36,"children":[ {"name":"Damien","age":14}]}]}
{"company":"company2","employees":[{"surname":"Jérémie","firstname":"Namie","nationality":"French","age":44,"children":[ {"name":"Patrice","age":25},{"name":"agathe","age":10},{"name":"Maélie","age":2}]}] }
{"company":"company3","employees":[{"surname":"Océane","firstname":"Dupont","nationality":"Spanish","age":41,"children":[ {"name":"Lilie","age":25},{"name":"Alice","age":10}]}]}
```

The following corresponding configuration file is named configuration.yml

```yaml
version: 1
seed: 42
frameSize: 1000
selectors:
  - group1:
    - employees.name.*
    - employees.firstname.*
  - group2:
    - employees.age.*
    - employees.nationality.*
  - employees.children.*
```

The permutation of the two groups will be done independently. The execution is done as follows

```console
< stream.jsonl | type
```

And the result will be the following

```json
{"company":"company1","employees":[{"lastname":"Jérémie","firstname":"Namie","nationality":"Spanish","age":41,"children":[ {"name":"Damien","age":14}]}]}
{"company":"company2","employees":[{"surname":"Océane","firstname":"Dupont","nationality":"italian","age":36,"children":[ {"name":"Patrice","age":25},{"name":"agathe","age":10},{"name":"Maélie","age":2}]}] }
{"company":"company3","employees":[{"lastname":"Martin","firstname":"Lebaron","nationality":"French","age":44,"children":[ {"name":"Lilie","age":25},{"name":"Alice","age":10}]}]}
```

Note that the age and nationality fields have been swapped consistently and independently of the surname and first name fields, which have also been swapped consistently.

## Contributors

- CGI France ✉[Contact support](mailto:LINO.fr@cgi.com)

## License

Copyright (C) 2023 CGI France

TIPO is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

TIPO is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with TIPO. If not, see http://www.gnu.org/licenses/.
