package internal

type User interface {
	Identifier() string
}

type LoggedUser struct {
	Username string
	Password string
}

func (l *LoggedUser) Identifier() string {
	return l.Username + " " + l.Password
}

type AnonymousUser struct {
	Username string
}

func (a *AnonymousUser) Identifier() string {
	return a.Username
}

type Mail struct {
	User User
	From string
	To   string
	Data string
}

type MailHandler interface {
	Handle(*Mail)
}

type MailHandlerComposite struct {
	handlers []MailHandler
}

func (m *MailHandlerComposite) Handle(mail *Mail) {
	for _, handler := range m.handlers {
		handler.Handle(mail)
	}
}

func NewMailHandlerComposite(handler MailHandler, handlers ...MailHandler) *MailHandlerComposite {
	return &MailHandlerComposite{
		handlers: append([]MailHandler{handler}, handlers...),
	}
}

type MailChanNotifier struct {
	chanNotif chan *Mail
}

func (m *MailChanNotifier) Handle(mail *Mail) {
	m.chanNotif <- mail
}

func NewMailChanNotifier(chanNotif chan *Mail) *MailChanNotifier {
	return &MailChanNotifier{
		chanNotif: chanNotif,
	}
}

type MailStorage struct {
	mails []*Mail
}

func (m *MailStorage) Handle(mail *Mail) {
	m.mails = append(m.mails, mail)
}

func (m *MailStorage) GetMails() []*Mail {
	return m.mails
}
