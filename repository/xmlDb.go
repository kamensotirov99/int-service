package repository

import (
	"encoding/xml"

	"io/ioutil"
)

type XMLType struct {
	extension string
}

func NewXML() DataManipulator {
	return &XMLType{
		extension: ".xml",
	}
}

func (x *XMLType) ReadData(fileName string, fileData interface{}) error {
	file, err := ioutil.ReadFile(fileName + x.extension)
	if err != nil {
		return err
	}

	err = xml.Unmarshal([]byte(file), fileData)
	if err != nil {
		return err
	}
	return err
}

func (x *XMLType) WriteFile(fileName string, fileData interface{}) error {
	data, err := xml.MarshalIndent(fileData, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fileName+x.extension, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
