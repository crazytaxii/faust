# Faust

Faust 是一款 Go 编写的轻量化工具，使用它将本地图片转换成 JPG 并上传至七牛云空间。

## 下载

**需要事先安装 Go 和配置开发环境！**

```bash
$ go get github.com/crazytaxii/faust
```

## 编译

```bash
$ go build
```

## 配置

七牛云所有的功能，都需要合法的授权。授权凭证的签算需要七牛账号下的一对有效的 Access Key 和 Secret Key，这对密钥可以通过如下步骤获得：

1. 点击[注册](https://portal.qiniu.com/signup?ref=developer.qiniu.com)开通七牛开发者帐号
2. 如果已有账号，直接登录七牛开发者后台，点击[这里](https://portal.qiniu.com/user/key)查看 Access Key 和 Secret Key

### 添加 Access Key

设置七牛账号的 Access Key

```bash
$ ./faust -a access_key
```

### 添加 Secret Key

设置七牛账号的 Secret Key

```bash
$ ./faust -a secret_key
```

### 添加 Bucket

Bucket 是存放图片的存储空间

```bash
$ ./faust -b bucket_name
```

### 添加域名（[域名接入七牛云存储](https://developer.qiniu.com/fusion/manual/4939/the-domain-name-to-access)）

已绑定存储空间的融合 CDN 加速域名

```bash
$ ./faust -d your_domain
```

## 使用

```bash
$ ./faust -i ./test/Go-Logo_Fuchsia.jpg
bucket: markdown
key: 19-02-17/94939921.jpg
file size: 71447
hash: FhpxfGzt6T241vme6_7j1CUEYw0k
public access url: pic.crazytaxii.com/19-02-17/94939921.jpg
```
