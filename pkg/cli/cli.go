package cli

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"
	"github.com/run-x/cloudgrep/pkg/api"
	"github.com/run-x/cloudgrep/pkg/command"
	"github.com/run-x/cloudgrep/pkg/util"
)

var (
	options command.Options
)

func exitWithMessage(message string) {
	fmt.Println("Error:", message)
	os.Exit(1)
}

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

	printVersion()
}

func printVersion() {
	chunks := []string{fmt.Sprintf("Cloudgrep v%s", command.Version)}

	if command.GitCommit != "" {
		chunks = append(chunks, fmt.Sprintf("(git: %s)", command.GitCommit))
	}

	if command.GoVersion != "" {
		chunks = append(chunks, fmt.Sprintf("(go: %s)", command.GoVersion))
	}

	if command.BuildTime != "" {
		chunks = append(chunks, fmt.Sprintf("(build time: %s)", command.BuildTime))
	}

	fmt.Println(strings.Join(chunks, " "))
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
	signal.Notify(c, os.Interrupt, os.Kill)
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

	exec.Command("open", url).Output()
}

func Run() {
	initOptions()
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
