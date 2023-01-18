package ll

import (
	"time"

	"git.kaiostech.com/cloud/common/limits"

	"git.kaiostech.com/cloud/common/cerrors"
	"git.kaiostech.com/cloud/common/model/core"
	"git.kaiostech.com/cloud/common/mq"
	common "git.kaiostech.com/cloud/common/utils/actions_common"
	"git.kaiostech.com/cloud/kc_rslbe/kc_rslbe_ll3/actions"

	"fmt"

	"git.kaiostech.com/cloud/common/config"
	l4g "git.kaiostech.com/cloud/thirdparty/code.google.com/p/log4go"
	// For profiling
	// "net/http"
	//	_ "net/http/pprof"
)

type T_ll_proc struct {
	srv             *mq.MQServer
	status          *mq.MQServerStatus
	subject_from_fe string
	automata        *actions.Automata
	stop            chan int
	req_exec_quota  chan int
	max_req         int
}

func New(conf *config.FEConfig) *T_ll_proc {
	var p T_ll_proc = T_ll_proc{}

	p.max_req = conf.LogicLayer.MaxActiveRequest
	p.req_exec_quota = make(chan int, p.max_req)

	if p.max_req == 0 {
		panic("Maximum number of active request is 0!")
	}

	// Filling with quota
	for i := 0; i < p.max_req; i++ {
		p.req_exec_quota <- 1
	}

	// Auth requests must be configured to use this queue. (see FrontLayer:Heavy_Load_Channel)
	if conf.LogicLayer.Process_Only_Heavy_Requests {
		l4g.Info("Processing heavy requests only")
		p.subject_from_fe = conf.QueueService.FE2LL_Heavy_Subject
	} else {
		l4g.Info("Processing regular requests")
		p.subject_from_fe = conf.QueueService.FE2LL_Subject
	}

	if config.GetFEConfig().Common.Percentage_Mem_Allowed <= 0 || config.GetFEConfig().Common.Percentage_Mem_Allowed >= 100 {
		l4g.Info("Memory throttling has been disabled.")
	} else {
		limits.Log_Actual_Max_Mem(uint64(config.GetFEConfig().Common.Percentage_Mem_Allowed))
	}

	nats_url := conf.GetNatsServer()

	// Server for communication from FE
	p.srv = mq.NewMQServer(nats_url, p.subject_from_fe)

	p.srv.Start()

	p.stop = make(chan int)

	p.automata = actions.NewAutomata()
	return &p
}

func (p *T_ll_proc) Stop(shutdown_timeout time.Duration) {
	l4g.Info("Worker %p: Shutdown requested ...", p)
	p.srv.Stop()
	select {
	case <-p.stop:
	case <-time.After(shutdown_timeout):
		l4g.Warn("Worker %p: Shutdown forced due to timeout ...", p)
	}
}

