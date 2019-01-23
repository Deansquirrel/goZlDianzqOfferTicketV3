package object

type TktCreateInfo struct {
	Queue         string
	MessageRoute  string
	TktInfo       []TktInfo
	TktYwInfo     YwInfo
	TktReturnInfo []TktReturnInfo
	MdFqTsInfo    []MdCreateTktTuis
	CzLx          int
	CzLxSm        string
}
