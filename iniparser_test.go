package iniparser

import (
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func validINIDataSample1() map[string]*section {
	return map[string]*section{
		"section 1": {data: map[string]string{"key key": "value value"}},
		"section 2": {data: map[string]string{}},
		"":          {data: map[string]string{"": ""}},
	}
}

var validINIStringSample1 = `
[section 1]
key key = value value
;comment
[section 2]

[]
 = 
`

func TestLoadData(t *testing.T) {
	t.Run("valid format", func(t *testing.T) {
		ini := INIParser{}
		expectedData := validINIDataSample1()

		err := ini.loadData(strings.NewReader(validINIStringSample1))
		assertNoErr(t, err)
		assertEqualData(t, ini.data, expectedData)
	})

	testCases := []struct {
		name  string
		input string
		err   error
	}{
		{
			name:  "invalid format: begin with comment",
			input: ";comment",
			err:   ErrNoGlobalDataAllowed,
		}, {
			name:  "invalid format: begin with global data",
			input: "key key = value value",
			err:   ErrNoGlobalDataAllowed,
		}, {
			name: "invalid format: section naming format",
			input: `
			[section 1
			key key = value value
			`,
			err: ErrInvalidFormat,
		}, {
			name: "invalid format: key value format",
			input: `
			[section 1]
			key key  value value
			`,
			err: ErrInvalidFormat,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ini := INIParser{}
			err := ini.loadData(strings.NewReader(tc.input))
			assertErrFound(t, err)
			if !errors.Is(err, tc.err) {
				t.Errorf("got %q want %q", err.Error(), tc.err.Error())
			}
		})
	}
}

func TestLoadFromString(t *testing.T) {
	t.Run("valid format", func(t *testing.T) {
		ini := INIParser{}
		expectedData := validINIDataSample1()

		err := ini.LoadFromString(validINIStringSample1)
		assertNoErr(t, err)
		assertEqualData(t, ini.data, expectedData)
	})

	testCases := []struct {
		name  string
		input string
		err   error
	}{
		{
			name:  "invalid format: begin with comment",
			input: ";comment",
			err:   ErrNoGlobalDataAllowed,
		}, {
			name:  "invalid format: begin with global data",
			input: "key key = value value",
			err:   ErrNoGlobalDataAllowed,
		}, {
			name: "invalid format: section naming format",
			input: `
			[section 1
			key key = value value
			`,
			err: ErrInvalidFormat,
		}, {
			name: "invalid format: key value format",
			input: `
			[section 1]
			key key  value value
			`,
			err: ErrInvalidFormat,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ini := INIParser{}
			err := ini.LoadFromString(tc.input)
			assertErrFound(t, err)
			if !errors.Is(err, tc.err) {
				t.Errorf("got %q want %q", err.Error(), tc.err.Error())
			}
		})
	}
}

