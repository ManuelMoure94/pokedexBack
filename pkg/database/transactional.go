package database

import (
	"context"

	"gorm.io/gorm"
)

type TransactionalKey struct{}

var TransactionalCtxKey = TransactionalKey{}

func Transactional(ctx context.Context, fn func(ctx context.Context) error) error {
	value := ctx.Value(TransactionalCtxKey)
	if value != nil {
		_, ok := value.(*gorm.DB)
		if ok {
			return fn(ctx)
		}
	}

	return orm.Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, TransactionalCtxKey, tx))
	})
}

func Orm(ctx context.Context) *gorm.DB {
	value := ctx.Value(TransactionalCtxKey)
	if value != nil {
		orm, ok := value.(*gorm.DB)
		if ok {
			return orm
		}
	}

	return orm
}
