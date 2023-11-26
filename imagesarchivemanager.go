package main

import (
	"crypto/rand"
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

func timeForImgName() string {
	t := time.Now().String()
	return t[:4] + t[5:7] + t[8:10]
}

func main() {
	op := os.Args[1]
	path1 := os.Args[2]
	switch op {
	case "renamefiles":
		var magickPngFileData []byte

		files, err := os.ReadDir(path1)

		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			fileName := file.Name()

			imgId := newImgId()

			newFileName := "src" + "_" + imgId + path.Ext(fileName)

			magickPngFileData = append(magickPngFileData, []byte("magick"+" "+newFileName+" "+"$(date +%y%m%d)"+"_"+imgId+".png")...)
			magickPngFileData = append(magickPngFileData, 0x0a)

			magickPngFileData = append(magickPngFileData, []byte("rm"+" "+newFileName)...)
			magickPngFileData = append(magickPngFileData, 0x0a)

			data, err := os.ReadFile(path.Join(path1, fileName))
			if err != nil {
				log.Fatal(err)
			}

			os.WriteFile(path.Join(path1, newFileName), data, 0666)

			fmt.Println("Copied " + fileName + " to " + newFileName)
		}

		os.WriteFile(path.Join(path1, "magickforpng"), magickPngFileData, 0666)
	}
}
