package routes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	assert.IsType(t, Routes{}, routes, "Unexpected type of routes")
	assert.Len(t, routes, 3, "Unexpected number of routes")
	assert.Equal(t, "Status", routes[0].Name, "Unexpected route name")
	assert.Equal(t, "Product", routes[1].Name, "Unexpected route name")
	assert.Equal(t, "Product", routes[2].Name, "Unexpected route name")
}
