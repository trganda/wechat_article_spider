# Wechat Article Crawer

微信公众号文章爬取的半自动化工具，当初是为了方便获取玄武实验室的每日推送而写的。并不适合大规模批量的获取公众号的文章内容，也不打算这么去做。

## Table of Contents
- [Introduction](#introduction)
   - [ChromeDriver](#chromedriver)
   - [Usage](#usage)  
- [Features](#features)
- [Principle](#Principle)
- [License](#license)

## Introduction

### ChromeDriver

为了方便登录和获取`cookie`，需要使用到`selenium`，请先下载对应`chrome`浏览器版本的[chromedriver](https://sites.google.com/a/chromium.org/chromedriver/)，并同时安装`chrome`浏览器。

将下载好的`brower driver`放在项目的`vendor`文件夹下，并修改配置文件中它的位置。当然，你也可以放在其它路径下。

### Usage

安装好浏览器和对应的`brower dirver`后，直接运行

```go
go run main.go
```

会打印出json格式的文章相关数据，及其链接。

## Features

## Principle

## License

