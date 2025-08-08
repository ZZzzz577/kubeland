package data

import (
	"api-server/api/v1/common"
	"context"
)

type Pageable[T any, R any] interface {
	Count(ctx context.Context) (int, error)
	Limit(limit int) T
	Offset(offset int) T
	All(ctx context.Context) ([]R, error)
}

func Page[T Pageable[T, R], R any](ctx context.Context, query Pageable[T, R], pageRequest *common.Page) (*common.Pagination, []R, error) {
	var current uint32 = 1
	var size uint32 = 20
	if pageRequest != nil && pageRequest.Current > 0 {
		current = pageRequest.Current
	}
	if pageRequest != nil && pageRequest.Size > 0 {
		size = pageRequest.Size
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, nil, err
	}
	pagination := &common.Pagination{
		Current:  current,
		Size:  size,
		Total: uint32(total),
	}
	if total <= 0 {
		return pagination, []R{}, nil
	}
	results, err := query.
		Limit(int(size)).
		Offset(int(current-1) * int(size)).
		All(ctx)
	if err != nil {
		return nil, nil, err
	}
	return pagination, results, nil
}
