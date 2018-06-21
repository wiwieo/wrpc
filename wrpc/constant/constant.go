package constant

import "time"

const (
	END_SIGN         = '\n'
	SUCCESS          = "0"
	SUCCESS_MSG      = "成功。"
	FAILED           = "1"
	FAILED_MSG       = "失败。"
	METHOD_NOT_EXIST = "方法不存在。"
	TIME_OUT = 1000*time.Second
)
