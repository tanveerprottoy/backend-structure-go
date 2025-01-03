package sqlext

import (
	"encoding/json"
	"fmt"
)

// JsonObject is a type for DB json array type
type JsonObject map[string]any

// Scan implements the Scanner interface for JsonObject
func (j *JsonObject) Scan(val any) error {
	var jsonObj map[string]any
	if val == nil {
		j = nil
		return nil
	}
	switch v := val.(type) {
	case []byte:
		err := json.Unmarshal(v, &jsonObj)
		if err != nil {
			return err
		}
		*j = jsonObj
		return nil
	}
	return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type map[string]any", val)
}

// JsonArray is a type for DB json array type
type JsonArray []map[string]any

// Scan implements the Scanner interface for JsonArray
func (j *JsonArray) Scan(val any) error {
	var jsonData []map[string]any
	switch v := val.(type) {
	case []byte:
		err := json.Unmarshal(v, &jsonData)
		if err != nil {
			return err
		}
		*j = jsonData
		return nil
	}
	return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type []map[string]any", val)
}

// JsonStringArray is a type for DB json string array type
type JsonStringArray []string

// Scan implements the Scanner interface for JsonStringArray
func (j *JsonStringArray) Scan(val any) error {
	var jsonData []string
	switch v := val.(type) {
	case []byte:
		err := json.Unmarshal(v, &jsonData)
		if err != nil {
			return err
		}
		*j = jsonData
		return nil

	case string:
		err := json.Unmarshal([]byte(val.(string)), &jsonData)
		if err != nil {
			return err
		}
		*j = jsonData
		return nil
	}
	return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type []string", val)
}

// JsonIntArray is a type for DB json int array type
type JsonIntArray []int

// Scan implements the Scanner interface for JsonIntArray
func (j *JsonIntArray) Scan(val any) error {
	var jsonData []int
	switch v := val.(type) {
	case []byte:
		err := json.Unmarshal(v, &jsonData)
		if err != nil {
			return err
		}
		*j = jsonData
		return nil

	case string:
		err := json.Unmarshal([]byte(val.(string)), &jsonData)
		if err != nil {
			return err
		}
		*j = jsonData
		return nil
	}
	return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type []int", val)
}
