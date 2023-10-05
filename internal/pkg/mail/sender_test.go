package mail

import (
	"testing"

	"github.com/stretchr/testify/require"

	"playground/internal/config"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	cfg := config.Get()

	sender := NewGmailSender(cfg.EmailSenderName, cfg.EmailSenderAddress, cfg.EmailSenderPassword)
	subject := "A test email"
	content := `
	<h1>Hello world</h1>
	<p>This is a test message from <a href="http://techschool.guru">Tech School</a></p>
	`
	to := []string{"example@example.com"}
	err := sender.Send(subject, content, to, nil, nil, nil)
	require.NoError(t, err)
}
