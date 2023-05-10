package jfilter

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Data Type for Hyper
const type_string = "string"
const type_integer = "integer"
const type_decimal = "decimal"
const type_number = "number"
const type_map = "map"
const type_array = "array"
const type_boolean = "boolean"

// fieldFilter is data model for filter some field
type fieldFilter struct {

	// Name and type of Field to filter
	Name string `structs:"name" json:"name" bson:"name"`
	Type string `structs:"type" json:"type" bson:"type"`

	// If all is true it's mean exposed all fields
	All bool `structs:"all" json:"all" bson:"all"`

	// Sub Fields
	Children []interface{} `structs:"children" json:"children" bson:"children"`

	// It is default value of field
	Default interface{} `structs:"default" json:"default" bson:"default"`

	// Validating Value
	Validated bool    `structs:"validated" json:"validated" bson:"validated"`
	Required  bool    `structs:"required" json:"required" bson:"required"`
	RegExp    string  `structs:"regexp" json:"regexp" bson:"regexp"`
	Min       float64 `structs:"min" json:"min" bson:"min"`
	Max       float64 `structs:"max" json:"max" bson:"max"`
}

// obj2FieldFilter is function to convert obj to fieldFilter data
func obj2FieldFilter(obj interface{}, fieldFilter *fieldFilter) (err error) {

	// Check null
	if obj == nil {
		err = errors.New("filter obj is null")
		return
	}
	if fieldFilter == nil {
		err = errors.New("fieldFilter is null")
		return
	}

	mapObj, ok := obj.(map[string]interface{})
	if !ok {
		err = errors.New("filter obj is not map")
		return
	}

	// Get value
	ffName, _ := mapObj["name"].(string)
	ffType, _ := mapObj["type"].(string)
	ffAll, _ := mapObj["all"].(bool)
	ffDefault := mapObj["default"]
	ffRequired, _ := mapObj["required"].(bool)
	ffChildren, _ := mapObj["children"].([]interface{})
	ffValidated, _ := mapObj["validated"].(bool)
	ffRegExp, _ := mapObj["regexp"].(string)
	ffMin := getDecimalValue(mapObj["min"])
	ffMax := getDecimalValue(mapObj["max"])

	// Check name and type is empty
	if ffName == "" {
		err = errors.New("filter name is empty")
		return
	}
	if ffType == "" {
		err = errors.New("filter type is empty")
		return
	}

	// Set value back
	fieldFilter.Name = ffName
	fieldFilter.Type = ffType
	fieldFilter.All = ffAll
	fieldFilter.Default = ffDefault
	fieldFilter.Required = ffRequired
	fieldFilter.Children = ffChildren
	fieldFilter.Validated = ffValidated
	fieldFilter.RegExp = ffRegExp
	fieldFilter.Min = ffMin
	fieldFilter.Max = ffMax

	return
}

// jsonFilterMap is function to filter data in dataBus, support root is Map type
func jsonFilterMap(dataBus map[string]interface{}, filter []interface{}) (fData map[string]interface{}, err error) {

	if len(filter) > 0 && dataBus != nil {

		// Init empty map
		fData = make(map[string]interface{})

		// Loop to filter
		for _, fObj := range filter {

			var fItem fieldFilter

			// Convert to fieldFilter
			err = obj2FieldFilter(fObj, &fItem)
			if err != nil {
				return
			}

			fName := fItem.Name
			fType := fItem.Type
			fAll := fItem.All
			fChilds := fItem.Children
			fDefault := fItem.Default
			fValidated := fItem.Validated

			// Value from dataBus
			value := dataBus[fName]

			// Set default value if it's null
			if value == nil {
				value = getDefaultValue(fDefault, fType)
			}

			// Check Validate before do anything
			if fValidated {
				err = validateValue(value, &fItem)
				if err != nil {
					return
				}
			}

			// Check Type
			if strings.EqualFold(type_map, fType) {

				// Check All Flag
				if fAll {

					// Exposed all fields
					fData[fName] = value

				} else {

					// Recursive to filter
					mapValue, ok := value.(map[string]interface{})
					if ok {

						// Call Recursive
						fMapValue, err := jsonFilterMap(mapValue, fChilds)
						if err != nil {
							return fData, err
						}

						// Set filter value
						if fMapValue != nil {
							fData[fName] = fMapValue
						}

					} else {
						msg := fmt.Sprintf("Field \"%s\" is not map", fName)
						err = errors.New(msg)
						return
					}
				}

			} else if strings.EqualFold(type_array, fType) {

				// Check All Flag
				if fAll {

					// Exposed all fields
					fData[fName] = value

				} else {

					// Recursive to filter
					arrayValue, ok := value.([]interface{})
					if ok {

						// Call Recursive
						fArrayValue, err := jsonFilterArray(arrayValue, fChilds)
						if err != nil {
							return fData, err
						}

						// Set filter value
						if fArrayValue != nil {
							fData[fName] = fArrayValue
						}

					} else {
						msg := fmt.Sprintf("Field \"%s\" is not array", fName)
						err = errors.New(msg)
						return
					}

				}

			} else {

				// Other Type
				// string, integer, decimal, boolean
				if value != nil {
					fData[fName] = value
				}

			}

		}
	} else {
		fData = map[string]interface{}{}
	}

	return fData, err

}

