package initloader

import (
	"archive/zip"
	"bufio"
	"log"
	"os"
	"strconv"

	jsoniter "github.com/json-iterator/go"

	"github.com/server-may-cry/highloadcup_18/accounthelper"
	"github.com/server-may-cry/highloadcup_18/db"
	"github.com/server-may-cry/highloadcup_18/structures"
)

type FileContent struct {
	Accounts []structures.Account `json:"accounts"`
}

func Load(db *db.Storage, path string) int {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	file, err := os.Open(path + "/options.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	if err != nil {
		log.Fatal(err)
	}
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	time, err := strconv.Atoi(lines[0])
	if err != nil {
		log.Fatal(err)
	}
	accounthelper.CurrentTime = time

	r, err := zip.OpenReader(path + "/data.zip")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	var i int
	for _, f := range r.File {
		func() {
			reader, err := f.Open()
			if err != nil {
				log.Fatal(err)
			}
			defer reader.Close()

			var accounts FileContent
			err = json.NewDecoder(reader).Decode(&accounts)
			if err != nil {
				log.Fatal(err)
			}

			for _, account := range accounts.Accounts {
				i++
				_ = db.AddAccount(account)
			}
		}()
	}
	return i
}
