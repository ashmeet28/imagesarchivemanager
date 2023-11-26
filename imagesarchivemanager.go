package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"
	"path"
	"time"
)

func newImgId() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(b)
}

func newImgTimePrefix() string {
	t := time.Now().String()
	return t[0:4] + t[5:7] + t[8:10]
}

func main() {
	op := os.Args[1]
	dir1 := os.Args[2]
	dir2 := os.Args[3]

	switch op {
	case "renamefiles":
		files, err := os.ReadDir(dir1)

		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			fName := file.Name()

			newFName := newImgTimePrefix() + "_" + newImgId() + path.Ext(fName)

			data, err := os.ReadFile(path.Join(dir1, fName))
			if err != nil {
				log.Fatal(err)
			}

			os.WriteFile(path.Join(dir2, newFName), data, 0666)
		}
	}
}
