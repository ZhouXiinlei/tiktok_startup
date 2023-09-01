package video

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tikstart/http/internal/logic/video"
	"tikstart/http/internal/svc"
	"tikstart/http/internal/types"
)

func CommentVideoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CommentRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := video.NewCommentVideoLogic(r.Context(), svcCtx)
		resp, err := l.CommentVideo(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
