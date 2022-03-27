package config

import (
	"bufio"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

func Read(fileName string, out interface{}) error {

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	return unmarshallYml(reader, out)
}

func unmarshallYml(reader io.Reader, out interface{}) error {

	file, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal([]byte(file), out)
	if err != nil {
		return err
	}

	return nil
}
