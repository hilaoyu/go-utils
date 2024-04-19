package utilFile

import (
	"encoding/json"
	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v3"
	"os"
)

func ReadJSON(path string, v interface{}) error {
	var (
		data []byte
		err  error
	)
	data, err = os.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, v)
	return err
}

func ReadYaml(path string, v interface{}) error {
	var (
		data []byte
		err  error
	)
	data, err = os.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, v)
	return err
}
func ReadIni(path string, v interface{}) error {
	err := ini.MapTo(v, path)
	return err
}
