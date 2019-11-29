# Faust

Faust 是一款 Go 编写的轻量化工具，使用它将本地图片转换成 JPG 并上传至七牛云空间。

## 下载

**需要事先配置好 Go (1.11+) 开发环境！**

```bash
$ go get github.com/crazytaxii/faust
```

## 编译

```bash
$ make build
```

## 安装

```bash
$ make install
```

## 配置

七牛云所有的功能，都需要合法的授权。授权凭证的签算需要七牛账号下的一对有效的 Access Key 和 Secret Key，这对密钥可以通过如下步骤获得：

1. 点击[注册](https://portal.qiniu.com/signup?ref=developer.qiniu.com)开通七牛开发者帐号
2. 如果已有账号，直接登录七牛开发者后台，点击[这里](https://portal.qiniu.com/user/key)查看 Access Key 和 Secret Key

- Access Key
- Secret Key
- Bucket
- Base URL (已绑定存储空间的融合 CDN 加速域名，比如 https://pic.crazytaxii.com)

**[域名接入七牛云存储](https://developer.qiniu.com/fusion/manual/4939/the-domain-name-to-access)**

```bash
$ ./faust \
  --access_key your_access_key \
  --secret_key your_secret_key \
  --bucket bucket_name \
  --base_url your_base_url \
  config
```

> 配置文件 config.yaml 默认生成在 ~/.faust 路径。

## 使用

```bash
$ ./faust --image ./test/Go-Logo_Fuchsia.jpg upload
bucket: markdown
key: 19-02-17/94939921.jpg
file size: 71447
hash: FhpxfGzt6T241vme6_7j1CUEYw0k
public access url: https://pic.crazytaxii.com/19-02-17/94939921.jpg
```
