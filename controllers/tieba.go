package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
	"net/http"
	"io"
	"strconv"
	"os"
)

type TieBaController struct {
	beego.Controller
}

func (tb *TieBaController) TieBa()  {
	working(1,10)
}

func HttpGet(url string) (result string,err error) {
	res,err1 := http.Get(url)
	if err1 != nil {
		err = err1
		return
	}

	defer  res.Body.Close()

	//循环读取数据
	buf := make([]byte,4096)
	for {
		n,err2 := res.Body.Read(buf)
		if n == 0 {
			//爬取完成
			break
		}
		if err2 != nil && err2 != io.EOF{
			err = err2
		}
		//累加每次读取到的数据
		result += string(buf[:n])
	}
	return
}
/**
https://tieba.baidu.com/f?kw=%E7%BB%9D%E5%9C%B0%E6%B1%82%E7%94%9F&ie=utf-8&pn=50
 */
func working(start,end int)  {
	fmt.Printf("正在爬取第%d页-%d页...\n",start,end)

	pageChan := make(chan int)

	for i := start; i <= end; i++  {
		//并发爬取
		go SpiderPage(i,pageChan)
	}

	//等待子go程完成
	for i := start; i <= end; i++  {
		fmt.Printf("第%d页面爬取完成\n", <- pageChan)
	}
}

func SpiderPage(i int,pageChan chan int)  {
	url := "https://tieba.baidu.com/f?kw=%E7%BB%9D%E5%9C%B0%E6%B1%82%E7%94%9F&ie=utf-8&pn="+strconv.Itoa((i-1)*50)

	result,err := HttpGet(url)
	if err != nil {
		fmt.Println("err : ",err)
		return
	}
	//fmt.Println("result=",result)
	//创建文件
	f,_ := os.Create("第" + strconv.Itoa(i) + "页" + ".html")
	f.WriteString(result)
	f.Close()

	//与主go程完成同步
	pageChan <- i
}