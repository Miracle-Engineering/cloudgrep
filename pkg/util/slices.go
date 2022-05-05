package util

//TODO use exp/slices in the future but as of now lint fails with this error
// no export data for \"golang.org/x/exp/slices
func Contains(sl []string, name string) bool {
	for _, value := range sl {
		if value == name {
			return true
		}
	}
	return false
}
