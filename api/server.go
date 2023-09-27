package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/okoroemeka/simple_bank/db/sqlc"
)

// Server serves http request for our banking app
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer creates a new http server and setup routing
func NewServer(store *db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.PUT("/accounts/:id", server.updateAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
