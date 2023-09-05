package union

import (
	"github.com/zeromicro/go-zero/core/logx"
	"tikstart/common/model"
	"tikstart/common/utils"
	"tikstart/rpc/user/internal/cache"
	"tikstart/rpc/user/internal/svc"
)

func IsFollow(svcCtx *svc.ServiceContext, userId int64, targetId int64) (bool, error) {
	val, err := svcCtx.RDS.Get(cache.GenFollowKey(userId, targetId))
	if err != nil {
		return false, utils.InternalWithDetails("(redis)error getting follow relation", err)
	}
	if v, hit := cache.TrueOrFalse(val); hit {
		return v, nil
	}

	var count int64
	err = svcCtx.DB.
		Model(&model.Follow{}).
		Where("follower_id = ? AND followed_id = ?", userId, targetId).
		Count(&count).
		Error
	if err != nil {
		return false, utils.InternalWithDetails("(mysql)error querying follow relation", err)
	}

	go func() {
		err = svcCtx.RDS.Set(cache.GenFollowKey(userId, targetId), cache.YesOrNo(count == 1))
		if err != nil {
			logx.Errorf("(redis)error setting follow relation", err)
		}
	}()
	return count == 1, nil
}
