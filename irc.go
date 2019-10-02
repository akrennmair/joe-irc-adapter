package irc

import (
	"net"
	"strings"

	"github.com/go-joe/joe"
	irc "gopkg.in/irc.v3"
)

// Config contains all the IRC configuration
// that you need. Required fields are Address,
// Nick and Channel.
type Config struct {
	Address string
	Nick    string
	User    string
	Pass    string
	Name    string
	Channel string
}

// Adapter accepts a Config object and returns a
// joe.Module that can then used when creating a
// new bot using joe.New.
func Adapter(cfg Config) joe.Module {
	return joe.ModuleFunc(func(joeConf *joe.Config) error {
		conn, err := net.Dial("tcp", cfg.Address)
		if err != nil {
			return err
		}

		a := &ircAdapter{
			conn: conn,
			cfg:  cfg,
		}

		user := cfg.User
		if user == "" {
			user = cfg.Nick
		}

		ircConfig := irc.ClientConfig{
			Nick:    cfg.Nick,
			User:    user,
			Pass:    cfg.Pass,
			Name:    cfg.Name,
			Handler: a,
		}

		client := irc.NewClient(conn, ircConfig)
		a.client = client

		go a.client.Run()

		joeConf.SetAdapter(a)

		return nil
	})
}

type ircAdapter struct {
	client *irc.Client
	brain  *joe.Brain
	cfg    Config
	conn   net.Conn
}

func (a *ircAdapter) RegisterAt(brain *joe.Brain) {
	a.brain = brain
}

func (a *ircAdapter) Handle(c *irc.Client, m *irc.Message) {
	if m.Command == "001" {
		c.Write("JOIN " + a.cfg.Channel)
	} else if m.Command == "PRIVMSG" && c.FromChannel(m) {
		msg := m.Trailing()
		myPrefix := "@" + a.cfg.Nick + " "
		if !strings.HasPrefix(msg, myPrefix) {
			return
		}

		msg = msg[len(myPrefix):]
		a.brain.Emit(joe.ReceiveMessageEvent{
			Text:    msg,
			Channel: m.Name,
		})
	}
}

func (a *ircAdapter) Send(text, user string) error {
	return a.client.WriteMessage(&irc.Message{
		Command: "PRIVMSG",
		Params: []string{
			a.cfg.Channel,
			"@" + user + " " + text,
		},
	})
}

func (a *ircAdapter) Close() error {
	a.client.WriteMessage(&irc.Message{
		Command: "PART",
		Params:  []string{a.cfg.Channel},
	})
	return a.conn.Close()
}
