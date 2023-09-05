package cache

import "fmt"

func GenFollowKey(userId int64, targetId int64) string {
	return fmt.Sprintf("follow_%d_%d", userId, targetId)
}

func GenFollowingCountKey(userId int64) string {
	return fmt.Sprintf("following_count_%d", userId)
}

func GenFollowerCountKey(userId int64) string {
	return fmt.Sprintf("follower_count_%d", userId)
}

func GenUserInfoKey(userId int64) string {
	return fmt.Sprintf("user_info_%d", userId)
}

func YesOrNo(cond bool) string {
	if cond {
		return "yes"
	}
	return "no"
}

func TrueOrFalse(str string) (val bool, match bool) {
	if str != "" {
		if str == "yes" {
			return true, true
		}
		return false, true
	}
	return false, false
}
