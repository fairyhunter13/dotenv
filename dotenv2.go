package dotenv

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fairyhunter13/envcompact/pkg/compacter"
)

// Option is the option struct.
type Option struct {
	Overload bool
	Paths    []string
}

// Default specifies the default values of the option struct.
func (o *Option) Default() {
	if len(o.Paths) <= 0 {
		o.Paths = append(o.Paths, ".env")
	}
}

// FnOption is a functional option for this package
type FnOption func(*Option)

// WithPaths specifies all paths of the environment file.
func WithPaths(paths ...string) FnOption {
	return func(o *Option) {
		o.Paths = paths
	}
}

// WithOverload overdrives the current set environment.
func WithOverload(overload bool) FnOption {
	return func(o *Option) {
		o.Overload = overload
	}
}

// Load2 is the version 2 of Load.
func Load2(opts ...FnOption) (err error) {
	optStruct := new(Option)
	for _, opt := range opts {
		opt(optStruct)
	}
	optStruct.Default()

	var isSuccess bool
	for _, path := range optStruct.Paths {
		err = loadFile2(path, optStruct.Overload)
		if err != nil {
			continue
		}

		isSuccess = true
	}

	if isSuccess {
		err = nil
	}
	return
}

// loadFile2 is the version 2 of loadFile.
func loadFile2(path string, overload bool) error {
	env, err := ReadFile2(path)
	if err != nil {
		return err
	}
	LoadMap(env, overload)
	return nil
}

// ReadFile2 is the version 2 of ReadFile.
func ReadFile2(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Read2(f)
}

// LoadReader2 is the version 2 of LoadReader.
func LoadReader2(r io.Reader) error {
	env, err := Read2(r)
	if err != nil {
		return err
	}
	LoadMap(env, false)
	return nil
}

// Read2 is the version 2 of Read.
func Read2(rd io.Reader) (map[string]string, error) {
	engine := compacter.NewParser(rd)
	envMap := make(map[string]string)
	var (
		line, k, v string
		err        error
	)
	for engine.Scan() {
		line = engine.Text()
		if regexVar.MatchString(line) {
			line = regexVar.ReplaceAllStringFunc(line, func(s string) string {
				return envMap[strings.Trim(s, "${}")]
			})
		}
		k, v, err = ParseString(line)
		if err == ErrInvalidln {
			return nil, fmt.Errorf("could not parse file: %v", err)
		}
		if err != nil {
			continue
		}
		envMap[k] = v
	}

	if err = engine.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return envMap, nil
}
