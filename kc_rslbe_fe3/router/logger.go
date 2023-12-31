/*
 * Remote Sim Lock API
 *
 * **Implementation of the Remote Sim Lock APIs**  **Vibe project documentation (including RSL) can be found using the following links:**     **1- DRAFT-RSL sequence flows pictures source ZY-20220225.docx**    * https://kaios.sharepoint.com/:w:/r/sites/vibe/Shared%20Documents/4.%20PRODUCT/06%20-%20RSL%20REQ%20and%20Designing/DRAFT-RSL%20sequence%20flows%20pictures%20source%20ZY-20220225.docx?d=w8cba978f89bc4cac8460d87c0e1053ba&csf=1&web=1&e=iQp9Eo *    This document includes squence diagrams which describe the use cases: “device initiate” , “daily ping”, and “user change SIM to another one”.    **2- DRAFT-Vibe Remote SIM Lock Operation Structure-20220307.docx**    * https://kaios.sharepoint.com/:w:/r/sites/vibe/Shared%20Documents/4.%20PRODUCT/06%20-%20RSL%20REQ%20and%20Designing/DRAFT-Vibe%20Remote%20SIM%20Lock%20Operation%20Structure-20220307.docx?d=w51d49ceb6ddc4329aaa18ce02ce22318&csf=1&web=1&e=orUFfm    This is a short document which includes an introduction to the RSL operation structure, and RSL portal.    **3- RSL flow charts for communicate.v0.2.docx**    * https://kaios.sharepoint.com/:w:/r/sites/vibe/Shared%20Documents/4.%20PRODUCT/06%20-%20RSL%20REQ%20and%20Designing/RSL%20flow%20charts%20for%20communicate.v0.2.docx?d=w8215af326530486ab8b84b14e3208967&csf=1&web=1&e=H2DhdF    **4- Vibe RSL Technical Keypoints.20220915.pptx**    * https://kaios.sharepoint.com/:p:/r/sites/vibe/Shared%20Documents/4.%20PRODUCT/06%20-%20RSL%20REQ%20and%20Designing/Vibe%20RSL%20Technical%20Keypoints%20.20220915.pptx?d=w58603c5ed4814c9086a799450262b97e&csf=1&web=1&e=czrBNc    **5- Vibe Product requirement_211118.pptx**    * https://kaios-my.sharepoint.com/:p:/r/personal/raffi_semerciyan_kaiostech_com/Documents/Documents/20211124-Vibe/Vibe%20Product%20Requirement_211118.pptx?d=w130eb4dd2a74481d93cb1fec26c0d068&csf=1&web=1&e=ivj2jW    **6- Vibe-Requirements-Analysis.docx**    * https://kaios-my.sharepoint.com/:w:/g/personal/raffi_semerciyan_kaiostech_com/EXlVgqdcF7pAii8-VIHeRXQB_VQEBIbjO_2aU3Ur_Rqk1w?e=R7eaOb     This document shows the detailed design of RSL function, after making a requirement summary.      **7- Vibe-Specification.docx**   * https://kaios-my.sharepoint.com/:w:/g/personal/raffi_semerciyan_kaiostech_com/EXlVgqdcF7pAii8-VIHeRXQB_VQEBIbjO_2aU3Ur_Rqk1w?e=R7eaOb    This document shows the detailed use case diagrams which show the interaction between Vive users and the system.
 *
 * API version: 1.0.0
 * Contact: maen.hammour@kaiostech.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package router

import (
	"log"
	"net/http"
	"time"
)

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s %s %s %s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
