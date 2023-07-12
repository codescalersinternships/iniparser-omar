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
var (
	ErrInvalidFormat           = errors.New("invalid ini file format")
	ErrNoGlobalDataAllowed     = errors.New("global data is not supported")
	ErrInvalidFileExtension    = errors.New("file extension should be .ini")
	ErrKeyCantBeEmpty          = errors.New("key can't be empty")
	ErrKeyMustBeUnique         = errors.New("key must be unique")
	ErrSectionNameCantBeEmpty  = errors.New("section name can't be empty")
	ErrSectionNameMustBeUnique = errors.New("section name must be unique")
)

type section map[string]string

// INIParser provides parser to ini files.
type INIParser struct {
	data map[string]section
}

func NewINIParser() INIParser {
	return INIParser{
		data: map[string]section{},
	}
}

func (ini *INIParser) loadData(reader io.Reader) error {
	ini.data = map[string]section{}

	fileScanner := bufio.NewScanner(reader)
	fileScanner.Split(bufio.ScanLines)

	currentSectionName := ""
	for fileScanner.Scan() {
		line := strings.TrimSpace(fileScanner.Text())

		// empty line or comment
		if len(line) == 0 || line[0] == ';' || line[0] == '#' {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			sectionName := line[1 : len(line)-1]

			if sectionName == "" {
				return ErrSectionNameCantBeEmpty
			}
			if ini.data[sectionName] != nil {
				return fmt.Errorf("%w, section = %q", ErrSectionNameMustBeUnique, sectionName)
			}

			ini.data[sectionName] = section{}
			currentSectionName = sectionName

			continue
		}

		if strings.Contains(line, "=") {
			if len(ini.data) == 0 { // if no sections add yet
				return ErrNoGlobalDataAllowed
			}
			// now i'm sure that 'currentSectionName' has a value

			splitRet := strings.SplitN(line, "=", 2)
			key := splitRet[0]
			value := splitRet[1]

			if key == "" {
				return fmt.Errorf("%w: at line %q", ErrKeyCantBeEmpty, line)
			}
			_, ok := ini.data[currentSectionName][key]
			if ok {
				return fmt.Errorf("%w, key = %q", ErrKeyMustBeUnique, key)
			}

			ini.data[currentSectionName][key] = value

			continue
		}

		return fmt.Errorf("%w: at line %q", ErrInvalidFormat, line)
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
	if fileExt := filepath.Ext(filePath); fileExt != ".ini" {
		return fmt.Errorf("%w: %s", ErrInvalidFileExtension, fileExt)
	}

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	return ini.loadData(f)
}

// GetSectionNames gets slice of section names of loaded data.
func (ini *INIParser) GetSectionNames() []string {
	keys := []string{}
	for k := range ini.data {
		keys = append(keys, k)
	}

	return keys
}

// GetSections gets ini data parsed as 'map[string]map[string]string'.
func (ini *INIParser) GetSections() map[string]map[string]string {
	serializedINI := map[string]map[string]string{}

	for sectionName, section := range ini.data {
		serializedINI[sectionName] = map[string]string{}
		for k, v := range section {
			serializedINI[sectionName][k] = v
		}
	}

	return serializedINI
}

// Get gets value given section name and key.
func (ini *INIParser) Get(sectionName, key string) (string, bool) {
	if ini.data[sectionName] == nil {
		return "", false
	}

	value, ok := ini.data[sectionName][key]
	return value, ok
}

// Set sets or updates given section name and key.
// creates new section and new key if not exist.
func (ini *INIParser) Set(sectionName, key, value string) error {
	if sectionName == "" {
		return ErrSectionNameCantBeEmpty
	}
	if key == "" {
		return ErrKeyCantBeEmpty
	}

	if ini.data[sectionName] == nil {
		ini.data[sectionName] = section{}
	}
	ini.data[sectionName][key] = value

	return nil
}

// String converts loaded data to string.
func (ini *INIParser) String() string {
	str := ""
	for sectionName, section := range ini.data {
		str += fmt.Sprintf("[%v]\n", sectionName)
		for k, v := range section {
			str += fmt.Sprintf("%v=%v\n", k, v)
		}
	}

	return str
}

// saveToFile saves loaded data to file.
func (ini *INIParser) SaveToFile(filePath string) error {
	if fileExt := filepath.Ext(filePath); fileExt != ".ini" {
		return fmt.Errorf("%w: %s", ErrInvalidFileExtension, fileExt)
	}

	return os.WriteFile(filePath, []byte(ini.String()), 0644)
}
