package video

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"tikstart/common/utils"
	"tikstart/http/internal/logic/video"
	"tikstart/http/internal/svc"
	"tikstart/http/internal/types"
	"tikstart/http/schema"
	rpcvideo "tikstart/rpc/video/video"

	"github.com/google/uuid"
	"github.com/h2non/filetype"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func PublishVideoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PublishVideoRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		userclaims, _ := utils.ParseToken(req.Token, svcCtx.Config.JwtAuth.Secret)
		UserId := userclaims.UserId
		file, fileHeader, err := r.FormFile("data")
		if err != nil {
			httpx.Error(w, schema.ServerError{
				ApiError: schema.ApiError{
					StatusCode: 400,
					Code:       40005,
					Message:    "文件上传失败",
				},
				Detail: err,
			})
			return
		}
		tmpFile, err := fileHeader.Open()
		if err != nil {
			httpx.Error(w, schema.ServerError{
				ApiError: schema.ApiError{
					StatusCode: 400,
					Code:       40005,
					Message:    "文件上传失败",
				},
				Detail: err,
			})
			return
		}
		defer tmpFile.Close()
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, tmpFile); err != nil {
			httpx.Error(w, schema.ServerError{
				ApiError: schema.ApiError{
					StatusCode: 400,
					Code:       40005,
					Message:    "文件上传失败",
				},
				Detail: err,
			})
			return
		}
		if !filetype.IsVideo(buf.Bytes()) {
			httpx.Error(w, schema.ServerError{
				ApiError: schema.ApiError{
					StatusCode: 400,
					Code:       40006,
					Message:    "无效的文件类型",
				},
				Detail: err,
			})
			return
		}

		fileName := uuid.New().String() + filepath.Ext(fileHeader.Filename)
		keyName := "/tiktok/" + fileName
		_, err = svcCtx.TengxunyunClient.Object.Put(context.Background(), keyName, file, nil)
		if err != nil {
			httpx.Error(w, schema.ServerError{
				ApiError: schema.ApiError{
					StatusCode: 500,
					Code:       50010,
					Message:    "文件上传cos失败",
				},
				Detail: err,
			})
			return
		}
		cdnUrl, err := url.Parse(svcCtx.Config.CDNBaseURL)
		if err != nil {
			httpx.Error(w, schema.ServerError{
				ApiError: schema.ApiError{
					StatusCode: 500,
					Code:       50011,
					Message:    "获取视频链接失败",
				},
				Detail: err,
			})
			return
		}
		cdnUrl.Path = keyName
		videoUrl := cdnUrl.String()

		// 获取视频封面的url，地址和视频地址一样，只是所在目录是/tiktok/cover/，文件名是视频文件名+.jpg后缀
		cdnUrl.Path = "/tiktok/cover/" + strings.Replace(fileName, ".mp4", ".jpg", -1)
		coverUrl := cdnUrl.String()

		_, err = svcCtx.VideoRpc.PublishVideo(r.Context(), &rpcvideo.PublishVideoRequest{
			Video: &rpcvideo.VideoInfo{
				AuthorId: UserId,
				Title:    req.Title,
				PlayUrl:  videoUrl,
				CoverUrl: coverUrl,
			},
		})
		if err != nil {
			httpx.Error(w, schema.ServerError{
				ApiError: schema.ApiError{
					StatusCode: 500,
					Code:       50000,
					Message:    "rpc service error",
				},
				Detail: err,
			})
			return
		}
		l := video.NewPublishVideoLogic(r.Context(), svcCtx)
		resp, err := l.PublishVideo(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
