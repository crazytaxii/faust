package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	help      bool
	srcFile   string
	accessKey string
	secretKey string
	bucket    string
	version   bool
	domain    string
)

type Config struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Bucket    string `json:"bucket"`
	Domain    string `json:"domain"`
}

var conf *Config

const (
	configFile = "./config.json"
)

func init() {
	flag.BoolVar(&help, "h", false, "show help")
	flag.StringVar(&srcFile, "i", "", "image path")
	flag.StringVar(&accessKey, "a", "", "access key")
	flag.StringVar(&secretKey, "s", "", "secret key")
	flag.StringVar(&bucket, "b", "", "bucket name")
	flag.BoolVar(&version, "v", false, "show version")
	flag.StringVar(&domain, "d", "", "your domain")

	conf = &Config{}
	err := loadConfig(configFile)
	if err != nil {
		fmt.Println("err:", err)
		os.Exit(1)
	}
}

func main() {
	flag.Parse()
	if help {
		flag.Usage()
		return
	}
	if version {
		fmt.Println("v0.9")
		return
	}
	defer func() {
		err := saveConfig(configFile)
		if err != nil {
			fmt.Println("err:", err.Error())
		}
	}()
	if len(accessKey) != 0 {
		conf.AccessKey = accessKey
		return
	} else if len(conf.AccessKey) == 0 {
		fmt.Println("input access key of your qiniu account")
		return
	}
	if len(secretKey) != 0 {
		conf.SecretKey = secretKey
		return
	} else if len(conf.SecretKey) == 0 {
		fmt.Println("input secret key of your qiniu account")
		return
	}
	if len(bucket) != 0 {
		conf.Bucket = bucket
		return
	} else if len(conf.Bucket) == 0 {
		fmt.Println("input bucket name of storage")
		return
	}
	if len(domain) != 0 {
		conf.Domain = domain
		return
	} else if len(conf.Domain) == 0 {
		fmt.Println("input your domain for making public URL")
		return
	}
	if len(srcFile) == 0 {
		fmt.Println("invalid src image path")
		return
	}

	err := upload(srcFile, conf.AccessKey, conf.SecretKey, conf.Bucket)
	if err != nil {
		fmt.Println("err:", err.Error())
	}
	return
}
