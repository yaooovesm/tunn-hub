package logging

import (
	log "github.com/cihub/seelog"
	"os"
	"time"
	"tunn-hub/version"
)

//
// getConfigString
// @Description:
// @param logPah
// @return string
//
func getConfigString(logPah string) string {
	file := ""
	if logPah != "" {
		file = "        <rollingfile type=\"size\" filename=\"" + logPah + "\" maxsize=\"102400\" maxrolls=\"5\"/>\n"
	}
	debug := ""
	if version.Develop {
		debug = ",debug"
	}
	return "" +
		"<seelog type=\"sync\">\n" +
		"    <outputs formatid=\"main\">\n" +
		"        <filter levels=\"trace,info" + debug + "\">\n" +
		"            <console formatid=\"colored-default\"/>\n" +
		"        </filter>\n" +
		"        <filter levels=\"warn\">\n" +
		"            <console formatid=\"colored-warn\"/>\n" +
		"        </filter>\n" +
		"        <filter levels=\"error,critical\">\n" +
		"            <console formatid=\"colored-error\"/>\n" +
		"        </filter>\n" +
		file +
		"    </outputs>\n" +
		"    <formats>\n" +
		"        <format id=\"main\" format=\"[%LEV][%Date(2006-01-02 15:04:05.999)][%File.%Line] %Msg%n\"/>\n" +
		"        <format id=\"colored-default\"  format=\"%EscM(34)[%LEV][%Date(2006-01-02 15:04:05.999)] %Msg%n%EscM(0)\"/>\n" +
		"        <format id=\"colored-warn\"  format=\"%EscM(33)[%LEV][%Date(2006-01-02 15:04:05.999)] %Msg%n%EscM(0)\"/>\n" +
		"        <format id=\"colored-error\"  format=\"%EscM(31)[%LEV][%Date(2006-01-02 15:04:05.999)] %Msg%n%EscM(0)\"/>\n" +
		"    </formats>\n" +
		"</seelog>"
}

//
// Initialize
// @Description:
//
func Initialize() {
	storagePath := "./log/"
	s, err := os.Stat(storagePath)
	if err != nil {
		_ = log.Warn("set log file error : ", err)
	}
	if os.IsNotExist(err) || !s.IsDir() {
		_ = os.MkdirAll(storagePath, 0644)
	}
	storageFile := storagePath + time.Now().Format("2006_01_02") + ".log"
	configString := getConfigString(storageFile)
	logger, err := log.LoggerFromConfigAsString(configString)
	if err != nil {
		return
	}
	err = log.ReplaceLogger(logger)
}
