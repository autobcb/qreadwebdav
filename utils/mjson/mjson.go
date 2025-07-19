package mjson

import (
	"encoding/json"
)

func Parse(jsonstr string) map[string]interface{} {
	var animals = map[string]interface{}{}
	err := json.Unmarshal([]byte(jsonstr), &animals)
	if err != nil {
		panic(err.Error())
	}
	return animals
}

func ParseHasErr(jsonstr string) (map[string]interface{}, error) {
	var animals = map[string]interface{}{}
	err := json.Unmarshal([]byte(jsonstr), &animals)
	if err != nil {
		return nil, err
	}
	return animals, nil
}

func ToJson(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err.Error())
	}
	return string(b)
}

func ParseByobj(jsonstr string, animals interface{}) interface{} {
	err := json.Unmarshal([]byte(jsonstr), animals)
	if err != nil {
		panic(err.Error())
	}
	return animals
}

func ParseByobjHasErr(jsonstr string, animals interface{}) error {
	err := json.Unmarshal([]byte(jsonstr), animals)
	if err != nil {
		return err
	}
	return nil
}
