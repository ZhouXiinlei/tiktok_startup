package utils

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
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
