package app

import (
	"fmt"
	"../../packs/gin"
	"../../packs/util"
)

type Index struct {
	Base
}

func (this *Base) Index(c *gin.Context) {
	totalNum :=	this.mysqlInstance().DefGetOne("select count(*) from `post` where status = 0", "0")

	page := util.Str2int(c.DefaultQuery("page", "0"))

	_, list, _ := this.mysqlInstance().GetAll(fmt.Sprintf("select * from `post` where status = 0 limit %d,%d", this.pageOffset(page), this.pageSize()))

	pagebar := util.NewPager(page, util.Str2int(totalNum), this.pageSize(), "/", true).ToString()

	base :=	getBaseData()

	this.display(c, map[string]interface{}{
		"pagebar"		:		pagebar,
		"list"			:		list,
		"base"			:		base,
	})
}

func (this *Base) Article(c *gin.Context) {
	var title, content string

	aid := c.DefaultQuery("aid", "-1")
	if "-1" != aid {
		row,  _ := this.mysqlInstance().GetRow("select title,content from `post` where id = " + aid)
		title  = (*row)["title"]
		content = (*row)["content"]

	}

	if "" == title || "" == content {
		this.errorShow(c, []string{"No Article!"})
		return
	}

	base :=	getBaseData()

	this.display(c, map[string]interface{}{
		"title"		:		title,
		"content"	:		content,
		"base"		:		base,
	})
}

func (this *Base) Category(c *gin.Context) {

	var(
		list map[string][]map[string]string
		tmplist []map[string]string
		where string
		tmp string
	)

	where = " where status = 0 "
	tag := c.DefaultQuery("tag", "default")

	if "default" != tag {
		where += fmt.Sprintf(" and tags = '%s' ", tag)
	}

	_, tmplist, _  = this.mysqlInstance().GetAll(fmt.Sprintf("select id,title,post_time,tags from `post` %s order by post_time desc", where))

	if len(tmplist) > 0 {
		list = make(map[string][]map[string]string)
		for _, row := range tmplist  {
			tmp = util.Unix2year(util.Date2unix(row["post_time"]))
			list[tmp] = append(list[tmp], row)
		}
	}

	base :=	getBaseData()

	this.display(c, map[string]interface{}{
		"list"		:		list,
		"base"		:		base,
	})
}

func getBaseData() map[string]interface{} {
	_, tags,_ := this.mysqlInstance().GetAll("select distinct tags from `post` where status = 0")

	options := map[string]string{
		"sitename"		:		"test",
		"subtitle"		:		"test",
		"siteurl"		:		"test",
		"stat"			:		"test",
	}

	return map[string]interface{} {
		"taglist"	:		tags,
		"options"	:		options,
	}
}