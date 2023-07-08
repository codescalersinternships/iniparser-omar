package iniparser

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

var errInvalidFormat = errors.New("INVALID INI FILE FORMAT")

type INIParser struct {
	sectionNames []string
	data         map[string]*section
}

func (ini *INIParser) init() {
	ini.data = map[string]*section{}
}

func (ini *INIParser) loadData(lines []string) error {
	// remove empty lines at the beginning of the file
	for i := range lines {
		if lines[i] != "" {
			lines = lines[i:]
			break
		}
	}

	ini.init()
	for _, line := range lines {
		if len(line) == 0 || line[0] != '[' {
			if len(ini.sectionNames) == 0 {
				return errInvalidFormat
			}
			currentSection := ini.sectionNames[len(ini.sectionNames)-1]
			err := ini.data[currentSection].readLine(line)
			if err != nil {
				return err
			}
		} else {
			line = strings.Trim(line, " ")
			if len(line) <= 2 || line[len(line)-1] != ']' {
				return errInvalidFormat
			}
			sectionName := line[1 : len(line)-1]
			sectionName = strings.Trim(sectionName, " ")
			if len(sectionName) == 0 {
				return errInvalidFormat
			}
			ini.sectionNames = append(ini.sectionNames, sectionName)
			ini.data[sectionName] = &section{}
			ini.data[sectionName].init()
		}
	}
	return nil
}

func (ini *INIParser) LoadFromString(str string) error {
	lines := strings.Split(str, "\n")
	err := ini.loadData(lines)
	return err
}

func (ini *INIParser) LoadFromFile(filePath string) error {
	readFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	var lines []string
	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}

	err = ini.loadData(lines)
	return err
}

func (ini *INIParser) GetSectionNames() []string {
	return ini.sectionNames
}

func (ini *INIParser) GetSections() map[string]map[string]string {
	serializedINI := map[string]map[string]string{}
	for _, name := range ini.sectionNames {
		serializedINI[name] = ini.data[name].getSection()
	}
	return serializedINI
}

func (ini *INIParser) Get(sectionName, key string) string {
	if ini.data[sectionName] == nil {
		return ""
	}
	return ini.data[sectionName].get(key)
}

func (ini *INIParser) Set(sectionName, key, value string) error {
	sectionName = strings.Trim(sectionName, " ")
	if len(sectionName) == 0 {
		return errInvalidFormat
	}
	if ini.data[sectionName] == nil {
		// add new section, key and value
		ini.sectionNames = append(ini.sectionNames, sectionName)
		ini.data[sectionName] = &section{}
		ini.data[sectionName].init()
	}
	err := ini.data[sectionName].set(key, value)
	return err
}

func (ini *INIParser) ToString() string {
	var str string
	for _, name := range ini.sectionNames {
		str += "[" + name + "]\n"
		str += strings.Join(ini.data[name].getSectionINI(), "\n") + "\n"
	}
	return str
}

func (ini *INIParser) SaveToFile(filePath string) error {
	err := os.WriteFile(filePath, []byte(ini.ToString()), 0644)
	return err
}
