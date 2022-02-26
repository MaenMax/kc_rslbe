package router

import (
	"git.kaiostech.com/cloud/common/utils/handlers_common"
	"git.kaiostech.com/cloud/kc_rslbe/kc_rslbe_fe/handlers"
)

func init_rslbe() (routes handlers_common.T_Routes) {

	routes = append(routes, handlers_common.T_Route{
		Name:          "Generic creation",
		Method:        "POST",
		Pattern:       "/kc_rslbe/v1.0/generic",
		Audiences:     DEFAULT_EP,
		CredChecker:   handlers_common.CredChecker,
		Pagination:    nil,
		Documentation: ``,
		Cors:          true,
		HandlerFunc:   handlers.Generic,
		Permissions:   "sys#accounts:c",
		InputType:     nil,
		OutputType:    nil,
		TC:            nil,
	})

	return
}
