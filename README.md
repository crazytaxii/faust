# faust ![CI Status](https://github.com/crazytaxii/faust/actions/workflows/ci.yaml/badge.svg)

faust 是一款将本地图片上传至*七牛云对象存储*的小工具。目前支持：

- jpg
- png
- webp
- gif
- [avif](https://aomediacodec.github.io/av1-avif/)
- svg

## 编译

**请先配置好 Go 开发环境！**

```bash
$ make build
```

## 安装

```bash
$ make install
```

## 配置

请先注册七牛云账号，并获取账号对应的 Access Key 和 Secret Key：

1. 点击[注册](https://portal.qiniu.com/signup?ref=developer.qiniu.com)开通七牛开发者帐号
2. 如果已有账号，直接登录七牛开发者后台，点击[这里](https://portal.qiniu.com/user/key)查看 Access Key 和 Secret Key

- Access Key & Secret Key
- Bucket
- Base URL (已绑定存储空间的融合 CDN 加速域名，比如 <https://pic.crazytaxii.com>)

**[域名接入七牛云存储](https://developer.qiniu.com/fusion/manual/4939/the-domain-name-to-access)**

```bash
$ cat <<EOF > ~/.faust/config.yaml
accessKey ${your_access_key} \
secretKey ${your_secret_key} \
bucket ${bucket_name} \
baseURL ${your_base_url}
EOF
```

> 配置文件 config.yaml 默认放置于 ~/.faust 路径

## 使用

1. 上传图片

    ```bash
    $ faust --image ./test/Go-Logo_Fuchsia.jpg upload
    ```

1. 上传证书（私钥 + 证书链）

    ```bash
    $ faust --key /path/to/private_key --cert /path/to/full_chan_certificates upload
    ```
