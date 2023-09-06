package utils

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/zeromicro/go-zero/core/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"regexp"
	"tikstart/http/schema"
)

func MatchError(err error, target *status.Status) (*status.Status, bool) {
	st, _ := status.FromError(err)
	if st.Message() == target.Message() {
		return st, true
	}
	return st, false
}

func MatchRegexp(pattern string, value string) bool {
	r := regexp.MustCompile(pattern)
	return r.MatchString(value)
}

func InternalWithDetails(msg string, items ...interface{}) error {
	fmt.Printf("InternalError: %s\n", msg)
	details := make([]proto.Message, 0, len(items))
	for index, item := range items {
		switch v := item.(type) {
		case error:
			details = append(details, &any.Any{
				Value: []byte(v.Error()),
			})
			fmt.Printf("%d: %v\n", index, v.Error())
		case string:
			details = append(details, &any.Any{
				Value: []byte(v),
			})
			fmt.Printf("%d: %v\n", index, v)
		default:
			// try fmt.Sprintf to stringify this item
			details = append(details, &any.Any{
				Value: []byte(fmt.Sprintf("%v", v)),
			})
			fmt.Printf("%d: %v\n", index, v)
		}
	}
	st, _ := status.New(codes.Internal, msg).WithDetails(details...)
	return st.Err()
}

func ReturnInternalError(ctx context.Context, st *status.Status, err error) error {
	for index, item := range st.Details() {
		detail := item.(*anypb.Any)
		fmt.Printf("%d: %s\n", index, string(detail.Value))
	}

	return schema.ServerError{
		ApiError: schema.ApiError{
			StatusCode: 500,
			Code:       50000,
			Message:    "Internal Server Error",
		},
		Detail:  err,
		TraceId: trace.TraceIDFromContext(ctx),
	}
}

func SortId(idA int64, idB int64) (int64, int64) {
	if idA < idB {
		return idA, idB
	}
	return idB, idA
}
