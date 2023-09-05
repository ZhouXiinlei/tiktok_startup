package union

import (
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"tikstart/common/model"
	"tikstart/common/utils"
	"tikstart/rpc/user/internal/cache"
)

func PickUserCounts(rds *redis.Redis, userId int64, field string, dbCount int64) (int64, error) {
	score, err := rds.Zscore(cache.GenUserCountsKey(field), strconv.FormatInt(userId, 10))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return dbCount, nil
		} else {
			return 0, utils.InternalWithDetails("(redis)error querying cache status", err)
		}
	}
	return score, nil
}

func ModifyUserCounts(db *gorm.DB, rds *redis.Redis, userId int64, field string, delta int64) error {
	// check cache first
	_, err := rds.Zscore(cache.GenUserCountsKey(field), strconv.FormatInt(userId, 10))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// cache miss, update database
			err = db.
				Clauses(clause.Locking{Strength: "UPDATE"}).
				Model(&model.User{}).
				Where("user_id = ?", userId).
				UpdateColumn(field, gorm.Expr(fmt.Sprintf("%s + ?", field), delta)).
				Error
			if err != nil {
				return utils.InternalWithDetails("error modifying user_counts", err)
			}
		} else {
			return utils.InternalWithDetails("(redis)error querying cache status", err)
		}
	}
	// cache hit in redis, then update cache only
	go func() {
		_, err = rds.Zincrby(cache.GenUserCountsKey(field), delta, strconv.FormatInt(userId, 10))
		if err != nil {
			logx.Errorf("(redis)error modifying user_counts", err)
		}
	}()

	// async manage cache
	go ManageCache(db, rds, userId, field)

	return nil
}

func ManageCache(db *gorm.DB, rds *redis.Redis, userId int64, field string) {
	// set or add heat
	heat, err := rds.Incr(cache.GenUserHeatKey(userId))
	if err != nil {
		logx.Errorf("(redis)error incrementing user_heat", err)
	}
	// reset expire time
	err = rds.Expire(cache.GenUserHeatKey(userId), 300)
	if err != nil {
		logx.Errorf("(redis)error resetting user_heat expire", err)
	}

	// popular user, cache field
	if heat > 0 {
		_, err := rds.Zscore(cache.GenUserCountsKey(field), strconv.FormatInt(userId, 10))
		if err != nil {
			if errors.Is(err, redis.Nil) {
				// not cached, setting new cache
				var count int64
				err = db.
					Model(&model.User{}).
					Where("user_id = ?", userId).
					Select(field).
					First(&count).
					Error
				if err != nil {
					logx.Errorf(fmt.Sprintf("(mysql)error querying %s", field), err)
				}

				_, err = rds.Zadd(cache.GenUserCountsKey(field), count, strconv.FormatInt(userId, 10))
				if err != nil {
					logx.Errorf(fmt.Sprintf("(redis)error setting %s", field), err)
				}
			} else {
				logx.Errorf("(redis)error querying cache status", err)
			}
		}

	}
}
