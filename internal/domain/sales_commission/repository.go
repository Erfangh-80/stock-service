package salescommission

type Repository interface {
	Save(sc *SalesCommission) error
	FindByID(id int64) (*SalesCommission, error)
	Delete(id int64) error
}
