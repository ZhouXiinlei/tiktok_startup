package cache

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
