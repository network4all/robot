package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
)

// encode files
func encode(filename string) (string, error) {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return "", fmt.Errorf("could not open file %s :%v", f.Name(), err)
	}

	reader := bufio.NewReader(f)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("could not read from file %s: %v", f.Name(), err)
	}

	encoded := base64.StdEncoding.EncodeToString(content)
	return encoded, nil
}

// decode and save files
func decode(filename string, data string) error {
	dec, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return fmt.Errorf("could not decode string %s :%v", filename, err)
	}
	f, err := os.Create(filename)

	if err != nil {
		return fmt.Errorf("could not create file %s :%v", filename, err)
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		return fmt.Errorf("could not write to file %s :%v", filename, err)
	}

	if err := f.Sync(); err != nil {
		return fmt.Errorf("could not sync file %s :%v", filename, err)
	}
	return nil
}
