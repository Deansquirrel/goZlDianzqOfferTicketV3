package object

import (
	"strings"
)

type ResponseCreateLittleTkt struct {
	TmpCol    int             `json:"tmpcol"`
	TktReturn []TktReturnInfo `json:"tktReturn"`

	ErrorModel  ErrorModel `json:"errormodel"`
	Description string     `json:"description"`

	HttpCode    int  `josn:"httpcode"`
	DBCommitted bool `json:"dbcommitted"`

	IsSuccess bool   `json:"issuccess"`
	FindAccid bool   `json:"findaccid"`
	HasCommit bool   `json:"hascommit"`
	Guid      string `json:"requestGuid"`
}

func (response *ResponseCreateLittleTkt) BaseFunc(req *RequestCreateLittleTkt, resp *ResponseCreateLittleTkt) {
	resp.Guid = req.Guid
}

func (response *ResponseCreateLittleTkt) Refresh() {
	//if (response.HttpCode == 503 || response.HttpCode == 403) && response.Description != "" {
	//	if strings.Contains(response.Description,"账号无法找到") {
	//		panic(response.Description)
	//	}
	//}

	if response.ErrorModel.ErrType == "" && response.ErrorModel.Desc == "" {
		response.DBCommitted = true
	} else {
		response.DBCommitted = false
	}

	if response.ErrorModel.ErrType != "" || response.ErrorModel.Desc != "" {
		response.IsSuccess = false
	} else if response.HttpCode != 200 {
		response.IsSuccess = false
	} else {
		response.IsSuccess = true
	}

	if response.HttpCode == 503 && strings.Contains(response.Description, "账号无法找到") {
		response.FindAccid = false
	} else {
		response.FindAccid = true
	}
}
