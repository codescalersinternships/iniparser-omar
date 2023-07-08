package iniparser

import (
	"reflect"
	"testing"
)

func Test_loadData(t *testing.T) {
	t.Run("valid format", func(t *testing.T) {
		ini := INIParser{}
		input := []string{"", "", "[section name]", "name=omar abdelghani", "age=-1", "[section 2]", "[section 3]", "hi hi = hello hello"}

		expectedData := map[string]*section{
			"section name": {data: map[string]string{
				"name": "omar abdelghani",
				"age":  "-1",
			}},
			"section 2": {data: map[string]string{}},
			"section 3": {data: map[string]string{
				"hi hi": "hello hello",
			}},
		}
		expectedSectionNames := []string{"section name", "section 2", "section 3"}

		got := ini.loadData(input)
		if got != nil {
			t.Fatal("not expected get an error, got: " + got.Error())
		}
		assertEqualMaps(t, ini.data, expectedData)
		if !reflect.DeepEqual(ini.sectionNames, expectedSectionNames) {
			t.Errorf("got %v want %v", ini.sectionNames, expectedSectionNames)
		}
	})

	testCases := []struct {
		name  string
		input []string
	}{
		{
			name:  "invalid format: begin with comment",
			input: []string{";comment", "[section name]", "name=omar abdelghani", " age=-1", "[section 2]", "[section 3]", "hi hi = hello hello"},
		}, {
			name:  "invalid format: begin with global data",
			input: []string{"age=1", "[section name]", "name=omar abdelghani", " age=-1", "[section 2]", "[section 3]", "hi hi = hello hello"},
		}, {
			name:  "invalid format: empty section name",
			input: []string{"[section name]", "name=omar abdelghani", " age=-1", "[section 2]", "[]", "hi hi = hello hello"},
		}, {
			name:  "invalid format: empty key",
			input: []string{"[section name]", "name=omar abdelghani", " age=-1", "[section 2]", "[section 3]", " = hello hello"},
		}, {
			name:  "invalid format: empty value",
			input: []string{"[section name]", "name=omar abdelghani", " age=-1", "[section 2]", "[section 3]", "hi hi = "},
		}, {
			name:  "invalid format: section naming format",
			input: []string{"[section name]", "name=omar abdelghani", " age=-1", "[section 2]", "[section 3", "hi hi = hello hello"},
		}, {
			name:  "invalid format: key value format",
			input: []string{"[section name]", "name=omar abdelghani", " age=-1", "[section 2]", "[section 3]", "hi hi  hello hello"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ini := INIParser{}
			got := ini.loadData(tc.input)
			if got == nil {
				t.Fatal("expected to get error.")
			}
			if got.Error() != errInvalidFormat.Error() {
				t.Errorf("got %q want %q", got.Error(), errInvalidFormat.Error())
			}
		})
	}
}

func TestGetSections(t *testing.T) {
	ini := INIParser{data: map[string]*section{
		"section name": {data: map[string]string{
			"name": "omar abdelghani",
			"age":  "-1",
		}},
		"section 2": {data: map[string]string{}},
		"section 3": {data: map[string]string{
			"hi  hi": "hello hello",
		}},
	},
		sectionNames: []string{"section name", "section 2", "section 3"}}

	got := ini.GetSections()
	want := map[string]map[string]string{
		"section name": {
			"name": "omar abdelghani",
			"age":  "-1",
		},
		"section 2": {},
		"section 3": {
			"hi  hi": "hello hello",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestGet(t *testing.T) {
	testCases := []struct {
		name     string
		section  string
		key      string
		expected string
		data     map[string]*section
	}{
		{
			name:     "normal get",
			section:  "section name2",
			key:      "name",
			expected: "omar abdelghani",
			data: map[string]*section{
				"section name": {data: map[string]string{
					"name": "omar",
					"age":  "-1",
				}},
				"section name2": {data: map[string]string{
					"name": "omar abdelghani",
					"age":  "-1",
				}},
			},
		}, {
			name:     "not found section",
			section:  "not found",
			key:      "not found",
			expected: "",
			data: map[string]*section{
				"section name": {data: map[string]string{
					"name": "omar abdelghani",
				}},
			},
		}, {
			name:     "not found key",
			section:  "section name",
			key:      "not found",
			expected: "",
			data: map[string]*section{
				"section name": {data: map[string]string{
					"name": "omar abdelghani",
				}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ini := INIParser{data: tc.data}
			got := ini.Get(tc.section, tc.key)
			want := tc.expected
			if got != want {
				t.Errorf("got %q want %q", got, want)
			}
		})
	}
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
		},
	}

	invalidTestCases := []struct {
		name    string
		data    map[string]*section
		section string
		key     string
		value   string
	}{
		{
			name:    "set empty section",
			section: "",
			key:     "new key",
			value:   "new value",
		},
		{
			name:    "set empty key",
			section: "new section",
			key:     "",
			value:   "new value",
		}, {
			name:    "set empty value",
			section: "new section",
			key:     "new key",
			value:   "",
		},
	}

	for _, tc := range validTestCases {
		t.Run(tc.name, func(t *testing.T) {
			ini := INIParser{data: tc.data}
			err := ini.Set(tc.section, tc.key, tc.value)
			if err != nil {
				t.Fatal("not expected get an error, got: " + err.Error())
			}
			expected := tc.expected
			assertEqualMaps(t, ini.data, expected)
		})
	}

	for _, tc := range invalidTestCases {
		ini := INIParser{data: map[string]*section{"new section": {data: map[string]string{
			"new key": "new value",
		}}}}
		err := ini.Set(tc.section, tc.key, tc.value)
		if err == nil {
			t.Fatal("expected to get error.")
		}
		if err.Error() != errInvalidFormat.Error() {
			t.Errorf("get %q want %q", err.Error(), errInvalidFormat.Error())
		}
	}
}

func TestToString(t *testing.T) {
	ini := INIParser{data: map[string]*section{
		"section name": {
			data: map[string]string{
				"name": "omar abdelghani",
				"age":  "-1",
			},
			lines: []string{"name", "age"},
		},
		"section 2": {data: map[string]string{}, lines: []string{}},
		"section 3": {
			data: map[string]string{
				"hi hi": "hello hello",
			},
			lines: []string{"hi hi"}},
	},
		sectionNames: []string{"section name", "section 2", "section 3"},
	}
	got := ini.ToString()
	want := "[section name]\nname = omar abdelghani\nage = -1\n[section 2]\n\n[section 3]\nhi hi = hello hello\n"
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func assertEqualMaps(t testing.TB, got, want map[string]*section) {
	t.Helper()
	if len(want) != len(got) {
		t.Fatal("two maps don't have the same number of sections")
	}
	for k := range got {
		if want[k] == nil {
			t.Fatal("two maps don't have the same number of sections")
		}
		if !reflect.DeepEqual(got[k].data, want[k].data) {
			t.Errorf("in section %q got %v want %v", k, got[k].data, want[k].data)
		}
	}
}
