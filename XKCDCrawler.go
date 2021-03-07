package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/antchfx/htmlquery"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
)

const (
	xkcdDomain = "https://xkcd.in"
)

type Comic struct {
	Id          string
	Title       string
	ImageUrl    string
	NextLink    string
	Description string
}

var saveDir = ""

func init() {
	flag.StringVar(&saveDir, "o", "", "下载图片和信息保存路径")
	flag.Parse()
}

func main() {
	if saveDir == "" {
		flag.Usage()
		os.Exit(1)
	}
	crawl("https://xkcd.in")
	log.Printf("crawl finish!")
}

// 抓取
func crawl(url string) {
	log.Printf("crawl url: %s", url)
	comic, err := fetchHtml(url)
	if err != nil {
		log.Fatal("%s Error: %s \n", url, err)
	} else {
		downloadImage(comic)
		saveComic(comic)
	}
	if comic.NextLink != "" {
		crawl(comic.NextLink)
	}
}

// 抓取并解析漫画HTML
func fetchHtml(link string) (*Comic, error) {
	doc, err := htmlquery.LoadURL(link)
	if err != nil {
		return nil, err
	}

	comicBodyNode := htmlquery.FindOne(doc, "//div[@class='comic-body']")
	if comicBodyNode == nil {
		return nil, errors.New("页面错误")
	}

	imgNode := htmlquery.FindOne(comicBodyNode, "//img")
	imageUrl := urlJoin(xkcdDomain, htmlquery.SelectAttr(imgNode, "src"))
	title := htmlquery.SelectAttr(imgNode, "alt")

	comicDetailsNode := htmlquery.FindOne(comicBodyNode, "//div[@class='comic-details']")
	description := htmlquery.InnerText(comicDetailsNode)

	nextLinkNode := htmlquery.FindOne(comicBodyNode, "//div[@class='nextLink']/a")
	nextLink := ""
	if nextLinkNode != nil {
		nextLink = htmlquery.SelectAttr(nextLinkNode, "href")
		nextLink = urlJoin(xkcdDomain, nextLink)
	}
	id := parseId(link, nextLink)
	comic := Comic{id, title, imageUrl, nextLink, description}

	return &comic, nil
}

// 拼接域名和相对URL
func urlJoin(base string, path string) string {
	u, err := url.Parse(path)
	if err != nil {
		log.Fatal(err)
	}
	b, err := url.Parse(base)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprint(b.ResolveReference(u))
}

// 解析漫画ID
func parseId(link string, nextLink string) string {
	u, _ := url.Parse(link)
	q := u.Query()
	id := q.Get("id")
	if id == "" {
		u, _ = url.Parse(nextLink)
		q = u.Query()
		iid, _ := strconv.Atoi(q.Get("id"))
		id = fmt.Sprintf("%d", iid+1)
	}
	return id
}

// 下载并保存漫画图片
func downloadImage(comic *Comic) {
	// 判断下载的图片是否存在，不存在去下载
	dir := saveDir + comic.Id
	mkDirs(dir)
	filename := dir + "/" + path.Base(comic.ImageUrl)
	if exists(filename) {
		return
	}

	//下载图片
	res, err := http.Get(comic.ImageUrl)
	if err != nil {
		log.Printf("%d-Error: %s \n", comic.Id, err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Printf("%d-Http StatusCode %d \n", comic.Id, res.StatusCode)
		return
	}

	// 写到文件
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("%d-Error: %s \n", comic.Id, err)
		return
	}
	err = ioutil.WriteFile(filename, body, 0644)
	if err != nil {
		log.Printf("%d-Error: %s \n", comic.Id, err)
		return
	}
}

// 将Comic保存为json
func saveComic(comic *Comic) {
	dir := path.Join(saveDir, comic.Id)
	mkDirs(dir)
	filename := path.Join(dir, comic.Id+".json")
	if exists(filename) {
		return
	}
	s, err := json.Marshal(comic)
	if err != nil {
		log.Printf("%d-Error: %s \n", comic.Id, err)
		return
	}

	err = ioutil.WriteFile(filename, s, 0644)
	if err != nil {
		log.Printf("%d-Error: %s \n", comic.Id, err)
		return
	}
}

// 创建目录
func mkDirs(dir string) {
	if !exists(dir) {
		os.MkdirAll(dir, 0777)
		os.Chmod(dir, 0777)
	}
}

// 判断文件或目录是否存在
func exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
