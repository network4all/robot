package main

import (
	"bufio"
	"encoding/base64"
	"io/ioutil"
	"os"
)

func encode(filename string) string {

	f, err := os.Open(filename)
    defer f.Close()
    if err != nil {
		panic(err)
	}
	
	reader := bufio.NewReader(f)
	content, err := ioutil.ReadAll(reader)
    if err != nil {
		panic(err)
	}

	encoded := base64.StdEncoding.EncodeToString(content)
	return encoded
}

func decode(filename string, data string) {
	dec, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		panic(err)
	}
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		panic(err)
	}
	if err := f.Sync(); err != nil {
		panic(err)
	}
}
