package werror

import (
	"fmt"
	"os"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println(fmt.Sprintf("%s", err))
		os.Exit(-1)
	}
}
