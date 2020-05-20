package server

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

// TODO add routes registration
// TODO add api version mechanism
// TODO add check user requests params

// TODO routes:
// - process/resize image
// - list prev resized images with inputted params
// - all api can return errors
// - return obj - links to orig and resized images
// - if requested prev image with prev params - just return already processed img link

// APIServer process http requests
// Using api version mechanism
type APIServer struct {
	logger  *logrus.Logger
	address string
	router  *mux.Router
}

// create new instance of APIServer
func NewAPIServer(bindAddr string, logger *logrus.Logger) *APIServer {
	return &APIServer{
		logger:  logger,
		address: bindAddr,
		router:  mux.NewRouter(),
	}
}

// Starting APIServer using port from config
func (s *APIServer) Start() error {
	s.logger.Infof("Starting api server. Port %s ", s.address)
	return http.ListenAndServe(s.address, s.router)
}
