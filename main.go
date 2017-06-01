package main

import (
	"log"
	"os"

	"os/exec"
	"path/filepath"

	"fmt"

	"github.com/cihub/seelog"
	"github.com/spf13/cobra"
	"github.com/sryanyuan/gocodesite/gocodecc"
)

// Default log config
const (
	defaultLogSetting = `
	<seelog minlevel="info">
    	<outputs formatid="main">
			<rollingfile namemode="postfix" type="date" filename="log/app.log" datepattern="060102" maxrolls="30"/>
       		<console />
    	</outputs>
    	<formats>
        	<format id="main" format="%Date/%Time [%LEV] %Msg (%File:%Line %FuncShort)%n"/>
    	</formats>
	</seelog>
	`
)

var (
	rootCommand = &cobra.Command{
		Use:   "gocodesite",
		Short: "gocodesite provides some commands to run the blog site",
	}
	setupCommand = &cobra.Command{
		Use:   "setup",
		Short: "setup initialize the site",
		Run:   siteSetupFunc,
	}
	runCommand = &cobra.Command{
		Use:   "run",
		Short: "run the site",
		Run:   siteRunFunc,
	}
)

// Run options
var (
	runConfigPath   string
	setupConfigPath string
)

func init() {
	runCommand.PersistentFlags().StringVar(&runConfigPath, "config", "", "site config file")
	setupCommand.PersistentFlags().StringVar(&setupConfigPath, "config", "", "site config file")

	rootCommand.AddCommand(setupCommand)
	rootCommand.AddCommand(runCommand)
}

func getModulePath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	dir := filepath.Dir(path)
	return dir
}

func main() {
	if err := rootCommand.Execute(); nil != err {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}

func siteSetupFunc(cmd *cobra.Command, args []string) {
	var err error

	// Load config
	if setupConfigPath == "" {
		cmd.Println("config mustn't be empty")
		cmd.Usage()
		return
	}
	config, err := gocodecc.ReadTOMLConfig(setupConfigPath)
	if nil != err {
		cmd.Println(err)
		return
	}

	site := gocodecc.NewSite(config)
	if err = site.Setup(true); nil != err {
		cmd.Println(err)
		return
	}

	// Create default user
	var password string
	if password, err = site.NewAdmin(); nil != err {
		cmd.Println(err)
		return
	}
	cmd.Println("Setup ok, account:admin password:", password)
}

func siteRunFunc(cmd *cobra.Command, args []string) {
	var err error

	// Load config
	if runConfigPath == "" {
		cmd.Println("config mustn't be empty")
		cmd.Usage()
		return
	}
	config, err := gocodecc.ReadTOMLConfig(runConfigPath)
	if nil != err {
		cmd.Println(err)
		return
	}

	// Clean up
	defer func() {
		e := recover()
		if nil != e {
			seelog.Error("Main routine quit with error:", e)
		} else {
			seelog.Info("Main routine quit normally")
		}

		seelog.Flush()
	}()

	// Initialize log
	logger, err := seelog.LoggerFromConfigAsFile(getModulePath() + "/static/conf/log.conf")
	if nil == logger {
		// Using default log config
		logger, err = seelog.LoggerFromConfigAsString(defaultLogSetting)
		if nil != err {
			log.Println("Failed to initialize log module, error:", err)
			os.Exit(1)
		}
	}
	seelog.ReplaceLogger(logger)

	site := gocodecc.NewSite(config)
	if err = site.Start(); nil != err {
		seelog.Error("gocodecc error:", err)
	}
}
