package commands

import (
	"github.com/bombsimon/logrusr/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/AyakuraYuki/go-toolkits/cmd/ffmpeg-make/internal/args"
	"github.com/AyakuraYuki/go-toolkits/cmd/ffmpeg-make/internal/exec"
)

var execute exec.Execute

func init() {
	logger := logrus.StandardLogger()
	logger.SetFormatter(&logrus.JSONFormatter{
		DisableHTMLEscape: true,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "@time",
			logrus.FieldKeyLevel: "@level",
			logrus.FieldKeyMsg:   "_message",
		},
	})

	options := exec.Options{Logger: logrusr.New(logger)}
	execute = exec.NewExecute(options)

	initM4rCmd()
}

var (
	m4rArgs args.M4RArgs
	M4rCmd  = &cobra.Command{
		Use:   "m4r",
		Short: "制作 iOS 铃声",

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return m4rArgs.Validate()
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return execute.MakeM4r(cmd.Context(), m4rArgs)
		},
	}
)

func initM4rCmd() {
	M4rCmd.Flags().StringVarP(&m4rArgs.Src, "input", "i", "", "源文件的路径")
	M4rCmd.Flags().StringVarP(&m4rArgs.Dst, "output", "o", "", "保存铃声的路径")
	M4rCmd.Flags().StringVarP(&m4rArgs.Start, "ss", "s", "00:00:00.000", "截取的开始时间，格式：00:00:00.000")
	M4rCmd.Flags().StringVarP(&m4rArgs.End, "to", "t", "", "截取的结束时间，格式：00:00:00.000；如不设置，截取到最大iOS铃声时长为止")
	M4rCmd.Flags().IntVar(&m4rArgs.FilterMode, "af", 0, "过滤器偏好。0：（默认）无增益；1：增益；2：高响度；-1：不使用过滤器")
}
