package router

import (
	"git.kaiostech.com/cloud/common/config"
	"git.kaiostech.com/cloud/common/utils/handlers_common"
)

var (
	DEFAULT_EP  []string
	LBS_EP      []string
	FXPUSH_EP   []string
	STORAGE_EP  []string
	_all_routes handlers_common.T_Routes
	_app_routes handlers_common.T_Routes
)

func Init() {

	DEFAULT_EP = config.GetFEConfig().FrontLayer.DefaultEP
	LBS_EP = config.GetFEConfig().FrontLayer.LBSEP
	FXPUSH_EP = config.GetFEConfig().FrontLayer.PushEP
	STORAGE_EP = config.GetFEConfig().FrontLayer.StorageEP

	handlers_common.Init_Cors_Domains()

	_app_routes = init_rslbe()

	_all_routes = init_rslbe()
	_all_routes = append(_all_routes, init_centreon()...)
	_all_routes = append(_all_routes, init_profiling()...)
}

func Routes() handlers_common.T_Routes {
	return _all_routes
}

func AppRoutes() handlers_common.T_Routes {
	return _app_routes
}
