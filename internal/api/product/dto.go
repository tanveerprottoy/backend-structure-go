package product

type CreateDTO struct {
	Name        string
	Description string
}

type UpdateDTO struct {
	Name        string
	Description string
	IsArchived  bool
}
