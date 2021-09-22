package defaultexplorerserver

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/bang9211/wire-jacket/internal/config"
	"github.com/bang9211/wire-jacket/wirejackettest/ossiconesblockchain"
)

type ExplorerServer interface {
	// Serve listens and serves the Explorer Server.
	Serve()
	// Close closes the Explorer Server.
	Close() error
}

const (
	defaultDHSHost = "0.0.0.0"
	defaultDHSPort = 3000
)

var dhs *DefaultExplorerServer
var once sync.Once

type DefaultExplorerServer struct {
	config     config.Config
	handler    *http.ServeMux
	blockchain ossiconesblockchain.Blockchain
	address    string
}

// GetOrCreate returns the existing singletone object of DefaultHTTPServer.
// Otherwise, it creates and returns the object.
func GetOrCreate(
	config config.Config,
	blocchain ossiconesblockchain.Blockchain) ExplorerServer {
	if dhs == nil {
		once.Do(func() {
			dhs = &DefaultExplorerServer{
				config:     config,
				handler:    http.NewServeMux(),
				blockchain: blocchain,
			}
		})
		err := dhs.init()
		if err != nil {
			dhs = nil
			return nil
		}
	}

	return dhs
}

func (d *DefaultExplorerServer) init() error {
	host := d.config.GetString("ossicones_explorer_server_host", defaultDHSHost)
	port := d.config.GetInt("ossicones_explorer_server_port", defaultDHSPort)
	d.address = host + ":" + strconv.Itoa(port)

	d.Serve()

	return nil
}

func (d *DefaultExplorerServer) Serve() {
	go func() {
		fmt.Printf("Listening Explorer Server on %s\n", d.address)
		log.Fatal(http.ListenAndServe(d.address, d.handler))
	}()
}

func (d *DefaultExplorerServer) Close() error {
	dhs = nil
	return nil
}
