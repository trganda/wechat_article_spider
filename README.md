# Wechat Article Spider

微信公众号文章爬取的半自动化工具，当初是为了方便获取玄武实验室的每日推送而写的。并不适合大规模批量的获取公众号的文章内容，也不打算这么去做。由于不确定反爬虫机制，所以每次请求之间会有一些时间间隔。

## Table of Contents
- [Introduction](#introduction)
   - [ChromeDriver](#chromedriver)
   - [Usage](#usage)  
- [Features](#features)
- [Principle](#Principle)

## Introduction

### ChromeDriver

为了方便登录和获取`cookie`，需要使用到`selenium`，请先下载对应`chrome`浏览器版本的[chromedriver](https://sites.google.com/a/chromium.org/chromedriver/)，并同时安装`chrome`浏览器。

将下载好的`brower driver`放在项目的`vendor`文件夹下，并修改配置文件中它的位置。当然，你也可以放在其它路径下。

### Usage

安装好浏览器和对应的`brower dirver`后，直接运行

```go
go run main.go
```

会打印出json格式的文章相关数据，如下：

```json
{
  "app_msg_list": [
    {
      "aid": "2651958330_1",
      "album_id": "0",
      "appmsgid": 2651958330,
      "checking": 0,
      "cover": "https://mmbiz.qlogo.cn/mmbiz_jpg/dWDic6IAXZsfZ5NKcSyULDMmjMncfAus29aTXCgabeiavgsebgt93sL07iahdxagl04wD6NwuJKCRalEXibDpghUwA/0?wx_fmt=jpeg",
      "create_time": 1648888670,
      "digest": "PHP Supply Chain Attack on PEAR；Go 语言将应用新 Mitigation 防御供应链攻击",
      "itemidx": 1,
      "link": "http://mp.weixin.qq.com/s?__biz=MzA5NDYyNDI0MA==\u0026mid=2651958330\u0026idx=1\u0026sn=a14fb5f431821a63dff80b219906e029\u0026chksm=8baecca5bcd945b3c1597d267fcd79304c7a5eac32dbb4be2a81f42aee7a9be6e15190e6d86d#rd",
      "title": "每日安全动态推送(04-02)",
      "update_time": 1648888670
    }
  ]
}
```

## Features

* ~~支持关键词搜索公众号~~
* 支持关键词搜索文章
* 支持根据时间获取文章

## Principle

