package mailer

type MsgType string

const (
	ForgotPassword MsgType = "forgotPassword"
	Welcome        MsgType = "welcome"
	VerifyEmail    MsgType = "verifyEmail"
	Notification   MsgType = "notification"
)

type ForgotPasswordData struct {
	Name string
	Code string
}

type WelcomeData struct {
	Name string
}

type VerifyEmailData struct {
	Name string
	Link string
}

type NotificationData struct {
	Title   string
	Message string
}

type EmailParams struct {
	Name      string
	Recipient string
	Code      string
	Link      string
	Title     string
	Message   string
	Type      MsgType
}

var templateRegistry = map[MsgType]struct {
	File      string
	Subject   string
	BuildData func(p EmailParams) any
}{
	ForgotPassword: {
		File:    forgotPasswordTemplate,
		Subject: "Reset Your Password",
		BuildData: func(p EmailParams) any {
			return ForgotPasswordData{
				Name: p.Name,
				Code: p.Code,
			}
		},
	},
	Welcome: {
		File:    welcomeTemplate,
		Subject: "Welcome to our platform 🎉",
		BuildData: func(p EmailParams) interface{} {
			return WelcomeData{Name: p.Name}
		},
	},
	VerifyEmail: {
		File:    verifyEmailTemplate,
		Subject: "Verify Your Email",
		BuildData: func(p EmailParams) interface{} {
			return VerifyEmailData{
				Name: p.Name,
				Link: p.Link,
			}
		},
	},
	Notification: {
		File:    notificationTemplate,
		Subject: "Notification",
		BuildData: func(p EmailParams) interface{} {
			return NotificationData{
				Title:   p.Title,
				Message: p.Message,
			}
		},
	},
}
