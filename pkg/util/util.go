package util

func SetLimit(limit int) int {
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return limit
}

func CalculateOffset(limit, page int) int {
	if page < 1 {
		page = 1
	}
	limit = SetLimit(limit)
	return limit * (page - 1)
}
