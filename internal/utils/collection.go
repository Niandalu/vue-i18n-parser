package utils

func Entries(someMap map[string]interface{}) ([]string, []interface{}) {
	keys := make([]string, 0, len(someMap))
	values := make([]interface{}, 0, len(someMap))

	for k, v := range someMap {
		keys = append(keys, k)
		values = append(values, v)
	}

	return keys, values
}

func Contains(collection []string, one string) bool {
	for _, a := range collection {
		if a == one {
			return true
		}
	}
	return false
}
