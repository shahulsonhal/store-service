package data

type MockRepo struct {
	MockGetStore func(int, string) ([]StoreDetails, error)
}

func (f *MockRepo) GetStore(max int, country string) ([]StoreDetails, error) {
	return f.MockGetStore(max, country)
}
