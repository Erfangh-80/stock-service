package store

type Repository interface {
	Save(store *Store) error
	FindByID(id int64) (*Store, error)
	FindAll() ([]*Store, error)
	Delete(id int64) error
}
