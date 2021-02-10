package database

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/datatypes"
)

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *JSONB) Scan(value interface{}) error {
	if err := json.Unmarshal([]byte(value.(string)), &j); err != nil {
		return err
	}
	return nil
}

func ToJSONB(v interface{}) datatypes.JSON {
	data, _ := json.Marshal(v)
	return datatypes.JSON(data)
}