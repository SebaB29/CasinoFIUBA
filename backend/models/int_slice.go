package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type IntSlice []int

func (s IntSlice) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *IntSlice) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("el tipo no es []byte")
	}
	return json.Unmarshal(bytes, s)
}
