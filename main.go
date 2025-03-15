package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/webtender/kuma-waybar/kuma"
)

const CONFIG_API_KEY = "UPTIME_KUMA_API_KEY"
const CONFIG_BASE_URL = "UPTIME_KUMA_BASE_URL"
const COMMAND = "kuma-waybar"
var format = ANSI

type Format int
const (
	PLAIN Format = 0
	ANSI Format = 1
	WAYBAR Format = 2
	JSON Format = 3
	JSONP Format = 4
)

func parseFormat(format string) (Format, error) {
	switch format {
	case "plain":
		return PLAIN, nil
	case "ansi":
		return ANSI, nil
	case "waybar":
		return WAYBAR, nil
	case "json":
		return JSON, nil
	case "jsonp":
		return JSONP, nil
	}
	return ANSI, errors.New("invalid format: " + format)
}

func main() {
	args := os.Args[1:]
	argsLen := len(args)

	var envFile string
	for _, arg := range args {
		if strings.HasPrefix(arg, "--env=") {
			envFile = strings.Split(arg, "=")[1]
			if (envFile[0] == '"' && envFile[len(envFile)-1] == '"') {
				envFile = envFile[1 : len(envFile)-1]
			}
			argsLen--
		}

		if strings.HasPrefix(arg, "--format=") {
			formatArg := strings.Split(arg, "=")[1]
			if (formatArg[0] == '"' && formatArg[len(formatArg)-1] == '"') {
				formatArg = formatArg[1 : len(formatArg)-1]
			}
			parsedFormat, err := parseFormat(formatArg)
			if (err != nil) {
				fmt.Println(err)
				os.Exit(1)
			}
			format = parsedFormat
			argsLen--
		}
		
		if arg == "--help" {
			showHelp()
			os.Exit(0)
		}
	}

	if (argsLen > 1) {
		fmt.Println("Usage: " + COMMAND + " [command] [--env=path]")
		fmt.Println("Use '" + COMMAND + " help' for more information")
		os.Exit(1)
	}

	command := "update"
	if (argsLen == 1) {
		command = strings.ToLower(args[0])
	}
	if (command == "help" ){
		showHelp()
		os.Exit(0)
	}

	dotEnv, _ := readEnv(envFile)

	kumaInstance, err := kuma.New(dotEnv[CONFIG_BASE_URL], dotEnv[CONFIG_API_KEY])
	if err != nil {
		panic(err)
	}

	switch(command) {
		case "update":
			run(kumaInstance)
		case "open":
			kumaInstance.Open()
		default:
			fmt.Println("Invalid argument.")
			fmt.Println("Use '" + COMMAND + " help' for more information")
			os.Exit(1)
	}
}

func showHelp() {
	fmt.Println("Usage: " + COMMAND + " [command] [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  help: Display this help message")
	fmt.Println("  open: open the Kuma dashboard in your default browser")
	fmt.Println("  update: (default) display a summary of the monitors")
	fmt.Println("\nOptions:")
	fmt.Println("  --env=path: specify the path to the .env file")
	fmt.Println("  --format=ansi|plain|waybar|json|jsonp: specify the output format")
	fmt.Println("\nExamples:")
	fmt.Println("  " + COMMAND + " update --env=.env")
	fmt.Println("  " + COMMAND + " open --env=.env")
	fmt.Println("  " + COMMAND + " --env=.env")
	fmt.Println("  " + COMMAND + " --format=waybar # Formats the output for Waybar")
	fmt.Println("  " + COMMAND + " --format=json # Shows verbose JSON output")
	fmt.Println("  " + COMMAND + " --format=jsonp # Shows verbose JSON output with indentation")
	fmt.Println("  " + COMMAND + " --help")
}

func run(kumaInstance *kuma.Kuma) {
	metrics, monitors, err := kumaInstance.GetMetrics()
	if err != nil {
		panic(err)
	}

	if (format == JSON) {
		jsonStr, err := json.Marshal(map[string]interface{}{
			"monitors": monitors,
			"metrics": metrics,
			"time": kumaInstance.LastUpdated,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(string(jsonStr))
		return
	}
	
	if (format == JSONP) {
		jsonStr, err := json.MarshalIndent(map[string]interface{}{
			"monitors": monitors,
			"metrics": metrics,
			"time": kumaInstance.LastUpdated,
		}, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(jsonStr))
		return
	}

	var countGreen uint64 = 0
	var countYellow uint64 = 0
	var countRed uint64 = 0

	for _, monitor := range monitors {
		if (monitor.Status == kuma.Up) {
			countGreen++
		} else if (monitor.Status == kuma.Pending) {
			countYellow++
		} else if (monitor.Status == kuma.Down) {
			countRed++
		}
	}

	if (countGreen > 0 && countYellow == 0 && countRed == 0) {
		// check mark
		fmt.Println(formatCounter("✓", Green))
		return;
	}

	out := ""
	if (countGreen > 0) {
		out += formatCounter(strconv.FormatUint(countGreen, 10), Green) + " "
	}
	if (countYellow > 0) {
		out += formatCounter(strconv.FormatUint(countYellow, 10), Yellow) + " "
	}
	if (countRed > 0) {
		out += formatCounter(strconv.FormatUint(countRed, 10), Red) + " "
	}
	fmt.Println(out)
}

type Color int
const (
	Red Color = 1
	Green Color = 2
	Yellow Color = 3
)

func (c Color) hex() string {
	switch c {
	case Red:
		return "#FF0000"
	case Green:
		return "#00FF00"
	case Yellow:
		return "#FFBF00"
	}
	return ""
}

func (c Color) ansi() string {
	switch c {
	case Red:
		return "1"
	case Green:
		return "2"
	case Yellow:
		return "3"
	}
	return ""
}

func formatCounter(str string, color Color) string {
	if (format == PLAIN) {
		return str
	}
	if (format == ANSI) {
		return "\033[38;5;" + color.ansi() + "m" + str + "\033[0m"
	}
	return "<span foreground=\"" + color.hex() + "\">" + str + "</span>"
}
