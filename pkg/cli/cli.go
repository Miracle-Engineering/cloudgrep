package cli

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"
	"github.com/run-x/cloudgrep/pkg/api"
	"github.com/run-x/cloudgrep/pkg/command"
	"github.com/run-x/cloudgrep/pkg/util"
)

var (
	options command.Options
)

func initClient() {
	//TODO init AWS client
}

func initOptions() {
	opts, err := command.ParseOptions(os.Args)
	if err != nil {
		switch err.(type) {
		case *flags.Error:
			// no need to print error, flags package already does that
		default:
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}
	command.Opts = opts
	options = opts
}

func startServer() {
	router := gin.Default()

	api.SetupRoutes(router)

	fmt.Println("Starting server...")
	go func() {
		err := router.Run(fmt.Sprintf("%v:%v", options.HTTPHost, options.HTTPPort))
		if err != nil {
			fmt.Println("Cant start server:", err)
			if strings.Contains(err.Error(), "address already in use") {
				openPage()
			}
			os.Exit(1)
		}
	}()
}

func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

func openPage() {
	url := fmt.Sprintf("http://%v:%v", options.HTTPHost, options.HTTPPort)
	fmt.Println("To view Cloudgrep UI", url, "in browser")

	if options.SkipOpen {
		return
	}

	_, err := exec.Command("which", "open").Output()
	if err != nil {
		return
	}

	err = exec.Command("open", url).Run()
	if err != nil {
		return
	}

}

func Run() {
	initOptions()

	if options.Version {
		print(command.Version)
		os.Exit(0)
	}

	initClient()

	if options.Debug {
		gin.SetMode("debug")
	} else {
		gin.SetMode("release")
	}
	// Print memory usage every 30 seconds with debug flag
	if options.Debug {
		util.StartProfiler()
	}

	startServer()
	openPage()
	handleSignals()
}
