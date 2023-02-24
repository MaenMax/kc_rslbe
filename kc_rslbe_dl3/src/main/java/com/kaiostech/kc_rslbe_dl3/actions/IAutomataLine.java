package com.kaiostech.kc_rslbe_dl3.actions;

/**
 * This class represents an abstraction of a line inside an Automata.
 *
 * <p>This interface is just used to initialize the Automata.
 *
 * <p>It is used in order to make the initialization process of the Automata more formal and easier.
 */
public interface IAutomataLine {
  public IAction[] getAutomataLine();
}
