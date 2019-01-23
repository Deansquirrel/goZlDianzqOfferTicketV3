package object

import "time"

type MdFqYwInfo struct {
	Sn      int       `json:"sn"`
	TktKind string    `json:"tktkind"`
	TktNum  int       `json:"tktnum"`
	XtTsr   time.Time `json:"xxtsr"`
	TsFs    int       `json:"tsfs"`
	TsNr    string    `json:"tsnr"`
}
