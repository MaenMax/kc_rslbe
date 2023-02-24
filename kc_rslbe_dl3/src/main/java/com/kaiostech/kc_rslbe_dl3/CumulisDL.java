package com.kaiostech.kc_rslbe_dl3;

import com.kaiostech.config.*;
import com.kaiostech.healthcheck.HealthCheckServer;
import com.kaiostech.healthcheck.HealthCheckServlet;
import com.kaiostech.mq.MQInit;
import com.kaiostech.mq.MQRequest;
import com.kaiostech.mq.MQRequestProvider;
import com.kaiostech.mq.MQServer;
import com.kaiostech.version.Version;
import java.util.Base64;
import java.util.LinkedList;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.apache.logging.log4j.core.LoggerContext;
import org.apache.logging.log4j.core.config.Configurator;

// import java.util.Base64.Decoder;

class ShutdownHook extends Thread {
  private static Logger _logger = LogManager.getLogger(ShutdownHook.class);
  private CumulisDL _dl;

  public ShutdownHook(CumulisDL dl) {
    this._dl = dl;
  }

  @Override
  public void run() {
    _logger.info("Shutdown hook triggered ... ");
    _dl.shutdown();
    _dl.wait_for_completion();
    _logger.info("Shutdown hook completed ... ");

    // Shutting down log4j2
    if (LogManager.getContext() instanceof LoggerContext) {
      _logger.info("Shutting down log4j2");
      Configurator.shutdown((LoggerContext) LogManager.getContext());
    }
  }
}

/**
 * This is the entry point and main thread of the Data Layer executable.
 *
 * <p>It implements the MQRequestProvider interface in order to be able to let worker threads to
 * execute the received requests.
 */
public class CumulisDL extends Thread implements MQRequestProvider {
  private static Logger _logger = LogManager.getLogger(CumulisDL.class);
  private static CumulisConfig _config = null;
  private static CumulisDL _dl = null;
  private static ShutdownHook _shutdown_hook;

  private static String Default_Config_Path = "conf/kc_rslbe.conf";
  private static String Default_Log_Path = "conf/kc_rslbe_dl3_log.xml";
  private static String Keys_Json_File = "keys.json";

  private LinkedList<MQRequest> _l_request;
  private Boolean _keep_running;
  private Boolean _shutdown = false;
  private HealthCheckServer _health_check_server = new HealthCheckServer();
  public MQServer _mqserver;

  public CumulisDL() {
    _keep_running = true;
    _shutdown = false;
    _l_request = new LinkedList<MQRequest>();
  }

  public static int Parse_Args(String args[]) {
    int err = 0;

    for (String s : args) {
      if (s.startsWith("--config")) {
        String[] parts = s.split("=");

        if (parts[1] == null) {
          System.err.println(
              "Syntax Error with option '" + s + "'. --{option}={value} form is expected.");
          err++;
          continue;
        }

        if (parts[1].length() == 0) {
          System.err.println(
              "Syntax Error with option '" + s + "'. --{option}={value} form is expected.");
          err++;
          continue;
        }

        Default_Config_Path = parts[1];
      } else if (s.startsWith("--log")) {
        String[] parts = s.split("=");

        if (parts[1] == null) {
          System.err.println(
              "Syntax Error with option '" + s + "'. --{option}={value} form is expected.");
          err++;
          continue;
        }

        if (parts[1].length() == 0) {
          System.err.println(
              "Syntax Error with option '" + s + "'. --{option}={value} form is expected.");
          err++;
          continue;
        }

        Default_Log_Path = parts[1];
        System.err.println(
            "Warning: --log not implemented yet. File '" + Default_Log_Path + "' will be ignored!");
      } else if (s.startsWith("--key_file")) {
        String[] parts = s.split("=");

        if (parts[1] == null) {
          System.err.println(
              "Syntax Error with option '" + s + "'. --{option}={value} form is expected.");
          err++;
          continue;
        }

        if (parts[1].length() == 0) {
          System.err.println(
              "Syntax Error with option '" + s + "'. --{option}={value} form is expected.");
          err++;
          continue;
        }

        Keys_Json_File = parts[1];
      } else if (s.startsWith("--version") || (s.startsWith("-version"))) {
        System.out.print(Version.getFullVersion());
        System.exit(0);
      }
    }

    return err;
  }

  public static void main(String args[]) {

    Parse_Args(args);

    _logger.debug(
        "Loading Config from: " + Default_Config_Path + " and keys json file " + Keys_Json_File);
    _config = CumulisConfig.loadConfig(Default_Config_Path, Keys_Json_File);

    if (_config == null) {
      System.err.println(
          "Failed to load configuration file '" + Default_Config_Path + "'! Aborting ...");
      System.exit(1);
    }

    _logger.info("Initialized ... ");

    // This is required to initialize the MQ module (we have
    // to call it after having read the configuration file).
    MQInit.init(_config);

    //
    if (CumulisConfig.get() != null) {
      if (CumulisConfig.get().getValue("fernet_key") == null) {
        System.err.println(
            "The fernet_key value used to encrypt database data cannot be read from the key file! Did you forget to provide a key file with fernet_key value defined? Aborting ...");
        System.exit(10);
      }
      if (CumulisConfig.get().getValue("fernet_iv") == null) {
        System.err.println(
            "The fernet_iv value used to encrypt database data cannot be read from the key file! Did you forget to provide a key file with fernet_iv value defined? Aborting ...");
        System.exit(10);
      }

      final byte iv_buff[] = Base64.getDecoder().decode(CumulisConfig.get().getValue("fernet_iv"));

      com.kaiostech.db.codec.DBCborFernetCodecFactory.setKey(
          CumulisConfig.get().getValue("fernet_key"), iv_buff);
      com.kaiostech.db.codec.DBCodecFactory.setKey(
          CumulisConfig.get().getValue("fernet_key"), iv_buff);
    }

    _dl = new CumulisDL();
    _shutdown_hook = new ShutdownHook(_dl);

    Runtime.getRuntime().addShutdownHook(_shutdown_hook);

    _logger.info("Start processing ... ");

    _dl.setName("CumulisDL3");
    _dl.start();

    while (true) {
      try {
        Thread.sleep(100);
      } catch (InterruptedException e) {
      }
    }
  }

