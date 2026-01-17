package exec

import (
	"context"
	"os"
	"strconv"

	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	ffmpeg "github.com/u2takey/ffmpeg-go"

	"github.com/AyakuraYuki/go-toolkits/cmd/ffmpeg-make/internal/args"
	"github.com/AyakuraYuki/go-toolkits/pkg/decimals"
)

const (
	maxIOSRingtoneSeconds = 40    // 默认 iOS 铃声时长（秒）
	fadeSeconds           = "1.5" // 淡入淡出标准时长（秒）
)

func (e *execute) MakeM4r(ctx context.Context, args args.M4RArgs) (err error) {
	keyTo, valueTo := "t", strconv.Itoa(maxIOSRingtoneSeconds)
	if args.End != "" {
		keyTo, valueTo = "to", args.End // 设置了结束时间点的
	}

	// 输入
	stream := ffmpeg.Input(args.Src, ffmpeg.KwArgs{
		"ss":  args.Start, // 开始时间点
		keyTo: valueTo,    // 结束时间点，或默认的40秒截取时长
	}).Audio()

	// 过滤器
	switch args.FilterMode {
	case 1:
		stream = m4rGainLimitFilter(stream, args)
	case 2:
		stream = m4rLoudnorm(stream, args)
	case -1:
		break // 不使用过滤器
	default:
		stream = m4rFadeInFadeOutFilter(stream, args)
	}

	// 输出
	stream = stream.Output(args.Dst, ffmpeg.KwArgs{
		"c:a":       "aac",     // 使用 AAC 编码（iOS 原生支持）
		"profile:a": "aac_low", // 符合 Apple 官方要求的 AAC-LC
		"b:a":       "256k",    // 高品质铃声，体积和音质的平衡选择
		"ar":        "44100",   // 更稳定播放的采样率
		"ac":        "2",       // 立体声双声道
	})

	err = stream.WithErrorOutput(os.Stderr).Run()
	if err != nil {
		e.logger.Error(err, "导出铃声失败")
		return err
	}

	return nil
}

func m4rGetFadeOutSt(args args.M4RArgs) string {
	dur := decimal.NewFromInt(maxIOSRingtoneSeconds)
	if args.GetTo() > 0 {
		dur = decimals.ParseString(cast.ToString(args.GetTo().Seconds()))
	}
	st := dur.Sub(decimals.ParseString(fadeSeconds))
	if st.IsNegative() {
		return "0"
	}
	return st.String()
}

// 不增益，淡入淡出
// "afade=t=in:st=0:d=1.5,afade=t=out:st=38.5:d=1.5"
func m4rFadeInFadeOutFilter(stream *ffmpeg.Stream, args args.M4RArgs) *ffmpeg.Stream {
	stream = stream.Filter("afade", nil, ffmpeg.KwArgs{
		"t":  "in",
		"st": "0",
		"d":  fadeSeconds,
	})

	stream = stream.Filter("afade", nil, ffmpeg.KwArgs{
		"t":  "out",
		"st": m4rGetFadeOutSt(args),
		"d":  fadeSeconds,
	})

	return stream
}

// 增益不失真
// "volume=1.4,alimiter=limit=0.98,afade=t=in:st=0:d=1.5,afade=t=out:st=38.5:d=1.5"
func m4rGainLimitFilter(stream *ffmpeg.Stream, args args.M4RArgs) *ffmpeg.Stream {
	// 增益
	stream = stream.Filter("volume", ffmpeg.Args{"1.4"})
	// 限制器，防止编码前削峰
	stream = stream.Filter("alimiter", nil, ffmpeg.KwArgs{"limit": "0.98"})
	// 淡入淡出
	return m4rFadeInFadeOutFilter(stream, args)
}

// 接近 Apple Music 的响度
// "loudnorm=I=-16:TP=-1.5:LRA=11,afade=t=in:st=0:d=1.5,afade=t=out:st=38.5:d=1.5"
func m4rLoudnorm(stream *ffmpeg.Stream, args args.M4RArgs) *ffmpeg.Stream {
	stream = stream.Filter("loudnorm", nil, ffmpeg.KwArgs{
		"I":   "-16",
		"TP":  "-1.5",
		"LRA": "11",
	})
	// 淡入淡出
	return m4rFadeInFadeOutFilter(stream, args)
}
