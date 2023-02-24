package com.kaiostech.kc_rslbe_dl3.actions;

import com.kaiostech.model.core.JWT;
import com.kaiostech.mq.MQRequest;
import com.kaiostech.mq.MQResponse;

/**
 * This class represents an abstraction of all possible actions we might receive from any other
 * layers.
 *
 * <p>So the arguments to the only function of this class is actually the union of all possible
 * arguments of all the possible functions.
 *
 * <p>This means not all arguments will be used for a given action. Some might be null or empty.
 *
 * <p>This abstraction has a very important impact on the architecture since it allows easily
 * splitting the source code by domain scope.
 *
 * <p>See the implementation of the com.kaiostech.kc_rslbe_dl3 .actions.TrackActions to see how both
 * the Database and Actions scopes are reduced to only one single one: Track.
 *
 * <p>This is the way teams will be able to work independently on their domain without affecting
 * others.
 *
 * <p>NOTE: A domain = a use case. So a use case may involves several tables.
 */
public interface IAction {
  public MQResponse execute(int id, JWT jwt, MQRequest request);
}
