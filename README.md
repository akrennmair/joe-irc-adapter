<h1 align="center">Joe Bot - IRC Adapter</h1>
<p align="center">IRC adapter. https://github.com/go-joe/joe</p>

IRC adapter for: https://github.com/go-joe/joe

This simple IRC adapter forwards all messages directed to it
on a particular channel to the brain. To direct a message to
the bot, the IRC message needs to contain the bot's IRC nick
prefixed with `@` at the beginning, e.g. `@thebot hello world`
makes the IRC adapter receive `hello world`. As the "channel",
it uses the IRC nick of the user who sent the message. The
"channel" will be used in responses to directly address the
same user, in the form `@user this is the response`.

### Example

```go
cfg := irc.Config{
	Address: "my-irc-server.com:6667",
	Nick: "my-irc-bot",
	Name: "My IRC Bot",
	Channel: "#my-irc-channel",
}

b := &ExampleBot{
	Bot: joe.New("my-irc-bot", irc.Adapter(cfg),
}
```
