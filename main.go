package main

import (
	"flag"
	"strings"
	"github.com/sirupsen/logrus"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/hetian9288/json2dart/app/convert"
	"os"
	"fmt"
)

var (
	ModelName   = flag.String("name", "Users", "-name=模块名 如Users")
	ApiUrl      = flag.String("api", "https://h5.cnganen.cn/test.json", "-api=你的API地址")
	Path        = flag.String("page", "./dart", "-path=文件保存位置")
	packagePath = flag.String("package", "./", "-package=引用包前缀")
)

func main() {
	flag.Parse()

	isRun := true
	flag.VisitAll(func(field *flag.Flag) {
		if strings.EqualFold(field.Value.String(), "") {
			isRun = false
			logrus.Error(field.Name, "必须设置, ", field.Usage)
		}
	})

	if !isRun {
		return
	}

	err := os.MkdirAll(*Path, 0777)
	if err != nil {
		logrus.Error("保存路径创建失败", err)
		os.Exit(101)
	}

	resp, err := http.Get(*ApiUrl)
	if err != nil {
		logrus.Error("发起请求失败, ", err)
		return
	}
	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("读取接口失败, ", err)
		return
	}
	var jsonMap map[string]interface{}
	err = json.Unmarshal(bodyByte, &jsonMap)
	if err != nil {
		logrus.Error("接口返回可能非正确的JSON格式, ", err)
		return
	}
	convert.NewConvert(*Path, *packagePath, *ModelName, jsonMap)
	fmt.Println("搞定收工")
}
