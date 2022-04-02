# -*- coding: utf-8 -*-
import requests
import time
import json


# 目标url
url = "https://mp.weixin.qq.com/cgi-bin/appmsg"

# 使用Cookie，跳过登陆操作
# headers = {
#   "Cookie": r"appmsglist_action_3912343346=card; ua_id=j0lbwaYANIsJh8K9AAAAAAuU0AwayCoaEQXZ8BlsjMM=; wxuin=48724447363676; uuid=5235998a046de218dcfacac397b0c516; rand_info=CAESIFBIv/W6IkVT1YdF9AfqTQHp0MhS7vaTrqDxnhfjrjj8; slave_bizuin=3912343346; data_bizuin=3912343346; bizuin=3912343346; data_ticket=lsM1ty+r5zyv0mTkODMVsZQC9vrhg3SwFBBKFXd1HN/O1JLUpA/fUiHd8FNq+im+; slave_sid=NldXQ0xRdHdWWHJrMGY4Sm9iaXhyQzhEX0ZwQkdjZm5IUVFVXzdWMnBuazJVbWgxV2syZ0dYVFlwN0Zfc3dGSEt0WUNCb0Z2TXc4MkxzcmluNVIyaTQ1dkZoZmtieUg5R21xUW5vX3BjdTJIZ2wyRWtsR2s2ZXA3aU5vN0xYam9ON2FMMXBhSU9sdW5FR0hQ; slave_user=gh_63f2d6498c0d; xid=589049049e73662ee710df1cf38d09aa; mm_lang=zh_CN",
#   "User-Agent": r"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.62 Safari/537.36",
# }

headers = {
  "Cookie": r"mm_lang=zh_CN; slave_user=gh_63f2d6498c0d; ua_id=EOY6zRld1fngTZcaAAAAAD-w73HIpTVaGndwiFNajcE=; xid=bf2db27c7025d0dbb6c361dc1cff2a7a; slave_sid=Z0dxbE54cWR2Z3pySVI5bk9KejFaWEswMjdQMkN0TGsyM3lidUZZd0hQTGJfRjJCcWZndnAwMmk3QUpZQ1dsSUs1MFh0U0hkd0ZETUZobXljSllIbjUzaUQzRzh3VU5KYlpJcUtGZEhMT0FURVIyNU9HRFZkTEU0ZkFWVGpaWm9kdjlhQ002emUwVHMzdWtB; uuid=102eee91d14d33477da5e75540d48d00; data_ticket=fxFOcR0KOjoq7ToxqSVHkA3uKIWQT3mXHec5eOZcqOWvud0r3QuSOpquzNNiwNQJ; bizuin=3912343346; data_bizuin=3912343346; slave_bizuin=3912343346; wxuin=48863846493347; rand_info=CAESIA3rnBNMk1g2VMrwAwD6QoXbhMmwBnE5B8+Mg8ICSx1X",
  "User-Agent": r"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.62 Safari/537.36",
}

"""
需要提交的data
以下个别字段是否一定需要还未验证。
注意修改yourtoken,number
number表示从第number页开始爬取，为5的倍数，从0开始。如0、5、10……
token可以使用Chrome自带的工具进行获取
fakeid是公众号独一无二的一个id，等同于后面的__biz
"""
data = {
    "token": "8138345",
    "lang": "zh_CN",
    "f": "json",
    "ajax": "1",
    "action": "list_ex",
    "begin": "0",
    "count": "5",
    "query": "",
    "fakeid": "MzA5NDYyNDI0MA==",
    "type": "9",
}

# 使用get方法进行提交
content_json = requests.get(url, headers=headers, params=data).json()
# 返回了一个json，里面是每一页的数据
for item in content_json["app_msg_list"]:
    # 提取每页文章的标题及对应的url
    print(item["title"], "url:", item["link"])
