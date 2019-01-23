package object

import "time"

type TktInfo struct {
	AppId    string    `json:"appid"`
	AccId    int       `json:"accid"`
	TktKind  string    `json:"tktKind"`
	EffDate  time.Time `json:"effDate"`
	Deadline time.Time `json:"deadline"`
	CrYwLsh  string    `json:"crYwlsh"`
	CrBr     string    `json:"crBr"`
	ChTs     int       `json:"chts"`
	TktNo    string    `json:"tktno"`
	AccIdEnc string    `json:"accidenc"`
	CashMy   float64   `json:"cashmy"`
	AddMy    float64   `json:"addmy"`
	TktName  string    `json:"tktname"`
	PCno     string    `json:"pcno"`
}

func (t *TktInfo) IsEmpty() bool {
	if t.PCno != "" {
		return false
	}
	return true
}
