package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/adirak/gojsonfilter/jfilter"
)

// ReadJsonFile is function to read json file
func ReadJsonFile(path string) (data interface{}, err error) {

	// Read file
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	// Convert to Map
	mapData := map[string]interface{}{}
	err = json.Unmarshal([]byte(file), &mapData)
	if err == nil {
		return mapData, nil
	}

	// Convert to Array
	arrData := []interface{}{}
	err = json.Unmarshal([]byte(file), &arrData)
	return arrData, err
}

func TestFilter1(t *testing.T) {

	obj, err := ReadJsonFile("./example/test_1.json")
	if err != nil {
		t.Error(err)
	}

	mapObj, _ := obj.(map[string]interface{})
	databus := mapObj["dataBus"]
	filter, _ := mapObj["filter"].([]interface{})

	out, err := jfilter.JsonFilter(databus, filter)
	if err != nil {
		t.Error(err)
	}

	t.Log(out)

}

func TestReqular(t *testing.T) {

	name := "name"
	val := "123456A"
	regExp := "^[ A-Za-z0-9]*$"

	err := jfilter.ValidateRegExp(name, val, regExp)
	if err != nil {
		t.Error(err)
	}

}
