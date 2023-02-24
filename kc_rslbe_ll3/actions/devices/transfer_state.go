package devices

import (
	"fmt"
	"strings"

	"git.kaiostech.com/cloud/common/cerrors"
	"git.kaiostech.com/cloud/common/config"
	"git.kaiostech.com/cloud/common/model/core"
	"git.kaiostech.com/cloud/common/mq"
	common "git.kaiostech.com/cloud/common/utils/actions_common"
	l4g "git.kaiostech.com/cloud/thirdparty/code.google.com/p/log4go"
)

var rsl_transfer_state_action_line []common.IAutomatizer

func init() {
	max_value := int(mq.MQRT_LAST)

	rsl_transfer_state_action_line = make([]common.IAutomatizer, max_value+1)

	// By default, we are just forwarding to data layer ...
	for i := 0; i < max_value; i++ {
		rsl_transfer_state_action_line[i] = nil
	}

	rsl_transfer_state_action_line[int(mq.MQRT_UPDATE)] = TransferState(common.ForwardRequest(nil))
}

func GetTransferStateAutomataLine() []common.IAutomatizer {
	return rsl_transfer_state_action_line
}

func TransferState(forward common.IAutomatizer) common.IAutomatizer {
	return common.AutomatizerFunc(func(dl_mqclient *mq.MQClient, handler_id uint64, jwt *core.JWT, request *mq.MQRequest) *mq.MQResponse {
		var rsp *mq.MQResponse
		if config.GetFEConfig().Common.Debug {
			l4g.Debug("req #%v: TransferState starts", request.ReqId)
			defer l4g.Debug("req #%v: TransferState ends", request.ReqId)
		}
		object_ids := strings.Split(request.ObjectId, ":")
		if len(object_ids) != 3 {
			l4g.Error("req #%v: TransferState: Invalid object ID:", request.ReqId)
			rsp = mq.NewMQResponse(request.Id, request.Type, request.Scope, nil, true, request.ReqId)
			return common.MakeErrorRsp(rsp, cerrors.ERROR_INTERNAL_SERVER_ERROR, "TransferState: Invalid object ID", "LL")
		}

		partnerID := object_ids[0]
		fromIMEI := object_ids[1]
		toIMEI := object_ids[2]

		l4g.Debug("req #%v: TransferState: PartnerID: %s", request.ReqId, partnerID)

		//Forwarding the request to DL
		req := mq.NewMQRequest(fmt.Sprintf("%s:%s:%s", partnerID, fromIMEI, toIMEI), mq.MQRT_UPDATE, mq.MQSCOPE_RSL_BE_TRANSFER_STATE, request.ReqId, request.ReqInfo)
		rsp, cerr := dl_mqclient.SendRecv(common.LL2DL, req)
		if cerr != nil {
			l4g.Error("req #%v: TransferState: transfer state of the device:'%s'.", request.ReqId, cerr.Error())
			rsp = mq.NewMQResponse(request.Id, request.Type, request.Scope, nil, true, request.ReqId)
			if cerr.Code == 404 {
				return common.MakeErrorRsp(rsp, cerrors.ERROR_NOT_FOUND, fmt.Sprintf("TransferState: Failed to transfer the state of device: '%s'.", cerr.Error()), "LL")
			}
			return common.MakeErrorRsp(rsp, cerrors.ERROR_INTERNAL_SERVER_ERROR, fmt.Sprintf("TransferState: Failed to update 2I data of the device: '%s'.", cerr.Error()), "LL")
		}
		//Returning the response to FE
		rsp = mq.NewMQResponse(request.Id, request.Type, request.Scope, nil, false, request.ReqId)
		return rsp
	})
}
