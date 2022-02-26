package main

import (
	"flag"
	"fmt"
	"git.kaiostech.com/cloud/common/config"
	"git.kaiostech.com/cloud/common/utils/handlers_common"
	"git.kaiostech.com/cloud/common/utils/raml/v2"
	"git.kaiostech.com/cloud/kc_rslbe/utils/version"
	"os"
	"time"

	l4g "git.kaiostech.com/cloud/thirdparty/code.google.com/p/log4go"

	router "git.kaiostech.com/cloud/kc_rslbe/kc_rslbe_fe/router"
)

type T_MicroService struct {
	name   string                   // Microservice Short Name  -- Maybe used for future documentation
	desc   string                   // Microservice Description -- Maybe used for future documentation
	routes handlers_common.T_Routes // Microservice Routes
}

var (
	out_file      *string = flag.String("out", "./kaicloud.raml", "Output file name containing the generated RAML doc.")
	log_file      *string = flag.String("log", "conf/gen_doc_log.xml", "Log L4G config file to use.")
	skeleton_only *bool   = flag.Bool("simple", false, "Skeleton only flag generating RAML file without details.")
)

func main() {
	var (
		exec_name string
		r         *raml.T_RAML
	)

	start := time.Now()

	// Allow l4g to flush its content into the log file before terminating so that we can have
	// all the logs even the last ones.
	defer time.Sleep(time.Second * 2)

	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Fprintf(os.Stderr, "[ERROR] gen_doc expects the Micro service name as argument.\n")
		os.Exit(1)
	}

	// Loading configuration file for the log output.
	l4g.LoadConfiguration(*log_file)

	// Modyfying default configuration to include Payment API
	conf := config.GetFEConfig()
	conf.Common.PaymentsEnabled = true

	router.Init()

	microservices := []T_MicroService{
		{name: exec_name, desc: fmt.Sprintf("KaiCloud (%s)", exec_name), routes: router.AppRoutes()},
	}

	protocols := []string{"https"}

	r = raml.New("https://api.kaiostech.com", version.Get_Version(), protocols, "KaiCloud Remote SimLock")

	route_count := 0
	for _, ms := range microservices {
		if len(ms.routes) == 0 {
			l4g.Error("No routes registered for: %s... . Skipping", ms.name)
			continue
		}

		for _, route := range ms.routes {
			l4g.Info("Processing %s %s", route.Method, route.Pattern)
			r.AddRoute(route)
			route_count++
		}
	}

	if err := r.GenDoc(*out_file, *skeleton_only); err != nil {
		l4g.Error("Failed to generate '%s' RAML file: '%s'.", *out_file, err)
		os.Exit(1)
	}

	spent := time.Since(start)
	fmt.Printf("RAML File: '%s' - Generated %v routes, %v end points and %v operations in %v.\n", *out_file, route_count, r.EP_Count(), r.OP_Count(), spent)
	l4g.Info("RAML File: '%s' - Generated %v routes, %v end points and %v operations in %v.", *out_file, route_count, r.EP_Count(), r.OP_Count(), spent)
}
