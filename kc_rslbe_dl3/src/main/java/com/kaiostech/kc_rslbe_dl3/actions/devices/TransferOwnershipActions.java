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
public class TransferOwnershipActions implements IAutomataLine {

  private IAction[] _actions;
  private static TransferOwnershipActions _instance = null;

  private TransferOwnershipActions() {
    int i;
    int max = MQRequestType.MQRT_LAST.getValue();
    _actions = new IAction[max + 1];

    for (i = 0; i < max + 1; i++) {
      _actions[i] = null;
    }

    _actions[MQRequestType.MQRT_CREATE.getValue()] = new TransferOwnership();
  }

  public static TransferOwnershipActions getInstance() {
    if (_instance == null) {
      _instance = new TransferOwnershipActions();
    }

    return _instance;
  }

  public IAction[] getAutomataLine() {
    return _actions;
  }
}

class TransferOwnership implements IAction {
  private static Logger _logger = LogManager.getLogger(TransferOwnership.class);
  private static VibeBeDb _db = VibeBeDb.getInstance();
  private static IMQCodec _codec = MQCodecFactory.getDefault();

  public MQResponse execute(int id, JWT jwt, MQRequest request) {

    if (_logger.isDebugEnabled()) {
      _logger.debug("Req #" + request.ReqId + ": " + "TransferOwnership starts");
    }
    byte[] data;
    CError err;
    String[] objectIDs = request.ObjectId.split(":", 5);

    if (objectIDs.length != 2) {
      _logger.error("Req #" + request.ReqId + " :TransferOwnership: Invalid object ID");
      err = CError.New(CError.ERROR_MALFORMED_PARAMETER, "Malformed data: Invalid object ID");
      data = _codec.encode(err);
      return new MQResponse(
          request.Id,
          "",
          MQRequestType.MQRT_CREATE,
          MQRequestScope.MQRS_RSL_BE_TRANSFER_OWNERSHIP,
          data,
          true,
          request.ReqId);
    }
    _logger.debug("MAEN object IDs: " + request.ObjectId);
    String partnerID = objectIDs[0];
    String Imei = objectIDs[1];

    if (partnerID.length() == 0) {
      _logger.error(
          "Req #" + request.ReqId + " :TransferOwnership: Malformed data, Partner ID not provided");
      err =
          CError.New(CError.ERROR_MALFORMED_PARAMETER, "Malformed data: Partner ID not provided!");
      data = _codec.encode(err);
      return new MQResponse(
          request.Id,
          "",
          MQRequestType.MQRT_CREATE,
          MQRequestScope.MQRS_RSL_BE_TRANSFER_OWNERSHIP,
          data,
          true,
          request.ReqId);
    }

    if (Imei.length() == 0) {
      _logger.error(
          "Req #"
              + request.ReqId
              + " :TransferOwnership: Malformed data, IMEI value is not provided");
      err =
          CError.New(
              CError.ERROR_MALFORMED_PARAMETER, "Malformed data:IMEI value is not provided!");
      data = _codec.encode(err);
      return new MQResponse(
          request.Id,
          "",
          MQRequestType.MQRT_CREATE,
          MQRequestScope.MQRS_RSL_BE_TRANSFER_OWNERSHIP,
          data,
          true,
          request.ReqId);
    }

    boolean result = false;

    try {
      result = _db.transferOwnership(partnerID, Imei, request.ReqId);
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
          MQRequestType.MQRT_CREATE,
          MQRequestScope.MQRS_RSL_BE_TRANSFER_OWNERSHIP,
          request.ReqId);
    }
    if (result == false) {
      err =
          CError.New(
              CError.ERROR_NOT_FOUND,
              "Failed to transfer Ownership of the device to the specified partner. Imei or partner was not found in DB",
              "DL:TransferOwnership:update");
      data = _codec.encode(err);
      return new MQResponse(
          request.Id,
          "",
          MQRequestType.MQRT_CREATE,
          MQRequestScope.MQRS_RSL_BE_TRANSFER_OWNERSHIP,
          data,
          false,
          request.ReqId);
    }
    if (_logger.isDebugEnabled()) {
      _logger.debug("Req #" + request.ReqId + ": " + "TransferOwnership ends");
    }
    return new MQResponse(
        request.Id, null, request.Type, request.Scope, null, false, request.ReqId);
  }
}
