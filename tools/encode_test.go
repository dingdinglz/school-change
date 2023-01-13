package tools

import (
	"fmt"
	"testing"
)

func TestMD5(t *testing.T) {
	fmt.Println(MD5("dinglz"))
	fmt.Println(MD5("adminadmin"))
}
