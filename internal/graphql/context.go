package graphql

import "context"

type ReqContextKey string

const DbContextKey = ReqContextKey("db")
const RequestCachesKey = ReqContextKey("requestCaches")

type CacheKey string
type RequestCache map[CacheKey]map[int]interface{}

const TransactionsCacheKey = CacheKey("transactions")
const DetailsCacheKey = CacheKey("details")

func (key CacheKey) getCache(ctx context.Context) map[int]interface{} {
	caches := ctx.Value(RequestCachesKey).(RequestCache)
	if caches[key] == nil {
		caches[key] = make(map[int]interface{})
	}
	return caches[key]
}
