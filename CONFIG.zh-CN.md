## 配置文件说明

```yaml
version: 1.5
addr: :10000                          # HTTP 服务监听的地址
token:                                # 鉴权Token，管理页面和客户端都需要通过该 Token 进行鉴权
domains: []                           # 反连平台绑定的域名
external_ip:                          # 反连平台的外部IP，需要确保你想测试的目标能通过该 IP 访问到平台
admin_path_prefix: "/revsuit"          # 管理页面的 http path 前缀，管理页面将位于：/admin_path_prefix/admin
# 数据库连接信息 支持MySQL、Postgres、SQLite3
database: "mysql://root:password@tcp(127.0.0.1:3306)/revsuit?charset=utf8mb4&parseTime=True&loc=Local"
# database: "postgres://host=127.0.0.1 user=root password=password dbname=revsuit port=5432 sslmode=disable TimeZone=Asia/Shanghai"
#database: revsuit.db # 在部分系统上使用 sqlite 数据库时，在并发场景下可能会出现 `SQLITE BUSY` 的问题，因此不推荐在正式环境中使用该类型数据库

log_level: info                       # 输出日志的级别，分为：debug、info、warning、error、fatal
check_upgrade: false                  # 是否自动检查更新

ip_location_database:                 # IP 数据库相关配置
  database: "qqwry"                   # qqwry 或者 geoip.
  geo_license_key: ""                 # 如果你使用 GeoIP 则该字段为必填
  
http:                                 # HTTP 反连相关配置
  ip_header:                          # 通过 HTTP 请求头获取来源IP，通常在有前置反向代理时需要配置
dns:                                  # DNS 反连相关配置
  enable: true                             
  addr: :53                           # DNS 服务监听的地址  
rmi:                                  # RMI 反连相关配置
  enable: true                                
  addr: :1099                         # RMI 服务监听的地址
ldap:
  enable: true
  addr: :1389                         # Ldap 服务监听的地址
mysql:                                # MySQL 反连相关配置
  enable: true                        
  addr: :3306                         # MySQL 服务监听的地址        
  version_string: 10.4.13-MariaDB-log # 伪装MySQL服务的版本
ftp:                                  # FTP 反连相关配置
  enable: true                        
  addr: :21                           # FTP 服务监听的地址
  pasv_port: 2020                     # FTP被动模式端口监听的地址
notice:
  dingtalk: https://oapi.dingtalk.com/robot/send?access_token={token}      # 钉钉机器人webhook地址
  lark: https://open.feishu.cn/open-apis/bot/v2/hook/{token}               # 飞书机器人webhook地址 
  weixin: https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key={key}       # 企业微信机器人webhook地址
  slack: https://hooks.slack.com/services/{id}/{token}                     # slack机器人webhook地址

```
