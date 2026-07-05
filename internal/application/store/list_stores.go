package store

import (
	"stock-service/internal/domain/store"
)

type ListStoresInput struct {
	UserID *int64
	Status *string
	Page   int
	Limit  int
}

type ListStoresOutput struct {
	Stores []*store.Store
	Total  int
	Page   int
	Limit  int
}

type ListStoresUseCase interface {
	Execute(input ListStoresInput) (*ListStoresOutput, error)
}

type listStoresInteractor struct {
	repo store.Repository
}

func NewListStoresUseCase(repo store.Repository) ListStoresUseCase {
	return &listStoresInteractor{repo: repo}
}

func (uc *listStoresInteractor) Execute(input ListStoresInput) (*ListStoresOutput, error) {
	all, err := uc.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var filtered []*store.Store
	for _, s := range all {
		if input.UserID != nil && s.UserID != *input.UserID {
			continue
		}
		if input.Status != nil && string(s.Status) != *input.Status {
			continue
		}
		filtered = append(filtered, s)
	}

	total := len(filtered)

	page := input.Page
	if page < 1 {
		page = 1
	}
	limit := input.Limit
	if limit < 1 {
		limit = 20
	}

	start := (page - 1) * limit
	if start >= total {
		return &ListStoresOutput{
			Stores: []*store.Store{},
			Total:  total,
			Page:   page,
			Limit:  limit,
		}, nil
	}

	end := start + limit
	if end > total {
		end = total
	}

	return &ListStoresOutput{
		Stores: filtered[start:end],
		Total:  total,
		Page:   page,
		Limit:  limit,
	}, nil
}
