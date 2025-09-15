package services

import "context"

// @provider
type test struct{}

func (t *test) Test(ctx context.Context) (string, error) {
	return "Test", nil
}
