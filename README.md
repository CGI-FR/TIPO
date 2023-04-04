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
version: 1          # version of the configuration, only 1 is allowed for now
seed: 42            # starting seed of the pseudo-random process, ensure coherence between executions
frameSize: 1000     # frameSize is the size of the processing window, should be as large as possible
selectors:          # each selector in this list will trigger a permutation between JSONLines
  - $.name          # a selector is defined by a jsonpath expression
  - $.surname
  - group:          # a group of selectors will swap attributes together
    - $.age
    - $.nationality
```

Notes :

- `seed` is optional, use it only if you need a reproducible execution (every execution gives the same result), change the value to get different results
- `frameSize` parameter is an essential element affecting the quality of the permutation, it is the size of the processing window. To ensure good permutation quality, its value must be large, in order to have a greater number of values ready to be permuted and to reduce the chance of having permutations with identical data at the origin.

## Example 1 : redistribute siblings to new parents

Suppose our input stream of type JSONLine is stored in a "stream.jsonl" file

```json
{"company":"acme","employees":[{"name":"one","children":[{"name":"child 1"},{"name":"child 2"}]},{"name":"two","children":[{"name":"child 3"},{"name":"child 4"},{"name":"child 5"}]}]}
{"company":"megacorp","employees":[{"name":"alpha","children":[{"name":"kid 1"}]},{"name":"beta","children":[{"name":"kid 2"},{"name":"kid 3"}]}]}
{"company":"dynatech","employees":[{"name":"first","children":[{"name":"offspring 1"},{"name":"offspring 2"}]},{"name":"second","children":[]}]}
```

and the following configuration file is named "swap.yml":

```YAML
version: 1
seed: 42
frameSize: 1000
selectors:
  - $.employees.*.children
```

In this case we want to swap children of the employees. Siblings will not be separated in this scenario because the JSONPath `$.employees.*.children` will select whole arrays of children. However children will be distributed to new parents.

### Execution (1st form)

```console
< stream.jsonl | tipo
```

The result will be the following

```json
{"company":"acme","employees":[{"name":"one","children":[{"name":"kid 2"},{"name":"kid 3"}]},{"name":"two","children":[{"name":"child 1"},{"name":"child 2"}]}]}
{"company":"megacorp","employees":[{"name":"alpha","children":[]},{"name":"beta","children":[{"name":"kid 1"}]}]}
{"company":"dynatech","employees":[{"name":"first","children":[{"name":"offspring 1"},{"name":"offspring 2"}]},{"name":"second","children":[{"name":"child 3"},{"name":"child 4"},{"name":"child 5"}]}]}
```

N.B.: The tipo command can use the path to the configuration file via the `-c` flag, if no path has been provided it will try to look by default for the swap.yml file which must be in the root of the project.

## Example 2 : mix the children together

With the same data file `stream.jsonl`, consider this TIPO configuration :

```YAML
version: 1
seed: 42
frameSize: 1000
selectors:
  - $.employees.*.children.*
```

In this case we want to mix children together. Siblings will be separated in this scenario because the JSONPath `$.employees.*.children.*` will select each child individualy. However each parent will keep its original number of children.

### Execution (2nd form)

```console
cat stream.jsonl | tipo
```

The result will be the following

```json
{"company":"acme","employees":[{"name":"one","children":[{"name":"offspring 2"},{"name":"child 2"}]},{"name":"two","children":[{"name":"child 5"},{"name":"kid 2"},{"name":"child 3"}]}]}
{"company":"megacorp","employees":[{"name":"alpha","children":[{"name":"kid 3"}]},{"name":"beta","children":[{"name":"kid 1"},{"name":"child 4"}]}]}
{"company":"dynatech","employees":[{"name":"first","children":[{"name":"child 1"},{"name":"offspring 1"}]},{"name":"second","children":[]}]}
```

## Example 3 : permutation of a group of attributes

Let's change our dataset example by adding a new information `childnumber`.

```json
{"company":"acme","employees":[{"name":"one","childnumber":2,"children":[{"name":"child 1"},{"name":"child 2"}]},{"name":"two","childnumber":3,"children":[{"name":"child 3"},{"name":"child 4"},{"name":"child 5"}]}]}
{"company":"megacorp","employees":[{"name":"alpha","childnumber":1,"children":[{"name":"kid 1"}]},{"name":"beta","childnumber":2,"children":[{"name":"kid 2"},{"name":"kid 3"}]}]}
{"company":"dynatech","employees":[{"name":"first","childnumber":2,"children":[{"name":"offspring 1"},{"name":"offspring 2"}]},{"name":"second","childnumber":0,"children":[]}]}
```

When the need is to permute a group of attributes in a coherent way, for example if `childnumber` and the `children` list must be swapped coherently then the configuration file named `swap.yml` will have the following content :

```yaml
version: 1
seed: 42
frameSize: 1000
selectors:
  - group:
    - $.employees.*.childnumber
    - $.employees.*.children
