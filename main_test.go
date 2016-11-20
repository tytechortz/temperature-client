package main

import (
	"io/ioutil"
	"testing"
)

var (
	testData = []byte(`3c 01 4b 46 7f ff 04 10 40 : crc=40 YES
3c 01 4b 46 7f ff 04 10 40 t=19750`)
)

func TestGetTemperature(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	f.Write(testData)
	f.Close()
	val := getTemperature(f.Name())
	if val != 19.750 {
		t.Fatal(val, "does not equal 19.750")
	}
}
