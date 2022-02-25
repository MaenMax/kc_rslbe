package router

import (
	"git.kaiostech.com/cloud/common/utils/handlers_common"
)

func init_centreon() {

	routes = append(routes, handlers_common.T_Route{
		Name:        "Global Performance Raw Counters for Centreon",
		Method:      "GET",
		Pattern:     "/centreon/v1.0/perfs/services/counters/raw",
		Audiences:   nil,
		CredChecker: nil,
		Pagination:  nil,
		HandlerFunc: handlers_common.Nagios_Raw_Perfs_Counters,
		Permissions: "",
		InputType:   nil,
		OutputType:  nil,
		TC:          nil,
	})

	routes = append(routes, handlers_common.T_Route{
		Name:        "Global Performance Counters Diff since last call for Centreon",
		Method:      "GET",
		Pattern:     "/centreon/v1.0/perfs/services/counters/diff",
		Audiences:   nil,
		CredChecker: nil,
		Pagination:  nil,
		HandlerFunc: handlers_common.Nagios_Diff_Perfs_Counters,
		Permissions: "",
		InputType:   nil,
		OutputType:  nil,
		TC:          nil,
	})

	routes = append(routes, handlers_common.T_Route{
		Name:        "Global Performance Counters Diff in percentage since last call for Centreon",
		Method:      "GET",
		Pattern:     "/centreon/v1.0/perfs/services/counters/rel",
		Audiences:   nil,
		CredChecker: nil,
		Pagination:  nil,
		HandlerFunc: handlers_common.Nagios_Diff_Percent_Perfs_Counters,
		Permissions: "",
		InputType:   nil,
		OutputType:  nil,
		TC:          nil,
	})

	routes = append(routes, handlers_common.T_Route{
		Name:        "Global Performance Counters Diff per API since last call for Centreon",
		Method:      "GET",
		Pattern:     "/v3.0/perfs/api/counters/diff",
		Audiences:   nil,
		CredChecker: nil,
		Pagination:  nil,
		HandlerFunc: handlers_common.Diff_Perfs_Counters_Per_API,
		Permissions: "",
		InputType:   nil,
		OutputType:  nil,
		TC:          nil,
	})

	routes = append(routes, handlers_common.T_Route{
		Name:        "Global Performance Counters Diff per Raw API since last call for Centreon",
		Method:      "GET",
		Pattern:     "/v3.0/perfs/raw_api/counters/diff",
		Audiences:   nil,
		CredChecker: nil,
		Pagination:  nil,
		HandlerFunc: handlers_common.Diff_Perfs_Counters_Per_Raw_API,
		Permissions: "",
		InputType:   nil,
		OutputType:  nil,
		TC:          nil,
	})

}
