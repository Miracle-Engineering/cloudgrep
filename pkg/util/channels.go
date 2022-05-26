package util

import "context"

// SendAllFromSlice iterates over `items` and sends each one over `out` unless `ctx` is cancelled.
// Returns ctx.Err() if not all items were written to the channel.
func SendAllFromSlice[T any](ctx context.Context, out chan<- T, items []T) error {
	for _, item := range items {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case out <- item:
		}
	}

	return nil
}
