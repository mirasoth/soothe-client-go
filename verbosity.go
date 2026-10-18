package soothe

// VerbosityLevel represents the user-configurable verbosity setting.
type VerbosityLevel string

const (
	VerbosityQuiet  VerbosityLevel = "quiet"
	VerbosityNormal VerbosityLevel = "normal"
	VerbosityDebug  VerbosityLevel = "debug"
)

// VerbosityTier is the visibility class for content, mirroring
// soothe_nano.core.verbosity's binary model (NORMAL=1 / INTERNAL=99).
type VerbosityTier int

const (
	TierNormal   VerbosityTier = 1  // client-visible (flows over the wire)
	TierInternal VerbosityTier = 99 // daemon-side only, never leaves the daemon
)

// WireTierForVerbosity maps the user verbosity knob to the daemon's
// loop_events wire_tier subscription param: quiet -> "compact", else "full".
func WireTierForVerbosity(verbosity VerbosityLevel) string {
	if verbosity == VerbosityQuiet {
		return "compact"
	}
	return "full"
}

// ShouldShow reports whether content at the given tier is visible to the
// client. Wire-level filtering is decided daemon-side by wire_tier; locally
// only INTERNAL-tier content is hidden.
func ShouldShow(tier VerbosityTier) bool {
	return tier != TierInternal
}

// IsValidVerbosityLevel checks whether a string is a valid verbosity level.
func IsValidVerbosityLevel(s string) bool {
	switch VerbosityLevel(s) {
	case VerbosityQuiet, VerbosityNormal, VerbosityDebug:
		return true
	default:
		return false
	}
}
