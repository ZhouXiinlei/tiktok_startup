package cache

import (
	"github.com/zeromicro/go-zero/core/logx"
	"tikstart/common/model"
	"tikstart/common/utils"
	"tikstart/rpc/video/internal/svc"
)

func IsFavorite(svcCtx *svc.ServiceContext, userId int64, videoId int64) (bool, error) {
	val, err := svcCtx.RDS.Get(GenFavoriteKey(userId, videoId))
	if err != nil {
		return false, utils.InternalWithDetails("(redis)error getting favorite relation", err)
	}
	if v, hit := TrueOrFalse(val); hit {
		return v, nil
	}

	var count int64
	err = svcCtx.DB.
		Model(&model.Favorite{}).
		Where("user_id = ? AND video_id = ?", userId, videoId).
		Count(&count).
		Error
	if err != nil {
		return false, utils.InternalWithDetails("(mysql)error querying favorite relation", err)
	}

	go func() {
		err = svcCtx.RDS.Set(GenFavoriteKey(userId, videoId), YesOrNo(count == 1))
		if err != nil {
			logx.Errorf("(redis)error setting favorite relation", err)
		}
	}()
	return count == 1, nil
}
