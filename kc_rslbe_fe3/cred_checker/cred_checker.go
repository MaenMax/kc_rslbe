package cred_checker

import (
	"fmt"
	"net/http"

	"git.kaiostech.com/cloud/common/limits"
	"git.kaiostech.com/cloud/common/utils/context"
	common "git.kaiostech.com/cloud/common/utils/handlers_common"
	"git.kaiostech.com/cloud/common/utils/handlers_common/cred"

	"git.kaiostech.com/cloud/common/config"
	l4g "git.kaiostech.com/cloud/thirdparty/code.google.com/p/log4go"
)

type T_RslCredChecker func(next http.Handler, accepted_aud_list []string) http.Handler

func RslCredChecker(api_version string) T_RslCredChecker {
	return T_RslCredChecker(func(next http.Handler, accepted_aud_list []string) http.Handler {
		var aud_list []string
		aud_list = accepted_aud_list

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if config.GetFEConfig().Common.Debug {
				// Output the received header from client
				for key, value := range r.Header {
					l4g.Finest("req #%v: %s: %s", context.GetReqId(r), key, value)
				}
			}

			///////////////////////////////////////////////////////////////////////
			// Read and check parameters.
			body, cerr := common.ReadBody(w, r, limits.MAX_POST_SIZE)

			if cerr != nil {
				// Error already managed in read_body.
				l4g.Error("req #%v: CredChecker:: Aborting due to error while reading body ...", context.GetReqId(r))

				// Bug 71463 - [Monitoring] API used for monitoring generate status code 0 in log file
				// To mark this request as a 'controlled' error, we mark the status code as 1000.
				//
				// NOTE: Despite responding with an error we will still set a return value into the
				//       context because we observe that when the request timeout, nothing is reported
				//       back to the request (w http.ResponseWriter might become invalid).

				context.SetStatusCode(r, 1000)
				common.RespondError(r, w, cerr)

				return
			}

			if config.GetFEConfig().Common.Debug {
				l4g.Fine("req #%v: CredChecker:: body size is %v.", context.GetReqId(r), len(body))
			}

			auth_head, cerr := cred.GetAuthHeader(r)

			if cerr != nil {

				common.RespondError(r, w, cerr)
				return
			}

			if config.GetFEConfig().Common.Debug {
				l4g.Fine("req #%v: %s", context.GetReqId(r), fmt.Sprintf("CredChecker:auth_head %v", auth_head))
			}

			if cerr := common.CheckCredential(r, auth_head, body, aud_list); cerr != nil {

				common.RespondError(r, w, cerr)
				return
			}

			next.ServeHTTP(w, r)

		})
	})
}
