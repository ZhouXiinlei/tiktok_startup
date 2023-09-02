package video

import (
	"bytes"
	"context"
	"github.com/google/uuid"
	"github.com/h2non/filetype"
	"github.com/zeromicro/go-zero/rest/httpx"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"tikstart/common/utils"
	"tikstart/http/internal/logic/video"
	"tikstart/http/internal/svc"
	"tikstart/http/internal/types"
	"tikstart/http/schema"
	rpcvideo "tikstart/rpc/video/video"
)

func PublishVideoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PublishVideoRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		Userclaims, _ := utils.ParseToken(req.Token, svcCtx.Config.JwtAuth.Secret)
		UserId := Userclaims.UserId
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
		Keyname := "/video/" + fileName
		_, err = svcCtx.TengxunyunClient.Object.Put(context.Background(), Keyname, file, nil)
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
		Url, err := url.Parse(svcCtx.Config.COS.Endpoint)
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
		Url.Path = "/" + Keyname
		VideoUrl := Url.String()

		_, err = svcCtx.VideoRpc.PublishVideo(r.Context(), &rpcvideo.PublishVideoRequest{
			Video: &rpcvideo.VideoInfo{
				AuthorId: UserId,
				Title:    req.Title,
				PlayUrl:  VideoUrl,
				CoverUrl: "11111",
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
