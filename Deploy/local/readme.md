#harbor本地调试部署环境搭建-使用SSL证书(registry.51yixiao.com)

##Setup1.获取本机IP - 此处IP地址 192.168.1.22
```
ifconfig | grep 'inet'

	inet6 ::1 prefixlen 128
	inet 127.0.0.1 netmask 0xff000000
	inet6 fe80::1%lo0 prefixlen 64 scopeid 0x1
	inet6 fe80::f65c:89ff:fec0:8007%en0 prefixlen 64 scopeid 0x4
	inet 192.168.1.22 netmask 0xffffff00 broadcast 192.168.1.255
	inet6 fe80::242a:8eff:fe21:eef7%awdl0 prefixlen 64 scopeid 0x8
```

##Setup2.更新本机shell环境配置

```
vi ~/.zshrc

#env begin
export MYSQL_HOST=127.0.0.1
export MYSQL_PORT=3306
export MYSQL_USR=root
export MYSQL_PWD=root

export REGISTRY_URL=http://registry:5000
export HARBOR_ADMIN_PASSWORD=abc123
export AUTH_MODE=db_auth
export LDAP_URL=ldaps://ldap.51yixiao.com
export LDAP_BASE_DN=uid=%s,ou=people,dc=mydomain,dc=com
export SELF_REGISTRATION=on
export LOG_LEVEL=debug
export GODEBUG=netdns=cgo
export UI_SECRET=anD1XUKCnYupirmQ
export VERIFY_REMOTE_CERT=on
export PRODUCTION=on

export HARBOR_REG_URL=registry.51yixiao.com
export HARBOR_URL=registry.51yixiao.com
export EXT_ENDPOINT=https://registry.51yixiao.com

#配置UI配置文件路径
export CONFIG_PATH=/works/goProject/src/github.com/vmware/harbor/Deploy/local/config/ui/app.conf

#配置JOB配置文件路径
export JOB_CONFIG_PATH=/works/goProject/src/github.com/vmware/harbor/Deploy/local/config/jobservice/app.conf

export JOB_SERVICE_URL=http://192.168.1.22:8000
export UI_URL=http://192.168.1.22:8080
export TOKEN_URL=http://192.168.1.22:8080

```

##Setup3.更新IP配置

```
vi /etc/hosts
192.168.1.22	registry.51yixiao.com registry ui jobservice
```

修改Nginx环境配置
```
vi ./config/nginx/nginx.conf

upstream ui {
    #server ui:80;
    server 192.168.1.22:8080;
  }
```

修改registry环境配置
```
vi ./config/registry/config.yml
endpoints:
  - name: harbor
    disabled: false
    url: http://192.168.1.22:8080/service/notifications
    timeout: 500ms
    threshold: 5
    backoff: 1s
```

##Setup4.启动基础容器

```
cd local
mkdir -p log
mkdir -p data/registry
mkdir -p data/mysql
mkdir -p data/job_logs

#docker-compose-local.yml 路径记得修改
docker-compose -f docker-compose-local.yml up -d

docker-compose -f docker-compose-local.yml down
```
##Setup5.修改代码配置
```
vi ../../job/config/config.go:83

#configPath := os.Getenv("CONFIG_PATH")
#configPath := os.Getenv("JOB_CONFIG_PATH")

vi ../../service/token/authutils.go:36
#privateKey = "/etc/ui/private_key.pem"
privateKey = "/works/goProject/src/github.com/vmware/harbor/Deploy/local/config/ui/private_key.pem"
```

##Setup6.启动harborUI/Job
```
cd ../../
cp ./ui/* .
bee run harborUI

cd jobservice
bee run jobService
```



