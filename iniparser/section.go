package iniparser

import (
	"errors"
	"fmt"
	"strings"
)

type Section struct {
	lines []string // key or comment or empty
	data  map[string]string
}

func (sec *Section) Init() {
	sec.data = map[string]string{}
}

func (sec *Section) ReadLine(line string) error {
	line = strings.Trim(line, " ")
	if len(line) == 0 || line[0] == ';' {
		sec.lines = append(sec.lines, line)
		return nil
	}

	if strings.Contains(line, "=") {
		splitRet := strings.Split(line, "=")
		key := strings.Trim(splitRet[0], " ")
		value := strings.Trim(splitRet[1], " ")

		sec.lines = append(sec.lines, key)
		sec.data[key] = value
		return nil
	}
	return errors.New("INVALID INI FILE FORMAT")
}

func (sec *Section) Get(key string) string {
	return sec.data[key]
}

func (sec *Section) Set(key, value string) {
	key = strings.Trim(key, " ")
	value = strings.Trim(value, " ")

	if _, ok := sec.data[key]; !ok {
		sec.lines = append(sec.lines, key)
	}
	sec.data[key] = value
}

func (sec *Section) GetSection() map[string]string {
	return sec.data
}

func (sec *Section) GetSectionINI() []string {
	var iniLines []string
	for _, line := range sec.lines {
		if len(line) == 0 || line[0] == ';' {
			// line is comment or empty
			iniLines = append(iniLines, line)
		} else {
			// line is key
			iniLines = append(iniLines, fmt.Sprintf("%s = %s", line, sec.data[line]))
		}
	}

	return iniLines
}

