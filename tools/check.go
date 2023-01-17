package tools

// IsAllEnglish 判断字符串中是否只含有英文
func IsAllEnglish(s string) bool {
	AllWord := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	for _, i := range s {
		var _flag = false
		for _, j := range AllWord {
			if i == j {
				_flag = true
			}
		}
		if !_flag {
			return false
		}
	}
	return true
}
