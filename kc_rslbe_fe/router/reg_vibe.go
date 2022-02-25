package router

import (
	"git.kaiostech.com/cloud/common/model/vibe"
	"git.kaiostech.com/cloud/common/utils/handlers_common"
	"git.kaiostech.com/cloud/kc_/kc_core_fe3/handlers/account"
)

func init_rslbe() {

	routes = append(routes, handlers_common.T_Route{
		Name:        "Account Creation by Admin",
		Method:      "POST",
		Pattern:     "/kc_vibe_be/v1.0/",
		Audiences:   DEFAULT_EP,
		CredChecker: handlers_common.CredChecker,
		Pagination:  nil,
        Documentation: &string(``),
        Cors:          true,
		HandlerFunc: account.Create_Account,
		Permissions: "sys#accounts:c",
		InputType:   vibe.T_Account{}.Post_Request(),
		OutputType:  vive.T_Account{}.Post_Response(),
		TC:          nil,
	})


}
