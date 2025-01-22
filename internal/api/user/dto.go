package user

type CreateDTO struct {
	Name      string
	Address   *string
	CreatedAt int64
	UpdatedAt int64
}

type UpdateDTO struct {
	Name       string
	Address    *string
	IsArchived bool
	UpdatedAt  int64
}
