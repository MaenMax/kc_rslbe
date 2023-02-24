package com.kaiostech.kc_rslbe_dl3.actions.devices;

import com.kaiostech.cerrors.CError;
import com.kaiostech.cerrors.CException;
import com.kaiostech.db.vibe.VibeBeDb;
import com.kaiostech.kc_rslbe_dl3.actions.IAction;
import com.kaiostech.kc_rslbe_dl3.actions.IAutomataLine;
import com.kaiostech.model.core.JWT;
import com.kaiostech.mq.*;
import com.kaiostech.mq.codec.IMQCodec;
import com.kaiostech.mq.codec.MQCodecFactory;
import com.kaiostech.utils.ErrorUtils;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

/** @author Maen Created - 11/09/2022 */
public class TransferStateActions implements IAutomataLine {

  private IAction[] _actions;
  private static TransferStateActions _instance = null;

  private TransferStateActions() {
    int i;
    int max = MQRequestType.MQRT_LAST.getValue();
    _actions = new IAction[max + 1];

    for (i = 0; i < max + 1; i++) {
      _actions[i] = null;
    }

    _actions[MQRequestType.MQRT_UPDATE.getValue()] = new TransferState();
  }

  public static TransferStateActions getInstance() {
    if (_instance == null) {
      _instance = new TransferStateActions();
    }

    return _instance;
  }

  public IAction[] getAutomataLine() {
    return _actions;
  }
}

class TransferState implements IAction {
  private static Logger _logger = LogManager.getLogger(TransferState.class);
  private static VibeBeDb _db = VibeBeDb.getInstance();
  private static IMQCodec _codec = MQCodecFactory.getDefault();

  public MQResponse execute(int id, JWT jwt, MQRequest request) {

    if (_logger.isDebugEnabled()) {
      _logger.debug("Req #" + request.ReqId + ": " + "TransferState starts");
    }
    byte[] data;
    CError err;
    String[] objectIDs = request.ObjectId.split(":", 5);

    if (objectIDs.length != 3) {
      _logger.error("Req #" + request.ReqId + " :TransferState: Invalid object ID");
      err = CError.New(CError.ERROR_MALFORMED_PARAMETER, "Malformed data: Invalid object ID");
      data = _codec.encode(err);
      return new MQResponse(
          request.Id,
          "",
          MQRequestType.MQRT_UPDATE,
          MQRequestScope.MQRS_RSL_BE_TRANSFER_STATE,
          data,
          true,
          request.ReqId);
    }
    _logger.debug("MAEN object IDs: " + request.ObjectId);
    String partnerID = objectIDs[0];
    String fromIMEI = objectIDs[1];
    String toIMEI = objectIDs[2];

    if (partnerID.length() == 0) {

      _logger.error(
          "Req #" + request.ReqId + " :TransferState: Malformed data, Partner ID not provided");
      err =
          CError.New(CError.ERROR_MALFORMED_PARAMETER, "Malformed data: Partner ID not provided!");
      data = _codec.encode(err);
      return new MQResponse(
          request.Id,
          "",
          MQRequestType.MQRT_UPDATE,
          MQRequestScope.MQRS_RSL_BE_TRANSFER_STATE,
          data,
          true,
          request.ReqId);
    }

    if (fromIMEI.length() == 0) {
      _logger.error(
          "Req #"
              + request.ReqId
              + " :TransferState: Malformed data, from IMEI value is not provided");
      err =
          CError.New(
              CError.ERROR_MALFORMED_PARAMETER, "Malformed data: from IMEI value is not provided!");
      data = _codec.encode(err);
      return new MQResponse(
          request.Id,
          "",
          MQRequestType.MQRT_UPDATE,
          MQRequestScope.MQRS_RSL_BE_TRANSFER_STATE,
          data,
          true,
          request.ReqId);
    }

    if (partnerID.length() == 0) {
      _logger.error(
          "Req #"
              + request.ReqId
              + " :TransferState: Malformed data, to IMEI value is not provided");
      err =
          CError.New(
              CError.ERROR_MALFORMED_PARAMETER, "Malformed data: to IMEI value is not provided!");
      data = _codec.encode(err);
      return new MQResponse(
          request.Id,
          "",
          MQRequestType.MQRT_UPDATE,
          MQRequestScope.MQRS_RSL_BE_TRANSFER_STATE,
          data,
          true,
          request.ReqId);
    }

    boolean result = false;

    try {
      result = _db.transferState(partnerID, fromIMEI, toIMEI, request.ReqId);
    } catch (CException e) {
      _logger.error(
          "Req #"
              + request.ReqId
              + ": "
              + "ReqProc #"
              + id
              + ": Received request "
              + request.Id
              + ": "
              + e.getMessage());
      return ErrorUtils.sendErrorResponse(
          e.getCError(),
          request.Id,
          MQRequestType.MQRT_UPDATE,
          MQRequestScope.MQRS_RSL_BE_TRANSFER_STATE,
          request.ReqId);
    }
    if (result == false) {
      err =
          CError.New(
              CError.ERROR_NOT_FOUND,
              "Failed to transfer device state to the other device. One or two of the devices were not found in DB",
              "DL:TransferState:update");
      data = _codec.encode(err);
      return new MQResponse(
          request.Id,
          "",
          MQRequestType.MQRT_UPDATE,
          MQRequestScope.MQRS_RSL_BE_TRANSFER_STATE,
          data,
          false,
          request.ReqId);
    }
    if (_logger.isDebugEnabled()) {
      _logger.debug("Req #" + request.ReqId + ": " + "TransferState ends");
    }
    return new MQResponse(
        request.Id, null, request.Type, request.Scope, null, false, request.ReqId);
  }
}
