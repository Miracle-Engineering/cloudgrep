package util

func Contains(sl []string, name string) bool {
	for _, value := range sl {
		if value == name {
			return true
		}
	}
	return false
}
