package trv

// SystemMode represents the system modes
type SystemMode string

// System modes
var (
	SystemModeAuto SystemMode = "auto"
	SystemModeHeat SystemMode = "heat"
	SystemModeOff  SystemMode = "off"
)

// ForceMode represents the force modes
type ForceMode string

// Force modes
var (
	ForceModeUnavailable ForceMode
	ForceModeNormal      ForceMode = "normal"
	ForceModeOpen        ForceMode = "open"
	ForceModeClose       ForceMode = "close"
)
