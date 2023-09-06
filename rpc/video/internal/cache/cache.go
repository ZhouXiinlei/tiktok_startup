package cache

import "fmt"

func GenFavoriteKey(userId int64, targetId int64) string {
	return fmt.Sprintf("favorite_%d_%d", userId, targetId)
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

func GenVideoCountsKey(topic string) string {
	return fmt.Sprintf("video_%ss", topic)
}

func GenVideoHeatKey(videoId int64) string {
	return fmt.Sprintf("video_heat_%d", videoId)
}
