package utils

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoPaginate struct {
	limit int64
	page  int64
}

func NewMongoPagination(limit, page int) *MongoPaginate {
	return &MongoPaginate{
		limit: int64(limit),
		page:  int64(page),
	}
}

func (mp *MongoPaginate) GetPaginatedOpts(field string, order int) *options.FindOptions {
	l := mp.limit
	skip := mp.page*mp.limit - mp.limit
	fOpt := options.FindOptions{
		Sort:  bson.D{{Key: field, Value: order}},
		Limit: &l,
		Skip:  &skip,
	}
	return &fOpt
}
