package routes

import (
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
	router := APIRouter()
	assert.IsType(t, mux.NewRouter(), router, "they should be equal")
}
