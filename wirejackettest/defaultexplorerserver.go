package wirejackettest

type ExplorerServer interface {
	// Serve listens and serves the Explorer Server.
	Serve()
	// Close closes the Explorer Server.
	Close() error
}
