package com.kaiostech.com.db.vibe;

import com.datastax.driver.core.Row;
import com.kaiostech.cerrors.CError;
import com.kaiostech.cerrors.CException;
import com.kaiostech.db.nosql.INoSqlDB_C;
import com.kaiostech.db.nosql.NoSqlDBFactory_C;
import com.kaiostech.model.financier.FinPartner;
import java.util.List;
import java.util.Map;
import java.util.TreeMap;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class RslPartnerDB {
  private static final String FIN_PARTNER_TABLE_NAME = "fin_partner";
  private static Logger _logger = LogManager.getLogger(RslPartnerDB.class);
  private INoSqlDB_C _nosqldb;
  private static RslPartnerDB _RslPartnerDB;

  private String EMAIL_D_COLUMN_NAME = "email_domain";
  private String FID_COLUMN_NAME = "fid";

  static {
    INoSqlDB_C nosqldb = NoSqlDBFactory_C.getDefault();
    _RslPartnerDB = new RslPartnerDB(nosqldb);
  }

  public static RslPartnerDB getInstance() {
    return _RslPartnerDB;
  }

  public RslPartnerDB(INoSqlDB_C nosqldb) {
    this._nosqldb = nosqldb;
  }

  public boolean connect() throws CException {
    try {
      this._nosqldb.connect();
    } catch (CException e) {
      _logger.error(
          "FinProfileDB: Connection error while connecting to database: '" + e.toString() + "'!");
      throw new CException(
          CError.New(
              CError.ERROR_INTERNAL_SERVER_ERROR, e.getMessage(), "DL:FinProfileDB:connect"));
    }
    return true;
  }

  public void close() {
    if (this._nosqldb != null) {
      this._nosqldb.close();
    }
  }

  public FinPartner createFinPartner(FinPartner fp, String req_id) throws CException {
    if (_logger.isDebugEnabled()) {
      _logger.debug("Req #" + req_id + ": DB: CreateFinPartner starts ");
    }
    Map<String, Object> query = new TreeMap<String, Object>();
    boolean ok;

    if (fp.getEmailDomain() != null) {
      query.put(EMAIL_D_COLUMN_NAME, fp.getEmailDomain());
    }

    if (fp.getFid() != null) {
      query.put(FID_COLUMN_NAME, fp.getFid());
    }

    ok = _nosqldb.setValue(req_id, FIN_PARTNER_TABLE_NAME, query);

    if (!ok) {
      return null;
    }
    return fp;
  }

  public FinPartner getFinPartner(String ed, String req_id) throws CException {
    if (_logger.isDebugEnabled()) {
      _logger.debug("Req #" + req_id + ": DB: GetFinPartner starts ");
    }

    FinPartner fp = new FinPartner();
    List<Row> results = null;
    Map<String, Object> query = new TreeMap<String, Object>();

    query.put(EMAIL_D_COLUMN_NAME, ed);

    results = _nosqldb.getValues(req_id, FIN_PARTNER_TABLE_NAME, query);
    if (results != null && !results.isEmpty()) {
      for (Row row : results) {
        fp.setEmailDomain(row.getString(EMAIL_D_COLUMN_NAME));
        fp.setFid(row.getString(FID_COLUMN_NAME));
      }
    } else {
      return null;
    }
    return fp;
  }

  public void deleteFinPartner(String ed, String req_id) throws CException {
    if (_logger.isDebugEnabled()) {
      _logger.debug(
          "Req #" + req_id + ": DB: DeleteFinPartner starts for email domain: '" + ed + "'!");
    }

    Map<String, Object> query = new TreeMap<String, Object>();
    boolean ok;
    query.put(EMAIL_D_COLUMN_NAME, ed);

    ok = _nosqldb.deleteRow(req_id, FIN_PARTNER_TABLE_NAME, query);
    if (!ok) {
      _logger.error(
          "Req #"
              + req_id
              + ": "
              + "DeleteFinPartner: Unable to delete Financier Partner with email domain: '"
              + ed
              + "'!");
    }
  }
}
