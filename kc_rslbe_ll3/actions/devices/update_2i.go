package devices

import (
	"fmt"

	"git.kaiostech.com/cloud/common/cerrors"
	"git.kaiostech.com/cloud/common/config"
	"git.kaiostech.com/cloud/common/model/core"
	"git.kaiostech.com/cloud/common/mq"
	common "git.kaiostech.com/cloud/common/utils/actions_common"
	l4g "git.kaiostech.com/cloud/thirdparty/code.google.com/p/log4go"
)

var rsl_update_2i_action_line []common.IAutomatizer

func init() {
	max_value := int(mq.MQRT_LAST)

	rsl_update_2i_action_line = make([]common.IAutomatizer, max_value+1)

	// By default, we are just forwarding to data layer ...
	for i := 0; i < max_value; i++ {
		rsl_update_2i_action_line[i] = nil
	}

	rsl_update_2i_action_line[int(mq.MQRT_UPDATE)] = Update2I(common.ForwardRequest(nil))
}

func GetUpdate2IAutomataLine() []common.IAutomatizer {
	return rsl_update_2i_action_line
}

func Update2I(forward common.IAutomatizer) common.IAutomatizer {
	return common.AutomatizerFunc(func(dl_mqclient *mq.MQClient, handler_id uint64, jwt *core.JWT, request *mq.MQRequest) *mq.MQResponse {
		var rsp *mq.MQResponse
		if config.GetFEConfig().Common.Debug {
			l4g.Debug("req #%v: Update2I starts", request.ReqId)
			defer l4g.Debug("req #%v: Update2I ends", request.ReqId)
		}
		partner_id := request.ObjectId
		data, err := request.GetData()
		if err != nil {
			l4g.Error("req #%v: Update2I: Failed to retrieve data from FE request: '%s'.", request.ReqId, err.Error())
			rsp = mq.NewMQResponse(request.Id, request.Type, request.Scope, nil, true, request.ReqId)
			return common.MakeErrorRsp(rsp, cerrors.ERROR_INTERNAL_SERVER_ERROR, fmt.Sprintf("Update2I: Failed to retrieve data from FE request: '%s'.", err.Error()), "LL")
		}

		l4g.Debug("req #%v: Update2I: PartnerID: %s", request.ReqId, partner_id)

		//Forwarding the request to DL
		req := mq.NewMQRequest(partner_id, mq.MQRT_UPDATE, mq.MQSCOPE_RSL_BE_UPDATE_2I, request.ReqId, request.ReqInfo)

		if err = req.SetData(data); err != nil {
			l4g.Error("req #%v: Update2I: Failed to set request data for updating financier device: '%s'.", request.ReqId, err.Error())
			return common.MakeErrorRsp(rsp, cerrors.ERROR_INTERNAL_SERVER_ERROR, fmt.Sprintf("Update2I: Error setting the body of MQ request: '%s'.", err), "LL")
		}
		rsp, cerr := dl_mqclient.SendRecv(common.LL2DL, req)
		if cerr != nil {
			l4g.Error("req #%v: Update2I: Failed to update 2I data of the device:'%s'.", request.ReqId, cerr.Error())
			rsp = mq.NewMQResponse(request.Id, request.Type, request.Scope, nil, true, request.ReqId)
			if cerr.Code == 404 {
				return common.MakeErrorRsp(rsp, cerrors.ERROR_NOT_FOUND, fmt.Sprintf("Update2I: Failed to update 2I data of the device: '%s'.", cerr.Error()), "LL")
			}
			return common.MakeErrorRsp(rsp, cerrors.ERROR_INTERNAL_SERVER_ERROR, fmt.Sprintf("Update2I: Failed to update 2I data of the device: '%s'.", cerr.Error()), "LL")
		}
		rsp_data, err := rsp.GetData()
		if err != nil {
			l4g.Error("req #%v: Update2I: Failed to retrieve data from response:'%s'.", request.ReqId, cerr.Error())
			rsp = mq.NewMQResponse(request.Id, request.Type, request.Scope, nil, true, request.ReqId)
			return common.MakeErrorRsp(rsp, cerrors.ERROR_INTERNAL_SERVER_ERROR, fmt.Sprintf("Update2I: Failed to retrieve data from response: '%s'.", cerr.Error()), "LL")
		}
		if rsp_data == nil {
			l4g.Error("req #%v: Update2I: Received nil data in the reponse body from DL:'%s'.", request.ReqId, nil)
			rsp = mq.NewMQResponse(request.Id, request.Type, request.Scope, nil, true, request.ReqId)
			return common.MakeErrorRsp(rsp, cerrors.ERROR_INTERNAL_SERVER_ERROR, fmt.Sprintf("Update2I: Received nil data in the reponse body from DL: '%s'.", cerr.Error()), "LL")
		}

		//Returning the response to FE
		rsp = mq.NewMQResponse(request.Id, request.Type, request.Scope, rsp_data, false, request.ReqId)
		return rsp
	})
}
