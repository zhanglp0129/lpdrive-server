package vo

type ResultVO struct {
	Code int     `json:"code,omitempty"`
	Data any     `json:"data,omitempty"`
	Msg  *string `json:"msg,omitempty"`
}

func Success(data any) ResultVO {
	return ResultVO{
		Code: 1000,
		Data: data,
	}
}

func Error(msg string) ResultVO {
	return ResultVO{
		Code: 2000,
		Msg:  &msg,
	}
}
