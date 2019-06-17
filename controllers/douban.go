package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
	"sync"
	"strconv"
	"regexp"
	"os"
)

/*
   @Time : 2019/6/17 15:42 
   @Author : ff
*/
var (
	wg sync.WaitGroup
)

type DouBanController struct {
	beego.Controller
}

func (db *DouBanController) DouBan()  {
	doubanWork(1,1)
}

func doubanWork(start,end int)  {
	fmt.Printf("正在爬取第%d页-%d页...\n",start,end)

	for i := start; i <= end; i++  {
		wg.Add(1)
		go spider(i)
	}

	wg.Wait()

	fmt.Printf("爬取第%d页-%d页任务完成\n",start,end)
}

func spider(index int)  {
	url := "https://movie.douban.com/top250?start="+strconv.Itoa((index-1)*25)+"&filter="
	result,err := HttpGet(url)
	if err != nil {
		fmt.Println("err : ",err)
		return
	}

	req := regexp.MustCompile(`<img width="100" alt="(?s:(.*?))"`)
	//
	filmName := req.FindAllStringSubmatch(result,-1)

	reqScore := regexp.MustCompile(`<span class="rating_num" property="v:average">(?s:(.*?))</span>`)
	filmScore := reqScore.FindAllStringSubmatch(result,-1)

	save2file(index,filmName,filmScore)
	wg.Done()
}

/**
  保存到文件
 */
func save2file(index int,filmName,filmScore [][]string)  {
	path := "D:/" + "第 " +strconv.Itoa(index)+ " 页.txt"
	f,_ := os.Create(path)

	defer f.Close()

	//获取爬取的电影数
	n := len(filmName)
	f.WriteString("电影名称" + "\t\t\t" + "电影评分" + "\n")
	for i := 0; i < n; i++ {
		f.WriteString(filmName[i][1] + "\t\t\t" + filmScore[i][1] + "\n")
	}

}