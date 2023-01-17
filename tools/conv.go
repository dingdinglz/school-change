package tools

import (
	"fmt"
	"strconv"
)

// StringToInt String转int
func StringToInt(i string) int {
	res, _ := strconv.Atoi(i)
	return res
}

// InterfaceToString Interface转string
func InterfaceToString(i interface{}) string {
	return fmt.Sprintf("%s", i)
}

// InterfaceToInt Interface转int
func InterfaceToInt(i interface{}) int {
	return StringToInt(InterfaceToString(i))
}
