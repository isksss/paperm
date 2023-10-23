package commands

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/google/subcommands"
	"github.com/isksss/paperm/commands/download"
	"github.com/isksss/paperm/config"
)

type StartCommand struct {
}

func (c *StartCommand) Name() string { return "start" }

func (c *StartCommand) Synopsis() string { return "start server." }

func (c *StartCommand) Usage() string { return "start" }

func (c *StartCommand) SetFlags(f *flag.FlagSet) {
}

func (c *StartCommand) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	// Download
	if err := download.ServerDownload(); err != nil {
		fmt.Printf("download error: %v\n", err)
		return subcommands.ExitFailure
	}

	// Plugin download
	if err := download.PluginDownload(); err != nil {
		fmt.Printf("plugin download error: %v\n", err)
		return subcommands.ExitFailure
	}

	// Get restart time.
	data, err := config.GetConfig()
	if err != nil {
		fmt.Printf("config read error: %v\n", err)
		return subcommands.ExitFailure
	}

	var parsedTime []time.Time
	for _, restartTime := range data.Server.RestartTime {
		parsed, err := time.Parse("15:04", restartTime)
		if err != nil {
			continue
		}

		parsedTime = append(parsedTime, parsed)
	}

	// exists java
	javaCmd := "java"
	javaBin, err := exec.LookPath(javaCmd)
	if err != nil {
		return subcommands.ExitFailure
	}

	//papermcの起動
	// opt
	min := strconv.Itoa(data.Server.MinMemory)
	max := strconv.Itoa(data.Server.MaxMemory)
	xms := "-Xms" + min + "M"
	xmx := "-Xmx" + max + "M"
	jar := data.JarName
	jarOpt := "nogui"

	for {
		// server
		cmd := exec.Command(javaBin, "-jar", xms, xmx, jar, jarOpt)

		// stdin取得
		stdin, err := cmd.StdinPipe()
		if err != nil {
			fmt.Println(err)
			return subcommands.ExitFailure
		}

		// server 起動
		err = cmd.Start()
		if err != nil {
			return subcommands.ExitFailure
		}

		// 再起動時間
		durations := getDurations(parsedTime)

		duration := findMinDuration(durations)

		announceTime := 1 * time.Minute
		announce := time.AfterFunc(duration-announceTime, func() {
			msg := fmt.Sprintf("say %s \015", data.Server.AnnounceMessage)
			fmt.Printf("Announce: %s\n", msg)
			io.WriteString(stdin, msg)
		})
		timer := time.AfterFunc(duration, func() {
			io.WriteString(stdin, "stop\015")
		})

		pid := cmd.Process.Pid
		fmt.Println("Restart:" + duration.String())
		fmt.Printf("Server start: %d\n", pid)
		// 標準入力を監視するゴルーチンを起動
		allstopchan := make(chan bool)
		go func() {
			reader := bufio.NewReader(os.Stdin)
			for {
				fmt.Printf("> ")
				input, _ := reader.ReadString('\n')

				if strings.TrimSpace(input) == "paperstop" {
					// サーバー停止
					fmt.Printf("Server stop: %d\n", pid)
					io.WriteString(stdin, "stop\015")
					return
				}

				if strings.TrimSpace(input) == "allstop" {
					// papermc-managerを停止
					fmt.Printf("Server stop: %d\n", pid)
					io.WriteString(stdin, "stop\015")
					allstopchan <- true
					return
				}

				// 残り時間を表示
				if strings.TrimSpace(input) == "time" {

					nowDuration := findMinDuration(getDurations(parsedTime))
					fmt.Printf("Remaining time: %s\n", nowDuration.String())
				}

				// サーバーにコマンドを送信
				if strings.TrimSpace(input) == "cmd" {
					fmt.Printf("cmd >>> ")
					input, _ := reader.ReadString('\n')
					io.WriteString(stdin, input)
				}
			}
		}()
		cmd.Wait()
		announce.Stop()
		timer.Stop()

		select {
		case <-allstopchan:
			fmt.Println("Done!")
			return subcommands.ExitSuccess
		default:
			fmt.Println("Restarting...")
		}
	}
}

func findMinDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	min := durations[0]
	for _, duration := range durations[1:] {
		if duration < min {
			min = duration
		}
	}

	return min
}

func getDurations(parsedTime []time.Time) []time.Duration {
	now := time.Now()
	var durations []time.Duration
	for _, restartTime := range parsedTime {
		n_h := now.Hour()
		n_m := now.Minute()
		r_h := restartTime.Hour()
		r_m := restartTime.Minute()
		var next time.Time
		y := now.Year()
		m := now.Month()
		d := now.Day()
		l := now.Location()

		if n_h >= r_h && n_m >= r_m {
			next = time.Date(y, m, d+1, r_h, r_m, 0, 0, l)
		} else {
			next = time.Date(y, m, d, r_h, r_m, 0, 0, l)
		}

		duration := next.Sub(now)
		durations = append(durations, duration)
	}

	return durations
}
