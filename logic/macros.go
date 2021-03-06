package logic

/**
宏管理
*/

import (
	"fmt"
	"strings"

	"github.com/chentaihan/NginxParse/util"
)

type Macros struct {
	Map map[string]*Macro
}

var macros *Macros = nil

func GetMacros() *Macros {
	if macros == nil {
		macros = &Macros{
			Map: make(map[string]*Macro, 1024),
		}
	}
	return macros
}

func (macros *Macros) Add(key string, mf *Macro) {
	macros.Map[key] = mf
}

//是否存在指定的宏
func (macros *Macros) Exist(macroName string) bool {
	if _, ok := macros.Map[macroName]; ok {
		return true
	}
	return false
}

//是否存在指定的宏
func (macros *Macros) Get(macroName string) *Macro {
	if m, ok := macros.Map[macroName]; ok {
		return m
	}
	return nil
}

func (macros *Macros) GetMacroValue(macroName string) string {
	index := strings.Index(macroName, "(")
	key := macroName
	if index > 0 {
		key = macroName[0:index]
	}

	macroInfo := GetMacros().Get(key)
	if macroInfo == nil {
		return ""
	}

	value := macroInfo.Value

	if index > 0 {
		actualName := macroName[index:]
		actualParams := macros.getMacroParams(actualName) //实参
		if len(actualParams) == 0 {
			util.Println("错误：宏 ", macroName, " 解析有误")
			return ""
		}
		formalParams := macros.getMacroParams(macroInfo.Name) //形参
		return macros.replaceParams(value, formalParams, actualParams)
	}

	return value
}

//解析宏参数
//将("secure_link",1)解析成[]string{"secure_link",1}
func (macros *Macros) getMacroParams(actualName string) []string {
	if actualName != "" {
		actualName = actualName[1 : len(actualName)-1]
		return strings.Split(actualName, ",")
	}
	return []string{}
}

//宏替换，实参代替形参
func (macros *Macros) replaceParams(value string, formalParams, actualParams []string) string {
	if value == "" {
		return value
	}
	minLen := len(formalParams)
	if minLen > len(actualParams) {
		minLen = len(actualParams)
	}
	outBuf := util.NewBufferWriter(len(value))
	for i := 0; i < minLen; i++ {
		for len(value) > 0 {
			index := strings.Index(value, formalParams[i])
			isMatch := true
			if index >= 0 {
				end := index+len(formalParams[i])
				if index >= 1 && util.IsLegalChar(value[index-1]) {
					isMatch = false
				} else if end < len(value) && util.IsLegalChar(value[end]) {
					isMatch = false
				}
				if isMatch {
					outBuf.WriteString(value[:index])
					value = value[index:]
					value = strings.Replace(value, formalParams[i], actualParams[i], 1)
				}else{
					outBuf.WriteString(value[:end])
					value = value[end:]
				}
			} else {
				outBuf.WriteString(value)
				break
			}
		}
	}
	return outBuf.ToString()
}

func (macros *Macros) Print() {
	fmt.Println()
	fmt.Println("----------------macro----------------")
	for key, macro := range macros.Map {
		fmt.Println("#define ", key, macro.Name, macro.Value)
	}
}
