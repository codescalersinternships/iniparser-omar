package iniparser

import (
	"errors"
	"fmt"
	"strings"
)

type section struct {
	lines []string // key or comment or empty
	data  map[string]string
}

func (sec *section) init() {
	sec.data = map[string]string{}
}

func (sec *section) readLine(line string) error {
	line = strings.Trim(line, " ")
	if len(line) == 0 || line[0] == ';' {
		sec.lines = append(sec.lines, line)
		return nil
	}

	if strings.Contains(line, "=") {
		splitRet := strings.Split(line, "=")
		key := strings.Trim(splitRet[0], " ")
		value := strings.Trim(splitRet[1], " ")
		if len(key) == 0 || len(value) == 0 {
			return errors.New("INVALID INI FILE FORMAT: KEY OR VALUE CAN'T BE EMPTY")
		}

		sec.lines = append(sec.lines, key)
		sec.data[key] = value
		return nil
	}
	return errors.New("INVALID INI FILE FORMAT")
}

func (sec *section) get(key string) string {
	return sec.data[key]
}

func (sec *section) set(key, value string) error {
	key = strings.Trim(key, " ")
	value = strings.Trim(value, " ")
	if len(key) == 0 || len(value) == 0 {
		return errors.New("INVALID INI FILE FORMAT: KEY OR VALUE CAN'T BE EMPTY")
	}

	if _, ok := sec.data[key]; !ok {
		sec.lines = append(sec.lines, key)
	}
	sec.data[key] = value

	return nil
}

func (sec *section) getSection() map[string]string {
	return sec.data
}

func (sec *section) getSectionINI() []string {
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
