package server

import (
	"fmt"
	"net/http"

	"github.com/CraigYanitski/mescli/client"
)

func CreateClient(name, password string) (*client.Client, error) {
    c := &client.Client{Name: name}
    err := c.HashPassword(password)
    if err != nil {
        return nil, err
    }
    return c, nil
}
