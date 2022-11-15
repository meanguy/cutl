package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v3"
)

type (
	Args struct {
		File       string
		FromFormat string
		ToFormat   string
	}

	FromHandler func(io.Reader, interface{}) error
	ToHandler   func(io.Writer, interface{}) error
)

func ParseArgs() Args {
	args := Args{}

	flag.StringVar(
		&args.File,
		"file",
		"tests/deeply-nested.json",
		"input file for translation - defaults to stdin if not specified",
	)
	flag.StringVar(&args.FromFormat, "f", "json", "format of the source input")
	flag.StringVar(&args.FromFormat, "from", "json", "format of the source input")
	flag.StringVar(&args.ToFormat, "t", "yaml", "format to translate the source input to")
	flag.StringVar(&args.ToFormat, "to", "yaml", "format to translate the source input to")
	flag.Parse()

	return args
}

func ReaderFromFile(path string) (io.Reader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %v", path, err)
	}

	buf, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %q: %v", path, err)
	}

	return bytes.NewBuffer(buf), nil
}

func JSONFromHandler(r io.Reader, out interface{}) error {
	return json.NewDecoder(r).Decode(out)
}

func JSONToHandler(w io.Writer, val interface{}) error {
	return json.NewEncoder(w).Encode(val)
}

func YAMLFromHandler(r io.Reader, out interface{}) error {
	return yaml.NewDecoder(r).Decode(out)
}

func YAMLToHandler(w io.Writer, val interface{}) error {
	return yaml.NewEncoder(w).Encode(val)
}

func TOMLFromHandler(r io.Reader, out interface{}) error {
	return toml.NewDecoder(r).Decode(out)
}

func TOMLToHandler(w io.Writer, val interface{}) error {
	return toml.NewEncoder(w).Encode(val)
}
func main() {
	args := ParseArgs()

	fromFormat := strings.ToLower(args.FromFormat)
	toFormat := strings.ToLower(args.ToFormat)

	var input io.Reader = os.Stdin
	if args.File != "" {
		r, err := ReaderFromFile(args.File)
		if err != nil {
			log.Fatal(err.Error())
		}

		input = r
	}

	var fromHandler FromHandler
	switch fromFormat {
	case "json":
		fromHandler = JSONFromHandler
	case "tml":
		fallthrough
	case "toml":
		fromHandler = TOMLFromHandler
	case "yml":
		fallthrough
	case "yaml":
		fromHandler = YAMLFromHandler
	default:
		log.Fatalf("cannot translate from %q format", fromFormat)
	}

	var toHandler ToHandler
	switch toFormat {
	case "json":
		toHandler = JSONToHandler
	case "tml":
		fallthrough
	case "toml":
		toHandler = TOMLToHandler
	case "yml":
		fallthrough
	case "yaml":
		toHandler = YAMLToHandler
	default:
		log.Fatalf("cannot translate to %q format", toFormat)
	}

	var value interface{}
	if err := fromHandler(input, &value); err != nil {
		log.Fatalf("failed to translate from %q format: %v", fromFormat, err)
	}

	if err := toHandler(os.Stdout, value); err != nil {
		log.Fatalf("failed to translate to %q format: %v", toFormat, err)
	}
}
