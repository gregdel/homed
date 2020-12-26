package main

// Boiler represents the central boiler for the house
type Boiler struct {
	State bool
}

// The boiler can be turned on or off but it should not change to much
// When the app starts, the boiler as no state, it should configure the state to the wanted state
