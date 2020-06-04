package main

import (
	"encoding/json"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	"github.com/sirupsen/logrus"
	"go.etcd.io/bbolt"
)

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

var (
	bucketName = []byte("passwords")
)

func commandHandler(ctx *context) func(string) {
	return func(input string) {
		cmds := strings.Split(input, " ")
		switch cmds[0] {
		case "get":
			name := cmds[1]
			ctx.db.View(func(tx *bbolt.Tx) error {
				payload := tx.Bucket(bucketName).Get([]byte(name))
				logrus.Println(payload)
				return nil
			})
		case "set":
			name := cmds[1]
			username := cmds[2]
			password := cmds[3]
			ctx.db.Update(func(tx *bbolt.Tx) error {
				bs, err := json.Marshal(map[string]string{
					"username": username,
					"password": password,
					"name":     name,
				})
				if err != nil {
					return err
				}
				bucket := tx.Bucket(bucketName)
				if bucket == nil {
					bucket, err = tx.CreateBucketIfNotExists(bucketName)
					if err != nil {
						return err
					}
				}
				err = bucket.Put([]byte(name), bs)
				if err != nil {
					return err
				}
				return nil
			})
		}
	}
}

type context struct {
	db *bbolt.DB
}

func openDatabase() (*bbolt.DB, error) {
	return bbolt.Open("passage.db", 0600, nil)
}

func main() {
	db, err := openDatabase()
	if err != nil {
		logrus.Fatalln(err)
	}
	logrus.Infoln("Database Loaded")
	defer db.Close()
	ctx := &context{db: db}
	p := prompt.New(commandHandler(ctx), completer)

	p.Run()
}
