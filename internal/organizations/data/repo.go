package data

import "errors"

type Organization struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	OwnerId int64  `json:"ownerId"`
}

func GetAll() ([]Organization, error) {
	return nil, errors.New("not implemented") // TODO
}

func GetOne(id int64) (*Organization, error) {
	return nil, errors.New("not implemented") // TODO
}

func Create() (*Organization, error) {
	return nil, errors.New("not implemented") // TODO
}

func Delete(id int64) error {
	return errors.New("not implemented") // TODO
}
