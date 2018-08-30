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
	"bytes"
)

var (
	ModelName   = flag.String("name", "users", "-name=模块名 如Users")
	ApiUrl      = flag.String("api", "http://127.0.0.1:3000/api2.0/app.users.base", "-api=你的API地址")
	Path        = flag.String("page", "./dart", "-path=文件保存位置")
	packagePath = flag.String("package", "./", "-package=引用包前缀")
	defDataName = flag.String("mapname", "data", "-mapname=默认的数据键名")
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
	var jsonMap interface{}
	//err = json.Unmarshal(bodyByte, &jsonMap)
	iJson := json.NewDecoder(bytes.NewBuffer(bodyByte))
	iJson.UseNumber()
	err = iJson.Decode(&jsonMap)

	if err != nil {
		logrus.Error("接口返回可能非正确的JSON格式, ", err)
		return
	}
	var defData interface{};
	if (*defDataName != "") {
		defData = jsonMap.(map[string]interface{})[*defDataName]
	}else{
		defData = jsonMap;
	}
	convert.NewConvert(*Path, *packagePath, *ModelName, defData)
	fmt.Println("搞定收工")
}