  public void run() {
    int i;
    AutomataRequestProcessor proc;
    LinkedList<AutomataRequestProcessor> l_processor;

    MQRequest req;

    _logger.info("Starting CumulisDL server ... ");

    _mqserver = new MQServer(_config.QueueService);

    l_processor = new LinkedList<>();

    for (i = 0; i < _config.DataLayer.NbRequestWorkerPerNode; i++) {
      proc =
          new AutomataRequestProcessor(i, this, this._mqserver, _config.QueueService.LLRspTimeOut);
      proc.start();
      proc.setName("ReqProc #" + (i + 1));
      l_processor.offer(proc);
    }

    if (_config.FrontLayer != null && _config.FrontLayer.HttpPort > 0) {
      HealthCheckServlet.SetHealthStatusVersionName("kc_rslbe_dl", Version.Get_Version());
      try {
        this._health_check_server.start(_config.FrontLayer.HttpPort + 3000);
      } catch (Exception e) {
        e.printStackTrace();
      }
    }

    while (_keep_running) {

      _logger.info("Server connecting to NATS ...");

      if (!_mqserver.connect()) {
        _logger.error("[FAILED] Attempting again in few seconds.");
        try {
          Thread.sleep(5000);
        } catch (InterruptedException e) {
        }
        continue;
      }

      _logger.info("Server successfully connected to NATS! Now listening ...");

      if (!_mqserver.listen()) {
        _mqserver.disconnect();
        _logger.error("[FAILED] Attempting again in few seconds.");
        try {
          Thread.sleep(5000);
        } catch (InterruptedException e) {
        }
        continue;
      }

      _logger.info("Server successfully listening for incoming requests!");

      while (_keep_running) {
        if (!_mqserver.isConnected()) {
          _logger.error(
              "Disconnection detected! Leaving main processing loop for reconnection ...");
          break;
        }

        req = _mqserver.getRequest();

        if (req != null) {
          deliver(req);
          continue;
        }

        // Since we had nothing to do during this loop, sleeping a while to yield CPU to
        // other threads.
        try {
          Thread.sleep(10);
        } catch (InterruptedException e) {
        }
      } // while(true) Level 2
    } // while(true) Level 1

    _logger.info("Shutting down  ...");

    // First unsubscribing from subject to avoid newly incoming requests.
    _mqserver.unsubscribe();

    // Second, making sure that all the pending requests have been processed (or
    // currently being processed).
    while (!is_empty()) {
      try {
        Thread.sleep(100);
      } catch (InterruptedException e) {
      }
    }

    // Third, sending shutdown signal to request processors
    for (AutomataRequestProcessor p : l_processor) {
      p.shutdown();
    }

    // Fourth, waiting for these request processors to complete their tasks and
    // shutdown.
    while (true) {
      Boolean all_down = true;

      try {
        Thread.sleep(100);
      } catch (InterruptedException e) {
      }

      for (AutomataRequestProcessor p : l_processor) {
        if (!p.is_stopped()) {
          all_down = false;
          break;
        }
      }

      if (all_down) {
        break;
      }
    }

    _logger.info("All requests completed  ...");

    _mqserver.disconnect();

    if (_config.FrontLayer != null && _config.FrontLayer.HttpPort > 0) {
      try {
        this._health_check_server.stop();
      } catch (Exception e) {
        e.printStackTrace();
      }
    }

    _logger.info("Daemon server shutdown ...");

    _shutdown = true;
  }

  private void deliver(MQRequest req) {
    synchronized (_l_request) {
      _l_request.offer(req);
    }
  }

  /**
   * Any worker in charge of executing request can use this function in order to retrieve a request
   * to process. If several workers are active (which is recommended), then one request will be
   * executed by only one worker.
   *
   * <p>Return: A pending request to process if successful else NULL if no request to process is
   * available.
   */
  public MQRequest get() {
    MQRequest res = null;

    synchronized (_l_request) {
      if (_l_request.size() > 0) {
        res = _l_request.poll();
      }
    }

    return res;
  }

  /** Check whether the request queue is empty or not. */
  public Boolean is_empty() {
    Boolean res = true;
    synchronized (_l_request) {
      if (_l_request.size() > 0) {
        res = false;
      }
    }
    return res;
  }

  public void shutdown() {
    _keep_running = false;
  }

  public Boolean is_shutdown() {
    return _shutdown;
  }

  public void wait_for_completion() {
    long cur_time;
    long ref_time;

    ref_time = System.currentTimeMillis();
    while (!_shutdown) {
      try {
        Thread.sleep(100);
      } catch (InterruptedException e) {
      }
      cur_time = System.currentTimeMillis();

      if (cur_time - ref_time > 30000) {
        _logger.warn("Forcing shutdown due to timeout ...");
        break;
      }
    }
  }
}
