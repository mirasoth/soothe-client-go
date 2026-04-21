// Package soothe provides a Go client for the Soothe daemon WebSocket API.
//
// The client implements the same protocol as the Python soothe-sdk, providing
// full access to the Soothe daemon's capabilities including thread management,
// event streaming, skills/models discovery, and daemon control.
//
// The package also defines verbosity types (VerbosityLevel, VerbosityTier) used
// for event classification and filtering based on user-configurable verbosity settings.
package soothe
