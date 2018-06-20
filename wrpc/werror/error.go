package werror

import (
	"fmt"
	"os"
)

func CheckErr(err error, line int) {
	if err != nil {
		fmt.Println(fmt.Sprintf("%d, %s", line, err))
		os.Exit(-1)
	}
}