package common

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUserAlreadyExists = status.New(codes.AlreadyExists, "user already exists")
	ErrUserNotFound      = status.New(codes.NotFound, "user not found")
	ErrCommentNotFound   = status.New(codes.NotFound, "comment not found")
	ErrVideoNotFound     = status.New(codes.NotFound, "video not found")
)
