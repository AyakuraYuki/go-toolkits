package dingtalk

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/AyakuraYuki/crypto-go"
	"github.com/go-resty/resty/v2"

	"github.com/AyakuraYuki/go-toolkits/pkg/cjson"
)

const signApi = "https://oapi.dingtalk.com/robot/send?access_token=%s&timestamp=%d&sign=%s"

type Webhook struct {
	queue chan *Message
	cli   *resty.Client
}

type Message struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	Auth  Auth   `json:"auth"`
}

type Auth struct {
	Token     string `json:"-"`
	Secret    string `json:"-"`
	timestamp int64
}

func (m *Auth) endpoint() string {
	m.timestamp = time.Now().UnixMilli()
	val := fmt.Sprintf("%d\n%s", m.timestamp, m.Secret)
	secretBS := []byte(m.Secret)
	valBS := []byte(val)
	hmacBS := crypto.HmacSHA256(valBS, secretBS)
	sign := base64.StdEncoding.EncodeToString(hmacBS)

	u := url.URL{
		Scheme: "https",
		Host:   "oapi.dingtalk.com",
		Path:   "/robot/send",
	}
	q := u.Query()
	q.Set("access_token", m.Token)
	q.Set("timestamp", fmt.Sprintf("%d", m.timestamp))
	q.Set("sign", sign)
	u.RawQuery = q.Encode()

	return u.String()
}

func (m *Message) String() string {
	return cjson.Stringify(m)
}

type markdown struct {
	MsgType  string   `json:"msgtype"`
	Markdown *Message `json:"markdown"`
	At       *at      `json:"at"`
}

func (m *markdown) String() string {
	return cjson.Stringify(m)
}

type at struct {
	AtMobiles []string `json:"atMobiles"`
	AtUserIds []string `json:"atUserIds"`
	IsAtAll   bool     `json:"isAtAll"`
}

var (
	instance *Webhook
	once     sync.Once
)

func init() {
	once.Do(func() {
		instance = &Webhook{
			queue: make(chan *Message, 32),
		}

		cli := resty.New()
		cli.SetTimeout(30 * time.Second)
		cli.SetRetryCount(0)
		cli.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
		instance.cli = cli

		go instance.start()
	})
}

func Notify(m *Message) {
	select {
	case instance.queue <- m:
	default:
	}
}

func (w *Webhook) start() {
	for {
		select {
		case m := <-w.queue:
			w.send(m)
		}
	}
}

func (w *Webhook) send(m *Message) {
	if m == nil {
		return
	}

	var now time.Time
	if loc, _ := time.LoadLocation("Asia/Shanghai"); loc != nil {
		now = time.Now().In(loc)
	} else {
		now = time.Now().In(time.UTC)
	}

	m.Text = fmt.Sprintf(`
### %s
#### 时间: %s

> %s`, m.Title, now.Format(time.DateTime), m.Text)

	md := &markdown{
		MsgType:  "markdown",
		Markdown: m,
		At: &at{
			IsAtAll: true,
		},
	}

	rsp, err := w.cli.R().
		SetHeader("Content-Type", "application/json").
		SetBody(md).
		Post(m.Auth.endpoint())
	if err != nil {
		log.Printf("[x] dingtalk notify failed: %v", err)
	}

	log.Printf(" <- dingtalk response: %s", rsp.String())
}
