package cache

import "fmt"

func GenFollowKey(userId int64, targetId int64) string {
	return fmt.Sprintf("follow_%d_%d", userId, targetId)
}

func GenUserCountsKey(topic string) string {
	return fmt.Sprintf("user_%ss", topic)
}

func GenUserHeatKey(userId int64) string {
	return fmt.Sprintf("user_heat_%d", userId)
}
