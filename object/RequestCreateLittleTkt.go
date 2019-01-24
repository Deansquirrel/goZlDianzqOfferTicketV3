package object

import (
	"errors"
)

type RequestCreateLittleTkt struct {
	Guid        string
	AppId       string
	ReturnTktNo int
	Body        requestBody
}

type requestBody struct {
	TmpCol      int
	TktInfo     []TktInfo
	CrmFqYwInfo CrmFqYwInfo
	CrmCardInfo []CrmCardInfo
	MdFqYwInfo  []MdFqYwInfo
	YwInfo      YwInfo
}

func (request *RequestCreateLittleTkt) CheckRequest() error {
	if request.Body.TktInfo == nil || request.Body.CrmFqYwInfo.CrmJhh == "" || request.Body.CrmCardInfo == nil || request.Body.YwInfo.OprYwSno == "" {
		return errors.New("传入的记录集为空")
	}
	for index := range request.Body.CrmCardInfo {
		if request.Body.CrmCardInfo[index].CardType != 5 {
			return errors.New("会员券目前仅支持按账户id加密值操作")
		}
	}
	if len(request.Body.CrmCardInfo) > 100 {
		return errors.New("券发放（立即生效）禁止超过100张")
	}
	return nil
}
