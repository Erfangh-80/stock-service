package salescommission

type CategoryCommissionRuleFilter struct {
	CategoryID *int32
	IsActive   *bool
	Page       int
	Limit      int
}

type CategoryCommissionRuleRepository interface {
	Save(rule *CategoryCommissionRule) error
	FindByID(id int64) (*CategoryCommissionRule, error)
	FindAll(filter CategoryCommissionRuleFilter) ([]*CategoryCommissionRule, int, error)
	Delete(id int64) error
}
