package structs

import "strings"

type tagOptions []string

func (t tagOptions) Has(opt string) bool {
	for _, tagOpt := range t {
		if tagOpt == opt {
			return true
		}
	}
	return false
}

func parseTag(tag string) (string, tagOptions) {
	res := strings.Split(tag, ",")
	return res[0], res[1:]
}

func parseGormColumn(tagName string) string {
	if !strings.Contains(tagName, "column") {
		return tagName
	}
	tags := strings.Split(tagName, ";")
	for _, tag := range tags {
		if strings.Contains(tagName, "column") {
			columns := strings.Split(tag, ":")
			if len(columns) > 0 {
				return columns[1]
			}
		}
	}
	return ""
}
