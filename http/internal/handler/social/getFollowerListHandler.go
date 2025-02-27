package social

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tikstart/http/internal/logic/social"
	"tikstart/http/internal/svc"
	"tikstart/http/internal/types"
)

func GetFollowerListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetFollowerListRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := social.NewGetFollowerListLogic(r.Context(), svcCtx)
		resp, err := l.GetFollowerList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
