package controller

import (
	"context"

	"github.com/Sunaga0709/go-api-template/internal/restapi/gen"
)

func (c *Controller) GetHealth(
	_ context.Context,
	_ gen.GetHealthRequestObject,
) (gen.GetHealthResponseObject, error) {
	return gen.GetHealth200JSONResponse{Status: "ok"}, nil
}
