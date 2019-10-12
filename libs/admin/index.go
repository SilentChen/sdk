package admin

import (
	"../../packs/gin"
	"os"
	"runtime"
)

type Index struct {
	Base
}

func (_ *Index) Index(c *gin.Context) {
	this.display(c, map[string]interface{}{
		"admin_name"	:		"admin",
		"version"		:		"1.0.1",
	})
}

func (_ *Index) Main(c *gin.Context) {
	this.display(c, map[string]interface{}{
		"app_ver"	:	"1.0.1",
		"hostname"	:	getHostName(),
		"go_ver"	:	runtime.Version(),
		"os"		:	runtime.GOOS,
		"cpu_num"	:	runtime.NumCPU(),
		"arch"		:	runtime.GOARCH,
		"postnum"	:	0,
		"tagnum"	:	0,
		"usernum"	:	0,
	})
}

func (this *Index) set(c *gin.Context) {

	this.display(c, map[string]interface{}{
		"webtitle"		:		"test",
		"websubtitle"	:		"test",
		"weburl"		:		"test",
		"webemail"		:		"test",
		"webpagenum"	:		"test",
		"webkeyword"	:		"test",
		"webdesc"		:		"test",
		"webtheme"		:		"test",
	})
}

/**
 * Private Begin.
 */
func getHostName() string {
	hname, err := os.Hostname()
	if nil != err {
		hname = "localhost"
	}

	return hname
}