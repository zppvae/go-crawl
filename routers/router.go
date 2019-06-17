package routers

import (
	"go-crawl/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
    beego.Router("/movie", &controllers.CrawlMovieController{}, "*:CrawlMovie")
	beego.Router("/tieba", &controllers.TieBaController{}, "*:TieBa")
	beego.Router("/douban", &controllers.DouBanController{}, "*:DouBan")
}
