package cache

import "fmt"

func GenFavoriteKey(userId int64, targetId int64) string {
	return fmt.Sprintf("favorite_%d_%d", userId, targetId)
}

func GenVideoCountsKey(topic string) string {
	return fmt.Sprintf("video_%ss", topic)
}

func GenVideoHeatKey(videoId int64) string {
	return fmt.Sprintf("video_heat_%d", videoId)
}
