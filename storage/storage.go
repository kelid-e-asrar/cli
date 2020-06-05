package storage

//PassageEntry each password in storage
type PassageEntry struct {
	Name     string `json:"name,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

//PassageStorage is the interface that each storage should implement.
type PassageStorage interface {
	Set(entry *PassageEntry) error
	Get(name string) (*PassageEntry, error)
	Close() error
}
