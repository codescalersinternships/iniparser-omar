package iniparser

import (
	"errors"
	"io"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func validINIDataSample1() map[string]section {
	return map[string]section{
		" section 1": map[string]string{"key key ": " value value"},
		"section 2":  map[string]string{},
		"section 3":  map[string]string{"key": ""},
	}
}

var validINIStringSample1 = `
;comment

[ section 1]
key key = value value
;comment
[section 2]

[section 3]
key= 
`

func TestLoadData(t *testing.T) {
	t.Run("valid format", func(t *testing.T) {
		ini := NewINIParser()
		expectedData := validINIDataSample1()

		err := ini.loadData(strings.NewReader(validINIStringSample1))
		assertErr(t, err, nil)
		if !reflect.DeepEqual(ini.data, expectedData) {
			t.Errorf("got %v want %v", ini.data, expectedData)
		}
	})

	testCases := []struct {
		name  string
		input string
		err   error
	}{
		{
			name:  "invalid format: begin with global data",
			input: "key key = value value",
			err:   ErrNoGlobalDataAllowed,
		}, {
			name: "invalid format: key value format",
			input: `
			[section 1]
			key key  value value
			`,
			err: ErrInvalidFormat,
		}, {
			name: "invalid format: empty section name",
			input: `
			[]
			key key = value value
			`,
			err: ErrSectionNameCantBeEmpty,
		}, {
			name: "invalid format: empty key",
			input: `
			[section name]
			   = value value
			`,
			err: ErrKeyCantBeEmpty,
		}, {
			name: "invalid format: section repeated",
			input: `
			[section name]
			[section name]
			`,
			err: ErrSectionNameMustBeUnique,
		}, {
			name: "invalid format: key repeated",
			input: `
			[section name]
			key = value
			key = value2
			`,
			err: ErrKeyMustBeUnique,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ini := NewINIParser()
			err := ini.loadData(strings.NewReader(tc.input))
			assertErr(t, err, tc.err)
			if !errors.Is(err, tc.err) {
				t.Errorf("got %q want %q", err.Error(), tc.err.Error())
			}
		})
	}
}

func TestLoadFromString(t *testing.T) {
	t.Run("valid format", func(t *testing.T) {
		ini := NewINIParser()
		expectedData := validINIDataSample1()

		err := ini.LoadFromString(validINIStringSample1)
		assertErr(t, err, nil)
		if !reflect.DeepEqual(ini.data, expectedData) {
			t.Errorf("got %v want %v", ini.data, expectedData)
		}
	})
}

func TestLoadFromFile(t *testing.T) {
	t.Run("file not exist", func(t *testing.T) {
		ini := NewINIParser()
		err := ini.LoadFromFile("notExist.ini")

		if !os.IsNotExist(err) {
			t.Errorf("expected to get not found error")
		}
	})

	t.Run("file extension is not valid", func(t *testing.T) {
		ini := NewINIParser()
		err := ini.LoadFromFile("invalid.in")
		assertErr(t, err, ErrInvalidFileExtension)
	})

	t.Run("valid format", func(t *testing.T) {
		f, err := os.CreateTemp("", "test_file*.ini")
		assertErr(t, err, nil)
		defer os.Remove(f.Name())

		_, err = f.WriteString(validINIStringSample1)
		assertErr(t, err, nil)
		_, err = f.Seek(0, io.SeekStart)
		assertErr(t, err, nil)

		ini := NewINIParser()
		expectedData := validINIDataSample1()

		err = ini.LoadFromFile(f.Name())
		assertErr(t, err, nil)

		if !reflect.DeepEqual(ini.data, expectedData) {
			t.Errorf("got %v want %v", ini.data, expectedData)
		}
	})
}

