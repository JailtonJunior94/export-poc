package excel

import (
	"context"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type Provider interface {
	NewFile(context.Context) File
}

type client struct{}

func NewProvider() Provider {
	return &client{}
}

func (c *client) NewFile(ctx context.Context) File {
	xls := excelize.NewFile()
	return newFile(xls)
}
