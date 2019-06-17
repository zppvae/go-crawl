package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
	"regexp"
	"net/http"
	"os"
	"strconv"
)

type DouYuController struct {
	beego.Controller
}

func (tb *DouYuController) DouYu()  {
	douyuWork(1,2)
}

func douyuWork(start,end int)  {
	fmt.Printf("正在爬取第%d页-%d页...\n",start,end)

	pageChan := make(chan int)

	for i := start; i <= end; i++  {
		//并发爬取
		go SpiderDouyu(i,pageChan)
	}

	//等待子go程完成
	for i := start; i <= end; i++  {
		fmt.Printf("第%d页爬取完成\n", <- pageChan)
	}
}

func SpiderDouyu(i int,pageChan chan int)  {
	url := "https://www.douyu.com/g_yz"

	result,err := HttpGet(url)
	if err != nil {
		fmt.Println("err : ",err)
		return
	}

	p := `"rs16":"(?s:(.*?))"`
	req := regexp.MustCompile(p)

	alls := req.FindAllStringSubmatch(result,-1)

	for index,imgUrl := range alls  {
		fmt.Println(imgUrl[1])
		go SaveImg(index,imgUrl[1])
	}

	//与主go程完成同步
	pageChan <- i
}

func SaveImg(index int,url string)  {
	path := "F:/goProject/src/go-crawl/static/img/"
	f,err1 := os.Create(path + strconv.Itoa(index+1)  + ".jpg")
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	defer f.Close()

	res,err2 := http.Get(url)
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	defer  res.Body.Close()

	//循环读取数据
	buf := make([]byte,4096)
	for {
		n,_ := res.Body.Read(buf)
		if n == 0 {
			//爬取完成
			break
		}
		f.Write(buf[:n])
	}
}