# cutl

`cutl` is a command-line utility for translating from one data format to another.

## Installation

```shell
$ go install github.com/meanguy/cutl/cmd/cutl@latest
```

## Usage

```shell
$ cutl -h                                        
Usage of cutl:                 
  -f string
        format of the source input (default "json")
  -file string
        input file for translation - defaults to stdin if not specified (default "tests/deeply-nested.json")
  -t string
        format to translate the source input to (default "yaml")
```

## Examples

```shell
$ cat tests/object.json
{
  "hello": "world",
  "abc": [1, 2, 3]
}

# translate json to yaml
$ cutl -file tests/object.json
abc:                                                       
    - 1
    - 2
    - 3
hello: world

# translate toml from yaml
$ cutl -file tests/object.json -t toml
abc = [1.0, 2.0, 3.0]
hello = "world"

# read from stdin
$ cutl -f toml -t json
abc = [1.0, 2.0, 3.0]
{"abc":[1,2,3]}

# read from pipe
$ cat tests/deeply-nested.json | cutl 
how:
    level: 2
    low:
        can:
            level: 4
            you:
                go:
                    - I
                    - wonder
```