func TestLoadFromFile(t *testing.T) {
	t.Run("file not exist", func(t *testing.T) {
		ini := INIParser{}
		err := ini.LoadFromFile("notExist.ini")
		assertErrFound(t, err)
		if !os.IsNotExist(err) {
			t.Errorf("expected to get not found error")
		}
	})

	t.Run("file extension is not valid", func(t *testing.T) {
		ini := INIParser{}
		err := ini.LoadFromFile("go.mod")
		assertErrFound(t, err)
		if !errors.Is(err, ErrInvalidFileExtension) {
			t.Errorf("got %q want %q", err.Error(), ErrInvalidFileExtension.Error())
		}
	})

	t.Run("valid format", func(t *testing.T) {
		testFilePath := "testdata/testfile.ini"
		ini := INIParser{}
		expectedData := validINIDataSample1()

		err := ioutil.WriteFile(testFilePath, []byte(validINIStringSample1), 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			// Clean up the test file after the test
			err := ioutil.WriteFile(testFilePath, []byte(""), 0644)
			if err != nil {
				t.Fatal(err)
			}
		}()

		err = ini.LoadFromFile(testFilePath)
		assertNoErr(t, err)
		assertEqualData(t, ini.data, expectedData)
	})

	invalidFormatTestCases := []struct {
		name  string
		input string
		err   error
	}{
		{
			name:  "invalid format: begin with comment",
			input: ";comment",
			err:   ErrNoGlobalDataAllowed,
		}, {
			name:  "invalid format: begin with global data",
			input: "key key = value value",
			err:   ErrNoGlobalDataAllowed,
		}, {
			name: "invalid format: section naming format",
			input: `
			[section 1
			key key = value value
			`,
			err: ErrInvalidFormat,
		}, {
			name: "invalid format: key value format",
			input: `
			[section 1]
			key key  value value
			`,
			err: ErrInvalidFormat,
		},
	}

	for _, tc := range invalidFormatTestCases {
		t.Run(tc.name, func(t *testing.T) {
			testFilePath := "testdata/testfile.ini"
			ini := INIParser{}

			err := ioutil.WriteFile(testFilePath, []byte(tc.input), 0644)
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				// Clean up the test file after the test
				err := ioutil.WriteFile(testFilePath, []byte(""), 0644)
				if err != nil {
					t.Fatal(err)
				}
			}()

			err = ini.LoadFromFile(testFilePath)
			assertErrFound(t, err)
			if !errors.Is(err, tc.err) {
				t.Errorf("got %q want %q", err.Error(), tc.err.Error())
			}
		})
	}
}

