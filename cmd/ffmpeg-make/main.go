package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/AyakuraYuki/go-toolkits/cmd/ffmpeg-make/internal/commands"
)

var (
	version = "0.0.0-internal"
	binName = "ffmpeg-make"

	cmd     *cobra.Command
	verbose int
)

func init() {
	cmd = &cobra.Command{
		Use:     binName,
		Short:   "一个利用FFMpeg制作媒体的工具集合",
		Version: version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			switch verbose {
			case 2:
				logrus.SetLevel(logrus.TraceLevel)
			case 1:
				logrus.SetLevel(logrus.DebugLevel)
			default:
				logrus.SetLevel(logrus.InfoLevel)
			}
		},
	}

	cmd.PersistentFlags().CountVarP(&verbose, "verbose", "v", "运行日志跟踪级别（0: (default) Info, 1: Debug, 2: Trace）")

	cmd.AddCommand(
		commands.M4rCmd, // ffmpeg-make m4r
	)
}

func main() {
	ctx, cancel := notifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := cmd.ExecuteContext(ctx); err != nil {
		log.Fatal(err)
	}
}

func notifyContext(parent context.Context, signals ...os.Signal) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)
	ch := make(chan os.Signal, 5)
	signal.Notify(ch, signals...)
	if ctx.Err() == nil {
		go func() {
			// 第一次取消上下文
			select {
			case <-ctx.Done():
			case <-ch:
				cancel()
			}
			// 第二次直接退出
			select {
			case <-ctx.Done():
			case <-ch:
				os.Exit(1)
			}
		}()
	}
	return ctx, cancel
}