// JsonFilterArray is function to filter data in dataBus, support root is Array type
func jsonFilterArray(dataBusArr []interface{}, filter []interface{}) (fDatas []interface{}, err error) {

	if len(filter) > 0 && dataBusArr != nil {

		// Init empty array
		fDatas = make([]interface{}, 0)

		// child of array is index 0 only
		var fItem fieldFilter

		// Convert to fieldFilter
		err = obj2FieldFilter(filter[0], &fItem)
		if err != nil {
			return
		}

		// field value to check
		fName := fItem.Name
		fType := fItem.Type
		fAll := fItem.All
		fChilds := fItem.Children
		fDefault := fItem.Default
		fValidated := fItem.Validated

		// Loop to filter
		for _, value := range dataBusArr {

			// Set default value if it's null
			if value == nil {
				value = getDefaultValue(fDefault, fType)
			}

			// Check Validate before do anything
			if fValidated {
				err = validateValue(value, &fItem)
				if err != nil {
					return
				}
			}

			// Check Type
			if strings.EqualFold(type_map, fType) {

				// Check All Flag
				if fAll {

					// Exposed all
					fDatas = append(fDatas, value)

				} else {

					// Recursive to filter
					mapValue, ok := value.(map[string]interface{})
					if ok {

						// Call Recursive
						fMapValue, err := jsonFilterMap(mapValue, fChilds)
						if err != nil {
							return fDatas, err
						}

						// Add filter value
						if fMapValue != nil {
							fDatas = append(fDatas, fMapValue)
						}

					} else {
						msg := fmt.Sprintf("Field \"%s\" is not map", fName)
						err = errors.New(msg)
						return
					}

				}

			} else if strings.EqualFold(type_array, fType) {

				// Check All Flag
				if fAll {

					// Exposed all
					fDatas = append(fDatas, value)

				} else {

					// Array Type
					// Recursive to filter
					arrayValue, ok := value.([]interface{})
					if ok {

						// Call Recursive
						fArrValue, err := jsonFilterArray(arrayValue, fChilds)
						if err != nil {
							return fDatas, err
						}

						// Add filter value
						if fArrValue != nil {
							fDatas = append(fDatas, fArrValue)
						}

					} else {
						msg := fmt.Sprintf("Field \"%s\" is not array", fName)
						err = errors.New(msg)
						return
					}
				}

			} else {

				// Other Type
				// string, integer, decimal, boolean
				if value != nil {
					fDatas = append(fDatas, value)
				}

			}

		}

	} else {
		fDatas = []interface{}{}
	}

	return fDatas, err
}

