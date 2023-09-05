package logic

//
//import (
//	"context"
//	"github.com/zeromicro/go-zero/core/logx"
//	"tikstart/rpc/video/internal/cache"
//	"tikstart/rpc/video/internal/svc"
//	"tikstart/rpc/video/video"
//)
//
//type GetCountByIdLogic struct {
//	ctx    context.Context
//	svcCtx *svc.ServiceContext
//	logx.Logger
//}
//
//func NewGetCountByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCountByIdLogic {
//	return &GetCountByIdLogic{
//		ctx:    ctx,
//		svcCtx: svcCtx,
//		Logger: logx.WithContext(ctx),
//	}
//}
//
//func (l *GetCountByIdLogic) GetCountById(in *video.GetCountByIdRequest) (*video.GetCountByIdResponse, error) {
//
//	var workCount, totalFavorited, userFavoriteCount int64
//	workCount, err := cache.GetWorkCount(l.svcCtx, in.UserId)
//	if err != nil {
//		logx.WithContext(l.ctx).Error(err)
//		return nil, err
//	}
//	totalFavorited, err = cache.GetTotalFavorited(l.svcCtx, in.UserId)
//	if err != nil {
//		logx.WithContext(l.ctx).Error(err)
//		return nil, err
//	}
//	userFavoriteCount, err = cache.GetUserFavoriteCount(l.svcCtx, in.UserId)
//	if err != nil {
//		logx.WithContext(l.ctx).Error(err)
//		return nil, err
//	}
//	return &video.GetCountByIdResponse{
//		TotalFavorited:    totalFavorited,
//		WorkCount:         workCount,
//		UserFavoriteCount: userFavoriteCount,
//	}, nil
//}
