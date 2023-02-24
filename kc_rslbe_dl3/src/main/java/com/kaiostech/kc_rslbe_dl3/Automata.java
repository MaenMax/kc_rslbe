package com.kaiostech.kc_rslbe_dl3;

// All the Actions should be imported since they need to be registered
// into the Automata!
import com.kaiostech.kc_rslbe_dl3.actions.IAction;
import com.kaiostech.kc_rslbe_dl3.actions.devices.TransferOwnershipActions;
import com.kaiostech.kc_rslbe_dl3.actions.devices.TransferStateActions;
import com.kaiostech.kc_rslbe_dl3.actions.devices.update_2i.Update2I;
import com.kaiostech.kc_rslbe_dl3.actions.partners.PartnerActions;
import com.kaiostech.mq.MQRequestScope;
import com.kaiostech.mq.MQRequestType;

/**
 * The Automata class allows the easy mapping of a request Type and Scope to an action. It is just a
 * two dimensional array (one dimension is the RequestType and the other dimension is the
 * RequestScope) that contain an action (represented by the IAction interface).
 *
 * <p>This is an important improvement compared to a double interlaced switch statements to check
 * for every Type and then for every Scope which would have used pages and pages of source code.
 *
 * <p>When adding new Scope or new Type, you will need to register the corresponding actions in this
 * file by initializing the mapping at the very end of the 'init' function.
 *
 * <p>NOTE: All the actions to register should be implemented into the
 * com.kaiostech.kc_rslbe_dl3.actions package.
 */
public class Automata {

  // First dimension: MQRequestScope,
  // Second dimension: MQRequestType
  // Cell Value: null or an IAction object
  private IAction[][] _action_mapping;

  private int _max_scope;
  private int _max_type;

  public Automata() {
    init();
  }

  public IAction get(MQRequestScope scope, MQRequestType type) {

    int type_val = type.getValue();
    int scope_val = scope.getValue();
    IAction[] scope_actions;

    if ((scope_val <= 0) || (scope_val >= _max_scope + 1)) {
      return null;
    }

    if ((type_val <= 0) || (type_val >= _max_type + 1)) {
      return null;
    }

    // Here we are sure the scope and type values are within the
    // admissible range of the automata.
    scope_actions = _action_mapping[scope_val];

    if (scope_actions != null) {
      return scope_actions[type_val];
    }
    return null;
  }

  private void init() {
    int i;

    _max_scope = MQRequestScope.MQRS_LAST.getValue();
    _max_type = MQRequestType.MQRT_LAST.getValue();

    _action_mapping = new IAction[_max_scope + 1][];

    // First making sure that by default all the actions (i.e. cells) are null.
    // Not all combination of Scope/Type should lead to an action!
    //
    // Example: MQRS_UNKNOWN and MQRS_NONE both have a Null registered which
    // is the expected behavior.
    for (i = 0; i < _max_scope + 1; i++) {
      _action_mapping[i] = null;
    }
    _action_mapping[MQRequestScope.MQRS_RSL_BE_REGISTER_3I.getValue()] =
        PartnerActions.getInstance().getAutomataLine();
    _action_mapping[MQRequestScope.MQRS_RSL_BE_UPDATE_2I.getValue()] =
        Update2I.getInstance().getAutomataLine();
    _action_mapping[MQRequestScope.MQRS_RSL_BE_TRANSFER_STATE.getValue()] =
        TransferStateActions.getInstance().getAutomataLine();
    _action_mapping[MQRequestScope.MQRS_RSL_BE_TRANSFER_OWNERSHIP.getValue()] =
        TransferOwnershipActions.getInstance().getAutomataLine();
  }
}
