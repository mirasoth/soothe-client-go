package soothe_test

import (
	"fmt"

	soothe "github.com/mirasoth/soothe-client-go"
)

// Example_shouldShow demonstrates content visibility filtering under nano's
// binary tier model: NORMAL-tier content is visible, INTERNAL-tier is not.
func Example_shouldShow() {
	fmt.Println(soothe.ShouldShow(soothe.TierNormal))   // true
	fmt.Println(soothe.ShouldShow(soothe.TierInternal)) // false
	// Output:
	// true
	// false
}

// Example_isValidVerbosityLevel checks whether a string is a recognized
// verbosity level ("quiet", "normal", or "debug").
func Example_isValidVerbosityLevel() {
	fmt.Println(soothe.IsValidVerbosityLevel("quiet"))  // true
	fmt.Println(soothe.IsValidVerbosityLevel("normal")) // true
	fmt.Println(soothe.IsValidVerbosityLevel("debug"))  // true
	fmt.Println(soothe.IsValidVerbosityLevel("trace"))  // false
	fmt.Println(soothe.IsValidVerbosityLevel(""))       // false
	// Output:
	// true
	// true
	// true
	// false
	// false
}

// Example_verbosityLevels shows the typed verbosity level and tier constants.
func Example_verbosityLevels() {
	// VerbosityLevel is a string type.
	level := soothe.VerbosityDebug
	fmt.Println(level) // debug

	// wire_tier mapping for loop_events subscriptions.
	fmt.Println(soothe.WireTierForVerbosity(soothe.VerbosityQuiet))  // compact
	fmt.Println(soothe.WireTierForVerbosity(soothe.VerbosityNormal)) // full
	// VerbosityTier is an int type.
	fmt.Println(int(soothe.TierNormal))   // 1
	fmt.Println(int(soothe.TierInternal)) // 99
	// Output:
	// debug
	// compact
	// full
	// 1
	// 99
}
