package configs

import (
	"fmt"
	"github.com/kataras/iris"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	Root, _  = os.Getwd()
	YamlConf = iris.YAML(filepath.Join(Root, "config", "config.yml")) // 加载配置文件
)

func GetConfigString(name string) string {

	if strings.Contains(name, ".") {
		if v, err := getValue(name, YamlConf.Other); err != nil {
			_, _ = fmt.Fprintf(os.Stdin, err.Error())
			return ""
		} else {
			return v
		}

		//return GetValue(name, Isc.Other)
	}

	if val, ok := YamlConf.Other[name]; ok {
		return fmt.Sprintf("%v", val)
	}

	return ""
}

func GetConfigInt(name string) int {

	if strings.Contains(name, ".") {
		if v, err := getValue(name, YamlConf.Other); err != nil {
			_, _ = fmt.Fprintf(os.Stdin, err.Error())
			return 0
		} else {
			if i, err := strconv.Atoi(v); err != nil {
				return i
			}
			return 0
		}

		//return GetValue(name, Isc.Other)
	}

	if val, ok := YamlConf.Other[name]; ok {
		if i, err := strconv.Atoi(fmt.Sprintf("%v", val)); err != nil {
			return i
		}
		//return int(fmt.Sprintf("%v", val))
	}

	return 0
}

func getValue(name string, conf map[string]interface{}) (string, error) {
	if strings.Contains(name, ".") {
		names := strings.Split(name, ".")
		firstName := names[0]
		subName := strings.Join(names[1:], ".")
		if val, ok := conf[firstName]; ok {
			//fmt.Printf("%s", val)
			//do something here
			switch v := val.(type) {
			case map[string]interface{}:
				return getValue(subName, v)
			case interface{}:
				value := conf[subName]
				return fmt.Sprintf("%v", value), nil
			}

		} else {
			return "", fmt.Errorf("没有找到配置项")
		}
		//GetValue(subName)
	}

	return fmt.Sprintf("%v", conf[name]), nil
}
