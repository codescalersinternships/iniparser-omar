package iniparser

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func validSectionDataSample1() map[string]string {
	return map[string]string{
		"key1": "value1",
		"key2": "value2",
		"":     "",
	}
}
func validSectionDataSample2() map[string]string {
	return map[string]string{
		"key1": "value1-2",
		"key2": "value2-2",
		"":     "",
	}
}

var validSectionDataString1 = `key1 = value1
key2 = value2
 = 
`

func TestSectionAddLine(t *testing.T) {
	t.Run("add correct data", func(t *testing.T) {
		sec := section{}
		err := sec.addLine("key1 = value1")
		assertNoErr(t, err)
		err = sec.addLine(";comment")
		assertNoErr(t, err)
		err = sec.addLine("key2 = value2")
		assertNoErr(t, err)
		err = sec.addLine("")
		assertNoErr(t, err)
		err = sec.addLine("=")
		assertNoErr(t, err)
		want := validSectionDataSample1()

		if !reflect.DeepEqual(sec.data, want) {
			t.Errorf("got %q want %q", sec.data, want)
		}
	})
	t.Run("add incorrect data", func(t *testing.T) {
		sec := section{}
		line := "key value"
		err := sec.addLine(line)
		assertErrFound(t, err)
		want := fmt.Errorf(errInvalidLine, ErrInvalidFormat, line)

		if err.Error() != want.Error() {
			t.Errorf("got %q want %q", err.Error(), want)
		}
	})
}

func TestSectionGet(t *testing.T) {
	t.Run("get value from key exist", func(t *testing.T) {
		dataSample := validSectionDataSample1()
		sec := section{data: dataSample}
		key := "key1"
		val, exist := sec.get(key)
		if !exist {
			t.Errorf("got %q not exist and should be exist", key)
		}
		if val != dataSample[key] {
			t.Errorf("got %q want %q", val, dataSample[key])
		}
	})
	t.Run("get value from key not exist", func(t *testing.T) {
		sec := section{data: validSectionDataSample1()}
		key := "not exist"
		_, exist := sec.get(key)
		if exist {
			t.Errorf("got %q exists and should not be exist", key)
		}
	})
	t.Run("get value from key not exist before load data", func(t *testing.T) {
		sec := section{}
		key := "not exist"
		_, exist := sec.get(key)
		if exist {
			t.Errorf("got %q exists and should not be exist", key)
		}
	})
}

func TestSectionSet(t *testing.T) {
	t.Run("set not exist key and value", func(t *testing.T) {
		dataSample := validSectionDataSample1()
		sec := section{}
		key1 := "key1"
		key2 := "key2"
		key3 := ""
		sec.set(key1, dataSample[key1])
		sec.set(key2, dataSample[key2])
		sec.set(key3, dataSample[key3])

		if !reflect.DeepEqual(sec.data, dataSample) {
			t.Errorf("got %v want %v", sec.data, dataSample)
		}
	})
	t.Run("update value", func(t *testing.T) {
		dataSample1 := validSectionDataSample1()
		dataSample2 := validSectionDataSample2()
		sec := section{data: dataSample1}
		key1 := "key1"
		key2 := "key2"
		sec.set(key1, dataSample2[key1])
		sec.set(key2, dataSample2[key2])

		if !reflect.DeepEqual(sec.data, dataSample2) {
			t.Errorf("got %v want %v", sec.data, dataSample2)
		}
	})
}

func TestSectionGetSection(t *testing.T) {
	dataSample := validSectionDataSample1()
	sec := section{data: dataSample}
	got := sec.getSection()
	if !reflect.DeepEqual(got, dataSample) {
		t.Errorf("got %v want %v", got, dataSample)
	}
}

func TestSectionString(t *testing.T) {
	dataSample := validSectionDataSample1()
	sec := section{data: dataSample}
	got := sec.string()

	sec2 := section{}
	for _, line := range strings.Split(got, "\n") {
		sec2.addLine(line)
	}
	if !reflect.DeepEqual(sec2.data, dataSample) {
		t.Errorf("got %q want %q", got, validSectionDataString1)
	}
}

func assertNoErr(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatal("expected not to get an error, got: " + err.Error())
	}
}

func assertErrFound(t testing.TB, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected to get an error, but got nothing")
	}
}
