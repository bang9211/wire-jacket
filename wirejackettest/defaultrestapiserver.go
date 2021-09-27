package wirejackettest

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/bang9211/wire-jacket/internal/config"
	"github.com/gorilla/mux"
)

type RESTAPIServer interface {
	// Serve listens and serves the REST API Server.
	Serve()
	// Close closes the REST API Server.
	Close() error
}

const (
	defaultDRSHost = "0.0.0.0"
	defaultDRSPort = 4000
)

var drs *DefaultRESTAPIServer
var drsOnce sync.Once

type defaultURL struct {
	address string
	path    string
}

func (u *defaultURL) String() string {
	return u.path
}

func (u *defaultURL) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("%s%s", u.address, u.path)
	return []byte(url), nil
}

type urlDescription struct {
	URL         defaultURL `json:"url"`
	Method      string     `json:"method"`
	Description string     `json:"description"`
	Payload     string     `json:"payload,omitempty"`
}

func (u urlDescription) String() string {
	return ""
}

type AddBlockBody struct {
	Message string
}

type ErrorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

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
		})
		if drs == nil {
			return nil
		}
		err := drs.init()
		if err != nil {
			drs = nil
			return nil
		}
	}

	return drs
}

func (d *DefaultRESTAPIServer) init() error {
	host := d.config.GetString("ossicones_rest_api_server_host", defaultDRSHost)
	port := d.config.GetInt("ossicones_rest_api_server_port", defaultDRSPort)
	d.address = host + ":" + strconv.Itoa(port)

	d.Serve()

	return nil
}

func (d *DefaultRESTAPIServer) Serve() {
	go func() {
		fmt.Printf("Listening REST API Server on %s\n", d.address)
		log.Fatal(http.ListenAndServe(d.address, d.handler))
	}()
}

func (d *DefaultRESTAPIServer) Close() error {
	// drs = nil
	return nil
}
