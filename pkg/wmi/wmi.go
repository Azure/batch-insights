package wmi

import (
	"context"
	"time"

	wmi_lib "github.com/StackExchange/wmi"
)

const Timeout = 3 * time.Second

func CreateQuery(src interface{}, where string) string {
	return wmi_lib.CreateQuery(src, where)
}

func QueryWithContext(ctx context.Context, query string, dst interface{}, connectServerArgs ...interface{}) error {
	if _, ok := ctx.Deadline(); !ok {
		ctxTimeout, cancel := context.WithTimeout(ctx, Timeout)
		defer cancel()
		ctx = ctxTimeout
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- wmi_lib.Query(query, dst, connectServerArgs...)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}
