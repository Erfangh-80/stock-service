package brand

type Repository interface {
	Save(b *Brand) error
	FindByID(id int64) (*Brand, error)
	FindAll() ([]*Brand, error)
	Delete(id int64) error
}
