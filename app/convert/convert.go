package convert

import (
	"reflect"
	"strings"
	"runtime"
	"github.com/hetian9288/json2dart/app/fields"
	"fmt"
	"os"
	"github.com/sirupsen/logrus"
	"encoding/json"
)

type Convert struct {
	ModelName   string
	Data        interface{}
	Path        string
	Fields      []fields.Fields
	FieldsSort  []string
	PackagePath string
	Imports     []string
}

func NewConvert(path string, ipackage string, name string, data interface{}) {
	iConvert := &Convert{
		ModelName:   name,
		Data:        data,
		PackagePath: ipackage,
		Path:        path,
	}
	switch {
	case strings.HasPrefix(reflect.TypeOf(data).String(), "map[string]"):
		iConvert.FromMap()
	case strings.HasPrefix(reflect.TypeOf(data).String(), "[]interface"):
		iConvert.FromArr()
	}

	iConvert.WriteModelContent()
	iConvert.WriteModelPartContent()
}

func (this *Convert) FromMap() {
	for name, value := range this.Data.(map[string]interface{}) {
		field := this.fieldConvert(name, value)
		this.Fields = append(this.Fields, field)
	}
}

func (this *Convert) fieldConvert(name string, value interface{}) fields.Fields {
	if value == nil {
		return fields.NewFields(name, "Null", false)
	}
	valueType := reflect.TypeOf(value).String()
	fmt.Println()
	switch {
	case valueType == "string":
		if fields.FieldIsDateTime(value.(string)) {
			return fields.NewFields(name, "DateTime", false)
		}
		return fields.NewFields(name, "String", false)

	case valueType == "bool":
		return fields.NewFields(name, "bool", false)

	case valueType == "json.Number":
		if strings.Contains(value.(json.Number).String(), ".") {
			return fields.NewFields(name, "double", false)
		}else{
			return fields.NewFields(name, "int", false)
		}

	case strings.HasPrefix(valueType, "int"):
		return fields.NewFields(name, "int", false)

	case strings.HasPrefix(valueType, "float"):
		return fields.NewFields(name, "double", false)

	case strings.HasPrefix(valueType, "map[string]"):
		this.Imports = append(this.Imports, fields.NameToFieldName(name))
		NewConvert(this.Path, this.PackagePath, fields.NameToFieldName(name), value)
		return fields.NewFields(name, strings.Title(name), true)

	case strings.HasPrefix(valueType, "[]interface"):
		one := value.([]interface{})[0]
		oneField := this.fieldConvert(name, one)
		return fields.NewFields(name, fmt.Sprintf("List<%s>", oneField.ValueType), false)
	}
	return fields.NewFields(name, "Null", false)
}

func (this *Convert) FromArr() {
	one := this.Data.([]interface{})[0]
	oneField := this.fieldConvert(this.ModelName, one)
	this.Fields = append(this.Fields, fields.NewFields(this.ModelName, fmt.Sprintf("List<%s>", oneField.ValueType), false))
}

func (this *Convert) GetFileName() string {
	delimiter := "/"
	if runtime.GOOS == "windows" {
		delimiter = "\\"
	}
	if strings.HasSuffix(this.Path, delimiter) {
		return this.Path + this.ModelName + ".dart"
	}
	return this.Path + delimiter + this.ModelName + ".dart"
}

func (this *Convert) GetFilePartName() string {
	delimiter := "/"
	if runtime.GOOS == "windows" {
		delimiter = "\\"
	}
	if strings.HasSuffix(this.Path, delimiter) {
		return this.Path + this.ModelName + ".g.dart"
	}
	return this.Path + delimiter + this.ModelName + ".g.dart"
}

func (this *Convert) Write(path string, content string) {
	var err error
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		logrus.Error("创建文件失败", err)
		return
	}
	_, err = f.WriteString(content)
	if err != nil {
		logrus.Error("写入文件失败", err)
		return
	}

	logrus.Info(path, "; 文件保存成功")
}

func (this *Convert) WriteModelContent() {
	modelName := fields.TypeToType(this.ModelName)
	fileContent := fmt.Sprintf(`import 'package:json_annotation/json_annotation.dart';
%s

part '%s.g.dart';

@JsonSerializable()
class %s {
  %s
  %s({%s});
  factory %s.fromJson(Map<String, dynamic> json) => _$%sFromJson(json);
  Map<String, dynamic> toJson() => _$%sToJson(this);
}
`, this.getImports(), this.ModelName, modelName, this.getFieldLines(), modelName, this.getFieldInitLines(), modelName, modelName, modelName)
	this.Write(this.GetFileName(), fileContent)
}

func (this *Convert) WriteModelPartContent() {
	modelName := fields.TypeToType(this.ModelName)
	fileContent := fmt.Sprintf(`part of '%s.dart';

%s _$%sFromJson(Map<String, dynamic> json) {
  return new %s(
      %s);
}

Map<String, dynamic> _$%sToJson(%s instance) => <String, dynamic>{
      %s
    };`, this.ModelName, modelName, modelName, modelName, this.getFieldJsonToValStr(), modelName, modelName, this.getFieldToDataStr())
	this.Write(this.GetFilePartName(), fileContent)
}

func (this *Convert) getFieldLines() (line string) {
	for _, item := range this.Fields {
		line += fmt.Sprintf("final %s %s;\n", item.ValueType, item.FieldName)
	}
	return
}

func (this *Convert) getFieldInitLines() (string) {
	var lines []string
	for _, item := range this.Fields {
		lines = append(lines, "this."+item.FieldName)
	}
	return strings.Join(lines, ", ")
}

func (this *Convert) getFieldJsonToValStr() (line string) {
	for _, item := range this.Fields {
		line += fmt.Sprintf("%s: %s,\n", item.FieldName, item.GetJsonToValStr())
	}
	return
}

func (this *Convert) getFieldToDataStr() (string) {
	var lines []string
	for _, item := range this.Fields {
		lines = append(lines, fmt.Sprintf("'%s': %s", item.FieldName, item.GetToDataStr()))
	}
	return strings.Join(lines, ", \n")
}

func (this *Convert) getImports() string {
	var lines []string
	for _, item := range this.Imports {
		lines = append(lines, fmt.Sprintf("import './%s.dart';", item))
	}
	return strings.Join(lines, "; \n")
}