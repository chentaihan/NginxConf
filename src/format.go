package main

import (
	"html/template"
	"io/ioutil"
)

const (
	COMMAND_HEADER = "<tr><th colspan=\"6\">{{.ModuleName}}</th></tr><tr><td colspan=\"3\">文件名：{{.FileName}}</td><td colspan=\"3\">struct：{{.StructName}}</td></tr>"
	COMMAND_TH     = "<tr><td width=\"16%\">name</td><td width=\"24%\">type</td><td width=\"17%\">set</td><td width=\"17%\">conf</td><td width=\"13%\">offset</td><td width=\"13%\">post</td></tr>"
	COMMAND_TD     = "<tr><td>{{.Name}}</td><td>{{.Type}}</td><td>{{.Set}}</td><td>{{.Conf}}</td><td>{{.Offset}}</td><td>{{.Post}}</td></tr>"
	COMMAND_HTML   = "command.html"
)

const (
	MODULE_HEADER = "<tr><th colspan=\"6\">{{.ModuleName}}</th></tr><tr><td colspan=\"3\">文件名：{{.FileName}}</td><td colspan=\"3\">struct：{{.StructName}}</td></tr>"
	MODULE_TH     = "<tr><td width=\"33%\" colspan=\"2\">Context</td><td width=\"34%\" colspan=\"2\">Command</td><td width=\"33%\" colspan=\"2\">Type</td></tr>"
	MODULE_TD     = "<tr><td colspan=\"2\">{{.Context}}</td><td colspan=\"2\">{{.Command}}</td><td colspan=\"2\">{{.Type}}</td></tr>"
	MODULE_HTML   = "module.html"
)

type formatStruct func(sct interface{}) []byte

func formatCommandList(list interface{}) []byte {
	buf := make([]byte, 0, 1024)
	cmdList := list.([]*CommandInfo)

	for _, cmdInfo := range cmdList {
		list := make([]interface{}, 0, len(cmdInfo.CmdList))
		for _, item := range cmdInfo.CmdList {
			list = append(list, item)
		}
		buf = append(buf, formatTable(COMMAND_HEADER, COMMAND_TH, COMMAND_TD, cmdInfo, list...)...)
	}
	return buf
}

func formatModuleList(list interface{}) []byte {
	buf := make([]byte, 0, 1024)
	moduleList := list.([]*ModuleInfo)

	for _, moduleInfo := range moduleList {
		list := make([]interface{}, 0, len(moduleInfo.ModuleList))
		for _, item := range moduleInfo.ModuleList {
			list = append(list, item)
		}
		buf = append(buf, formatTable(MODULE_HEADER, MODULE_TH, MODULE_TD, moduleInfo, list...)...)
	}
	return buf
}

func formatTable(headerFormat, thFormat, tdFormat string, structInfo interface{}, list ...interface{}) []byte {
	writer := NewBufferWriter(0)
	writer.Write([]byte("<table>"))
	tmpl := template.New("tmpl1")
	tmpl.Parse(headerFormat)
	tmpl.Execute(writer, structInfo)
	writer.Write([]byte(thFormat))

	for _, cmd := range list {
		tmpl = template.New("tmpl1")
		tmpl.Parse(tdFormat)
		tmpl.Execute(writer, cmd)
	}
	writer.Write([]byte("</table><br/><br/>"))

	return writer.buffer
}

func outPutCommand(cmdList []*CommandInfo) {
	outputFile(FILE_CONFIG_FORMAT, COMMAND_HTML, formatCommandList, cmdList)
}

func outPutModule(moduleList []*ModuleInfo) {
	outputFile(FILE_CONFIG_FORMAT, MODULE_HTML, formatModuleList, moduleList)
}

func outputFile(formatFile, outputFile string, formatFunc formatStruct, list interface{}) {
	configFormat := getConfigFile(formatFile)
	content, err := ioutil.ReadFile(configFormat)
	if err != nil {
		return
	}
	var index int = -1
	for i := 0; i < len(content); i++ {
		if content[i] == '$' && content[i+1] == '$' {
			index = i
		}
	}

	outputFile = getOutPutFile(outputFile)
	ioutil.WriteFile(outputFile, content[:index], 0666)

	buf := formatFunc(list)
	WriteFileAppend(outputFile, buf, 0666)

	index += 2
	WriteFileAppend(outputFile, content[index:], 0666)
}