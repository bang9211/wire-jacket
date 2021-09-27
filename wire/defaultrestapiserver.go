package wire

import (
	"strconv"
	"sync"

	"github.com/bang9211/wire-jacket/internal/config"
	"github.com/gorilla/mux"
)

type RESTAPIServer interface {
	// Serve listens and serves the REST API Server.
	Serve()
	// Get gets routing paths
	GetPaths() []string
	// Close closes the REST API Server.
	Close() error
}

const (
	defaultDRSHost = "0.0.0.0"
	defaultDRSPort = 4000
)

var drs *DefaultRESTAPIServer
var drsOnce sync.Once

type DefaultRESTAPIServer struct {
	config     config.Config
	handler    *mux.Router
	blockchain Blockchain
	homePath   string
	address    string
}

// GetOrCreate returns the existing singletone object of DefaultAPIServer.
// Otherwise, it creates and returns the object.
func GetOrCreateDefaultRESTAPIServer(
	config config.Config,
	blocchain Blockchain) RESTAPIServer {
	if drs == nil {
		drsOnce.Do(func() {
			drs = &DefaultRESTAPIServer{
				config:     config,
				handler:    mux.NewRouter(),
				blockchain: blocchain,
			}
			host := drs.config.GetString("ossicones_restapi_server_host", defaultDRSHost)
			port := drs.config.GetInt("ossicones_restapi_server_port", defaultDRSPort)
			drs.address = host + ":" + strconv.Itoa(port)
			drs.Serve()
		})
	}

	return drs
}

func (d *DefaultRESTAPIServer) Serve() {
	// go func() {
	// 	fmt.Printf("Listening REST API Server on %s\n", d.address)
	// 	log.Fatal(http.ListenAndServe(d.address, d.handler))
	// }()
}

func (d *DefaultRESTAPIServer) GetPaths() []string {
	return []string{"/"}
}

func (d *DefaultRESTAPIServer) Close() error {
	// drs = nil
	return nil
}
