package helpers

func Contains(slice []string, element map[string]interface{}) (contains bool) {
	for _, v := range slice {
		if v == element["Value"] {
			return true
		}
	}
	return
}
