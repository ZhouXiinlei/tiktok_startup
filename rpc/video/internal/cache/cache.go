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
func GenWorkCountKey(userId int64) string {
	return fmt.Sprintf("getworkcount_%d", userId)

}
func GenUserFavoriteCountKey(userId int64) string {
	return fmt.Sprintf("getuserfavoritecount_%d", userId)
}
func GenTotalFavoritedKey(userId int64) string {
	return fmt.Sprintf("gettotalfavorited_%d", userId)

}

func GenVideoCountsKey(topic string) string {
	return fmt.Sprintf("video_%ss", topic)
}

func GenVideoHeatKey(videoId int64) string {
	return fmt.Sprintf("video_heat_%d", videoId)
}
