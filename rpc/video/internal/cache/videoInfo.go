package cache

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"tikstart/common"
	"tikstart/common/model"
	"tikstart/rpc/video/video"
)

func GetVideoInfoById(db *gorm.DB, rds *redis.Redis, videoId int64) (*video.VideoInfo, error) {
	videoRecord := model.Video{}
	err := db.Where("video_id = ?", videoId).First(&videoRecord).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrUserNotFound.Err()
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	CommentCount, err := PickVideoCounts(db, rds, videoId, "comment_count", videoRecord.CommentCount)
	if err != nil {
		return nil, err
	}
	FavoriteCount, err := PickVideoCounts(db, rds, videoId, "favorite_count", videoRecord.FavoriteCount)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return &video.VideoInfo{
		Id:            videoRecord.VideoId,
		AuthorId:      videoRecord.AuthorId,
		Title:         videoRecord.Title,
		PlayUrl:       videoRecord.PlayUrl,
		CoverUrl:      videoRecord.CoverUrl,
		FavoriteCount: FavoriteCount,
		CommentCount:  CommentCount,
		CreateTime:    videoRecord.CreatedAt.Unix(),
	}, nil
}
