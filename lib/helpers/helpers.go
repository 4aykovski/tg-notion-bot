package helpers

import "fmt"

func ErrWrapIfNotNil(msg string, err error) error {
	if err != nil {
		return fmt.Errorf("%s: %w", msg, err)
	}
	return nil
}
