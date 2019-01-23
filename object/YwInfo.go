package object

import "time"

type YwInfo struct {
	OprBrid   string    `json:"oprbrid"`
	OprYwSno  string    `json:"oprywsno"`
	OprPpid   string    `json:"oprppid"`
	OprTime   time.Time `json:"oprtime"`
	OprHsDate time.Time `json:"oprhsdate"`
	OprId     string    `josn:"oprid"`
	OprAccId  int       `json:"opraccid"`
	XtTsr     time.Time `json:"xxtsr"`
	TsFs      int       `json:"tsfs"`
	TsNr      string    `json:"tsnr"`
}
