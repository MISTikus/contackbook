package common

func BuildUrl(routeData ...string) string {
	routeData = removeEmpty(routeData)
	result, separator := "", "/"
	result += separator
	for i, route := range routeData {
		result += route
		if i+1 < len(routeData) {
			result += separator
		}
	}
	return result
}
func removeEmpty(routeData []string) []string {
	for i, route := range routeData {
		if route == "" {
			routeData = append(routeData[:i], routeData[i+1:]...)
			return removeEmpty(routeData)
		}
	}
	return routeData
}
