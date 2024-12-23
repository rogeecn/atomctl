package {{.ModuleName}}

import (
	log "github.com/sirupsen/logrus"
)

// @provider
type Controller struct {
	svc *Service
	log *log.Entry `inject:"false"`
}

func (c *Controller) Prepare() error {
	c.log = log.WithField("module", "{{.ModuleName}}.Controller")
	return nil
}

// Test godoc
//
//	@Summary		Test
//	@Description	Test
//	@Tags			Test
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"AccountID"
//	@Param			queryFilter	query		dto.AlarmListQuery		true	"AlarmListQueryFilter"
//	@Param			pageFilter	query		common.PageQueryFilter	true	"PageQueryFilter"
//	@Param			sortFilter	query		common.SortQueryFilter	true	"SortQueryFilter"
//	@Success		200			{object}	common.PageDataResponse{list=dto.AlarmItem}
//	@Router			/v1/test/:id<int> [get]