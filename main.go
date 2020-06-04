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

type entry struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (e *entry) AsBytes() ([]byte, error) {
	//TODO: encrypt data
	return json.Marshal(e)
}
func commandHandler(ctx *context) func(string) {
	return func(input string) {
		cmds := strings.Split(input, " ")
		switch cmds[0] {
		case "get":
			name := cmds[1]
			ctx.db.View(func(tx *bbolt.Tx) error {
				payload := tx.Bucket(bucketName).Get([]byte(name))
				m := entry{}
				err := json.Unmarshal(payload, &m)
				if err != nil {
					return err
				}
				logrus.Printf("%+v", m)
				return nil
			})
		case "set":
			e := &entry{
				Name:     cmds[1],
				Username: cmds[2],
				Password: cmds[3],
			}
			ctx.db.Update(func(tx *bbolt.Tx) error {
				bs, err := e.AsBytes()
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
				err = bucket.Put([]byte(e.Name), bs)
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
