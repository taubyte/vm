package service

import (
	"context"
)

func (r *runtime) Close(ctx context.Context) error {
	if err := r.primitive.Close(ctx); err != nil {
		return err
	}

	close(r.wasiStartDone)

	return nil
}
