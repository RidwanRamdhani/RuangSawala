package repositories

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type MatchmakingRepository struct {
	rdb *redis.Client
}

func NewMatchmakingRepository(rdb *redis.Client) *MatchmakingRepository {
	return &MatchmakingRepository{rdb: rdb}
}

const PoolKey = "matchmaking_pool"

func (r *MatchmakingRepository) AddToPool(ctx context.Context, userID int) error {
	return r.rdb.SAdd(ctx, PoolKey, userID).Err()
}

func (r *MatchmakingRepository) RemoveFromPool(ctx context.Context, userID int) error {
	return r.rdb.SRem(ctx, PoolKey, userID).Err()
}

func (r *MatchmakingRepository) GetAllCandidates(ctx context.Context) ([]string, error) {
	return r.rdb.SMembers(ctx, PoolKey).Result()
}

func (r *MatchmakingRepository) PopRandomCandidates(ctx context.Context, count int64) ([]string, error) {
	return r.rdb.SPopN(ctx, PoolKey, count).Result()
}

func (r *MatchmakingRepository) AddToPoolMulti(ctx context.Context, userIDs []int) error {
	if len(userIDs) == 0 {
		return nil
	}
	members := make([]interface{}, len(userIDs))
	for i, id := range userIDs {
		members[i] = id
	}
	return r.rdb.SAdd(ctx, PoolKey, members...).Err()
}
