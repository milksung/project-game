package helper

import (
	"fmt"
	"strconv"
)

func IncrementNumber(number string) string {
	conv, _ := strconv.Atoi(number)
	return fmt.Sprintf("%06d", conv+1)
}
