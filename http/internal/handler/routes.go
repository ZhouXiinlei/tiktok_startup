// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	app "tikstart/http/internal/handler/app"
	social "tikstart/http/internal/handler/social"
	user "tikstart/http/internal/handler/user"
	video "tikstart/http/internal/handler/video"
	"tikstart/http/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/ping",
				Handler: app.PingHandler(serverCtx),
			},
		},
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/douyin/feed",
				Handler: video.GetVideoListHandler(serverCtx),
			},
		},
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.JwtAuth},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/douyin/video",
					Handler: video.PublishVideoHandler(serverCtx),
				},
			}...,
		),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/register",
				Handler: user.RegisterHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/login",
				Handler: user.LoginHandler(serverCtx),
			},
		},
		rest.WithPrefix("/douyin/user"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.JwtAuth},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/",
					Handler: user.GetUserInfoHandler(serverCtx),
				},
			}...,
		),
		rest.WithPrefix("/douyin/user"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/action",
				Handler: social.FollowHandler(serverCtx),
			},
		},
		rest.WithPrefix("/douyin/relation"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.JwtAuth},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/action",
					Handler: video.FavoriteHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/list",
					Handler: video.GetFavoriteListHandler(serverCtx),
				},
			}...,
		),
		rest.WithPrefix("/douyin/favorite"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/friend/list",
				Handler: social.GetFriendListHandler(serverCtx),
			},
		},
		rest.WithPrefix("/douyin/relation"),
	)
}
