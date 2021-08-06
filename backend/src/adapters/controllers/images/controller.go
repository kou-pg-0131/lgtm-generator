package images

import (
	"github.com/kou-pg-0131/lgtm-generator/backend/src/adapters/controllers"
	"github.com/kou-pg-0131/lgtm-generator/backend/src/entities"
	"github.com/pkg/errors"
)

type Controller struct {
	config *ControllerConfig
}

type ControllerConfig struct {
	Renderer         controllers.Renderer
	ImagesRepository controllers.ImagesRepository
}

func NewController(cfg *ControllerConfig) *Controller {
	return &Controller{config: cfg}
}

func (ctrl *Controller) Search(ctx controllers.Context) {
	q := ctx.Query("q")
	if q == "" {
		ctrl.config.Renderer.BadRequest(ctx, entities.ErrCodeQueryIsEmpty, errors.New("query is empty"))
		return
	}

	imgs, err := ctrl.config.ImagesRepository.Search(q)
	if err != nil {
		ctrl.config.Renderer.InternalServerError(ctx, errors.WithStack(err))
		return
	}

	ctrl.config.Renderer.OK(ctx, imgs)
}