// validateValue is function to validate value
func validateValue(val interface{}, fieldFilter *fieldFilter) (err error) {

	// init
	name := fieldFilter.Name
	objType := fieldFilter.Type
	required := fieldFilter.Required
	checkMinMax := true
	if fieldFilter.Max == 0 && fieldFilter.Min == 0 {
		checkMinMax = false
	}

	// Check Required Field
	if val == nil && required {
		msg := fmt.Sprintf("Field \"%s\" is required", name)
		err = errors.New(msg)
		return
	}

	// Check type for default
	switch objType {
	case type_string:

		// Check type
		sVal, ok := val.(string)
		if !ok {
			msg := fmt.Sprintf("Field \"%s\" is not string", name)
			err = errors.New(msg)
			return
		}
		// Ignore if not requiered as empty string
		if !required && sVal == "" {
			return
		}
		// Check pattern
		if fieldFilter.RegExp != "" {
			return validateRegExp(name, sVal, fieldFilter.RegExp)
		}

		// Get Value
		length := len(sVal)

		// Check length Max
		if checkMinMax && length > int(fieldFilter.Max) {
			msg := fmt.Sprintf("Field \"%s\" more than maximum length", name)
			err = errors.New(msg)
			return
		}
		// Check length Min
		if checkMinMax && length < int(fieldFilter.Min) {
			msg := fmt.Sprintf("Field \"%s\" less than minimum length", name)
			err = errors.New(msg)
			return
		}

		return nil

	case type_boolean:

		// Check Type
		_, ok := val.(bool)
		if !ok {
			msg := fmt.Sprintf("Field \"%s\" is not boolean", name)
			err = errors.New(msg)
			return
		}

		return nil

	case type_decimal, type_integer, type_number:

		// Ignore if not requiered as null
		if !required && val == nil {
			return
		}

		// Check Type
		isNum := isTypeNumber(val)
		if !isNum {
			msg := fmt.Sprintf("Field \"%s\" is not number", name)
			err = errors.New(msg)
			return
		}

		// Get Value
		value := getDecimalValue(val)

		// Check Valule Max
		if checkMinMax && value > fieldFilter.Max {
			msg := fmt.Sprintf("Field \"%s\" more than maximum value", name)
			err = errors.New(msg)
			return
		}
		// Check Value Min
		if checkMinMax && value < fieldFilter.Min {
			msg := fmt.Sprintf("Field \"%s\" less than minimum value", name)
			err = errors.New(msg)
			return
		}

		return nil

	case type_map:

		// Ignore if not requiered as null
		if !required && val == nil {
			return
		}

		// Check Type
		isMap := isTypeMap(val)
		if !isMap {
			msg := fmt.Sprintf("Field \"%s\" is not map", name)
			err = errors.New(msg)
			return
		}
		// Check length Max
		mapObj, _ := val.(map[string]interface{})
		length := len(mapObj)
		if fieldFilter.Max > 0 && length > int(fieldFilter.Max) {
			msg := fmt.Sprintf("Field \"%s\" more than maximum size", name)
			err = errors.New(msg)
			return
		}
		// Check length Min
		if fieldFilter.Min >= 0 && length < int(fieldFilter.Min) {
			msg := fmt.Sprintf("Field \"%s\" less than minimum size", name)
			err = errors.New(msg)
			return
		}

		return nil

	case type_array:

		// Ignore if not requiered as null
		if !required && val == nil {
			return
		}

		// Check Type
		isArr := isTypeArray(val)
		if !isArr {
			msg := fmt.Sprintf("Field \"%s\" is not array", name)
			err = errors.New(msg)
			return
		}
		// Check length Max
		arrObj, _ := val.([]interface{})
		length := len(arrObj)
		if fieldFilter.Max > 0 && length > int(fieldFilter.Max) {
			msg := fmt.Sprintf("Field \"%s\" more than maximum length", name)
			err = errors.New(msg)
			return
		}
		// Check length Min
		if fieldFilter.Min >= 0 && length < int(fieldFilter.Min) {
			msg := fmt.Sprintf("Field \"%s\" less than minimum length", name)
			err = errors.New(msg)
			return
		}

		return nil

	}

	return nil
}

// validateRegExp is function to validate by regular expression
func validateRegExp(name, sVal string, regExp string) error {
	match, err := regexp.MatchString(regExp, sVal)
	if err != nil {
		msg := fmt.Sprintf("Regular expression invalid for field \"%s\"", name)
		return errors.New(msg)
	}
	if !match {
		msg := fmt.Sprintf("Value pattern not match for field \"%s\"", name)
		return errors.New(msg)
	}
	return nil
}

// ValidateRegExp is function to validate by regular expession
func ValidateRegExp(name, sVal string, regExp string) error {
	return validateRegExp(name, sVal, regExp)
}

// Get default value and cast type
func getDefaultValue(dfVal interface{}, objType string) interface{} {

	// Check empty
	empty := false
	if dfVal == nil {
		empty = true
	}
	sVal := fmt.Sprintf("%v", dfVal)
	if sVal == "" {
		empty = true
	}

	switch objType {

	case type_integer:
		if empty {
			return nil
		}
		f := getDecimalValue(dfVal)
		return int64(f)

	case type_decimal, type_number:
		if empty {
			return nil
		}
		f := getDecimalValue(dfVal)
		return f

	case type_boolean:
		if empty {
			return nil
		}
		b, _ := strconv.ParseBool(sVal)
		return b

	case type_map:
		return map[string]interface{}{}

	case type_array:
		return []interface{}{}

	}

	return dfVal
}

// getDecimalValue is function to get decimal value from object
func getDecimalValue(val interface{}) float64 {

	if val == nil {
		return 0
	}
	sVal := fmt.Sprintf("%v", val)
	f, _ := strconv.ParseFloat(sVal, 64)
	return f
}

// isNumber is function to check number in databus
func isTypeNumber(val interface{}) bool {
	switch val.(type) {
	case float64:
		return true
	case float32:
		return true
	case int64:
		return true
	case int32:
		return true
	case int:
		return true
	}
	return false
}

// isTypeArray is function to check val whether that is array
func isTypeArray(val interface{}) bool {
	rt := reflect.TypeOf(val)
	switch rt.Kind() {
	case reflect.Slice:
		return true
	case reflect.Array:
		return true
	}
	return false
}

// isTypeMap is function to check val whether that is map
func isTypeMap(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Map
}
