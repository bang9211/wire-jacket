package wirejackettest

type RESTAPIServer interface {
	// Serve listens and serves the REST API Server.
	Serve()
	// Close closes the REST API Server.
	Close() error
}
