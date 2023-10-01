package mail

import (
	"testing"

	"playground/config"

	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	cfg, err := config.NewConfig()
	require.NoError(t, err)

	sender := NewGmailSender(cfg.EmailSenderName, cfg.EmailSenderAddress, cfg.EmailSenderPassword)

	subject := "A test email"
	content := `
	<h1>Hello world</h1>
	<p>This is a test message from <a href="http://techschool.guru">Tech School</a></p>
	`
	to := []string{"gleamsoda99@gmail.com"}
	attachFiles := []string{"../../README.md"}

	err = sender.Send(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
