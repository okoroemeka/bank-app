package mail

import (
	"github.com/okoroemeka/simple_bank/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	subject := "Test Email"
	content := `
				<h1>Test Email</h1>
				<p>This is a test email</p>
			`
	to := []string{"meka.okoro@gmail.com"}

	attachFiles := []string{"../start.sh"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)

	require.NoError(t, err)
}
