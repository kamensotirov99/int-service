package repository

import (
	"encoding/json"

	"io/ioutil"
)

type JSONType struct {
	extension string
}

func NewJSON() DataManipulator {
	return &JSONType{
		extension: ".json",
	}
}

func (j *JSONType) ReadData(fileName string, fileData interface{}) error {
	content, err := ioutil.ReadFile(fileName + j.extension)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(content), fileData)
}

func (j *JSONType) WriteFile(fileName string, fileData interface{}) error {
	data, err := json.MarshalIndent(fileData, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName+j.extension, data, 0644)
}
