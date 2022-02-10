package output

import "strings"

func GetOutputFilename(filename string) string {
	if strings.HasSuffix(filename, ".hz") {
		return filename
	} else {
		return filename + ".hz"
	}
}
