package category

type Repository interface {
	Save(c *Category) error
	FindByID(id int64) (*Category, error)
	FindByParentID(parentID int64) ([]*Category, error)
	FindAll() ([]*Category, error)
	Delete(id int64) error
}
