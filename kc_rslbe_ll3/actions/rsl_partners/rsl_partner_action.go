package rsl_partners

import (
	"fmt"

	"git.kaiostech.com/cloud/common/cerrors"
	"git.kaiostech.com/cloud/common/config"
	"git.kaiostech.com/cloud/common/model/core"
	"git.kaiostech.com/cloud/common/mq"
	common "git.kaiostech.com/cloud/common/utils/actions_common"
	l4g "git.kaiostech.com/cloud/thirdparty/code.google.com/p/log4go"
)

var rsl_partners_actions_line []common.IAutomatizer

func init() {
	max_value := int(mq.MQRT_LAST)

	rsl_partners_actions_line = make([]common.IAutomatizer, max_value+1)

	// By default, we are just forwarding to data layer ...
	for i := 0; i < max_value; i++ {
		rsl_partners_actions_line[i] = nil
	}

	rsl_partners_actions_line[int(mq.MQRT_CREATE)] = Register3I(common.ForwardRequest(nil))
}

func GetAutomataLine() []common.IAutomatizer {
	return rsl_partners_actions_line
}

func Register3I(forward common.IAutomatizer) common.IAutomatizer {
	return common.AutomatizerFunc(func(dl_mqclient *mq.MQClient, handler_id uint64, jwt *core.JWT, request *mq.MQRequest) *mq.MQResponse {
		var rsp *mq.MQResponse
		if config.GetFEConfig().Common.Debug {
			l4g.Debug("req #%v: Register3I starts", request.ReqId)
			defer l4g.Debug("req #%v: Register3I ends", request.ReqId)
		}
		partner_id := request.ObjectId
		data, err := request.GetData()
		if err != nil {
			l4g.Error("req #%v: Register3I: Failed to retrieve data from FE request: '%s'.", request.ReqId, err.Error())
			rsp = mq.NewMQResponse(request.Id, request.Type, request.Scope, nil, true, request.ReqId)
			return common.MakeErrorRsp(rsp, cerrors.ERROR_INTERNAL_SERVER_ERROR, fmt.Sprintf("Register3I: Failed to retrieve data from FE request: '%s'.", err.Error()), "LL")
		}

		l4g.Debug("req #%v: Register3I: PartnerID: %s", request.ReqId, partner_id)

		//Forwarding the request to DL
		req := mq.NewMQRequest(partner_id, mq.MQRT_CREATE, mq.MQSCOPE_RSL_BE_REGISTER_3I, request.ReqId, request.ReqInfo)

		if err = req.SetData(data); err != nil {
			l4g.Error("req #%v: UpdateFinDevice: Failed to set request data for updating financier device: '%s'.", request.ReqId, err.Error())
			return common.MakeErrorRsp(rsp, cerrors.ERROR_INTERNAL_SERVER_ERROR, fmt.Sprintf("Register3I: Error setting the body of MQ request: '%s'.", err), "LL")
		}
		rsp, cerr := dl_mqclient.SendRecv(common.LL2DL, req)
		if cerr != nil || rsp.IsError {
			l4g.Error("req #%v: Register3I: Error getting response from DL while registering 3I data of the device:'%s'.", request.ReqId, cerr.Error())
			rsp = mq.NewMQResponse(request.Id, request.Type, request.Scope, nil, true, request.ReqId)
			return common.MakeErrorRsp(rsp, cerrors.ERROR_INTERNAL_SERVER_ERROR, fmt.Sprintf("Register3I: Error getting response from DL while registering 3i data: '%s'.", cerr.Error()), "LL")

		}
		rsp_data, err := rsp.GetData()
		if err != nil {
			l4g.Error("req #%v: Register3I: Failed to retrieve data from response:'%s'.", request.ReqId, cerr.Error())
			rsp = mq.NewMQResponse(request.Id, request.Type, request.Scope, nil, true, request.ReqId)
			return common.MakeErrorRsp(rsp, cerrors.ERROR_INTERNAL_SERVER_ERROR, fmt.Sprintf("Register3I: Failed to retrieve data from response: '%s'.", cerr.Error()), "LL")
		}
		//Returning the response to FE
		rsp = mq.NewMQResponse(request.Id, request.Type, request.Scope, rsp_data, false, request.ReqId)
		return rsp
	})
}
