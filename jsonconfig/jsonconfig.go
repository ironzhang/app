package jsonconfig

import (
	"encoding/json"
	"fmt"

	"github.com/ironzhang/x-pearls/config"
)

func Load(filename string, values map[string]interface{}) (err error) {
	var m map[string]json.RawMessage
	if err = config.JSON.LoadFromFile(filename, &m); err != nil {
		return err
	}
	for k, raw := range m {
		if v, ok := values[k]; ok {
			if err = json.Unmarshal(raw, v); err != nil {
				return fmt.Errorf("unmarshal value of %q: %v", k, err)
			}
		}
	}
	return nil
}

func Write(filename string, values map[string]interface{}) error {
	return config.JSON.WriteToFile(filename, values)
}
