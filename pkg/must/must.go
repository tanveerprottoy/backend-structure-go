package must

import "log"

func Must[T any](v T, err error) T {
	if err != nil {
		log.Fatalf("fatal error: %v", err)
	}

	return v
}
