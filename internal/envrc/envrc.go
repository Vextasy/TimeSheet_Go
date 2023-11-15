package envrc

// Simple .envrc file handling.
// Envrc reads a .envrc file and expects to find line of the form:
//
//	export {varname}={varvalue}
//
// With such lines it builds a key value map.
//
// A richer .env processing library can be found here: http://godoc.org/github.com/joho/godotenv

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type EnvRc map[string]string

// Open a ".envrc" file containing lines of the form:  export {varname}={varvalue}
// Look at the given path, and if that fails try again locally.
func NewEnvRc(path string) EnvRc {
	rc, err := mkEnvRc(path)
	if err != nil {
		return make(map[string]string)
	}
	return rc
}

// Open a ".envrc" file containing lines of the form:  export {varname}={varvalue}
// Look at the given path, and if that fails try again locally.
func MustEnvRc(path string) (EnvRc, error) {
	return mkEnvRc(path)
}

// Open a ".envrc" file containing lines of the form:  export {varname}={varvalue}
// Look at the given path, and if that fails try again locally.
func mkEnvRc(path string) (EnvRc, error) {
	const filename = ".envrc"
	const prefix = "export " // including trailing space
	const lenPrefix = len(prefix)
	// Try to read lines from the
	fullname := filepath.Join(path, filename)
	lines, err := readLines(fullname)
	if err != nil {
		if os.IsNotExist(err) {
			lines, err = readLines(filename)
			if err != nil {
				return nil, fmt.Errorf("unable to read env files: '%v' or '%v'", fullname, filename)
			}
		} else {
			return nil, fmt.Errorf("unable to open env file: '%v'", fullname)
		}
	}
	vars := make(map[string]string, 0)
	for _, line := range lines {
		if len(line) < lenPrefix || !strings.HasPrefix(line, prefix) {
			continue
		}
		line = line[lenPrefix:]
		key, val, found := strings.Cut(line, "=")
		if !found {
			continue
		}
		val = unquoteIfQuoted(val)
		vars[key] = val
	}
	return vars, nil
}

func (e EnvRc) Get(name string) string {
	return e[name]
}
func (e EnvRc) Try(name string, def string) string {
	if v, found := e[name]; found {
		return v
	} else {
		return def
	}
}

func readLines(path string) (lines []string, err error) {
	var (
		file   *os.File
		part   []byte
		prefix bool
	)
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))
	for {
		if part, prefix, err = reader.ReadLine(); err != nil { // err must be an unwrapped error to be compared with EOF below.
			break
		}
		buffer.Write(part)
		if !prefix {
			lines = append(lines, buffer.String())
			buffer.Reset()
		}
	}
	if err == io.EOF {
		err = nil
	}
	return lines, err
}

func unquoteIfQuoted(value string) string {
	var bytes []byte
	bytes = []byte(value)
	if (len(bytes) > 2 && bytes[0] == '"' && bytes[len(bytes)-1] == '"') ||
		(len(bytes) > 2 && bytes[0] == '\'' && bytes[len(bytes)-1] == '\'') {
		bytes = bytes[1 : len(bytes)-1]
	}
	return string(bytes)
}