func TestGetSectionNames(t *testing.T) {
	t.Run("get before load data", func(t *testing.T) {
		ini := INIParser{}
		got := ini.GetSectionNames()
		want := []string{}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %q want %q", got, want)
		}
	})
	t.Run("get data after load data", func(t *testing.T) {
		dataSample := validINIDataSample1()
		ini := INIParser{data: dataSample}
		got := ini.GetSectionNames()
		sort.Strings(got)
		want := []string{"", "section 1", "section 2"}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestGetSections(t *testing.T) {
	t.Run("get after load data", func(t *testing.T) {
		dataSample := validINIDataSample1()
		ini := INIParser{data: dataSample}
		got := ini.GetSections()
		want := map[string]map[string]string{
			"section 1": {"key key": "value value"},
			"section 2": {},
			"":          {"": ""},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("get before load data", func(t *testing.T) {
		ini := INIParser{}
		got := ini.GetSections()

		if len(got) != 0 {
			t.Errorf("expected to get empty map")
		}
	})
}

func TestGet(t *testing.T) {
	testCases := []struct {
		name          string
		section       string
		key           string
		expectedValue string
		isFound       bool
	}{
		{
			name:          "get exist value",
			section:       "section 1",
			key:           "key key",
			expectedValue: "value value",
			isFound:       true,
		}, {
			name:          "not found section",
			section:       "not found",
			key:           "not found",
			expectedValue: "",
			isFound:       false,
		}, {
			name:          "not found key",
			section:       "section 1",
			key:           "not found",
			expectedValue: "",
			isFound:       false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dataSample := validINIDataSample1()
			ini := INIParser{data: dataSample}
			gotValue, gotIsFound := ini.Get(tc.section, tc.key)

			if gotValue != tc.expectedValue {
				t.Errorf("value: got %q want %q", gotValue, tc.expectedValue)
			}
			if gotIsFound != tc.isFound {
				t.Errorf("is found: got is %t want %t", gotIsFound, tc.isFound)
			}
		})
	}
	t.Run("get before load data", func(t *testing.T) {
		ini := INIParser{}
		gotValue, gotIsFound := ini.Get("anything", "anything")

		if gotValue != "" {
			t.Errorf("value: got %q want empty string", gotValue)
		}
		if gotIsFound != false {
			t.Errorf("is found: got is %t want %t", gotIsFound, false)
		}
		
	})
}

func TestSet(t *testing.T) {
	validTestCases := []struct {
		name     string
		data     map[string]*section
		section  string
		key      string
		value    string
		expected map[string]*section
	}{
		{
			name:    "set new section and key before load data",
			section: "new section",
			key:     "new key",
			value:   "new value",
			expected: map[string]*section{
				"new section": {data: map[string]string{
					"new key": "new value",
				}},
			},
		},
		{
			name:    "set new section and key",
			data:    map[string]*section{},
			section: "new section",
			key:     "new key",
			value:   "new value",
			expected: map[string]*section{
				"new section": {data: map[string]string{
					"new key": "new value",
				}},
			},
		}, {
			name:    "set new key",
			data:    map[string]*section{"new section": {data: map[string]string{}}},
			section: "new section",
			key:     "new key",
			value:   "new value",
			expected: map[string]*section{
				"new section": {data: map[string]string{
					"new key": "new value",
				}},
			},
		}, {
			name: "update value",
			data: map[string]*section{"new section": {data: map[string]string{
				"new key": "new value",
			}}},
			section: "new section",
			key:     "new key",
			value:   "new value",
			expected: map[string]*section{
				"new section": {data: map[string]string{
					"new key": "new value",
				}},
			},
		}, {
			name:    "set empty section name, key and value",
			data:    map[string]*section{},
			section: "",
			key:     "",
			value:   "",
			expected: map[string]*section{
				"": {data: map[string]string{
					"": "",
				}},
			},
		},
	}

	for _, tc := range validTestCases {
		t.Run(tc.name, func(t *testing.T) {
			ini := INIParser{data: tc.data}
			ini.Set(tc.section, tc.key, tc.value)
			expected := tc.expected
			assertEqualData(t, ini.data, expected)
		})
	}
}

func TestString(t *testing.T) {
	t.Run("get string after load data", func(t *testing.T) {
		dataSample := validINIDataSample1()
		ini := INIParser{data: dataSample}
		got := ini.String()
	
		ini2 := INIParser{}
		ini2.LoadFromString(got)
	
		assertEqualData(t, ini2.data, dataSample)
	})
	t.Run("get string before load data", func(t *testing.T) {
		ini := INIParser{}
		got := ini.String()
	
		if got != "" {
			t.Errorf("got %q want empty string", got)
		}
	})
}

func TestSaveToFile(t *testing.T) {
	t.Run("save to file before load data", func(t *testing.T) {
		testFilePath := "testdata/testfile.ini"
		ini := INIParser{}
		err := ini.SaveToFile(testFilePath)
		assertNoErr(t, err)

		ini2 := INIParser{}
		ini2.LoadFromFile(testFilePath)

		assertEqualData(t, ini2.data, map[string]*section{})
	})
	
	t.Run("save to file", func(t *testing.T) {
		testFilePath := "testdata/testfile.ini"
		dataSample := validINIDataSample1()
		ini := INIParser{data: dataSample}
		err := ini.SaveToFile(testFilePath)
		assertNoErr(t, err)

		ini2 := INIParser{}
		ini2.LoadFromFile(testFilePath)

		assertEqualData(t, ini2.data, dataSample)
	})

	t.Run("save to file has not ini extension", func(t *testing.T) {
		testFilePath := "README.md"
		dataSample := validINIDataSample1()
		ini := INIParser{data: dataSample}
		err := ini.SaveToFile(testFilePath)
		assertErrFound(t, err)
		if !errors.Is(err, ErrInvalidFileExtension) {
			t.Errorf("got %q want %q", err.Error(), ErrInvalidFileExtension.Error())
		}
	})
}

func assertEqualData(t testing.TB, got, want map[string]*section) {
	t.Helper()
	if len(want) != len(got) {
		t.Fatal("two maps don't have the same number of sections")
	}
	for k := range got {
		if want[k] == nil {
			t.Fatal("key found in 'got map' and is not found in 'want map'")
		}
		if len(got[k].data) == len(want[k].data) && len(want[k].data) == 0 {
			continue
		}
		if !reflect.DeepEqual(got[k].data, want[k].data) {
			t.Errorf("in section %q got %v want %v", k, got[k].data, want[k].data)
		}
	}
}
