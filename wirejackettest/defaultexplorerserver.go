package wirejackettest

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/bang9211/wire-jacket/internal/config"
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
var dhsOnce sync.Once
var templates *template.Template

type homeData struct {
	PageTitle string
	Blocks    []Block
}

type DefaultExplorerServer struct {
	config     config.Config
	handler    *http.ServeMux
	blockchain Blockchain
	homePath   string
	address    string
}

// GetOrCreate returns the existing singletone object of DefaultHTTPServer.
// Otherwise, it creates and returns the object.
func GetOrCreateExplorerServer(
	config config.Config,
	blocchain Blockchain) ExplorerServer {
	if dhs == nil {
		dhsOnce.Do(func() {
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
	var err error
	// d.homePath, err = utils.GetOrSetHomePath()
	if err != nil {
		return err
	}
	host := d.config.GetString("ossicones_explorer_server_host", defaultDHSHost)
	port := d.config.GetInt("ossicones_explorer_server_port", defaultDHSPort)
	d.address = host + ":" + strconv.Itoa(port)

	// templates = template.Must(template.ParseGlob(d.homePath + "/templates/pages/*.gohtml"))
	// templates = template.Must(templates.ParseGlob(d.homePath + "/templates/partials/*.gohtml"))

	// d.handler.HandleFunc("/", d.home)
	// d.handler.HandleFunc("/add", d.add)

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
