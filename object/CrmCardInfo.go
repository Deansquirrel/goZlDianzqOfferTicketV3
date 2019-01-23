package object

type CrmCardInfo struct {
	OprYwSno string `json:"oprywsno"`
	Sn       int    `json:"sn"`
	CardNo   string `json:"cardno"`
	CardType int    `json:"cardtype"`
	AccId    int    `json:"accid"`
}
