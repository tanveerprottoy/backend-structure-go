package product

type CreateDTO struct {
	Name        string
	Description *string
	CreatedAt   int64
	UpdatedAt   int64
}

type UpdateDTO struct {
	Name        string
	Description *string
	IsArchived  bool
	UpdatedAt   int64
}
