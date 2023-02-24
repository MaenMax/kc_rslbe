package actions

import (
	"git.kaiostech.com/cloud/common/mq"
	common "git.kaiostech.com/cloud/common/utils/actions_common"
	"git.kaiostech.com/cloud/kc_rslbe/kc_rslbe_ll3/actions/devices"
	"git.kaiostech.com/cloud/kc_rslbe/kc_rslbe_ll3/actions/rsl_partners"
)

type Automata struct {
	// First dimension: MQRequestScope,
	// Second dimension: MQRequestType
	// Cell Value: null or an IAutomatizer
	action_mapping [][]common.IAutomatizer

	_max_scope int
	_max_type  int
}

func NewAutomata() *Automata {
	var automata *Automata = &Automata{}

	automata.Init()
	return automata
}

func (a *Automata) Get(scope mq.MQScope, req_type mq.MQRequestType) common.IAutomatizer {
	var scope_actions []common.IAutomatizer
	var type_val int
	var scope_val int

	scope_val = int(scope)
	type_val = int(req_type)

	if scope_val <= 0 || scope_val >= (a._max_scope+1) {
		return nil
	}

	if type_val <= 0 || type_val >= (a._max_type+1) {
		return nil
	}

	scope_actions = a.action_mapping[scope_val]

	if scope_actions == nil {
		return nil
	}

	return scope_actions[type_val]

}

func (a *Automata) Init() {

	a._max_scope = int(mq.MQSCOPE_LAST)
	a._max_type = int(mq.MQRT_LAST)

	a.action_mapping = make([][]common.IAutomatizer, a._max_scope+1)

	for i := 0; i < a._max_scope+1; i++ {
		a.action_mapping[i] = nil
	}
	a.action_mapping[int(mq.MQSCOPE_RSL_BE_REGISTER_3I)] = rsl_partners.GetAutomataLine()
	a.action_mapping[int(mq.MQSCOPE_RSL_BE_UPDATE_2I)] = devices.GetUpdate2IAutomataLine()
	a.action_mapping[int(mq.MQSCOPE_RSL_BE_TRANSFER_STATE)] = devices.GetTransferStateAutomataLine()
	a.action_mapping[int(mq.MQSCOPE_RSL_BE_TRANSFER_OWNERSHIP)] = devices.GetTransferOwnershipAutomataLine()
}
