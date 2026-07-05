package promotion

type Repository interface {
	Save(p *Promotion) error
	FindByID(id int64) (*Promotion, error)
	Delete(id int64) error
}