```

### Execution (3rd form)

```console
tipo < stream.jsonl
```

The result will be the following

```json
{"company":"acme","employees":[{"name":"one","childnumber":2,"children":[{"name":"kid 2"},{"name":"kid 3"}]},{"name":"two","childnumber":2,"children":[{"name":"child 1"},{"name":"child 2"}]}]}
{"company":"megacorp","employees":[{"name":"alpha","childnumber":0,"children":[]},{"name":"beta","childnumber":1,"children":[{"name":"kid 1"}]}]}
{"company":"dynatech","employees":[{"name":"first","childnumber":2,"children":[{"name":"offspring 1"},{"name":"offspring 2"}]},{"name":"second","childnumber":3,"children":[{"name":"child 3"},{"name":"child 4"},{"name":"child 5"}]}]}
```

## Example 4 : mix multiple groups and attributes

Suppose the following incoming stream is stored in a file named stream.jsonl

```json
{"company":"acme","employees":[{"name":"one","surname":"ONE","age":20,"nationality":"Kenyan"},{"name":"two","surname":"TWO","age":30,"nationality":"Icelandic"}]}
{"company":"megacorp","employees":[{"name":"alpha","surname":"ALPHA","age":40,"nationality":"Colombian"},{"name":"beta","surname":"BETA","age":50,"nationality":"Malaysian"}]}
{"company":"dynatech","employees":[{"name":"first","surname":"FIRST","age":60,"nationality":"Belgian"},{"name":"second","surname":"SECOND","age":70,"nationality":"Egyptian"}]}
```

The following corresponding configuration file is named configuration.yml

```yaml
version: 1
seed: 42
frameSize: 1000
selectors:
  - group1: # name and surname will be swapped together
    - employees.*.name
    - employees.*.surname
  - group2: # age and nationality will be swapped together
    - employees.*.age
    - employees.*.nationality
  - company
```

The permutation of the two groups will be done independently. The execution is done as follows

```console
< stream.jsonl | type
```

And the result will be the following

```json
{"company":"acme","employees":[{"name":"beta","surname":"BETA","age":50,"nationality":"Malaysian"},{"name":"one","surname":"ONE","age":30,"nationality":"Icelandic"}]}
{"company":"megacorp","employees":[{"name":"second","surname":"SECOND","age":70,"nationality":"Egyptian"},{"name":"alpha","surname":"ALPHA","age":20,"nationality":"Kenyan"}]}
{"company":"dynatech","employees":[{"name":"first","surname":"FIRST","age":60,"nationality":"Belgian"},{"name":"two","surname":"TWO","age":40,"nationality":"Colombian"}]}
```

Note that the age and nationality fields have been swapped consistently and independently of the surname and first name fields, which have also been swapped consistently.

## Different ways to configure a group

In previous examples, this kind of group configuration has been presented.

```yaml
selectors:
  - groupname:
    - employees.*.name
    - employees.*.surname
```

But there is other ways to configure a group.

1- **Inline array**

```yaml
selectors:
  - ["employees.*.name", "employees.*.surname"]
```

2- **Inline array with map**

```yaml
selectors:
  - groupname: ["employees.*.name", "employees.*.surname"]
```

## Contributors

- CGI France ✉[Contact support](mailto:LINO.fr@cgi.com)

## License

Copyright (C) 2023 CGI France

TIPO is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

TIPO is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with TIPO. If not, see http://www.gnu.org/licenses/.