func TestGetSectionNames(t *testing.T) {
	t.Run("get data", func(t *testing.T) {
		dataSample := validINIDataSample1()
		ini := INIParser{data: dataSample}

		got := ini.GetSectionNames()
		sort.Strings(got)
		want := []string{" section 1", "section 2", "section 3"}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestGetSections(t *testing.T) {
	t.Run("get sections", func(t *testing.T) {
		dataSample := validINIDataSample1()
		ini := INIParser{data: dataSample}

		got := ini.GetSections()
		want := map[string]map[string]string{
			" section 1": {"key key ": " value value"},
			"section 2":  {},
			"section 3":  {"key": ""},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestGet(t *testing.T) {
	testCases := []struct {
		name          string
		section       string
		key           string
		expectedValue string
		isOk          bool
	}{
		{
			name:          "get exist value",
			section:       " section 1",
			key:           "key key ",
			expectedValue: " value value",
			isOk:          true,
		}, {
			name:          "not found section",
			section:       "not found",
			key:           "not found",
			expectedValue: "",
			isOk:          false,
		}, {
			name:          "not found key",
			section:       "section 1",
			key:           "not found",
			expectedValue: "",
			isOk:          false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dataSample := validINIDataSample1()
			ini := INIParser{data: dataSample}

			gotValue, gotOk := ini.Get(tc.section, tc.key)

			if gotValue != tc.expectedValue {
				t.Errorf("value: got %q want %q", gotValue, tc.expectedValue)
			}
			if gotOk != tc.isOk {
				t.Errorf("is found: got is %t want %t", gotOk, tc.isOk)
			}
		})
	}
}

func TestSet(t *testing.T) {
	testCases := []struct {
		name     string
		data     map[string]section
		section  string
		key      string
		value    string
		expected map[string]section
		err      error
	}{
		{
			name:    "set new section and key",
			data:    map[string]section{},
			section: "new section",
			key:     "new key",
			value:   "new value",
			expected: map[string]section{
				"new section": map[string]string{
					"new key": "new value",
				},
			},
		}, {
			name:    "set new key",
			data:    map[string]section{"new section": map[string]string{}},
			section: "new section",
			key:     "new key",
			value:   "new value",
			expected: map[string]section{
				"new section": map[string]string{
					"new key": "new value",
				},
			},
		}, {
			name: "update value",
			data: map[string]section{"new section": map[string]string{
				"new key": "new value",
			}},
			section: "new section",
			key:     "new key",
			value:   "new value",
			expected: map[string]section{
				"new section": map[string]string{
					"new key": "new value",
				},
			},
		}, {
			name:    "set empty section name",
			section: "",
			key:     "key",
			value:   "value",
			err:     ErrSectionNameCantBeEmpty,
		},
		{
			name:    "set empty key",
			section: "section name",
			key:     "",
			value:   "value",
			err:     ErrKeyCantBeEmpty,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ini := INIParser{data: tc.data}
			err := ini.Set(tc.section, tc.key, tc.value)
			assertErr(t, err, tc.err)

			if tc.err != nil {
				if !reflect.DeepEqual(ini.data, tc.expected) {
					t.Errorf("got %v want %v", ini.data, tc.expected)
				}
			}
		})
	}
}

func TestString(t *testing.T) {
	t.Run("get string", func(t *testing.T) {
		dataSample := validINIDataSample1()
		ini := INIParser{data: dataSample}
		got := ini.String()

		ini2 := NewINIParser()
		err := ini2.LoadFromString(got)
		assertErr(t, err, nil)

		if !reflect.DeepEqual(ini2.data, dataSample) {
			t.Errorf("got %v want %v", ini2.data, dataSample)
		}
	})
}

func TestSaveToFile(t *testing.T) {
	t.Run("save to file", func(t *testing.T) {
		f, err := os.CreateTemp("", "test_file*.ini")
		assertErr(t, err, nil)
		defer os.Remove(f.Name())

		dataSample := validINIDataSample1()
		ini := INIParser{data: dataSample}

		err = ini.SaveToFile(f.Name())
		assertErr(t, err, nil)

		ini2 := NewINIParser()
		err = ini2.LoadFromFile(f.Name())
		assertErr(t, err, nil)

		if !reflect.DeepEqual(ini2.data, dataSample) {
			t.Errorf("got %v want %v", ini2.data, dataSample)
		}
	})

	t.Run("save to file has not ini extension", func(t *testing.T) {
		f, err := os.CreateTemp("", "test_file")
		assertErr(t, err, nil)
		defer os.Remove(f.Name())

		dataSample := validINIDataSample1()
		ini := INIParser{data: dataSample}

		err = ini.SaveToFile(f.Name())
		assertErr(t, err, ErrInvalidFileExtension)
	})
}

func assertErr(t testing.TB, got, want error) {
	t.Helper()
	if !errors.Is(got, want) {
		if got == nil {
			t.Errorf("got nil want %q", want.Error())
		} else if want == nil {
			t.Errorf("got %q want nil", got.Error())
		} else {
			t.Errorf("got %q want %q", got.Error(), want.Error())
		}
	}
}
