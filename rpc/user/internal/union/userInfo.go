package union

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"tikstart/common"
	"tikstart/common/model"
	"tikstart/rpc/user/user"
)

func GetUserInfoById(db *gorm.DB, rds *redis.Redis, userId int64) (*user.UserInfo, error) {
	userRecord := model.User{}
	err := db.Where("user_id = ?", userId).First(&userRecord).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrUserNotFound.Err()
		} else {
			return nil, status.Error(codes.Internal, err.Error())
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
			return nil, status.Error(codes.Internal, err.Error())
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
