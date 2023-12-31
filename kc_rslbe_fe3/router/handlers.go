/*
 * Remote Sim Lock API
 *
 * **Implementation of the Remote Sim Lock APIs**  **Vibe project documentation (including RSL) can be found using the following links:**     **1- DRAFT-RSL sequence flows pictures source ZY-20220225.docx**    * https://kaios.sharepoint.com/:w:/r/sites/vibe/Shared%20Documents/4.%20PRODUCT/06%20-%20RSL%20REQ%20and%20Designing/DRAFT-RSL%20sequence%20flows%20pictures%20source%20ZY-20220225.docx?d=w8cba978f89bc4cac8460d87c0e1053ba&csf=1&web=1&e=iQp9Eo *    This document includes squence diagrams which describe the use cases: “device initiate” , “daily ping”, and “user change SIM to another one”.    **2- DRAFT-Vibe Remote SIM Lock Operation Structure-20220307.docx**    * https://kaios.sharepoint.com/:w:/r/sites/vibe/Shared%20Documents/4.%20PRODUCT/06%20-%20RSL%20REQ%20and%20Designing/DRAFT-Vibe%20Remote%20SIM%20Lock%20Operation%20Structure-20220307.docx?d=w51d49ceb6ddc4329aaa18ce02ce22318&csf=1&web=1&e=orUFfm    This is a short document which includes an introduction to the RSL operation structure, and RSL portal.    **3- RSL flow charts for communicate.v0.2.docx**    * https://kaios.sharepoint.com/:w:/r/sites/vibe/Shared%20Documents/4.%20PRODUCT/06%20-%20RSL%20REQ%20and%20Designing/RSL%20flow%20charts%20for%20communicate.v0.2.docx?d=w8215af326530486ab8b84b14e3208967&csf=1&web=1&e=H2DhdF    **4- Vibe RSL Technical Keypoints.20220915.pptx**    * https://kaios.sharepoint.com/:p:/r/sites/vibe/Shared%20Documents/4.%20PRODUCT/06%20-%20RSL%20REQ%20and%20Designing/Vibe%20RSL%20Technical%20Keypoints%20.20220915.pptx?d=w58603c5ed4814c9086a799450262b97e&csf=1&web=1&e=czrBNc    **5- Vibe Product requirement_211118.pptx**    * https://kaios-my.sharepoint.com/:p:/r/personal/raffi_semerciyan_kaiostech_com/Documents/Documents/20211124-Vibe/Vibe%20Product%20Requirement_211118.pptx?d=w130eb4dd2a74481d93cb1fec26c0d068&csf=1&web=1&e=ivj2jW    **6- Vibe-Requirements-Analysis.docx**    * https://kaios-my.sharepoint.com/:w:/g/personal/raffi_semerciyan_kaiostech_com/EXlVgqdcF7pAii8-VIHeRXQB_VQEBIbjO_2aU3Ur_Rqk1w?e=R7eaOb     This document shows the detailed design of RSL function, after making a requirement summary.      **7- Vibe-Specification.docx**   * https://kaios-my.sharepoint.com/:w:/g/personal/raffi_semerciyan_kaiostech_com/EXlVgqdcF7pAii8-VIHeRXQB_VQEBIbjO_2aU3Ur_Rqk1w?e=R7eaOb    This document shows the detailed use case diagrams which show the interaction between Vive users and the system.
 *
 * API version: 1.0.0
 * Contact: maen.hammour@kaiostech.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.kaiostech.com/cloud/common/cerrors"
	"git.kaiostech.com/cloud/common/config"
	"git.kaiostech.com/cloud/common/model/vibe/vibe_v1_0"
	"git.kaiostech.com/cloud/common/mq"
	"git.kaiostech.com/cloud/common/security/jwttoken"
	"git.kaiostech.com/cloud/common/utils/context"
	"git.kaiostech.com/cloud/common/utils/handlers_common"
	l4g "git.kaiostech.com/cloud/thirdparty/code.google.com/p/log4go"
	"git.kaiostech.com/cloud/thirdparty/github.com/gorilla/mux"
)

func VibeBeDevicesRegisterPost(w http.ResponseWriter, r *http.Request) {
	if config.GetFEConfig().Common.Debug {
		l4g.Debug("req #%v: handlers: VibeBeDevicesRegisterPost starts", context.GetReqId(r))
		defer l4g.Debug("req #%v: handlers:VibeBeDevicesRegisterPost ends", context.GetReqId(r))
	}
	// Getting partner ID form the route.
	vars := mux.Vars(r)
	partnerID := vars["id"]

	body := context.GetPayload(r)

	l4g.Debug("req #%v: handlers: VibeBeDevicesRegisterPost: Received request body: %s", context.GetReqId(r), body)
	l4g.Debug("req #%v: handlers: VibeBeDevicesRegisterPost: Received partner ID in the request: %s", context.GetReqId(r), partnerID)

	if len(body) == 0 {
		cerr := cerrors.New(cerrors.ERROR_MALFORMED_POST_DATA, "No data found in request payload", "FE", true)
		l4g.Error("req #%v: handlers: VibeBeDevicesRegisterPost: '%s'", context.GetReqId(r), cerr.Cause)
		handlers_common.RespondError(r, w, cerr)
		return
	}
	var err error

	device_info, err := vibe_v1_0.JsonToDeviceInfo(body)
	if err != nil {
		cerr := cerrors.New(cerrors.ERROR_MALFORMED_POST_DATA, fmt.Sprintf("unmarshal failed: '%s'", err.Error()), "FE", true)
		l4g.Error("req #%v: handlers: VibeBeDevicesRegisterPost: Error while unmarshalling body data: '%s'. offending payload: '%s'", context.GetReqId(r), err.Error(), body)
		handlers_common.RespondError(r, w, cerr)
		return
	}

	if err := device_info.SelfCheck(); err != nil {
		cerr := cerrors.New(cerrors.ERROR_INVALID_REQUEST, fmt.Sprintf("selfcheck failed: '%s'", err.Error()), "FE", true)
		l4g.Error("req #%v: handlers: VibeBeDevicesRegisterPost: Error in object data SelfCheck: '%s'.", context.GetReqId(r), err.Error())
		handlers_common.RespondError(r, w, cerr)
		return
	}

	jwt_claims_str := context.GetClaims(r)
	if jwt_claims_str == "" {
		cerr := cerrors.New(cerrors.ERROR_INVALID_REQUEST, "Failed to recover Token Claims from Context!", "FE", true)
		l4g.Error("req #%v: handlers: VibeBeDevicesRegisterPost: %s", context.GetReqId(r), cerr.Cause)
		handlers_common.RespondError(r, w, cerr)
		return
	}

	var jwt_claims jwttoken.Claims

	if jwt_claims, err = jwttoken.ParseClaims(jwt_claims_str); err != nil {
		cerr := cerrors.New(cerrors.ERROR_INVALID_REQUEST, fmt.Sprintf("Failed to parse Token Claims: '%s'!", err), "FE", true)
		l4g.Error("req #%v: handlers: VibeBeDevicesRegisterPost: Failed to parse claims: %s. Claims content: '%s'", context.GetReqId(r), err.Error(), jwt_claims)
		handlers_common.RespondError(r, w, cerr)
		return
	}

	// Sending to Messaging Queue for request processing
	request := mq.NewMQRequest(partnerID, mq.MQRT_CREATE, mq.MQSCOPE_RSL_BE_REGISTER_3I, context.GetReqId(r), mq.NewReqInfo(r))

	//Set body of the request
	if err = request.SetData(body); err != nil {
		cerr := cerrors.New(cerrors.ERROR_INTERNAL_SERVER_ERROR, fmt.Sprintf("Failed to set request data: '%s'", err.Error()), "FE", true)
		l4g.Error("req #%v: handlers: VibeBeDevicesRegisterPost: Failed to set request data: '%s'.", context.GetReqId(r), err.Error())
		handlers_common.RespondError(r, w, cerr)
		return
	}
	response, cerr := handlers_common.MQSendRecv(request)
	if cerr != nil {
		handlers_common.RespondError(r, w, cerr)
	}
	if response == nil {
		cerr := cerrors.New(cerrors.ERROR_INTERNAL_SERVER_ERROR, "Received nil response from LL", "FE", true)
		l4g.Error("req #%v: handlers: VibeBeDevicesRegisterPost: Received nil response from LL", context.GetReqId(r))
		handlers_common.RespondError(r, w, cerr)
		return
	}
	w.WriteHeader(http.StatusCreated)
	handlers_common.HTTPFeedbackGetResult(r, w, fmt.Sprintf("/kc_rsl_be/v1.0/partners/%s/3i", partnerID), response, nil)
	return
}
func VibeBeDevicesRslImeiPost(w http.ResponseWriter, r *http.Request) {
	if config.GetFEConfig().Common.Debug {
		l4g.Debug("req #%v: handlers: VibeBeDevicesRslImeiPost starts", context.GetReqId(r))
		defer l4g.Debug("req #%v: handlers:VibeBeDevicesRslImeiPost ends", context.GetReqId(r))
	}

	var err error
	// Getting the IMEI from the request
	vars := mux.Vars(r)
	imei := vars["imei"]

	jwt_claims_str := context.GetClaims(r)
	if jwt_claims_str == "" {
		cerr := cerrors.New(cerrors.ERROR_INVALID_REQUEST, "Failed to recover Token Claims from Context!", "FE", true)
		l4g.Error("req #%v: handlers: VibeBeDevicesRslImeiPost: %s", context.GetReqId(r), cerr.Cause)
		handlers_common.RespondError(r, w, cerr)
		return
	}

	var jwt_claims jwttoken.Claims

	if jwt_claims, err = jwttoken.ParseClaims(jwt_claims_str); err != nil {
		cerr := cerrors.New(cerrors.ERROR_INVALID_REQUEST, fmt.Sprintf("Failed to parse Token Claims: '%s'!", err), "FE", true)
		l4g.Error("req #%v: handlers: VibeBeDevicesRslImeiPost: Failed to parse claims: %s. Claims content: '%s'", context.GetReqId(r), err.Error(), jwt_claims)
		handlers_common.RespondError(r, w, cerr)
		return
	}
	partnerID := jwt_claims.Get("pid")
	// Sending to Messaging Queue for request processing
	request := mq.NewMQRequest(fmt.Sprintf("%s:%s", partnerID, imei), mq.MQRT_CREATE, mq.MQSCOPE_RSL_BE_COMMAND, context.GetReqId(r), mq.NewReqInfo(r))

	response, cerr := handlers_common.MQSendRecv(request)
	if cerr != nil {
		handlers_common.RespondError(r, w, cerr)
	}
	handlers_common.HTTPFeedbackGetResult(r, w, fmt.Sprintf("/kc_rsl_be/v1.0/devices/rsl/%s", imei), response, nil)
	return

}

func VibeBeDevicesTransferOwnershipImeiPartnerIdPost(w http.ResponseWriter, r *http.Request) {
	if config.GetFEConfig().Common.Debug {
		l4g.Debug("req #%v: handlers: VibeBeDevicesTransferOwnershipImeiPartnerIdPost starts", context.GetReqId(r))
		defer l4g.Debug("req #%v: handlers:VibeBeDevicesTransferOwnershipImeiPartnerIdPost ends", context.GetReqId(r))
	}

	var err error
	// Getting the IMEI and partner_id from the request
	vars := mux.Vars(r)
	imei := vars["imei"]
	partnerID := vars["partner_id"]

	jwt_claims_str := context.GetClaims(r)
	if jwt_claims_str == "" {
		cerr := cerrors.New(cerrors.ERROR_INVALID_REQUEST, "Failed to recover Token Claims from Context!", "FE", true)
		l4g.Error("req #%v: handlers: VibeBeDevicesTransferOwnershipImeiPartnerIdPost: %s", context.GetReqId(r), cerr.Cause)
		handlers_common.RespondError(r, w, cerr)
		return
	}

	var jwt_claims jwttoken.Claims

	if jwt_claims, err = jwttoken.ParseClaims(jwt_claims_str); err != nil {
		cerr := cerrors.New(cerrors.ERROR_INVALID_REQUEST, fmt.Sprintf("Failed to parse Token Claims: '%s'!", err), "FE", true)
		l4g.Error("req #%v: handlers: VibeBeDevicesTransferOwnershipImeiPartnerIdPost: Failed to parse claims: %s. Claims content: '%s'", context.GetReqId(r), err.Error(), jwt_claims)
		handlers_common.RespondError(r, w, cerr)
		return
	}
	// Sending to Messaging Queue for request processing
	request := mq.NewMQRequest(fmt.Sprintf("%s:%s", partnerID, imei), mq.MQRT_CREATE, mq.MQSCOPE_RSL_BE_TRANSFER_OWNERSHIP, context.GetReqId(r), mq.NewReqInfo(r))

	response, cerr := handlers_common.MQSendRecv(request)
	if cerr != nil {
		handlers_common.RespondError(r, w, cerr)
	}
	handlers_common.HTTPFeedbackCreateResult(r, w, fmt.Sprintf("kc_rsl_be/v1.0/devices/transfer_ownership/%s/%s", imei, partnerID), response, nil)
	return
}

func VibeBeDevicesTransferStateImeiPartnerIdPost(w http.ResponseWriter, r *http.Request) {
	if config.GetFEConfig().Common.Debug {
		l4g.Debug("req #%v: handlers: VibeBeDevicesTransferStateImeiPartnerIdPost starts", context.GetReqId(r))
		defer l4g.Debug("req #%v: handlers:VibeBeDevicesTransferStateImeiPartnerIdPost ends", context.GetReqId(r))
	}

	var err error
	// Getting both IMEIs from the request
	vars := mux.Vars(r)
	from_imei := vars["from_imei"]
	to_imei := vars["to_imei"]

	jwt_claims_str := context.GetClaims(r)
	if jwt_claims_str == "" {
		cerr := cerrors.New(cerrors.ERROR_INVALID_REQUEST, "Failed to recover Token Claims from Context!", "FE", true)
		l4g.Error("req #%v: handlers: VibeBeDevicesTransferStateImeiPartnerIdPost: %s", context.GetReqId(r), cerr.Cause)
		handlers_common.RespondError(r, w, cerr)
		return
	}

	var jwt_claims jwttoken.Claims

	if jwt_claims, err = jwttoken.ParseClaims(jwt_claims_str); err != nil {
		cerr := cerrors.New(cerrors.ERROR_INVALID_REQUEST, fmt.Sprintf("Failed to parse Token Claims: '%s'!", err), "FE", true)
		l4g.Error("req #%v: handlers: VibeBeDevicesTransferStateImeiPartnerIdPost: Failed to parse claims: %s. Claims content: '%s'", context.GetReqId(r), err.Error(), jwt_claims)
		handlers_common.RespondError(r, w, cerr)
		return
	}

	partnerID := jwt_claims.Get("pid")

	// Sending to Messaging Queue for request processing
	request := mq.NewMQRequest(fmt.Sprintf("%s:%s:%s", partnerID.(string), from_imei, to_imei), mq.MQRT_UPDATE, mq.MQSCOPE_RSL_BE_TRANSFER_STATE, context.GetReqId(r), mq.NewReqInfo(r))

	response, cerr := handlers_common.MQSendRecv(request)
	if cerr != nil {
		handlers_common.RespondError(r, w, cerr)
	}
	handlers_common.HTTPFeedbackCreateResult(r, w, fmt.Sprintf("/kc_rsl_be/v1.0/devices/transfer_state/%s/%s", from_imei, to_imei), response, nil)
	return
}

func VibeBeDevicesUnleashImeiPost(w http.ResponseWriter, r *http.Request) {
	if config.GetFEConfig().Common.Debug {
		l4g.Debug("req #%v: handlers: VibeBeDevicesUnleashImeiPost starts", context.GetReqId(r))
		defer l4g.Debug("req #%v: handlers:VibeBeDevicesUnleashImeiPost ends", context.GetReqId(r))
	}

	var err error
	// Getting both IMEI
	vars := mux.Vars(r)
	imei := vars["imei"]

	jwt_claims_str := context.GetClaims(r)
	if jwt_claims_str == "" {
		cerr := cerrors.New(cerrors.ERROR_INVALID_REQUEST, "Failed to recover Token Claims from Context!", "FE", true)
		l4g.Error("req #%v: handlers: VibeBeDevicesUnleashImeiPost: %s", context.GetReqId(r), cerr.Cause)
		handlers_common.RespondError(r, w, cerr)
		return
	}

	var jwt_claims jwttoken.Claims

	if jwt_claims, err = jwttoken.ParseClaims(jwt_claims_str); err != nil {
		cerr := cerrors.New(cerrors.ERROR_INVALID_REQUEST, fmt.Sprintf("Failed to parse Token Claims: '%s'!", err), "FE", true)
		l4g.Error("req #%v: handlers: VibeBeDevicesUnleashImeiPost: Failed to parse claims: %s. Claims content: '%s'", context.GetReqId(r), err.Error(), jwt_claims)
		handlers_common.RespondError(r, w, cerr)
		return
	}
	// Sending to Messaging Queue for request processing
	request := mq.NewMQRequest(imei, mq.MQRT_CREATE, mq.MQSCOPE_RSL_BE_UNLEASH, context.GetReqId(r), mq.NewReqInfo(r))

	response, cerr := handlers_common.MQSendRecv(request)
	if cerr != nil {
		handlers_common.RespondError(r, w, cerr)
	}
	handlers_common.HTTPFeedbackCreateResult(r, w, fmt.Sprintf("/kc_rsl_be/v1.0/devices/unleash/%s", imei), response, nil)
	return
}

func VibeBeDevicesUpdateImeiPut(w http.ResponseWriter, r *http.Request) {
	if config.GetFEConfig().Common.Debug {
		l4g.Debug("req #%v: handlers: VibeBeDevicesUpdateImeiPut starts", context.GetReqId(r))
		defer l4g.Debug("req #%v: handlers:VibeBeDevicesUpdateImeiPut ends", context.GetReqId(r))
	}

	var err error
	// Getting both IMEI
	vars := mux.Vars(r)
	imei := vars["imei"]

	body := context.GetPayload(r)

	if len(body) == 0 {
		cerr := cerrors.New(cerrors.ERROR_MALFORMED_POST_DATA, "No data found in request payload", "FE", true)
		l4g.Error("req #%v: handlers: VibeBeDevicesUpdateImeiPut: '%s'", context.GetReqId(r), cerr.Cause)
		handlers_common.RespondError(r, w, cerr)
		return
	}

	device_info, err := vibe_v1_0.JsonToDeviceInfo(body)
	if err != nil {
		cerr := cerrors.New(cerrors.ERROR_MALFORMED_POST_DATA, fmt.Sprintf("unmarshal failed: '%s'", err.Error()), "FE", true)
		l4g.Error("req #%v: handlers: VibeBeDevicesUpdateImeiPut: Error while unmarshalling body data: '%s'. offending payload: '%s'", context.GetReqId(r), err.Error(), body)
		handlers_common.RespondError(r, w, cerr)
		return
	}

	device_info.Imei = imei

	if err := device_info.SelfCheck(); err != nil {
		cerr := cerrors.New(cerrors.ERROR_INVALID_REQUEST, fmt.Sprintf("selfcheck failed: '%s'", err.Error()), "FE", true)
		l4g.Error("req #%v: handlers: VibeBeDevicesRegisterPost: Error in object data SelfCheck: '%s'.", context.GetReqId(r), err.Error())
		handlers_common.RespondError(r, w, cerr)
		return
	}

	jwt_claims_str := context.GetClaims(r)
	if jwt_claims_str == "" {
		cerr := cerrors.New(cerrors.ERROR_INVALID_REQUEST, "Failed to recover Token Claims from Context!", "FE", true)
		l4g.Error("req #%v: handlers: VibeBeDevicesUpdateImeiPut: %s", context.GetReqId(r), cerr.Cause)
		handlers_common.RespondError(r, w, cerr)
		return
	}

	var jwt_claims jwttoken.Claims

	if jwt_claims, err = jwttoken.ParseClaims(jwt_claims_str); err != nil {
		cerr := cerrors.New(cerrors.ERROR_INVALID_REQUEST, fmt.Sprintf("Failed to parse Token Claims: '%s'!", err), "FE", true)
		l4g.Error("req #%v: handlers: VibeBeDevicesUpdateImeiPut: Failed to parse claims: %s. Claims content: '%s'", context.GetReqId(r), err.Error(), jwt_claims)
		handlers_common.RespondError(r, w, cerr)
		return
	}

	partner_id := jwt_claims.Get("pid")
	// Sending to Messaging Queue for request processing
	request := mq.NewMQRequest(partner_id.(string), mq.MQRT_UPDATE, mq.MQSCOPE_RSL_BE_UPDATE_2I, context.GetReqId(r), mq.NewReqInfo(r))

	newBody, err := json.Marshal(device_info)
	if err != nil {
		cerr := cerrors.New(cerrors.ERROR_INTERNAL_SERVER_ERROR, fmt.Sprintf("Failed to marshal request body: %s", err), "FE", true)
		l4g.Error("req #%v: handlers: VibeBeDevicesUpdateImeiPut: Failed to parse claims: %s. Claims content: '%s'", context.GetReqId(r), err.Error(), jwt_claims)
		handlers_common.RespondError(r, w, cerr)
		return
	}
	//Set body of the request
	if err = request.SetData(newBody); err != nil {
		cerr := cerrors.New(cerrors.ERROR_INTERNAL_SERVER_ERROR, fmt.Sprintf("Failed to set request data: '%s'", err.Error()), "FE", true)
		l4g.Error("req #%v: handlers: VibeBeDevicesRegisterPost: Failed to set request data: '%s'.", context.GetReqId(r), err.Error())
		handlers_common.RespondError(r, w, cerr)
		return
	}

	response, cerr := handlers_common.MQSendRecv(request)
	if cerr != nil {
		handlers_common.RespondError(r, w, cerr)
		return
	}
	handlers_common.HTTPFeedbackCreateResult(r, w, fmt.Sprintf("/kc_rsl_be/v1.0/devices/update/%s", imei), response, nil)
	return
}
