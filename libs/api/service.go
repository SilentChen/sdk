package api

import (
	"../../packs/gin"
	"../../packs/util"
)

type Service struct {
	Base
}

func (this *Service)Index(c *gin.Context)  {
	this.ResponseJson(c, util.HTTPSTATUSCODE_OK, util.JSONSTATUSCODE_Success, "success")
}

func (this *Service) Server_list(c *gin.Context) {
	channel  := c.DefaultQuery("channel", "")
	platform := c.DefaultQuery("platform", "")
	page     := c.DefaultQuery("page", "")
	accname  := c.DefaultQuery("accname", "")

	this.ResponseJson(c, util.HTTPSTATUSCODE_OK, util.JSONSTATUSCODE_Success, []string{
		channel,
		platform,
		page,
		accname,
	})
}