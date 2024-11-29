package common

func RemoveDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	var list []string
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func FindIndex(strSLice []string, name string) int {
	for i, item := range strSLice {
		if item == name {
			return i
		}
	}
	return -1
}
