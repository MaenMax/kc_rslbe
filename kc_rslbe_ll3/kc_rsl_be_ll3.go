package main

import (
	"log"
	"net/http"
	"runtime"
	"runtime/pprof"
	"time"

	"git.kaiostech.com/cloud/kc_rslbe/kc_rslbe_ll3/actions"

	common "git.kaiostech.com/cloud/common/utils/actions_common"

	"git.kaiostech.com/cloud/common/db/redisdb"

	"git.kaiostech.com/cloud/common/mq"
	"git.kaiostech.com/cloud/common/security/credkey"
	"git.kaiostech.com/cloud/common/security/jwttoken"
	"git.kaiostech.com/cloud/common/security/session"
	"git.kaiostech.com/cloud/common/security/token_blacklist"
	"git.kaiostech.com/cloud/common/utils/handlers_common"

	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"git.kaiostech.com/cloud/kc_rslbe/kc_rslbe_ll3/ll"

	"git.kaiostech.com/cloud/common/config"
	l4g "git.kaiostech.com/cloud/thirdparty/code.google.com/p/log4go"

	"git.kaiostech.com/cloud/common/healthcheck"
	"git.kaiostech.com/cloud/common/utils/health_check"
	"git.kaiostech.com/cloud/kc_rslbe/kc_rslbe_ll3/version"

	// For profiling
	//"net/http"
	_ "net/http/pprof"
)

const (
	NB_OF_WORKER                = 1
	NB_OF_CONN_HANDLERS         = 10
	SHUTDOWN_TIMEOUT            = 30
	NATS_RECONNECT_RETRY_PERIOD = 5 // seconds
)

var (
	_conf                           *config.FEConfig
	automata                        *actions.Automata
	ErrMaxActiveProcessCountReached = errors.New("Max Active Process Count reached")

	// Used for graceful shutdown
	c_kill      chan os.Signal
	c_int       chan os.Signal
	c_usr1      chan os.Signal
	l_c_stop_ll chan chan int
)

func init() {
	c_kill = make(chan os.Signal, 1)
	c_int = make(chan os.Signal, 1)
	c_usr1 = make(chan os.Signal, 1)
}

func health_check_callback() {

}

func Init() error {
	var err error

	// Making sure to have the latest configuration.
	_conf = config.GetFEConfig()

	common.Init(_conf.QueueService.Nodes, _conf.QueueService.FE2CP_Subject, _conf.QueueService.FE2LL_Subject, _conf.QueueService.LL2DL_Subject,
		_conf.QueueService.LL2SA_Subject)

	if err = common.Init_Idempotency(_conf.RedisService.Nodes, _conf.RedisService.MaxRedirects); err != nil {
		return errors.New(fmt.Sprintf("Common Init_Idempotency error - %s", err))
	}

	// Initializing handlers_common for Service center cache and filtering
	handlers_common.Init(_conf.QueueService.Nodes, _conf.QueueService.LL2DL_Subject)

	automata = actions.NewAutomata()

	mq.Init(_conf.QueueService.FlushTimeOut, _conf.QueueService.MaxMsgSize, _conf.QueueService.RetryNb, _conf.QueueService.LLRspTimeOut, _conf.Common.Debug, _conf)

	if err = redisdb.Init(_conf.Common.Debug); err != nil {
		return errors.New(fmt.Sprintf("RedisDB Init error - %s", err))
	}

	if err = session.Init(_conf.RedisService.Nodes, _conf.RedisService.MaxRedirects); err != nil {
		return errors.New(fmt.Sprintf("Session Init error - %s", err))
	}

	if err = credkey.Init(_conf.RedisService.Nodes, _conf.RedisService.MaxRedirects); err != nil {
		return errors.New(fmt.Sprintf("CredKey Init error - %s", err))
	}

	if err = jwttoken.Init(_conf); err != nil {
		return errors.New(fmt.Sprintf("JWTToken Init error - %s", err))
	}

	if err = token_blacklist.Init(_conf.RedisService.Nodes, _conf.RedisService.MaxRedirects); err != nil {
		return errors.New(fmt.Sprintf("Token BlackList Init error - %s", err))
	}

	if _conf.FrontLayer.HttpPort != 0 {
		healthcheck.InitHttpHealthCheck()
		healthcheck.SetHealthStatusVersionName("kc_rslbe_ll", version.Get_Version())
	}

	return nil
}

var (
	opt_version *bool   = flag.Bool("version", false, "Display current tool version and exit.")
	config_file *string = flag.String("config", "conf/kc_rslbe.conf", "Config file to use.")
	log_file    *string = flag.String("log", "conf/kc_rslbe_ll3_log.xml", "Log L4G config file to use.")
	key_file    *string = flag.String("key_file", "keys.json", "Keys used to boostrap fe node.")

	print_conf *bool = flag.Bool("print-conf", false, "Print current active configuration and exit.")

	cpu_prof           = flag.String("cpu-prof", "", "write cpu profile to `file`")
	mem_prof           = flag.String("mem-prof", "", "write memory profile to `file`")
	net_addr_port_prof = flag.String("port-prof", "localhost:6062", "The port of net/http/pprof")
)

