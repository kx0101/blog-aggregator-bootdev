package middlewares

import (
	"github.com/kx0101/blog-aggregator-bootdev/internal/database"
)

type APIConfig struct {
	DBQueries *database.Queries
}
