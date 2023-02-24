package com.kaiostech.kc_rslbe_dl3.actions.devices.update_2i;

import com.kaiostech.cerrors.CError;
import com.kaiostech.cerrors.CException;
import com.kaiostech.db.vibe.VibeBeDb;
import com.kaiostech.kc_rslbe_dl3.actions.IAction;
import com.kaiostech.kc_rslbe_dl3.actions.IAutomataLine;
import com.kaiostech.model.core.JWT;
import com.kaiostech.model.vibe.DeviceInfo;
import com.kaiostech.mq.*;
import com.kaiostech.mq.codec.IMQCodec;
import com.kaiostech.mq.codec.MQCodecFactory;
import com.kaiostech.utils.ErrorUtils;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

/** @author Maen Created - 11/09/2022 */
public class Update2I implements IAutomataLine {

  private IAction[] _actions;
  private static Update2I _instance = null;

  private Update2I() {
    int i;
    int max = MQRequestType.MQRT_LAST.getValue();
    _actions = new IAction[max + 1];

    for (i = 0; i < max + 1; i++) {
      _actions[i] = null;
    }

    _actions[MQRequestType.MQRT_UPDATE.getValue()] = new Update2IData();
  }

  public static Update2I getInstance() {
    if (_instance == null) {
      _instance = new Update2I();
    }

    return _instance;
  }

  public IAction[] getAutomataLine() {
    return _actions;
  }
}

class Update2IData implements IAction {
  private static Logger _logger = LogManager.getLogger(Update2IData.class);
  private static VibeBeDb _db = VibeBeDb.getInstance();
  private static IMQCodec _codec = MQCodecFactory.getDefault();

  public MQResponse execute(int id, JWT jwt, MQRequest request) {

    String partnerID = request.ObjectId;
    CError err;
    byte[] data;

    if (partnerID.length() == 0) {
      _logger.error(
          "Req #" + request.ReqId + " :Update2IData: Malformed data, Partner ID not provided");
      err =
          CError.New(CError.ERROR_MALFORMED_PARAMETER, "Malformed data: Partner ID not provided!");
      data = _codec.encode(err);
      return new MQResponse(
          request.Id,
          "",
          MQRequestType.MQRT_UPDATE,
          MQRequestScope.MQRS_RSL_BE_UPDATE_2I,
          data,
          true,
          request.ReqId);
    }

    if (_logger.isDebugEnabled()) {
      _logger.debug("Req #" + request.ReqId + ": " + "Update2IData starts");
    }

    MQDataResult dr = request.Data();

    if (dr.err != null) {
      _logger.error(
          "Req #"
              + request.ReqId
              + ": Update2IData aborted due to error while fetching request body data: '"
              + dr.err.toString()
              + "'!");
      data = _codec.encode(dr.err);
      return new MQResponse(request.Id, "", request.Type, request.Scope, data, true, request.ReqId);
    }
    DeviceInfo deviceInfo = new DeviceInfo();
    deviceInfo = _codec.decode(dr.data, DeviceInfo.class);
    if (deviceInfo == null) {
      _logger.error(
          "Req #"
              + request.ReqId
              + ": "
              + "ReqProc #"
              + id
              + ": Received request "
              + request.Id
              + " cannot be decoded.");
      if (_logger.isDebugEnabled()) {
        _logger.debug(
            "Req #"
                + request.ReqId
                + ": "
                + "ReqProc #"
                + id
                + ": Received request "
                + request.Id
                + " Data Dump: '"
                + new String(dr.data)
                + "'");
      }
      err =
          CError.New(CError.ERROR_MALFORMED_PARAMETER, "Malformed data: Cannot decode DeviceInfo!");
      data = _codec.encode(err);
      return new MQResponse(
          request.Id,
          "",
          MQRequestType.MQRT_UPDATE,
          MQRequestScope.MQRS_RSL_BE_UPDATE_2I,
          data,
          true,
          request.ReqId);
    }

    DeviceInfo result = new DeviceInfo();
    try {
      result = _db.update2IData(partnerID, deviceInfo, request.ReqId);
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
          MQRequestScope.MQRS_RSL_BE_UPDATE_2I,
          request.ReqId);
    }

    if (result == null) {
      err =
          CError.New(
              CError.ERROR_INTERNAL_SERVER_ERROR,
              "Failed to update 2I data of the device.",
              "DL:Update2I:update");
      data = _codec.encode(err);
      return new MQResponse(
          request.Id,
          "",
          MQRequestType.MQRT_UPDATE,
          MQRequestScope.MQRS_RSL_BE_UPDATE_2I,
          data,
          false,
          request.ReqId);
    }
    data = _codec.encode(result);

    if (_logger.isDebugEnabled()) {
      _logger.debug("Req #" + request.ReqId + ": " + "Update2IData ends");
    }
    return new MQResponse(
        request.Id, null, request.Type, request.Scope, data, false, request.ReqId);
  }
}
