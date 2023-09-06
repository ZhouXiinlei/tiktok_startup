package logic

import (
	"context"
	"tikstart/common/model"
	"tikstart/common/utils"
	"tikstart/rpc/video/internal/svc"
	"tikstart/rpc/video/internal/union"
	"tikstart/rpc/video/video"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVideoListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVideoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoListLogic {
	return &GetVideoListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetVideoListLogic) GetVideoList(in *video.GetVideoListRequest) (*video.GetVideoListResponse, error) {
	var videos []model.Video
	err := l.svcCtx.DB.
		Where("created_at < ?", time.Unix(in.LatestTime, 0)).
		Order("created_at desc").
		Limit(int(in.Num)).
		Find(&videos).
		Error
	if err != nil {
		return nil, utils.InternalWithDetails("error querying video list", err)
	}

	var videoList []*video.VideoInfo
	for _, v := range videos {
		videoInfo, err := union.GetVideoInfoById(l.svcCtx.DB, l.svcCtx.RDS, v.VideoId)
		if err != nil {
			return nil, err
		}
		videoList = append(videoList, videoInfo)
	}
	var nextTime int64 = 0
	if len(videoList) != 0 {
		nextTime = videoList[len(videoList)-1].CreateTime
	}
	return &video.GetVideoListResponse{
		VideoList: videoList,
		NextTime:  nextTime,
	}, nil
}
