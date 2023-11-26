package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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

func createHashMapForExistingFiles(d string) map[string]bool {
	files, err := os.ReadDir(d)
	if err != nil {
		log.Fatal(err)
	}

	existingFilesHashMap := map[string]bool{}

	for _, f := range files {
		fName := f.Name()
		data, err := os.ReadFile(path.Join(d, fName))
		if err != nil {
			log.Fatal(err)
		}
		sum := sha256.Sum256(data)
		existingFilesHashMap[hex.EncodeToString(sum[0:32])] = true
	}
	return existingFilesHashMap
}

func main() {
	op := os.Args[1]

	switch op {
	case "renamefiles":
		dir1 := os.Args[2]
		dir2 := os.Args[3]

		files, err := os.ReadDir(dir1)

		if err != nil {
			log.Fatal(err)
		}

		nameTimePrefix := newImgTimePrefix()

		for _, f := range files {
			fName := f.Name()

			newFName := nameTimePrefix + "_" + newImgId() + "_0" + path.Ext(fName)

			data, err := os.ReadFile(path.Join(dir1, fName))
			if err != nil {
				log.Fatal(err)
			}

			os.WriteFile(path.Join(dir2, newFName), data, 0666)
		}
	case "copynewfilesonly":
		dir1 := os.Args[2]
		dir2 := os.Args[3]

		existingFilesHashMap := createHashMapForExistingFiles(dir2)

		files, err := os.ReadDir(dir1)
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range files {
			fName := f.Name()
			data, err := os.ReadFile(path.Join(dir1, fName))
			if err != nil {
				log.Fatal(err)
			}
			sum := sha256.Sum256(data)

			if !existingFilesHashMap[hex.EncodeToString(sum[0:32])] {
				os.WriteFile(path.Join(dir2, fName), data, 0666)
				existingFilesHashMap[hex.EncodeToString(sum[0:32])] = true
			}
		}
	case "createmagickfile":
		dir1 := os.Args[2]
		dir2 := os.Args[3]

		files, err := os.ReadDir(dir1)

		if err != nil {
			log.Fatal(err)
		}

		fNamesIdHashMap := map[string]bool{}

		for _, f := range files {
			fName := f.Name()
			fNamesIdHashMap[fName[9:43]] = true
		}

		var magickFileData []byte

		for _, f := range files {
			fName := f.Name()
			if fName[42:43] == "0" && (!fNamesIdHashMap[fName[9:41]+"_1"]) {
				magickFileData = append(magickFileData, []byte("magick"+" "+fName+" "+fName[0:41]+"_1.png")...)
				magickFileData = append(magickFileData, 0x0a)
			}
		}

		os.WriteFile(path.Join(dir2, "magickfile"), magickFileData, 0666)

	default:
		fmt.Println("Invaild Operation")
	}
}
