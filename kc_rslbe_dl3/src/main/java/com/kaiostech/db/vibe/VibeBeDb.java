package com.kaiostech.db.vibe;

import com.datastax.driver.core.Row;
import com.kaiostech.cerrors.CError;
import com.kaiostech.cerrors.CException;
import com.kaiostech.db.nosql.INoSqlDB_C;
import com.kaiostech.db.nosql.NoSqlDBFactory_C;
import com.kaiostech.model.vibe.DeviceInfo;
import java.nio.ByteBuffer;
import java.util.List;
import java.util.Map;
import java.util.TreeMap;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class VibeBeDb {
  // Table names
  private static final String VIBE_PARTNER_TABLE_NAME = "vibe_partner";
  private static final String VIBE_DEVICES_TABLE_NAME = "devices";
  private static final String VIBE_IMEI2PARTNER_TABLE_NAME = "imei2partner";

  private static Logger _logger = LogManager.getLogger(VibeBeDb.class);
  private INoSqlDB_C _nosqldb;
  private static VibeBeDb _VibeBeDb;

  // Column names
  private String PARTNER_ID_COLUMN_NAME = "partner_id";
  private String IMEI_COLUMN_NAME = "imei";
  private String ISDN_COLUMN_NAME = "isdn";
  private String IMSI_COLUMN_NAME = "imsi";
  private String DEVICE_ID_COLUMN_NAME = "id";
  private String DEVICE_ACTUAL_RSL_STATE_COLUMN_NAME = "actual_rsl_state";

  static {
    INoSqlDB_C nosqldb = NoSqlDBFactory_C.getDefault();
    _VibeBeDb = new VibeBeDb(nosqldb);
  }

  public static VibeBeDb getInstance() {
    return _VibeBeDb;
  }

  public VibeBeDb(INoSqlDB_C nosqldb) {
    this._nosqldb = nosqldb;
  }

  public boolean connect() throws CException {
    try {
      this._nosqldb.connect();
    } catch (CException e) {
      _logger.error(
          "FinProfileDB: Connection error while connecting to database: '" + e.toString() + "'!");
      throw new CException(
          CError.New(CError.ERROR_INTERNAL_SERVER_ERROR, e.getMessage(), "DL:VibeBeDb:connect"));
    }
    return true;
  }

  public void close() {
    if (this._nosqldb != null) {
      this._nosqldb.close();
    }
  }

  public boolean registerImeiForPartner(String pid, String imei, String req_id) throws CException {
    if (_logger.isDebugEnabled()) {
      _logger.debug("Req #" + req_id + ": DB: registerImeiForPartner starts ");
    }
    Map<String, Object> query = new TreeMap<String, Object>();
    boolean ok;

    query.put(PARTNER_ID_COLUMN_NAME, pid);
    query.put(IMEI_COLUMN_NAME, imei);

    ok = _nosqldb.setValue(req_id, VIBE_IMEI2PARTNER_TABLE_NAME, query);

    if (!ok) {
      return false;
    }
    return true;
  }

  public DeviceInfo register3IData(String partner_id, DeviceInfo dev_info, String req_id)
      throws CException {
    if (_logger.isDebugEnabled()) {
      _logger.debug("Req #" + req_id + ": DB: Register3IData starts ");
    }
    boolean result = false;
    Map<String, Object> query = new TreeMap<String, Object>();
    boolean ok;

    query.put(PARTNER_ID_COLUMN_NAME, partner_id);
    query.put(IMEI_COLUMN_NAME, dev_info.getImei());
    query.put(IMSI_COLUMN_NAME, dev_info.getImsi());
    query.put(ISDN_COLUMN_NAME, dev_info.getIsdn());

    ok = _nosqldb.setValue(req_id, VIBE_PARTNER_TABLE_NAME, query);

    if (!ok) {
      return null;
    }

    // Now, registr the record in imei2partner tabel

    try {
      result = registerImeiForPartner(partner_id, dev_info.getImei(), req_id);
    } catch (CException e) {
      _logger.error(
          "Req #"
              + req_id
              + ": "
              + "Error while adding a record to imei2partner table"
              + e.getMessage());
    }

    if (result == false) {
      _logger.error("Req #" + req_id + ": " + "Could not add a record in imei2partner table");
      return null;
    }
    return dev_info;
  }

  public DeviceInfo update2IData(String partner_id, DeviceInfo dev_info, String req_id)
      throws CException {
    if (_logger.isDebugEnabled()) {
      _logger.debug("Req #" + req_id + ": DB: Register3IData starts ");
    }

    if (dev_info == null) {
      _logger.error(
          "Req #"
              + req_id
              + ": "
              + "#"
              + Thread.currentThread().getId()
              + ": '"
              + "Missing or empty device_info"
              + "'!");
      throw new CException(
          CError.New(
              CError.ERROR_INVALID_PARAMETER_VALUE,
              "Missing or empty device_info.",
              "DL:VibeBeDB:update2IData"));
    }

    List<Row> results = null;
    results = map2List(req_id, VIBE_PARTNER_TABLE_NAME, PARTNER_ID_COLUMN_NAME, partner_id);
    if (results == null || results.isEmpty()) {
      _logger.error(
          "Req #"
              + req_id
              + ": "
              + "#"
              + Thread.currentThread().getId()
              + ": '"
              + "Device_Id '"
              + dev_info.getImei()
              + "'Partner ID is not found in DB"
              + "'!");
      throw new CException(
          CError.New(
              CError.ERROR_NOT_FOUND,
              "partner ID '" + partner_id + "'Is not found in DB: ",
              "DL:VibeBeDB:update2IData"));
    }

    // At this point, the device is found in DB, we can now proceed to update the 2I data for the
    // device.

    Map<String, Object> query = new TreeMap<String, Object>();
    boolean ok;

    query.put(PARTNER_ID_COLUMN_NAME, partner_id);
    query.put(IMEI_COLUMN_NAME, dev_info.getImei());
    query.put(IMSI_COLUMN_NAME, dev_info.getImsi());
    query.put(ISDN_COLUMN_NAME, dev_info.getIsdn());

    ok = _nosqldb.setValue(req_id, VIBE_PARTNER_TABLE_NAME, query);

    if (!ok) {
      return null;
    }
    return dev_info;
  }

  public boolean transferState(String partner_id, String fromImei, String toIMEI, String req_id)
      throws CException {
    if (_logger.isDebugEnabled()) {
      _logger.debug("Req #" + req_id + ": DB: transferState starts ");
    }

    Map<String, Object> query = new TreeMap<String, Object>();
    boolean ok;

    List<Row> results = null;
    results = map2List(req_id, VIBE_DEVICES_TABLE_NAME, DEVICE_ID_COLUMN_NAME, fromImei);
    if (results == null || results.isEmpty()) {
      _logger.error(
          "Req #"
              + req_id
              + ": "
              + "#"
              + Thread.currentThread().getId()
              + ": '"
              + "Device_Id '"
              + fromImei
              + "'Device to transfer the state from, was not found in DB"
              + "'!");
      throw new CException(
          CError.New(
              CError.ERROR_NOT_FOUND,
              "Device ID: '" + fromImei + "'Is not found in DB: ",
              "DL:VibeBeDB:transferState"));
    }

    // At this point, the device/IMEI to transfer the state from, exists in DB. We now can proceed
    // to check for the second IMEI.

    List<Row> results2 = null;
    results2 = map2List(req_id, VIBE_DEVICES_TABLE_NAME, DEVICE_ID_COLUMN_NAME, toIMEI);
    if (results2 == null || results.isEmpty()) {
      _logger.error(
          "Req #"
              + req_id
              + ": "
              + "#"
              + Thread.currentThread().getId()
              + ": '"
              + "Device_Id '"
              + toIMEI
              + "'Device to transfer the state to, was not found in DB"
              + "'!");
      throw new CException(
          CError.New(
              CError.ERROR_NOT_FOUND,
              "IMEI '" + toIMEI + "'Does not exist in DB: ",
              "DL:VibeBeDB:transferState"));
    }
    // byte[] device_actual_rsl_state_bytes;

    // Processing results:
    for (Row row : results) {
      try {
        ByteBuffer bb_a = row.getBytes(DEVICE_ACTUAL_RSL_STATE_COLUMN_NAME);
        byte[] device_actual_rsl_state_bytes = bb_a.array();

        // Trasfering the state to the second IMEI
        query.put(DEVICE_ID_COLUMN_NAME, toIMEI);
        query.put(DEVICE_ACTUAL_RSL_STATE_COLUMN_NAME, device_actual_rsl_state_bytes);

        ok = _nosqldb.setValue(req_id, VIBE_DEVICES_TABLE_NAME, query);
        if (!ok) {
          _logger.debug(
              "Req #"
                  + req_id
                  + ": DB: transferState: IMEI: "
                  + fromImei
                  + "to transfer the state from, does not exist in devices table!");
          return false;
        }
        return true;
      } catch (Exception e) {
        _logger.debug("JSONException e: " + e);
        return false;
      }
    }
    return false;
  }

  public boolean transferOwnership(String target_partner_id, String imei, String req_id)
      throws CException {
    if (_logger.isDebugEnabled()) {
      _logger.debug("Req #" + req_id + ": DB: transferOwnership starts ");
    }

    boolean ok;
    List<Row> results = null;
    Map<String, Object> query = new TreeMap<String, Object>();

    // [1] check if the partner exists in DB for that IMEI by quering imei2partner table.
    results = map2List(req_id, VIBE_IMEI2PARTNER_TABLE_NAME, IMEI_COLUMN_NAME, imei);
    if (results == null || results.isEmpty()) {
      _logger.error(
          "Req #"
              + req_id
              + ": "
              + "#"
              + Thread.currentThread().getId()
              + ": '"
              + "Device_Id '"
              + imei
              + "'No partner found for the imei: "
              + imei);
      throw new CException(
          CError.New(
              CError.ERROR_NOT_FOUND,
              "IMEI '" + imei + "'Does not exist in DB: ",
              "DL:VibeBeDB:transferOwnership"));
    }
    // byte[] device_actual_rsl_state_bytes;

    // Processing results:
    for (Row row : results) {
      try {
        String original_partner_id = row.getString(PARTNER_ID_COLUMN_NAME);
        // [2] Get the device Info from devices table

        DeviceInfo dev_info = readDeviceInfo(original_partner_id, imei, req_id);

        if (dev_info == null) {
          return false;
        }

        // [3] Add device info to tagret partner.
        dev_info = register3IData(target_partner_id, dev_info, req_id);
        if (dev_info == null) {
          return false;
        }

        // [4] Delete the IMEI from original partner

        query.put(PARTNER_ID_COLUMN_NAME, original_partner_id);
        query.put(IMEI_COLUMN_NAME, imei);
        ok = _nosqldb.deleteRow(req_id, VIBE_PARTNER_TABLE_NAME, query);

        // [5] Modify record in imei2partner to reflect on the new change.
        query.put(IMEI_COLUMN_NAME, imei);
        query.put(PARTNER_ID_COLUMN_NAME, target_partner_id);
        ok = _nosqldb.setValue(req_id, VIBE_IMEI2PARTNER_TABLE_NAME, query);
        if (!ok) {
          return false;
        }
        return true;
      } catch (Exception e) {
        _logger.debug("Error while transfering device ownership: " + e);
        return false;
      }
    }
    return false;
  }

  public DeviceInfo readDeviceInfo(String partner_id, String imei, String req_id)
      throws CException {
    if (_logger.isDebugEnabled()) {
      _logger.debug("Req #" + req_id + ": DB: readDeviceInfo starts ");
    }

    DeviceInfo dev_info = new DeviceInfo();

    List<Row> results = null;
    results =
        cmpdKey2List(
            req_id,
            VIBE_PARTNER_TABLE_NAME,
            PARTNER_ID_COLUMN_NAME,
            partner_id,
            IMEI_COLUMN_NAME,
            imei);
    if (results == null || results.isEmpty()) {
      _logger.error(
          "Req #"
              + req_id
              + ": "
              + "#"
              + Thread.currentThread().getId()
              + ": '"
              + "Device_Id '"
              + imei
              + "'Does not belong to a partner"
              + "'!");
      throw new CException(
          CError.New(
              CError.ERROR_NOT_FOUND,
              "Device_Id '" + imei + "'Does not belong to a petner",
              "DL:readDeviceInfo"));
    }
    // Processing results:
    for (Row row : results) {
      try {
        String imsi = row.getString(IMSI_COLUMN_NAME);
        String isdn = row.getString(ISDN_COLUMN_NAME);
        dev_info.setImei(imei);
        dev_info.setImsi(imsi);
        dev_info.setIsdn(isdn);

        return dev_info;
      } catch (Exception e) {
        _logger.debug("Error while reading device info from devices table: " + e);
        return null;
      }
    }
    return null;
  }

  public List<Row> map2List(
      String req_id, String TABLE_NAME, String column_name, String column_value) throws CException {
    Map<String, Object> query = new TreeMap<String, Object>();
    query.put(column_name, column_value);
    List<Row> results = _nosqldb.getValues(req_id, TABLE_NAME, query);
    return results;
  }

  public List<Row> cmpdKey2List(
      String req_id,
      String TABLE_NAME,
      String partition_key_name,
      String partition_key_value,
      String clustering_key_name,
      String clustering_key_value)
      throws CException {
    Map<String, Object> query = new TreeMap<String, Object>();
    query.put(partition_key_name, partition_key_value);
    query.put(clustering_key_name, clustering_key_value);
    List<Row> results = _nosqldb.getValues(req_id, TABLE_NAME, query);
    return results;
  }
}
