# 抓取XKDC中文网站的漫画图片

XKCD中文站，是一个关于浪漫、隐喻、数字、以及语言的线上漫画。该站有好多有趣的漫画。我们的目标是将漫画图片下载并保存到本地目录，并且将漫画信息保存为json格式

## 编译运行

在项目目录执行：
go run XKCDCrawler.go  -o 下载图片和漫画Json的保存路径

## 运行结果

下载图片和漫画Json到指定目录

保存的漫画信息Json文件如下：
```json
{ 
  "Id":"2427",
  "Title":"毅力号麦克风",
  "ImageUrl":"https://xkcd.in/resources/compiled_cn/1020d8737d5e905092716e6fa1c60395.jpg",
  "NextLink":"https://xkcd.in/comic?lg=cn\u0026id=2426",
  "Description":"如果他们首先接收到的是探测器穿过大气层时的音频，那么我们大概除了探测器的尖鸣声，什么也听不见吧。"
}
```
