package fields

import (
	"regexp"
	"strings"
	"fmt"
)

type Fields struct {
	Name      string
	ValueType string
	FieldName string
	IsAuto    bool
}

func (this Fields) GetToDataStr() string {
	switch this.ValueType {
	case "String":
		return fmt.Sprintf("instance.%s", this.FieldName)
	case "bool":
		return fmt.Sprintf("instance.%s", this.FieldName)
	case "int":
		return fmt.Sprintf("instance.%s", this.FieldName)
	case "double":
		return fmt.Sprintf("instance.%s", this.FieldName)
	case "List<String>":
		return fmt.Sprintf("instance.%s", this.FieldName)
	case "List<bool>":
		return fmt.Sprintf("instance.%s", this.FieldName)
	case "List<int>":
		return fmt.Sprintf("instance.%s", this.FieldName)
	case "List<double>":
		return fmt.Sprintf("instance.%s", this.FieldName)
	case "DateTime":
		return fmt.Sprintf("instance.%s.toIso8601String()", this.FieldName)
	case "Null":
		return fmt.Sprintf("instance.%s", this.FieldName)
	default:
		if strings.HasPrefix(this.ValueType, "List") {
			return fmt.Sprintf(`instance.%s != null ? instance.%s.map((v) => v.toJson()).toList() : null`, this.FieldName, this.FieldName)
		}else if this.IsAuto{
			return fmt.Sprintf(`instance.%s != null ? instance.%s.toJson() : null`, this.FieldName, this.FieldName)
		}
		return ""
	}
}

func (this Fields) GetJsonToValStr() string {
	switch this.ValueType {
	case "String":
		return fmt.Sprintf("json['%s'] as String", this.Name)
	case "bool":
		return fmt.Sprintf("json['%s'] as bool", this.Name)
	case "int":
		return fmt.Sprintf("json['%s'] as int", this.Name)
	case "double":
		return fmt.Sprintf("json['%s'] as double", this.Name)
	case "List<String>":
		return fmt.Sprintf("json['%s'] as List<String>", this.Name)
	case "List<bool>":
		return fmt.Sprintf("json['%s'] as List<bool>", this.Name)
	case "List<int>":
		return fmt.Sprintf("json['%s'] as List<int>", this.Name)
	case "List<double>":
		return fmt.Sprintf("json['%s'] as List<double>", this.Name)
	case "DateTime":
		return fmt.Sprintf("DateTime.parse(json['%s'] as String)", this.Name)
	case "Null":
		return fmt.Sprintf("json['%s']", this.Name)
	default:
		if strings.HasPrefix(this.ValueType, "List") {
			return fmt.Sprintf(`(){
  List<%s>  %s = new List<%s>();
  if (json['%s'] != null) {
      json['%s'].forEach((v) {
        %s.add(new %s.fromJson(v));
      });
    }
      return userlist;
}()`, TypeToType(this.Name), this.Name, TypeToType(this.Name), this.Name, this.Name, this.Name, TypeToType(this.Name))
		}else if this.IsAuto{
			return fmt.Sprintf(`json['%s'] != null ? new %s.fromJson(json['%s']) : null`, this.Name, TypeToType(this.Name), this.Name)
		}
		return ""
	}
}

func NewFields(name string, valueType string, isauto bool) Fields {
	return Fields{
		Name: name, ValueType: TypeToType(valueType), IsAuto: isauto, FieldName: NameToFieldName(name),
	}
}

func TypeToType(ty string) (r string) {
	if ty == "int" || ty == "bool" || ty == "double" {
		return ty
	}
	tyArr := strings.Split(ty, "_")
	for _, item := range tyArr {
		r += strings.Title(item)
	}
	return
}

func NameToFieldName(n string) (name string) {
	tyArr := strings.Split(n, "_")
	for i, item := range tyArr {
		if i == 0 {
			name += strings.ToLower(item)
		}else {
			name += strings.Title(item)
		}
	}
	return
}

func FieldIsDateTime(date string) bool {
	dateReg, _ := regexp.Compile(`^[1-9]\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1])$`)
	dateTimeReg, _ := regexp.Compile(`^[1-9]\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1])\s+(20|21|22|23|[0-1]\d):[0-5]\d:[0-5]\d$`)

	if dateTimeReg.MatchString(date) {
		return true
	} else if dateReg.MatchString(date) {
		return true
	}
	return false
}
