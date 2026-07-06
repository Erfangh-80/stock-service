package category

import "time"

type Category struct {
	ID          int64
	Name        string
	Slug        string
	ParentID    *int64
	Description *string
	ImageFileID *int64
	SortOrder   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewCategory(name, slug string, parentID *int64) (*Category, error) {
	if name == "" {
		return nil, ErrNameRequired
	}
	if slug == "" {
		return nil, ErrSlugRequired
	}

	now := time.Now()
	return &Category{
		Name:      name,
		Slug:      slug,
		ParentID:  parentID,
		SortOrder: 0,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (c *Category) UpdateName(name string) error {
	if name == "" {
		return ErrNameRequired
	}
	c.Name = name
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Category) UpdateSlug(slug string) error {
	if slug == "" {
		return ErrSlugRequired
	}
	c.Slug = slug
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Category) UpdateDescription(desc string) {
	c.Description = &desc
	c.UpdatedAt = time.Now()
}

func (c *Category) UpdateImage(fileID int64) {
	c.ImageFileID = &fileID
	c.UpdatedAt = time.Now()
}

func (c *Category) UpdateSortOrder(order int) {
	c.SortOrder = order
	c.UpdatedAt = time.Now()
}

func (c *Category) Reparent(parentID int64) {
	c.ParentID = &parentID
	c.UpdatedAt = time.Now()
}

func (c *Category) RemoveParent() {
	c.ParentID = nil
	c.UpdatedAt = time.Now()
}
