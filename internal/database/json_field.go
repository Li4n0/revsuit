package database

import (
	"database/sql/driver"
	"encoding/json"
)

type MapField map[string]string

func (f MapField) Value() (driver.Value, error) {
	return json.Marshal(f)
}

func (f *MapField) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), f)
}


type ListField []string

func (f ListField) Value() (driver.Value, error) {
	return json.Marshal(f)
}

func (f *ListField) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), f)
}