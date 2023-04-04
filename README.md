![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/CGI-FR/TIPO/ci.yml?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/cgi-fr/tipo)](https://goreportcard.com/report/github.com/cgi-fr/tipo)
![GitHub all releases](https://img.shields.io/github/downloads/CGI-FR/TIPO/total)
![GitHub](https://img.shields.io/github/license/CGI-FR/TIPO)
![GitHub Repo stars](https://img.shields.io/github/stars/CGI-FR/TIPO)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/CGI-FR/TIPO)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/CGI-FR/TIPO)

# TIPO: Tidy Input Permuted Output

This tool enables reproducible data shuffling for JSONLines streams.

## Configuration File

Here is an example YAML configuration file

```YAML
version: 1          # Version of the configuration, only 1 is allowed for now
seed: 42            # Starting seed for the pseudo-random process, ensures consistency between executions
frameSize: 1000     # Frame size is the size of the processing window; should be as large as possible
selectors:          # Each selector in this list triggers a permutation between JSONLines
  - $.name          # A selector is defined by a JSONPath expression
  - $.surname
  - group:          # A group of selectors swaps attributes together
    - $.age
    - $.nationality
```

Notes :

- The `seed` parameter is optional. Use it only if you need a reproducible execution (every execution gives the same result). Change the value to obtain different results.
- The `frameSize` parameter is a crucial element affecting the quality of the permutation, as it defines the size of the processing window. To ensure good permutation quality, set its value as large as possible. This allows for a greater number of values to be permuted and reduces the likelihood of permutations with identical data at the origin.

## Example 1: Redistribute Siblings to New Parents

Suppose our input stream of type JSONLines is stored in a "stream.jsonl" file:

```json
{"company":"acme","employees":[{"name":"one","children":[{"name":"child 1"},{"name":"child 2"}]},{"name":"two","children":[{"name":"child 3"},{"name":"child 4"},{"name":"child 5"}]}]}
{"company":"megacorp","employees":[{"name":"alpha","children":[{"name":"kid 1"}]},{"name":"beta","children":[{"name":"kid 2"},{"name":"kid 3"}]}]}
{"company":"dynatech","employees":[{"name":"first","children":[{"name":"offspring 1"},{"name":"offspring 2"}]},{"name":"second","children":[]}]}
```

The following configuration file, named `swap.yml`, is used:

```yaml
version: 1
seed: 42
frameSize: 1000
selectors:
  - $.employees.*.children
```

In this example, we want to swap the children of the employees. Siblings will not be separated, as the JSONPath `$.employees.*.children` selects entire arrays of children. However, children will be redistributed to new parents.

### Execution (1st form)

```console
< stream.jsonl | tipo
```

The result will be the following:

```json
{"company":"acme","employees":[{"name":"one","children":[{"name":"kid 2"},{"name":"kid 3"}]},{"name":"two","children":[{"name":"child 1"},{"name":"child 2"}]}]}
{"company":"megacorp","employees":[{"name":"alpha","children":[]},{"name":"beta","children":[{"name":"kid 1"}]}]}
{"company":"dynatech","employees":[{"name":"first","children":[{"name":"offspring 1"},{"name":"offspring 2"}]},{"name":"second","children":[{"name":"child 3"},{"name":"child 4"},{"name":"child 5"}]}]}
```

Note: The `tipo` command can use the path to the configuration file with the `-c` flag. If no path is provided, it will look for the `swap.yml` file by default, which must be located in the project's root directory.

## Example 2: Mix the Children Together

Using the same data file `stream.jsonl`, consider this TIPO configuration:

```yaml
version: 1
seed: 42
frameSize: 1000
selectors:
  - $.employees.*.children.*
```

In this example, we want to mix the children together. Siblings will be separated because the JSONPath `$.employees.*.children.*` selects each child individually. However, each parent will retain their original number of children.

### Execution (2nd form)

```console
cat stream.jsonl | tipo
```

The result will be the following:

```json
{"company":"acme","employees":[{"name":"one","children":[{"name":"offspring 2"},{"name":"child 2"}]},{"name":"two","children":[{"name":"child 5"},{"name":"kid 2"},{"name":"child 3"}]}]}
{"company":"megacorp","employees":[{"name":"alpha","children":[{"name":"kid 3"}]},{"name":"beta","children":[{"name":"kid 1"},{"name":"child 4"}]}]}
{"company":"dynatech","employees":[{"name":"first","children":[{"name":"child 1"},{"name":"offspring 1"}]},{"name":"second","children":[]}]}
```

## Example 3: Permutation of a Group of Attributes

Let's modify our dataset example by adding a new attribute, `childnumber`.

```json
{"company":"acme","employees":[{"name":"one","childnumber":2,"children":[{"name":"child 1"},{"name":"child 2"}]},{"name":"two","childnumber":3,"children":[{"name":"child 3"},{"name":"child 4"},{"name":"child 5"}]}]}
{"company":"megacorp","employees":[{"name":"alpha","childnumber":1,"children":[{"name":"kid 1"}]},{"name":"beta","childnumber":2,"children":[{"name":"kid 2"},{"name":"kid 3"}]}]}
{"company":"dynatech","employees":[{"name":"first","childnumber":2,"children":[{"name":"offspring 1"},{"name":"offspring 2"}]},{"name":"second","childnumber":0,"children":[]}]}
```

When there is a need to permute a group of attributes coherently, such as ensuring that `childnumber` and the `children` list are swapped together, the configuration file named `swap.yml` should have the following content:

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

The result will be the following:

```json
{"company":"acme","employees":[{"name":"one","childnumber":2,"children":[{"name":"kid 2"},{"name":"kid 3"}]},{"name":"two","childnumber":2,"children":[{"name":"child 1"},{"name":"child 2"}]}]}
{"company":"megacorp","employees":[{"name":"alpha","childnumber":0,"children":[]},{"name":"beta","childnumber":1,"children":[{"name":"kid 1"}]}]}
{"company":"dynatech","employees":[{"name":"first","childnumber":2,"children":[{"name":"offspring 1"},{"name":"offspring 2"}]},{"name":"second","childnumber":3,"children":[{"name":"child 3"},{"name":"child 4"},{"name":"child 5"}]}]}
```

## Example 4: Mix Multiple Groups and Attributes

Suppose the following input stream is stored in a file named `stream.jsonl`:

```json
{"company":"acme","employees":[{"name":"one","surname":"ONE","age":20,"nationality":"Kenyan"},{"name":"two","surname":"TWO","age":30,"nationality":"Icelandic"}]}
{"company":"megacorp","employees":[{"name":"alpha","surname":"ALPHA","age":40,"nationality":"Colombian"},{"name":"beta","surname":"BETA","age":50,"nationality":"Malaysian"}]}
{"company":"dynatech","employees":[{"name":"first","surname":"FIRST","age":60,"nationality":"Belgian"},{"name":"second","surname":"SECOND","age":70,"nationality":"Egyptian"}]}
```

The corresponding configuration file is named `configuration.yml`:

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

The permutation of the two groups will be performed independently, as for the single attribute `company`. The execution is done as follows:

```console
< stream.jsonl | type
```

And the result will be the following:

```json
{"company":"acme","employees":[{"name":"beta","surname":"BETA","age":50,"nationality":"Malaysian"},{"name":"one","surname":"ONE","age":30,"nationality":"Icelandic"}]}
{"company":"megacorp","employees":[{"name":"second","surname":"SECOND","age":70,"nationality":"Egyptian"},{"name":"alpha","surname":"ALPHA","age":20,"nationality":"Kenyan"}]}
{"company":"dynatech","employees":[{"name":"first","surname":"FIRST","age":60,"nationality":"Belgian"},{"name":"two","surname":"TWO","age":40,"nationality":"Colombian"}]}
```

Note that the `age` and `nationality` fields have been swapped consistently and independently of the `surname` and `name` fields, which have also been swapped consistently.

## Different Ways to Configure a Group

In previous examples, the following group configuration has been presented:

```yaml
selectors:
  - groupname:
    - employees.*.name
    - employees.*.surname
```

However, there are other ways to configure a group:

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
