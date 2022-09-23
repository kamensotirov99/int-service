package repository

type DataManipulator interface {
	ReadData(fileName string, fileData interface{}) error
	WriteFile(fileName string, result interface{}) error
}