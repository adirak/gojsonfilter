package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/adirak/gojsonfilter/jfilter"
	"github.com/buger/jsonparser"
	"github.com/mitchellh/mapstructure"
	"github.com/pquerna/ffjson/ffjson"
)

func main() {
	testSimpleCase()
	testBigJson()
	testUnmarshalSpeed()
	testBugerJsonParserSpeed()
	testFFJsonParserSpeed()
}

func testSimpleCase() {

	fmt.Println("*********** Test Filter Simple ************")

	// Read Databus
	databus, err := readJsonMap("./example/databus_simple.json")
	if err != nil {
		panic(err)
	}

	// Read Filter
	filter, err := readJsonArr("./example/filter_simple.json")
	if err != nil {
		panic(err)
	}

	snsec := (time.Now().UnixNano())
	fmt.Printf("start=%d\n", snsec)

	result, err := jfilter.JsonFilter(databus, filter)
	if err != nil {
		fmt.Println("------- validate error --------")
		fmt.Println(err)
		fmt.Println("-------------------------------")
	}
	fmt.Println(result)

	ensec := (time.Now().UnixNano())
	fmt.Printf("end=%d\n", ensec)
	fmt.Printf("usetime=%d usec\n", ((ensec - snsec) / 1000))
	fmt.Println("****************************************")
}

func testBigJson() {

	fmt.Println("*********** Test Filter Big Json ************")

	// Read Databus
	databus, err := readJsonMap("./example/databus_big.json")
	if err != nil {
		panic(err)
	}

	// Read Filter
	filter, err := readJsonArr("./example/filter_big.json")
	if err != nil {
		panic(err)
	}

	snsec := (time.Now().UnixNano())
	fmt.Printf("start=%d\n", snsec)

	result, err := jfilter.JsonFilter(databus, filter)
	if err != nil {
		fmt.Println("------- validate error --------")
		fmt.Println(err)
		fmt.Println("-------------------------------")
	}

	ensec := (time.Now().UnixNano())
	fmt.Printf("end=%d\n", ensec)
	fmt.Printf("usetime=%d usec\n", ((ensec - snsec) / 1000))
	err = writeJson("./example/databus_big_result.json", result)
	if err != nil {
		panic(err)
	}
	fmt.Println("****************************************")
}

func testUnmarshalSpeed() {

	fmt.Println("*********** Test Unmarshal Speed ************")

	bData, err := readFile("./example/databus_big.json")
	if err != nil {
		panic(err)
	}

	t1 := (time.Now().UnixNano())
	fmt.Printf("t1=%d\n", t1)

	obj := map[string]interface{}{}
	err = json.Unmarshal(bData, &obj)
	if err != nil {
		panic(err)
	}
	t2 := (time.Now().UnixNano())
	d1 := t2 - t1
	fmt.Printf("t2=%d\n", t2)
	fmt.Printf("unmarshal to map =%d usec\n", d1/1000)

	obj2 := bigData{}
	err = mapstructure.Decode(obj, &obj2)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("obj2=%v\n", obj2)
	t3 := (time.Now().UnixNano())
	d2 := t3 - t2
	fmt.Printf("t3=%d\n", t3)
	fmt.Printf("mapstructure decode =%d usec\n", d2/1000)

	obj3 := bigData{}
	err = json.Unmarshal(bData, &obj3)
	if err != nil {
		panic(err)
	}
	t4 := (time.Now().UnixNano())
	d3 := t4 - t3
	fmt.Printf("t4=%d\n", t4)
	fmt.Printf("unmarshal to struct =%d usec\n", d3/1000)

	fmt.Println("****************************************")
}

func testBugerJsonParserSpeed() {

	fmt.Println("*********** Test BuggerJsonParser Speed ************")

	bData, err := readFile("./example/databus_big.json")
	if err != nil {
		panic(err)
	}

	t1 := (time.Now().UnixNano())
	fmt.Printf("t1=%d\n", t1)

	obj := bigData{}

	bR, _, _, _ := jsonparser.Get(bData, "data", "result")
	err = json.Unmarshal(bR, &obj.Data.Result)
	if err != nil {
		panic(err)
	}
	bR2, _, _, _ := jsonparser.Get(bData, "data", "result2")
	err = json.Unmarshal(bR2, &obj.Data.Result2)
	if err != nil {
		panic(err)
	}
	bR3, _, _, _ := jsonparser.Get(bData, "data", "result3")
	err = json.Unmarshal(bR3, &obj.Data.Result3)
	if err != nil {
		panic(err)
	}
	bR4, _, _, _ := jsonparser.Get(bData, "data", "result4")
	err = json.Unmarshal(bR4, &obj.Data.Result4)
	if err != nil {
		panic(err)
	}

	t2 := (time.Now().UnixNano())
	d1 := t2 - t1
	fmt.Printf("t2=%d\n", t2)
	fmt.Printf("time parser to struct=%d usec\n", d1/1000)

	fmt.Println("****************************************")
}

func testFFJsonParserSpeed() {

	fmt.Println("*********** Test FFJsonParser Speed ************")

	bData, err := readFile("./example/databus_big.json")
	if err != nil {
		panic(err)
	}

	t1 := (time.Now().UnixNano())
	fmt.Printf("t1=%d\n", t1)

	obj := bigData{}

	err = ffjson.Unmarshal(bData, &obj)
	if err != nil {
		panic(err)
	}

	t2 := (time.Now().UnixNano())
	d1 := t2 - t1
	fmt.Printf("t2=%d\n", t2)
	fmt.Printf("time parser to struct=%d usec\n", d1/1000)

	obj2 := map[string]interface{}{}
	err = ffjson.Unmarshal(bData, &obj2)
	if err != nil {
		panic(err)
	}

	t3 := (time.Now().UnixNano())
	d2 := t3 - t2
	fmt.Printf("t3=%d\n", t3)
	fmt.Printf("time parser to map=%d usec\n", d2/1000)

	fmt.Println("****************************************")
}

func readJsonArr(path string) ([]interface{}, error) {

	bData, err := readFile(path)
	if err != nil {
		return nil, err
	}

	arr := []interface{}{}
	err = json.Unmarshal([]byte(bData), &arr)

	return arr, err
}

func readJsonMap(path string) (map[string]interface{}, error) {

	bData, err := readFile(path)
	if err != nil {
		return nil, err
	}

	obj := map[string]interface{}{}
	err = json.Unmarshal([]byte(bData), &obj)

	return obj, err
}

func readFile(path string) (bData []byte, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}

	defer file.Close()
	bData, err = ioutil.ReadAll(file)

	return
}

func writeJson(path string, data interface{}) error {
	bData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return writeFile(bData, path)
}

func writeFile(bData []byte, path string) error {

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(bData)
	if err != nil {
		return err
	}

	return nil
}

type bigData struct {
	Data struct {
		Result  interface{} `structs:"result" json:"result" bson:"result"`
		Result2 interface{} `structs:"result2" json:"result2" bson:"result2"`
		Result3 interface{} `structs:"result3" json:"result3" bson:"result3"`
		Result4 interface{} `structs:"result4" json:"result4" bson:"result4"`
	} `structs:"data" json:"data" bson:"data"`
}
