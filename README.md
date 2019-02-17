# Faust

Faust 是一款 Go 编写的轻量化工具，使用它将本地图片转换成 JPG 并上传至七牛云空间。

## 下载

**需要事先安装 Go 和配置开发环境**

```bash
$ go get github.com/crazytaxii/faust
```

## 编译

```bash
$ go build
```

## 配置

* `-a access_key` 用户七牛账号的 Access Key
* `-s secret key` 用户七牛账号的 Secret Key
* `-b bucket name` 存放图片的存储空间
* `-d your domain` 已绑定存储空间的融合 CDN 加速域名

### 使用

```bash
$ ./faust -i ./test/Go-Logo_Fuchsia.jpg
bucket: markdown
key: 19-02-17/78885642.jpg
file size: 71447
hash: FhpxfGzt6T241vme6_7j1CUEYw0k
public access url: /19-02-17/78885642.jpg
```
