package main

import (
	"io/ioutil"
)

func writeFile(b []byte) {
	c := loadConfig()
	// write the whole body at once
	err := ioutil.WriteFile(c.CredFilePath, b, 0644)
	if err != nil {
		panic(err)
	}
}
