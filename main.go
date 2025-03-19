package main

import (
	"encoding/json"
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

func main() {
	args := os.Args[1:]
	argsLen := len(args)

	var envFile string
	for _, arg := range args {
		if strings.HasPrefix(arg, "--env=") {
			envFile = strings.Split(arg, "=")[1]
			if envFile[0] == '"' && envFile[len(envFile)-1] == '"' {
				envFile = envFile[1 : len(envFile)-1]
			}
			argsLen--
		}

		if strings.HasPrefix(arg, "--format=") {
			formatArg := strings.Split(arg, "=")[1]
			if formatArg[0] == '"' && formatArg[len(formatArg)-1] == '"' {
				formatArg = formatArg[1 : len(formatArg)-1]
			}
			parsedFormat, err := parseFormat(formatArg)
			if err != nil {
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

	if argsLen > 1 {
		fmt.Println("Usage: " + COMMAND + " [command] [--env=path]")
		fmt.Println("Use '" + COMMAND + " help' for more information")
		os.Exit(1)
	}

	command := "status"
	if argsLen == 1 {
		for _, arg := range args {
			if ! strings.HasPrefix(arg, "--") {
		command = strings.ToLower(arg)
			}
		}
	}
	if command == "help" {
		showHelp()
		os.Exit(0)
	}

	dotEnv, _ := readEnv(envFile)

	kumaInstance, err := kuma.New(dotEnv[CONFIG_BASE_URL], dotEnv[CONFIG_API_KEY])
	if err != nil {
		panic(err)
	}

	switch command {
	case "list":
		handleList(kumaInstance)
	case "status":
		run(kumaInstance)
	case "open":
		kumaInstance.Open()
	default:
		fmt.Println("Invalid argument.")
		fmt.Println("Use '" + COMMAND + " help' for more information")
		os.Exit(1)
	}
}

func handleList(kumaInstance *kuma.Kuma) error {
	_, monitors, err := kumaInstance.GetMetrics()

	if err != nil {
		println("Unable to get monitors")
		panic(err)
	}

	for _, monitor := range monitors {
		if monitor.Status == kuma.Up {
			fmt.Printf("✅ %s - %s\n", monitor.Name, monitor.Type)
		}
	}

	for _, monitor := range monitors {
		if monitor.Status == kuma.Pending {
			fmt.Printf("⚠️ %s - %s\n", monitor.Name, monitor.Type)
		}
	}

	for _, monitor := range monitors {
		if monitor.Status == kuma.Down {
			fmt.Printf("❌ %s - %s\n", monitor.Name, monitor.Type)
		}
	}

	return nil
}

func showHelp() {
	fmt.Println("Usage: " + COMMAND + " [command] [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  help   # Display this help message")
	fmt.Println("  open   # open the Kuma dashboard in your default browser")
	fmt.Println("  status # (default) display a summary of the monitors")
	fmt.Println("  list   # list all monitors with their status")
	fmt.Println("\nOptions:")
	fmt.Println("  --env=path                            # specify the path to the .env file")
	fmt.Println("  --format=ansi|plain|waybar|json|jsonp # specify the output format")
	fmt.Println("\nExamples:")
	fmt.Println("  " + COMMAND + "                 # Shows a summary of Uptime Kuma's status")
	fmt.Println("  " + COMMAND + " open            # Opens Uptime Kuma dashboard in your default browser")
	fmt.Println("  " + COMMAND + " --env=/path/to/.custom.env")
	fmt.Println("  " + COMMAND + " --format=waybar # Formats the output for Waybar")
	fmt.Println("  " + COMMAND + " --format=json   # Shows verbose JSON output")
	fmt.Println("  " + COMMAND + " --format=jsonp  # Shows verbose JSON output with indentation")
	fmt.Println("  " + COMMAND + " --help")
}

func run(kumaInstance *kuma.Kuma) {
	metrics, monitors, err := kumaInstance.GetMetrics()
	if err != nil {
		panic(err)
	}

	if format == JSON {
		jsonStr, err := json.Marshal(map[string]interface{}{
			"monitors": monitors,
			"metrics":  metrics,
			"time":     kumaInstance.LastUpdated,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(string(jsonStr))
		return
	}

	if format == JSONP {
		jsonStr, err := json.MarshalIndent(map[string]interface{}{
			"monitors": monitors,
			"metrics":  metrics,
			"time":     kumaInstance.LastUpdated,
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
		if monitor.Status == kuma.Up {
			countGreen++
		} else if monitor.Status == kuma.Pending {
			countYellow++
		} else if monitor.Status == kuma.Down {
			countRed++
		}
	}

	if countGreen > 0 && countYellow == 0 && countRed == 0 {
		// check mark
		fmt.Println(formatCounter("✓", Green))
		return
	}

	out := ""
	if countGreen > 0 {
		out += formatCounter(strconv.FormatUint(countGreen, 10), Green) + " "
	}
	if countYellow > 0 {
		out += formatCounter(strconv.FormatUint(countYellow, 10), Yellow) + " "
	}
	if countRed > 0 {
		out += formatCounter(strconv.FormatUint(countRed, 10), Red) + " "
	}
	fmt.Println(out)
}
