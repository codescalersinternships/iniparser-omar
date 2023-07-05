package iniparser

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type INIParser struct {
	sectionNames []string
	data         map[string]*Section
}

func (ini *INIParser) Init() {
	ini.data = map[string]*Section{}
}

func (ini *INIParser) loadData(lines []string) error {
	// remove empty lines at the beginning of the file
	for i := range lines {
		if lines[i] != "" {
			lines = lines[i:]
			break
		}
	}

	for _, line := range lines {
		if len(line) == 0 || line[0] != '[' {
			if len(ini.sectionNames) == 0 {
				return errors.New("INVALID INI FILE FORMAT: MUST BEGIN WITH SECTION HEADER")
			}
			currentSection := ini.sectionNames[len(ini.sectionNames)-1]
			err := ini.data[currentSection].ReadLine(line)
			if err != nil {
				return err
			}
		} else {
			line = strings.Trim(line, " ")
			if len(line) <= 2 || line[len(line)-1] != ']' {
				return errors.New("INVALID INI FILE FORMAT")
			}
			sectionName := line[1 : len(line)-1]
			sectionName = strings.Trim(sectionName, " ")
			if len(sectionName) == 0 {
				return errors.New("INVALID INI FILE FORMAT: SECTION HEADER CAN'T BE EMPTY")
			}
			ini.sectionNames = append(ini.sectionNames, sectionName)
			ini.data[sectionName] = &Section{}
			ini.data[sectionName].Init()
		}
	}
	return nil
}

func (ini *INIParser) LoadFromString(str string) {
	lines := strings.Split(str, "\n")
	err := ini.loadData(lines)
	check(err)
}

func (ini *INIParser) LoadFromFile(filePath string) {
	readFile, err := os.Open(filePath)
	check(err)
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	var lines []string
	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}

	err = ini.loadData(lines)
	check(err)
}

func (ini *INIParser) GetSectionNames() []string {
	return ini.sectionNames
}

func (ini *INIParser) GetSections() map[string]map[string]string {
	serializedINI := map[string]map[string]string{}
	for _, name := range ini.sectionNames {
		serializedINI[name] = ini.data[name].GetSection()
	}
	return serializedINI
}

func (ini *INIParser) Get(sectionName, key string) string {
	if ini.data[sectionName] == nil {
		return ""
	}
	return ini.data[sectionName].Get(key)
}

func (ini *INIParser) Set(sectionName, key, value string) {
	sectionName = strings.Trim(sectionName, " ")
	if len(sectionName) == 0 {
		panic("INVALID INI FILE FORMAT: SECTION HEADER CAN'T BE EMPTY")
	}
	if ini.data[sectionName] == nil {
		ini.data[ini.sectionNames[len(ini.sectionNames)-1]].ReadLine("") // add new line
		// add new section, key and value
		ini.sectionNames = append(ini.sectionNames, sectionName)
		ini.data[sectionName] = &Section{}
		ini.data[sectionName].Init()
	}
	err := ini.data[sectionName].Set(key, value)
	check(err)
}

func (ini *INIParser) ToString() string {
	var str string
	for _, name := range ini.sectionNames {
		str += "[" + name + "]\n"
		str += strings.Join(ini.data[name].GetSectionINI(), "\n") + "\n"
	}
	return str
}

func (ini *INIParser) SaveToFile(filePath string) {
	err := os.WriteFile(filePath, []byte(ini.ToString()), 0644)
	check(err)
}
