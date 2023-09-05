package cache

import "fmt"

func GenFollowKey(userId int64, targetId int64) string {
	return fmt.Sprintf("follow_%d_%d", userId, targetId)
}

func GenPopularUsersKey() string {
	return "popular_users"
}

func GenUserCountsKey(topic string) string {
	return fmt.Sprintf("user_%ss", topic)
}

func GenUserHeatKey(userId int64) string {
	return fmt.Sprintf("user_heat_%d", userId)
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
