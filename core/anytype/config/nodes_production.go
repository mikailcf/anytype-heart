//go:build !envnetworkcustom

package config

// Offline-only build: no production nodes configured.
// The embedded config is empty to prevent any network connections.
var nodesConfYmlBytes []byte
