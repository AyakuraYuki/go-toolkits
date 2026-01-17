package args

import (
	"errors"
	"fmt"
	"time"

	ayfile "github.com/AyakuraYuki/go-aybox/file"
	aytime "github.com/AyakuraYuki/go-aybox/time"
)

type M4RArgs struct {
	Src        string // 源文件的路径
	Dst        string // 保存铃声的路径
	Start      string // 截取的开始时间，格式：00:00:00.000
	End        string // 截取的结束时间，格式：00:00:00.000；如不设置，截取到最大iOS铃声时长为止
	FilterMode int    // 过滤器偏好。0：（默认）无增益；1：增益；2：高响度；-1：不使用过滤器

	ss time.Duration
	to time.Duration
}

func (args M4RArgs) Validate() (err error) {
	if args.Src == "" || args.Dst == "" {
		return errors.New("必须完整提供源文件和输出文件的路径")
	}

	if args.Src == args.Dst {
		return errors.New("出于对源文件的保护，输出文件不能覆盖源文件")
	}

	if !ayfile.PathExist(args.Src) {
		return fmt.Errorf("%v 文件不存在", args.Src)
	}

	if ayfile.IsDir(args.Src) {
		return fmt.Errorf("%v 不是一个有效的文件", args.Src)
	}

	if args.Start != "" {
		if _, err = aytime.ParseFlexibleDuration(args.Start); err != nil {
			return fmt.Errorf("开始时间不正确: %w", err)
		}
		if args.GetSS() < 0 {
			return fmt.Errorf("开始时间不正确: %v 小于 0", args.GetSS())
		}
	}

	if args.End != "" {
		if _, err = aytime.ParseFlexibleDuration(args.End); err != nil {
			return fmt.Errorf("结束时间不正确: %w", err)
		}
		if args.GetTo() < args.GetSS() {
			return fmt.Errorf("结束时间早于开始时间: %v, %v", args.GetSS(), args.GetTo())
		}
	}

	return nil
}

func (args M4RArgs) GetSS() time.Duration {
	if args.Start == "" {
		return args.ss
	}
	if args.ss > 0 {
		return args.ss
	}
	args.ss, _ = aytime.ParseFlexibleDuration(args.Start)
	return args.ss
}

func (args M4RArgs) GetTo() time.Duration {
	if args.End == "" {
		return args.to
	}
	if args.to > 0 {
		return args.to
	}
	args.to, _ = aytime.ParseFlexibleDuration(args.End)
	return args.to
}
