package storage

import (
	"encoding/json"
	"io/ioutil"
)

//JSONFile storage for passage.
type JSONFile struct {
	virtualJson map[string]*PassageEntry
	path        string
}

const (
	_DefaultJSONPath = "passage.json"
)

func NewJSONFile(path string) (*JSONFile, error) {
	if path == "" {
		path = _DefaultJSONPath
	}
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	vj := map[string]*PassageEntry{}
	err = json.Unmarshal(bs, &vj)
	if err != nil {
		return nil, err
	}
	return &JSONFile{vj, path}, nil
}

func (j *JSONFile) Set(entry *PassageEntry) error {
	j.virtualJson[entry.Name] = entry
	return nil
}

func (j *JSONFile) Get(name string) (*PassageEntry, error) {
	return j.virtualJson[name], nil
}

func (j *JSONFile) Close() error {
	bs, err := json.Marshal(j.virtualJson)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(j.path, bs, 0644)
}
