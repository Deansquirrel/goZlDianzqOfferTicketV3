package object

import "time"

type CrmFqYwInfo struct {
	FfDx   int       `json:"ffdx"`
	CrmJhh string    `json:"crmjhh"`
	FqZs   int       `json:"fqzs"`
	XtTsr  time.Time `json:"xxtsr"`
	TsFs   int       `json:"tsfs"`
	TsNr   string    `json:"tsnr"`
}