func (p *T_ll_proc) Execute() {
	var err error
	var req *mq.MQRequest
	var stop bool

	l4g.Info("Worker %p: Starting ...", p)

main_loop:
	for {

		if req, stop, err = p.srv.Listen(); err != nil {
			l4g.Error("Worker %p: Failed to listen: %s.\n", p, err.Error())
			continue
		}

		if stop {
			l4g.Info("Worker %p: Shutdown requested ...", p)
			break main_loop
		}

		// TO DO: replace with the request id
		request_id := req.ReqId

		select {
		case <-p.req_exec_quota:
			l4g.Debug("Worker %p: Processing Request #%v with response subject '%s' ...", p, request_id, req.Id)
			go p.process_request(request_id, req)
			continue

		case <-time.After(time.Millisecond * 1):
			l4g.Warn("Worker %p: Maximum Number of request execution quota (%v) exceeded.", p, p.max_req)
			cerr := cerrors.New(cerrors.ERROR_TOO_MANY_REQUESTS, fmt.Sprintf("Action for scope (%v) and request type (%v) failed to execute due to lack of server resource!", req.Scope, req.Type), "LL", true)
			rsp := mq.NewMQResponse(req.Id, req.Type, req.Scope, []byte(cerr.ToJson()), true, req.ReqId)
			p.Send(request_id, rsp)
		}

	} // for { level 1

	// Second, waiting for completion of all ongoing requests.
	// NOTE: since each single request has an internal timeout, we are sure they are going to complete
	//       at some point of time so we don't need watchdog here.
	for {
		c := len(p.req_exec_quota)

		if c == p.max_req {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}

	l4g.Info("Worker %p: All active requests are completed ... ", p)

	// Finally, fully disconnecting from NATS and signalling back to main thread the termination of the processing.
	//p.srv.Disconnect()
	p.stop <- 1

	l4g.Info("Worker #%v: Shutdown ...", p)
}

func (p *T_ll_proc) process_request(request_id string, req *mq.MQRequest) {
	var rsp *mq.MQResponse
	var rsp_id string

	defer func() { p.req_exec_quota <- 1 }()

	dl_client, err := common.MQCliMan.Client()

	if err != nil {
		l4g.Error("Worker %p, request #%v: MQClientManager returned error '%s'.", p, request_id, err.Error())
		cerr := cerrors.New(cerrors.ERROR_INTERNAL_SERVER_ERROR, fmt.Sprintf("Request #%v: Action for scope (%v) and request type (%v) failed due to error '%s'!", request_id, req.Scope, req.Type, err.Error()), "LL", true)
		rsp := mq.NewMQResponse(req.Id, req.Type, req.Scope, []byte(cerr.ToJson()), true, req.ReqId)

		if err := p.srv.Send(rsp); err != nil {
			// Here we failed to send back a message maybe because of the client.
			l4g.Error("Worker %p, request #%v: Failed to send back error response: %s.", p, request_id, err.Error())
		}
		return
	}

	if dl_client == nil {
		l4g.Error("Worker %p, request #%v: Nil NATS Client returned by client manager.", p, request_id)
		cerr := cerrors.New(cerrors.ERROR_INTERNAL_SERVER_ERROR, fmt.Sprintf("Request #%v: Action for scope (%v) and request type (%v) failed due to unavailable message queue connection!", request_id, req.Scope, req.Type), "LL", true)
		rsp := mq.NewMQResponse(req.Id, req.Type, req.Scope, []byte(cerr.ToJson()), true, req.ReqId)

		if err := p.srv.Send(rsp); err != nil {
			// Here we failed to send back a message maybe because of the client.
			l4g.Error("Worker %p, request #%v: Failed to send back error response: %s.", p, request_id, err.Error())
		}
		return
	}

	// Making sure the response id matches the request id else routing is going to fail.
	// We are going to save the req.id into a temporary variable in the case this ID
	// is changed within the Automata processing.
	rsp_id = req.Id

	if 0 < config.GetFEConfig().Common.Percentage_Mem_Allowed && config.GetFEConfig().Common.Percentage_Mem_Allowed < 100 {
		cerr := limits.Check_Mem(req.Id, uint64(config.GetFEConfig().Common.Percentage_Mem_Allowed), "LL")
		if cerr != nil {
			l4g.Error("Worker %p, request #%v:  (%v,%v): Server is busy error: %s.", p, request_id, req.Scope, req.Type, cerr.Error())
			rsp = mq.NewMQResponse(req.Id, req.Type, req.Scope, []byte(cerr.ToJson()), true, req.ReqId)
		}
	}

	if config.GetFEConfig().Common.Debug {
		l4g.Debug("Worker %p, request #%v: Processing request with response subject '%s' ...", p, request_id, rsp_id)
	}

	// we don't check the error from the jwt conversion yet. We will
	// let the action do it. This is because some actions require JWT
	// while some others don't.
	jwt, _ := core.JSON2JWT([]byte(req.ReqInfo.JWT))

	action := p.automata.Get(req.Scope, req.Type)

	if action == nil {
		l4g.Error("Worker %p, request #%v:  (%v,%v) Failed! Not linked with any action.", p, request_id, req.Scope, req.Type)
		cerr := cerrors.New(cerrors.ERROR_INVALID_REQUEST, fmt.Sprintf("Request #%v: Invalid scope (%v), request type (%v) pair as no action linked.", request_id, req.Scope, req.Type), "LL", true)
		rsp = mq.NewMQResponse(req.Id, req.Type, req.Scope, []byte(cerr.ToJson()), true, req.ReqId)

	} else {
		start := time.Now()
		//rsp = action.Execute(dl_client, request_id, jwt, req)
		rsp = action.Execute(dl_client, 0, jwt, req)
		duration := time.Since(start)

		// Once done, recycling the dl_client if no I/O error occured!
		if !dl_client.IsIOError() {
			common.MQCliMan.Recycle(dl_client)
		}

		if rsp == nil {
			l4g.Warn("Request #%v: ==METRICS (%v,%v) Null. Processing time %d.", request_id, req.Scope, req.Type, duration.Nanoseconds())
			cerr := cerrors.New(cerrors.ERROR_INTERNAL_SERVER_ERROR, fmt.Sprintf("Action for scope (%v) and request type (%v) pair leads to nil response!", req.Scope, req.Type), "LL", true)
			rsp = mq.NewMQResponse(req.Id, req.Type, req.Scope, []byte(cerr.ToJson()), true, req.ReqId)
		} else {
			// Making sure to assign same ID as the origin request in order to route properly back to the sender.
			l4g.Info("Worker %p, request #%v: ==METRICS (%v,%v) Success. Processing time %d.", p, request_id, req.Scope, req.Type, duration.Nanoseconds())
		}
	}

	rsp.Id = rsp_id
	p.Send(request_id, rsp)

}

func (p *T_ll_proc) Send(request_id string, rsp *mq.MQResponse) {
	if err := p.srv.Send(rsp); err != nil {
		// Here we failed to send back a message maybe because of the client.
		l4g.Error("Worker %p, request #%v: Failed to send back error response: %s.", p, request_id, err.Error())
	} else {
		if config.GetFEConfig().Common.Debug {
			l4g.Debug("Worker %p, request #%v: Processing request completed.", p, request_id)
		}
	}
}
