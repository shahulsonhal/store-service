package data

type Repo interface {
	GetStore(max int, country string) ([]StoreDetails, error)
}

type DB struct {
	// TODO: db struct for operations
}

func NewRepo(inMemStorage bool) Repo {
	if inMemStorage {
		s()
		return Store
	}
	return &DB{}
}

func (db *DB) GetStore(int, string) ([]StoreDetails, error) { return nil, nil }
