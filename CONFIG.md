## Configuration file description

RevSuit will generate a profile template the first time it is run, and you can edit its contents to implement some customization requirements.
```yaml
version: 1.5
addr: :10000                          # Address of the HTTP service will listen
token:                                # Authentication Token, both the admin page and the client need to be authenticated by this Token
domains: []                           # The domain names used by the platform
external_ip:                          # The external IP of the platform, you need to make sure that the target you want to test can access the platform through this IP
admin_path_prefix: "/revsuit"          # The http path prefix for the admin page, the page will be located at: /admin_path_prefix/admin
# Database connection information, support using MySQL, Postgres, Sqlite3
database: "mysql://root:password@tcp(127.0.0.1:3306)/revsuit?charset=utf8mb4&parseTime=True&loc=Local"
# database: "postgres://host=127.0.0.1 user=root password=password dbname=revsuit port=5432 sslmode=disable TimeZone=Asia/Shanghai"
#database: revsuit.db # The use of sqlite databases on some systems may cause `SQLITE BUSY` problems in concurrent scenarios, so it is not recommended for use in a formal environment

log_level: info                       # Output log levels, divided into: debug, info, warning, error, fatal
check_upgrade: true                   # Whether to automatically check for updates
  
ip_location_database:                 # IP database related configuration
  database: "qqwry"                   # qqwry or geoip.
  geo_license_key: ""                 # Mandatory field, if you choose to use GeoIP.
  
http:                                 # HTTP receive connection related configuration
  ip_header:                          # Specify the source IP via the HTTP Header, which is usually required when there is a predecessor reverse proxy
dns:                                  # DNS receive connection related configuration
  enable: true                             
  addr: :53                           # Address of the DNS service will listen  
rmi:                                  # RMI receive connection related configuration
  enable: true                                
  addr: :1099                         # Address of the RMI service will listen  
ldap:
  enable: true
  addr: :1389                         # Address of the Ldap service will listen 
mysql:                                # MySQL receive connection related configuration
  enable: true                        
  addr: :3306                         # Address of the MySQL service will listen         
  version_string: 10.4.13-MariaDB-log # The version of the MySQL service in disguise
ftp:                                  # FTP receive connection related configuration
  enable: true                        
  addr: :21                           # Address of the FTP service will listen 
  pasv_port: 2020                     # FTP passive mode port listening address
notice:
  dingtalk: https://oapi.dingtalk.com/robot/send?access_token={token}      # Webhook of DingTalk Bot
  lark: https://open.feishu.cn/open-apis/bot/v2/hook/{token}               # Webhook of Lark Bot 
  weixin: https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key={key}       # Webhook of Weixin Bot
  slack: https://hooks.slack.com/services/{id}/{token}                     # Webhook of Slack Bot
```
