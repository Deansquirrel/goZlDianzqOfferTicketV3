package object

import "time"

type TktModel struct {
	AppId     string
	TktNo     string
	AccId     int
	TktName   string
	TktKind   string
	AddMy     float64
	CashMy    float64
	PcNo      string
	EffTime   time.Time
	DeadLine  time.Time
	DbConnKey string
}
