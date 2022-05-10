package logging

import (
	log "github.com/cihub/seelog"
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
		file = "        <rollingfile type=\"size\" filename=\"" + file + "\" maxsize=\"102400\" maxrolls=\"5\"/>\n"
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
		"        <format id=\"colored-default\"  format=\"[%LEV][%Date(2006-01-02 15:04:05.999)] %Msg%n\"/>\n" +
		"        <format id=\"colored-warn\"  format=\"[%LEV][%Date(2006-01-02 15:04:05.999)] %Msg%n\"/>\n" +
		"        <format id=\"colored-error\"  format=\"[%LEV][%Date(2006-01-02 15:04:05.999)] %Msg%n\"/>\n" +
		"    </formats>\n" +
		"</seelog>"
}

//
// Initialize
// @Description:
//
func Initialize() {
	configString := getConfigString("")
	logger, err := log.LoggerFromConfigAsString(configString)
	if err != nil {
		return
	}
	err = log.ReplaceLogger(logger)
}
