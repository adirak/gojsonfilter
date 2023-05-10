package jfilter

import "errors"

// JsonFilter is function to filter data from databus
func JsonFilter(dataBus interface{}, filter []interface{}) (out interface{}, err error) {

	// Check is map dataBus
	dataBusMap, isMap := dataBus.(map[string]interface{})
	if isMap {
		return jsonFilterMap(dataBusMap, filter)
	}

	// Check is array databus
	dataBusArray, isArray := dataBus.([]interface{})
	if isArray {
		return jsonFilterArray(dataBusArray, filter)
	}

	return nil, errors.New("dataBus type is invalid")
}
