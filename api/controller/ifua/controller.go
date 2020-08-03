package ifua

import (
	"github.com/gorilla/mux"
)

// Controller struct
type Controller struct {
}

// NewController func
func NewController() *Controller {
	return &Controller{}
}

// Route func
func (c *Controller) Route(r *mux.Router) {
	//routes grouping
	s := r.PathPrefix("/dummy/bri").Subrouter()
	s.HandleFunc("/test", c.Test).Methods("GET")
	s.HandleFunc("/reqifua", c.sendDataIfuaReq).Methods("POST")
}
