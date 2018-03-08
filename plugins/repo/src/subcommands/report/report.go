package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/dokku/dokku/plugins/common"
)

func reportSingleApp(appName, infoFlag string) {
	infoFlags := map[string]string{
		"--repo-container-copy-folder": common.PropertyGet("repo", appName, "container-copy-folder"),
		"--repo-host-copy-folder":      common.PropertyGet("repo", appName, "host-copy-folder"),
	}

	if len(infoFlag) == 0 {
		common.LogInfo2Quiet(fmt.Sprintf("%s repo information", appName))
		for k, v := range infoFlags {
			key := common.UcFirst(strings.Replace(strings.TrimPrefix(k, "--"), "-", " ", -1))
			common.LogVerbose(fmt.Sprintf("%s: %s", key, v))
		}
		return
	}

	for k, v := range infoFlags {
		if infoFlag == k {
			fmt.Fprintln(os.Stdout, v)
			return
		}
	}

	keys := reflect.ValueOf(infoFlags).MapKeys()
	strkeys := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		strkeys[i] = keys[i].String()
	}
	common.LogFail(fmt.Sprintf("Invalid flag passed, valid flags: %s", strings.Join(strkeys, ", ")))
}

// displays a repo report for one or more apps
func main() {
	flag.Parse()
	appName := flag.Arg(1)
	infoFlag := flag.Arg(2)

	if strings.HasPrefix(appName, "--") {
		infoFlag = appName
		appName = ""
	}

	if len(appName) == 0 {
		apps, err := common.DokkuApps()
		if err != nil {
			return
		}
		for _, appName := range apps {
			reportSingleApp(appName, infoFlag)
		}
		return
	}

	reportSingleApp(appName, infoFlag)
}
