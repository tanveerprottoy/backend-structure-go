package product

import "errors"

type Product struct {
	ID          string
	Name        string
	Description *string
	IsArchived  bool
	CreatedAt   int64
	UpdatedAt   int64
}

func NewProduct(id, name string, description string, createdAt, updatedAt int64) *Product {
	p := &Product{
		ID:        id,
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	p.SetDescription(description)
	return p
}

func (p Product) Validate() error {
	if p.Name == "" {
		return errors.New("name required")
	}
	return nil
}

func (p *Product) SetDescription(val string) {
	if val != "" {
		v := val
		p.Description = &v
	}
}
