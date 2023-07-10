package iniparser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Generic file system errors.
// Errors returned by file systems can be tested against these errors
// using errors.Is.
var (
	ErrInvalidFormat        = errors.New("invalid ini file format")
	ErrNoGlobalDataAllowed  = errors.New("global data is not supported")
	ErrInvalidFileExtension = errors.New("ErrInvalidFileExtension")
)

// An INIParser provides parser to ini files.
type INIParser struct {
	data map[string]*section
}

func (ini *INIParser) loadData(reader io.Reader) error {
	fileScanner := bufio.NewScanner(reader)
	fileScanner.Split(bufio.ScanLines)
	var lines []string
	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}

	// remove empty lines at the beginning of the file
	for i := range lines {
		if lines[i] != "" {
			lines = lines[i:]
			break
		}
	}

	ini.data = map[string]*section{}
	var lastSectionName string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] != '[' { // not section definition line
			if len(ini.data) == 0 {
				return ErrNoGlobalDataAllowed
			}
			err := ini.data[lastSectionName].addLine(line)
			if err != nil {
				return err
			}
		} else {
			if len(line) < 2 || line[len(line)-1] != ']' {
				return fmt.Errorf(errInvalidLine, ErrInvalidFormat, line)
			}
			sectionName := line[1 : len(line)-1]
			sectionName = strings.TrimSpace(sectionName)
			ini.data[sectionName] = &section{}
			lastSectionName = sectionName
		}
	}
	return nil
}

// LoadFromString loads the given string to the parser.
func (ini *INIParser) LoadFromString(str string) error {
	reader := strings.NewReader(str)
	return ini.loadData(reader)
}

// LoadFromFile loads the given string to the parser.
func (ini *INIParser) LoadFromFile(filePath string) error {
	_, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	if fileExt := filepath.Ext(filePath); fileExt != ".ini" {
		return fmt.Errorf("%w: %s", ErrInvalidFileExtension, fileExt)
	}

	reader, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer reader.Close()
	return ini.loadData(reader)
}

// GetSectionNames gets slice of section names of loaded data.
func (ini *INIParser) GetSectionNames() []string {
	if ini.data == nil {
		ini.data = map[string]*section{}
	}

	keys := []string{}
	for k := range ini.data {
		keys = append(keys, k)
	}
	return keys
}

// GetSections gets ini data parsed as 'map[string]map[string]string'.
func (ini *INIParser) GetSections() map[string]map[string]string {
	if ini.data == nil {
		ini.data = map[string]*section{}
	}

	serializedINI := map[string]map[string]string{}
	for k, v := range ini.data {
		serializedINI[k] = v.getSection()
	}
	return serializedINI
}

// Get gets value given section name and key.
// returns tuple of (value, true) incase value found.
// returns tuple of (value, false) incase value is not found.
func (ini *INIParser) Get(sectionName, key string) (string, bool) {
	if ini.data == nil {
		ini.data = map[string]*section{}
	}

	if ini.data[sectionName] == nil {
		return "", false
	}
	return ini.data[sectionName].get(key)
}

// Set sets or updates given section name and key.
// creates new section and new key if not exist.
func (ini *INIParser) Set(sectionName, key, value string) {
	if ini.data == nil {
		ini.data = map[string]*section{}
	}

	sectionName = strings.TrimSpace(sectionName)
	if ini.data[sectionName] == nil {
		ini.data[sectionName] = &section{}
	}
	ini.data[sectionName].set(key, value)
}

// String converts loaded data to string.
func (ini *INIParser) String() string {
	if ini.data == nil {
		ini.data = map[string]*section{}
	}

	var str string
	for k, v := range ini.data {
		str += "[" + k + "]\n"
		str += v.string()
	}
	return str
}

// saveToFile saves loaded data to file.
func (ini *INIParser) SaveToFile(filePath string) error {
	if fileExt := filepath.Ext(filePath); fileExt != ".ini" {
		return fmt.Errorf("%w: %s", ErrInvalidFileExtension, fileExt)
	}

	if ini.data == nil {
		ini.data = map[string]*section{}
	}

	return os.WriteFile(filePath, []byte(ini.String()), 0644)
}
