package user

type CreateDTO struct {
	Name    string
	Address *string
}

type UpdateDTO struct {
	Name       string
	Address    *string
	IsArchived bool
	UpdatedAt  int64
}
