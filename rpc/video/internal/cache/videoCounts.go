package cache

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
)

func PickVideoCounts(db *gorm.DB, rds *redis.Redis, videoId int64, field string, dbCount int64) (int64, error) {
	score, err := rds.Zscore(GenVideoCountsKey(field), strconv.FormatInt(videoId, 10))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return dbCount, nil
		} else {
			return 0, utils.InternalWithDetails("(redis)error querying cache status", err)
		}
	}
	go ManageCache(db, rds, videoId, field)
	return score, nil
}

func ModifyVideoCounts(db *gorm.DB, rds *redis.Redis, videoId int64, field string, delta int64) error {
	// check cache first
	_, err := rds.Zscore(GenVideoCountsKey(field), strconv.FormatInt(videoId, 10))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// cache miss, update database
			err = db.
				Clauses(clause.Locking{Strength: "UPDATE"}).
				Model(&model.Video{}).
				Where("video_id = ?", videoId).
				UpdateColumn(field, gorm.Expr(fmt.Sprintf("%s + ?", field), delta)).
				Error
			if err != nil {
				return utils.InternalWithDetails("error adding video_counts", err)
			}
		} else {
			return utils.InternalWithDetails("(redis)error querying cache status", err)
		}
	}
	// cache hit in redis, then update cache only
	go func() {
		_, err = rds.Zincrby(GenVideoCountsKey(field), delta, strconv.FormatInt(videoId, 10))
		if err != nil {
			logx.Errorf("(redis)error incrementing video_counts", err)
		}
	}()

	// async manage cache
	go ManageCache(db, rds, videoId, field)

	return nil
}

func ManageCache(db *gorm.DB, rds *redis.Redis, videoId int64, field string) {
	// set or add heat
	heat, err := rds.Incr(GenVideoHeatKey(videoId))
	if err != nil {
		logx.Errorf("(redis)error incrementing video_heat", err)
	}
	// reset expire time
	err = rds.Expire(GenVideoHeatKey(videoId), 300)
	if err != nil {
		logx.Errorf("(redis)error resetting video_heat expire", err)
	}

	// popular video, cache field
	if heat > 0 {
		_, err := rds.Zscore(GenVideoCountsKey(field), strconv.FormatInt(videoId, 10))
		if err != nil {
			if errors.Is(err, redis.Nil) {
				// not cached, setting new cache
				var count int64
				err = db.
					Model(&model.Video{}).
					Where("video_id = ?", videoId).
					Select(field).
					First(&count).
					Error
				if err != nil {
					logx.Errorf(fmt.Sprintf("(mysql)error querying %s", field), err)
				}

				_, err = rds.Zadd(GenVideoCountsKey(field), count, strconv.FormatInt(videoId, 10))
				if err != nil {
					logx.Errorf(fmt.Sprintf("(redis)error setting %s", field), err)
				}
			} else {
				logx.Errorf("(redis)error querying cache status", err)
			}
		}

	}
}
