package brand

import "time"

type Brand struct {
	ID        int64
	Name      string
	Slug      string
	LogoFileID *int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewBrand(name, slug string) (*Brand, error) {
	if name == "" {
		return nil, ErrNameRequired
	}
	if slug == "" {
		return nil, ErrSlugRequired
	}

	now := time.Now()
	return &Brand{
		Name:      name,
		Slug:      slug,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (b *Brand) UpdateName(name string) error {
	if name == "" {
		return ErrNameRequired
	}
	b.Name = name
	b.UpdatedAt = time.Now()
	return nil
}

func (b *Brand) UpdateSlug(slug string) error {
	if slug == "" {
		return ErrSlugRequired
	}
	b.Slug = slug
	b.UpdatedAt = time.Now()
	return nil
}

func (b *Brand) UpdateLogo(fileID int64) {
	b.LogoFileID = &fileID
	b.UpdatedAt = time.Now()
}

func (b *Brand) RemoveLogo() {
	b.LogoFileID = nil
	b.UpdatedAt = time.Now()
}
