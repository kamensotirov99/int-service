package repository

import (
	"context"
	"int-service/dto"

	"github.com/pkg/errors"
)

const fileName = "clothes"

type FileDatabase struct {
	ReadWrite DataManipulator
}

func NewFileDatabase(format DataManipulator) Repository {
	return &FileDatabase{
		ReadWrite: format,
	}
}

func (f *FileDatabase) CreateClothing(ctx context.Context, clothing *dto.ClothingDTO) (*dto.ClothingDTO, error) {
	clothes := dto.ClothesDTO{}
	err := f.ReadWrite.ReadData(fileName, &clothes)

	if err != nil {
		return nil, errors.Wrap(err, "Error while reading the file")
	}

	clothes.Clothes = append(clothes.Clothes, *clothing)
	err = f.ReadWrite.WriteFile(fileName, &clothes)
	if err != nil {
		return nil, errors.Wrap(err, "Error while writing the file")
	}
	return clothing, nil
}

func (f *FileDatabase) DeleteClothing(ctx context.Context, ID string) error {
	clothes := dto.ClothesDTO{}
	err := f.ReadWrite.ReadData(fileName, &clothes)
	if err != nil {
		return errors.Wrap(err, "Error while reading the file")
	}

	found := false
	for i, c := range clothes.Clothes {
		if c.ID == ID {
			clothes.Clothes = append(clothes.Clothes[:i], clothes.Clothes[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return errors.New("Error while deleting clothing by id")
	}
	return f.ReadWrite.WriteFile(fileName, &clothes)
}

func (f *FileDatabase) GetAll(ctx context.Context) (*dto.ClothesDTO, error) {
	clothes := dto.ClothesDTO{}
	err := f.ReadWrite.ReadData(fileName, &clothes)
	if err != nil {
		return nil, errors.Wrap(err, "Error while reading the file")
	}
	return &clothes, nil
}
