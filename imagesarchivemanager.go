package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func newImgTimePrefix() string {
	return time.Now().Format("20060102")
}

func main() {
	op := os.Args[1]

	switch op {
	case "archive":
		filePathsArg := os.Args[2]
		archiveDirArg := os.Args[3]

		files, err_1 := os.ReadDir(archiveDirArg)
		if err_1 != nil {
			log.Fatal(err_1)
		}

		existingFilesHash := map[string]bool{}

		for _, f := range files {
			fName := f.Name()
			existingFilesHash[fName[9:25]] = true
		}

		data, err_2 := os.ReadFile(filePathsArg)
		if err_2 != nil {
			log.Fatal(err_2)
		}

		var withoutNewLinesData []byte

		for _, c := range data {
			if c == 0x0a {
				withoutNewLinesData = append(withoutNewLinesData, 0x3a)
			} else if c == 0x3a {
				fmt.Println("error: found colon in file paths")
				os.Exit(1)
			} else {
				withoutNewLinesData = append(withoutNewLinesData, c)
			}
		}

		filePaths := strings.Split(string(withoutNewLinesData), ":")

		var c1 int
		var c2 int

		for _, p := range filePaths {
			ext := path.Ext(p)
			if ext == ".jpeg" || ext == ".jpg" || ext == ".png" {
				data, err := os.ReadFile(p)
				if err != nil {
					log.Fatal(err)
				}
				sum := sha256.Sum256(data)
				sumHex := hex.EncodeToString(sum[0:8])
				if !existingFilesHash[sumHex] {
					os.WriteFile(path.Join(archiveDirArg, newImgTimePrefix()+"_"+sumHex+ext), data, 0666)
					existingFilesHash[sumHex] = true
					c2++
				}
				c1++
			}
		}
		fmt.Println("Total files:", strconv.FormatInt(int64(len(filePaths)), 10))
		fmt.Println("Images found:", strconv.FormatInt(int64(c1), 10))
		fmt.Println("Images archived:", strconv.FormatInt(int64(c2), 10))

	case "check":
		archiveDirArg := os.Args[2]

		files, err_1 := os.ReadDir(archiveDirArg)
		if err_1 != nil {
			log.Fatal(err_1)
		}

		var c int

		for _, f := range files {
			data, err := os.ReadFile(path.Join(archiveDirArg, f.Name()))
			if err != nil {
				log.Fatal(err)
			}
			sum := sha256.Sum256(data)
			sumHex := hex.EncodeToString(sum[0:8])

			if f.Name()[9:25] == sumHex {
				fmt.Println("OK", f.Name())
			} else {
				fmt.Println("ERROR:", f.Name())
				c++
			}

		}

		fmt.Println("Total files checked:", strconv.FormatInt(int64(len(files)), 10))
		fmt.Println("Total errors:", strconv.FormatInt(int64(c), 10))
	case "magickcheck":
		archiveDirArg := os.Args[2]
		magickFileDirArg := os.Args[3]

		files, err_1 := os.ReadDir(archiveDirArg)
		if err_1 != nil {
			log.Fatal(err_1)
		}

		var d []byte
		for _, f := range files {
			d = append(d, []byte("magick "+f.Name()+" /mnt/t/temp.rgba")...)
			d = append(d, 0x0a)
			d = append(d, []byte("rm /mnt/t/temp.rgba")...)
			d = append(d, 0x0a)
		}
		os.WriteFile(path.Join(magickFileDirArg, "magickcheckbashfile"), d, 0666)

	default:
		fmt.Println("Invaild Operation")
	}
}