func main() {
	// Sleep 1 second before the end in order to give time to l4g to sync its output
	// to the log file. Failing to wait will prevent the latest logged messages to be
	// output and flushed.
	defer time.Sleep(time.Second * 1)

	var err error
	var procs []*ll.T_ll_proc

	flag.Parse()

	if *cpu_prof != "" {
		f, err := os.Create(*cpu_prof)
		if err != nil {
			panic(fmt.Sprintf("could not create CPU profile: '%s'", err))
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			panic(fmt.Sprintf("could not start CPU profile: '%s'", err))
		}
		defer pprof.StopCPUProfile()
	}

	if *mem_prof != "" {
		f, err := os.Create(*mem_prof)
		if err != nil {
			panic(fmt.Sprintf("could not create memory profile: '%s'", err))
		}
		defer f.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			panic(fmt.Sprintf("could not write memory profile: '%s'", err))
		}
	}

	// TODO: add a configuration for net/http/pprof port
	go func() {
		err := http.ListenAndServe(*net_addr_port_prof, nil)
		if err != nil {
			panic(fmt.Sprintf("net/http/pprof listen and serve failed: '%s'", err))
		} else {
			log.Println("net/http/pprof is listening on 6062")
		}
	}()

	if *opt_version {
		fmt.Printf(version.String())
		os.Exit(0)
	}

	l4g.LoadConfiguration(*log_file)

	if file, err := os.OpenFile(*config_file, os.O_RDONLY, 0666); err == nil {

		_conf, err = config.Load_FEConfig(*config_file, *key_file)
		if err == nil {
			l4g.Info("Read config from '%s' ...", *config_file)
		} else {
			l4g.Warn("Failed to Read config from '%s': '%s'.", *config_file, err)
		}
		file.Close()
	}

	// Making sure to have the latest configuration.
	_conf = config.GetFEConfig()

	if *print_conf {
		fmt.Printf(config.GetFEConfig().String())
		os.Exit(0)
	}

	err = Init()

	if err != nil {
		fmt.Printf("Error during initialization: '%s'! Aborting ...", err)
		return
	}

	l4g.Info("Daemon server initialized ...")

	if _conf.LogicLayer.HealthCheck.HealthCheckEnable {
		health_check.Set_Status("Startup", "Starting services")
	}

	l_c_stop_ll = make(chan chan int, NB_OF_WORKER)

	signal.Notify(c_kill, syscall.SIGTERM)
	signal.Notify(c_int, syscall.SIGINT)
	signal.Notify(c_usr1, syscall.SIGUSR1)

	l4g.Info("Daemon server starting ...")

	procs = make([]*ll.T_ll_proc, NB_OF_WORKER)

	for i := 0; i < NB_OF_WORKER; i++ {
		ll_process := ll.New(config.GetFEConfig())
		procs = append(procs, ll_process)
		go ll_process.Execute()
	}

	err_chan := start_servers()

	if _conf.LogicLayer.HealthCheck.HealthCheckEnable {
		health_check.Set_Status("Up", "Started successfully")
	}

main_loop:
	for {
		select {
		case <-c_kill:
			if _conf.LogicLayer.HealthCheck.HealthCheckEnable {
				health_check.Set_Status("Shutdown", "SIGTERM received")
			}
			l4g.Info("SIGTERM signal received! Shutting down ...")
			break main_loop

		case <-c_int:
			if _conf.LogicLayer.HealthCheck.HealthCheckEnable {
				health_check.Set_Status("Shutdown", "SIGINT received")
			}
			l4g.Info("SIGINT (Ctrl+c) detected! Shutting down ...")
			break main_loop

		case err = <-err_chan:
			l4g.Error("Server graceful shutdown due to error '%s'", err.Error())
			stop_servers()
			break main_loop

		case <-c_usr1:
			l4g.LoadConfiguration(*log_file)
			l4g.Info("Reread L4G configuration file on reception of SIGUSR1 signal!")
		}
	}
	l4g.Info("Daemon shutdown completed ...")
}

func stop(n int, procs []*ll.T_ll_proc) {

	defer time.Sleep(time.Millisecond * 300)

	for _, proc := range procs {
		// Since we have NB_OF_WORKER to shutdown, we will only give
		// the allowed shutdown time divided by NB_OF_WORKER+1
		// We use NB_OF_WORKER+1 (instead of NB_OF_WORKER) to take in consideration
		// the overhead.
		proc.Stop(SHUTDOWN_TIMEOUT * time.Second / (NB_OF_WORKER + 1))
	}
}

func stop_servers() {

	time.Sleep(time.Second * 1)

	if _conf.FrontLayer.HttpPort != 0 {
		healthcheck.StopHttpHealthCheck()
	}

}

func start_servers() chan error {

	errs := make(chan error)

	if _conf.FrontLayer.HttpPort != 0 {
		healthcheck.StartHttpHealthCheck(_conf.FrontLayer.HttpPort+2000, &errs)
	}
	l4g.Info("Daemon Server started ...")
	return errs
}
