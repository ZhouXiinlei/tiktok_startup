package union

import (
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
	"tikstart/common"
	"tikstart/common/model"
	"tikstart/common/utils"
	"tikstart/rpc/user/user"
)

func GetUserInfoById(db *gorm.DB, rds *redis.Redis, userId int64) (*user.UserInfo, error) {
	userRecord := model.User{}
	err := db.Where("user_id = ?", userId).First(&userRecord).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrUserNotFound.Err()
		} else {
			return nil, utils.InternalWithDetails("error querying user", err)
		}
	}

	followingCount, err := PickUserCounts(db, rds, userId, "following_count", userRecord.FollowingCount)
	if err != nil {
		return nil, err
	}
	followerCount, err := PickUserCounts(db, rds, userId, "follower_count", userRecord.FollowerCount)
	if err != nil {
		return nil, err
	}
	totalFavorited, err := PickUserCounts(db, rds, userId, "total_favorited", userRecord.TotalFavorited)
	if err != nil {
		return nil, err
	}
	favoriteCount, err := PickUserCounts(db, rds, userId, "favorite_count", userRecord.FavoriteCount)
	if err != nil {
		return nil, err
	}

	return &user.UserInfo{
		UserId:         userRecord.UserId,
		Username:       userRecord.Username,
		FollowingCount: followingCount,
		FollowerCount:  followerCount,
		TotalFavorited: totalFavorited,
		WorkCount:      userRecord.WorkCount,
		FavoriteCount:  favoriteCount,
		Password:       userRecord.Password,
		CreatedAt:      userRecord.CreatedAt.Unix(),
		UpdatedAt:      userRecord.UpdatedAt.Unix(),
	}, nil
}

func GetUserInfoByName(db *gorm.DB, rds *redis.Redis, username string) (*user.UserInfo, error) {
	userRecord := model.User{}
	err := db.Where("username = ?", username).First(&userRecord).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrUserNotFound.Err()
		} else {
			return nil, utils.InternalWithDetails("error querying user", err)
		}
	}

	followingCount, err := PickUserCounts(db, rds, userRecord.UserId, "following_count", userRecord.FollowingCount)
	if err != nil {
		return nil, err
	}
	followerCount, err := PickUserCounts(db, rds, userRecord.UserId, "follower_count", userRecord.FollowerCount)
	if err != nil {
		return nil, err
	}
	totalFavorited, err := PickUserCounts(db, rds, userRecord.UserId, "total_favorited", userRecord.TotalFavorited)
	if err != nil {
		return nil, err
	}
	favoriteCount, err := PickUserCounts(db, rds, userRecord.UserId, "favorite_count", userRecord.FavoriteCount)
	if err != nil {
		return nil, err
	}

	return &user.UserInfo{
		UserId:         userRecord.UserId,
		Username:       userRecord.Username,
		FollowingCount: followingCount,
		FollowerCount:  followerCount,
		TotalFavorited: totalFavorited,
		WorkCount:      userRecord.WorkCount,
		FavoriteCount:  favoriteCount,
		Password:       userRecord.Password,
		CreatedAt:      userRecord.CreatedAt.Unix(),
		UpdatedAt:      userRecord.UpdatedAt.Unix(),
	}, nil
}

func GetManyUserInfos(db *gorm.DB, rds *redis.Redis, condition *gorm.DB) ([]*user.UserInfo, error) {
	userIdList := make([]int64, 0)
	err := db.
		Model(&model.User{}).
		Where("user_id in (?)", condition).
		Select("user_id").
		Find(&userIdList).
		Error
	if err != nil {
		return nil, utils.InternalWithDetails("error querying many users", err)
	}

	order := make(map[int64]int, len(userIdList))
	userInfoList, err := mr.MapReduce(func(source chan<- int64) {
		for i, userId := range userIdList {
			source <- userId
			order[userId] = i
		}
	}, func(item int64, writer mr.Writer[*user.UserInfo], cancel func(error)) {
		userInfo, err := GetUserInfoById(db, rds, item)
		if err != nil {
			cancel(err)
			return
		}
		writer.Write(userInfo)
	}, func(pipe <-chan *user.UserInfo, writer mr.Writer[[]*user.UserInfo], cancel func(error)) {
		list := make([]*user.UserInfo, len(userIdList))
		for item := range pipe {
			userInfo := item
			i, _ := order[item.UserId]
			list[i] = userInfo
		}
		writer.Write(list)
	})

	if err != nil {
		return nil, err
	}
	return userInfoList, nil
}
