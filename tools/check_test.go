package tools

import (
	"fmt"
	"testing"
)

func TestIsAllEnglish(t *testing.T) {
	fmt.Println(IsAllEnglish("dingdasdslz"))
	fmt.Println(IsAllEnglish("d5ingl8z"))
	fmt.Println(IsAllEnglish("sdasd"))
}
