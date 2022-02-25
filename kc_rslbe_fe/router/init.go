package router

import (
	"git.kaiostech.com/cloud/common/config"
	"git.kaiostech.com/cloud/common/utils/handlers_common"
)

var routes handlers_common.T_Routes

var (
	DEFAULT_EP []string
	LBS_EP     []string
	FXPUSH_EP  []string
	STORAGE_EP []string
)

func Init() {

	DEFAULT_EP = config.GetFEConfig().FrontLayer.DefaultEP
	LBS_EP = config.GetFEConfig().FrontLayer.LBSEP
	FXPUSH_EP = config.GetFEConfig().FrontLayer.PushEP
	STORAGE_EP = config.GetFEConfig().FrontLayer.StorageEP

	handlers_common.Init_Cors_Domains()

    init_rslbe()

	init_centreon()
	init_profiling()
}
