package union

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
	"tikstart/common/model"
	"tikstart/common/utils"
	"tikstart/rpc/user/internal/cache"
)

func IsFollow(db *gorm.DB, rds *redis.Redis, userId int64, targetId int64) (bool, error) {
	val, err := rds.Get(cache.GenFollowKey(userId, targetId))
	if err != nil {
		return false, utils.InternalWithDetails("(redis)error getting follow relation", err)
	}
	if v, hit := cache.TrueOrFalse(val); hit {
		return v, nil
	}

	var count int64
	err = db.
		Model(&model.Follow{}).
		Where("follower_id = ? AND followed_id = ?", userId, targetId).
		Count(&count).
		Error
	if err != nil {
		return false, utils.InternalWithDetails("(mysql)error querying follow relation", err)
	}

	go func() {
		err = rds.Set(cache.GenFollowKey(userId, targetId), cache.YesOrNo(count == 1))
		if err != nil {
			logx.Errorf("(redis)error setting follow relation", err)
		}
	}()
	return count == 1, nil
}
