package helpers

func ContainsJson(slice []string, element map[string]interface{}) (contains bool) {
	for _, v := range slice {
		if v == element["Value"] {
			return true
		}
	}
	return
}

func ContainsString(slice []string, element string) (contains bool) {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return
}
