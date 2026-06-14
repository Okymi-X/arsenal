package registry

// Source abstracts where a registry comes from and how it is refreshed.
//
// Load returns the currently available registry. Sync refreshes the local
// copy from its upstream origin. Implementations are injected so that tests
// and offline modes can substitute a static or fake source.
type Source interface {
	// Load returns the currently available registry.
	Load() (*Registry, error)
	// Sync refreshes the local copy from its upstream origin.
	Sync() error
}
