package cache

import (
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
	"tikstart/common/model"
	"tikstart/common/utils"
	"tikstart/rpc/video/internal/svc"
)

func GetWorkCount(svcCtx *svc.ServiceContext, userId int64) (int64, error) {
	val, err := svcCtx.RDS.Get(GenWorkCountKey(userId))
	if err != nil {
		return 0, utils.InternalWithDetails("(redis)error getting workCount", err)
	}
	var workCount int64
	if val != "" {
		// string到int64
		workCount, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			_, err = svcCtx.RDS.Del(GenUserFavoriteCountKey(userId))
		}
		return workCount, nil
	}
	err = svcCtx.DB.Model(&model.Video{}).
		Where("author_id = ?", userId).
		Count(&workCount).Error
	if err != nil {
		return 0, utils.InternalWithDetails("error querying WorkCount ", err)
	}
	workCache := strconv.FormatInt(workCount, 10)
	go func() {
		err = svcCtx.RDS.Set(GenWorkCountKey(userId), workCache)
		if err != nil {
			logx.Errorf("(redis)error setting WorkCount", err)
		}
	}()
	return workCount, nil
}
func GetTotalFavorited(svcCtx *svc.ServiceContext, userId int64) (int64, error) {
	val, err := svcCtx.RDS.Get(GenTotalFavoritedKey(userId))
	if err != nil {
		return 0, utils.InternalWithDetails("(redis)error getting totalFavorited relation", err)
	}
	var totalFavorited int64
	if val != "" {
		// string到int64
		totalFavorited, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			_, err = svcCtx.RDS.Del(GenTotalFavoritedKey(userId))
		}
		return totalFavorited, nil
	}
	err = svcCtx.DB.Model(&model.Favorite{}).
		Where("video_id in (?)",
			svcCtx.DB.Model(&model.Video{}).
				Select("video_id").
				Where("author_id = ?", userId)).
		Count(&totalFavorited).Error
	if err != nil {
		return 0, utils.InternalWithDetails("error querying totalFavorited ", err)
	}
	totalcache := strconv.FormatInt(totalFavorited, 10)
	go func() {
		err = svcCtx.RDS.Set(GenTotalFavoritedKey(userId), totalcache)
		if err != nil {
			logx.Errorf("(redis)error setting totalFavorited", err)
		}
	}()
	return totalFavorited, nil
}
func GetUserFavoriteCount(svcCtx *svc.ServiceContext, userId int64) (int64, error) {
	val, err := svcCtx.RDS.Get(GenUserFavoriteCountKey(userId))
	if err != nil {
		return 0, utils.InternalWithDetails("(redis)error getting userFavorite relation", err)
	}
	var userFavoriteCount int64
	if val != "" {
		// string到int64
		userFavoriteCount, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			_, err = svcCtx.RDS.Del(GenUserFavoriteCountKey(userId))
		}
		return userFavoriteCount, nil
	}
	err = svcCtx.DB.Model(&model.Favorite{}).
		Where("user_id = ?", userId).
		Count(&userFavoriteCount).Error
	if err != nil {
		return 0, utils.InternalWithDetails("error querying UserFavoriteCount ", err)
	}
	userFavoriteCache := strconv.FormatInt(userFavoriteCount, 10)
	go func() {
		err = svcCtx.RDS.Set(GenUserFavoriteCountKey(userId), userFavoriteCache)
		if err != nil {
			logx.Errorf("(redis)error setting userFavorite", err)
		}
	}()
	return userFavoriteCount, nil
}
