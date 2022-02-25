package router

import (
	"net/http"
	"strings"
	"sync"

	"git.kaiostech.com/cloud/common/utils/handlers_common"
	"git.kaiostech.com/cloud/thirdparty/github.com/gorilla/mux"

	"git.kaiostech.com/cloud/common/config"
	l4g "git.kaiostech.com/cloud/thirdparty/code.google.com/p/log4go"
)

var mutex sync.Mutex

func Routes() handlers_common.T_Routes {
	return routes
}

func NewRouter(max_req_config int, max_mem_percentage int) *mux.Router {
    var (
        handler http.Handler
        aud_list []string
    )

	mutex.Lock()
	handlers_common.Init_Throttle(max_req_config)
	handlers_common.Init_Mem_Throttle(max_mem_percentage)

	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		handler = route.HandlerFunc

		/*
			Wrapping with throttle controls.
		*/
		handler = handlers_common.HTTP_Throttle(handler)

		if 0 < max_mem_percentage && max_mem_percentage < 100 {
			handler = handlers_common.Memory_Throttle(handler, "FE")
		}

		// Audiences are list of End Point that are expected and accepted
		// to be found into the audience field of the JWT token.
		aud_list = route.Audiences

		if route.Pagination != nil {
			handler = route.Pagination(handler)
		}

		/*
		   Wrapping Permission check handler if permission not empty.
		*/
		if len(route.Permissions) > 0 {
            if config.GetFEConfig().Common.Debug {
                l4g.Debug("NewRouter: installing route " + route.Method + " " + route.Pattern+ ": Permission handler wrapped ")
            }
			handler = handlers_common.PermChecker(handler, route.Permissions)
		}

		/*
		   Wrapping Credential check only for the routes that require it.
		*/
		if route.CredChecker != nil {
            if config.GetFEConfig().Common.Debug {
                l4g.Debug("NewRouter: installing route " + route.Method + " " + route.Pattern+ ": CredChecker '%+v' handler wrapped ",route.CredChecker)
            }
			handler = route.CredChecker(handler, aud_list)
		}

		/*
		   Wrapping logging for ALL routes to see their execution trace in log file.
		*/
		handler = handlers_common.HTTPLogger(handler, route.Name)

		// Stats API calls should not be counted into the statistics
		if !strings.HasPrefix(route.Name, "Global Performance") {
			/*
			   Accumulating performance statistics/counters for Centreon monitoring
			*/
            if config.GetFEConfig().Common.Debug {
                l4g.Debug("NewRouter: installing route " + route.Method + " " + route.Pattern+ ": Global Perf handler wrapped ")
            }
			handler = handlers_common.HTTPStats(handler, route.Method, route.Pattern)
		}

        /*
            Cors handler chaining should be done AFTER credential checks because 
            we want to generate the CORS headers even if the authentication fails. 
            We want CORS sensitive client being able to receive 401 or 403 responses
            with proper CORS headers.
        */
        if route.Cors {
            if config.GetFEConfig().Common.Debug {
                l4g.Debug("NewRouter: installing route " + route.Method + " " + route.Pattern+ ": Cors handler wrapped ")
            }
            handler = handlers_common.Allow_Cors_Json_Response_Handler(handler)
        }

		/*
		   Single Request Trail feature request to assign a unique ID to each incoming request.
		   This will be done in this handler below.
		*/
		handler = handlers_common.AssignReqId(handler)

		/* Wrapping context cleaning to avoid memory leak */
		handler = handlers_common.ClearContext(handler)

		/*
		   Wrapping handler with execution timeout to limit resource consumption of micro service.
		*/
        /*
           /!\/!\/!\/!\/!\/!\/!\/!\/!\/!\/!\/!\/!\/!\/!\/!\/!\/!\/!\/!\/!\/!\/!\/!\/!\/!\
           Bug 129368 - [PROD] The FE process crashes randomly after some running time
           
           We found that we cannot use the timeout handler of the standard library if our 
           handler implementation hasn't been implementated to collaborate with it in order 
           to detect the abort/timeout signal. So such mechanism cannot be used without heavy
           modifications of all of our handlers. 

           NOTE: We mistakenly believed that a preemptive like implementation would suffice
                 but we were wrong.

		if route.ExecTimeout!=nil { 
			l4g.Debug("HandlerWithTimeout: %s %s: using timeout of %v from config.",route.Method,route.Pattern,*route.ExecTimeout)
			handler = handlers_common.HandlerWithTimeout(handler, *route.ExecTimeout)
		} else {
			l4g.Debug("HandlerWithTimeout: %s %s: using timeout of %v from config.",route.Method,route.Pattern,config.GetFEConfig().FrontLayer.ExecTimeout.GetDuration())
			handler = handlers_common.HandlerWithTimeout(handler, config.GetFEConfig().FrontLayer.ExecTimeout.GetDuration())
		}
           /!\/!\/!\/!\/!\/!\/!\/!\/!\/!\/!\/!\/!\/!\/!\/!\/!\/!\/!\/!\/!\/!\/!\/!\/!\/!\
        */

		if config.GetFEConfig().Common.Debug {
			l4g.Debug("NewRouter: installing route " + route.Method + " " + route.Pattern)
		}

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	mutex.Unlock()

	return router
}
