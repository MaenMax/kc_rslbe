package router

import (

	//"runtime/pprof"
	"git.kaiostech.com/cloud/common/utils/handlers_common"
	"net/http/pprof"
)

// net/http/pprof only registers its handlers with http.DefaultServeMux.
// Since we are not using the default mux, we need to remap the default endpoints.
func init_profiling() (routes handlers_common.T_Routes) {

	routes = append(routes, handlers_common.T_Route{
		Name:        "Global Performance - Profling Index Page",
		Method:      "GET",
		Pattern:     "/debug/pprof",
		Audiences:   nil,
		CredChecker: nil,
		Pagination:  nil,
		HandlerFunc: pprof.Index,
		Permissions: "",
		InputType:   nil,
		OutputType:  nil,
		TC:          nil,
	})

	routes = append(routes, handlers_common.T_Route{
		Name:        "Global Performance - Profiling CmdLine",
		Method:      "GET",
		Pattern:     "/debug/pprof/cmdline",
		Audiences:   nil,
		CredChecker: nil,
		Pagination:  nil,
		HandlerFunc: pprof.Cmdline,
		Permissions: "",
		InputType:   nil,
		OutputType:  nil,
		TC:          nil,
	})

	routes = append(routes, handlers_common.T_Route{
		Name:        "Global Performance - Profiling Profile",
		Method:      "GET",
		Pattern:     "/debug/pprof/profile",
		Audiences:   nil,
		CredChecker: nil,
		Pagination:  nil,
		HandlerFunc: pprof.Profile,
		Permissions: "",
		InputType:   nil,
		OutputType:  nil,
		TC:          nil,
	})

	routes = append(routes, handlers_common.T_Route{
		Name:        "Global Performance - Profiling Trace",
		Method:      "GET",
		Pattern:     "/debug/pprof/trace",
		Audiences:   nil,
		CredChecker: nil,
		Pagination:  nil,
		HandlerFunc: pprof.Trace,
		Permissions: "",
		InputType:   nil,
		OutputType:  nil,
		TC:          nil,
	})

	routes = append(routes, handlers_common.T_Route{
		Name:        "Global Performance - Profiling Symbols",
		Method:      "GET",
		Pattern:     "/debug/pprof/symbol",
		Audiences:   nil,
		CredChecker: nil,
		Pagination:  nil,
		HandlerFunc: pprof.Symbol,
		Permissions: "",
		InputType:   nil,
		OutputType:  nil,
		TC:          nil,
	})

	routes = append(routes, handlers_common.T_Route{
		Name:        "Global Performance - Profiling Goroutine",
		Method:      "GET",
		Pattern:     "/debug/pprof/goroutine",
		Audiences:   nil,
		CredChecker: nil,
		Pagination:  nil,
		HandlerFunc: pprof.Handler("goroutine").ServeHTTP,
		Permissions: "",
		InputType:   nil,
		OutputType:  nil,
		TC:          nil,
	})

	routes = append(routes, handlers_common.T_Route{
		Name:        "Global Performance - Profiling Heap",
		Method:      "GET",
		Pattern:     "/debug/pprof/heap",
		Audiences:   nil,
		CredChecker: nil,
		Pagination:  nil,
		HandlerFunc: pprof.Handler("heap").ServeHTTP,
		Permissions: "",
		InputType:   nil,
		OutputType:  nil,
		TC:          nil,
	})

	routes = append(routes, handlers_common.T_Route{
		Name:        "Global Performance - Profiling Thread Create",
		Method:      "GET",
		Pattern:     "/debug/pprof/threadcreate",
		Audiences:   nil,
		CredChecker: nil,
		Pagination:  nil,
		HandlerFunc: pprof.Handler("threadcreate").ServeHTTP,
		Permissions: "",
		InputType:   nil,
		OutputType:  nil,
		TC:          nil,
	})

	routes = append(routes, handlers_common.T_Route{
		Name:        "Global Performance - Profiling Block",
		Method:      "GET",
		Pattern:     "/debug/pprof/block",
		Audiences:   nil,
		CredChecker: nil,
		Pagination:  nil,
		HandlerFunc: pprof.Handler("block").ServeHTTP,
		Permissions: "",
		InputType:   nil,
		OutputType:  nil,
		TC:          nil,
	})

	routes = append(routes, handlers_common.T_Route{
		Name:        "Global Performance - Profiling Mutex",
		Method:      "GET",
		Pattern:     "/debug/pprof/mutex",
		Audiences:   nil,
		CredChecker: nil,
		Pagination:  nil,
		HandlerFunc: pprof.Handler("mutex").ServeHTTP,
		Permissions: "",
		InputType:   nil,
		OutputType:  nil,
		TC:          nil,
	})

	routes = append(routes, handlers_common.T_Route{
		Name:        "Global Performance - Profiling Allocations",
		Method:      "GET",
		Pattern:     "/debug/pprof/allocs",
		Audiences:   nil,
		CredChecker: nil,
		Pagination:  nil,
		HandlerFunc: pprof.Handler("allocs").ServeHTTP,
		Permissions: "",
		InputType:   nil,
		OutputType:  nil,
		TC:          nil,
	})

	return
}
