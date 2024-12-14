package report

import (
	"fmt"
	"github.com/go-lark/lark"
)

type LarkWriter struct {
	bot *lark.Bot
}

func NewLarkWriter(token string) *LarkWriter {
	bot := lark.NewNotificationBot(token)
	return &LarkWriter{
		bot: bot,
	}
}

func (l *LarkWriter) Write(p []byte) (n int, err error) {
	builder := lark.NewCardBuilder()
	card := builder.Card(builder.Markdown(fmt.Sprintf("```json \n%s", string(p)))).Yellow().Title("错误日志")
	msg := lark.NewMsgBuffer(lark.MsgInteractive)
	om := msg.Card(card.String()).Build()
	if _, err = l.bot.PostNotificationV2(om); err != nil {
		return 0, err
	}
	return len(p), nil

}
func (l *LarkWriter) Sync() error {
	return nil
}
