package iniparser

import (
	"fmt"
	"strings"
)

var errInvalidLine = "%w: at line %q"

type section struct {
	data map[string]string
}

func (sec *section) addLine(line string) error {
	if sec.data == nil {
		sec.data = map[string]string{}
	}

	line = strings.TrimSpace(line)
	if len(line) == 0 || line[0] == ';' {
		return nil
	}

	if strings.Contains(line, "=") {
		splitRet := strings.Split(line, "=")
		key := strings.TrimSpace(splitRet[0])
		value := strings.TrimSpace(splitRet[1])

		sec.data[key] = value
		return nil
	}
	return fmt.Errorf(errInvalidLine, ErrInvalidFormat, line)
}

func (sec *section) get(key string) (string, bool) {
	if sec.data == nil {
		sec.data = map[string]string{}
	}

	val, exist := sec.data[key]
	return val, exist
}

func (sec *section) set(key, value string) {
	if sec.data == nil {
		sec.data = map[string]string{}
	}

	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)
	sec.data[key] = value
}

func (sec *section) getSection() map[string]string {
	return sec.data
}

func (sec *section) string() string {
	var str string
	for k, v := range sec.data {
		str += k + " = " + v + "\n"
	}
	return str
}
