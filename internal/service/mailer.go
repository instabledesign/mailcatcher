package service

import (
	"github.com/instabledesign/mailcatcher/internal"
	"sync"
)

var mailHandlerOnce sync.Once

func (container *Container) GetMailHandler() internal.MailHandler {
	mailHandlerOnce.Do(func() {
		container.mailHandler = internal.NewMailHandlerComposite(
			container.GetMailChanNotifier(),
			container.GetMailStorage(),
		)

	})
	return container.mailHandler
}

func (container *Container) GetMailChan() chan *internal.Mail {
	return container.mailChan
}

var mailChanNotifierOnce sync.Once

func (container *Container) GetMailChanNotifier() *internal.MailChanNotifier {
	mailChanNotifierOnce.Do(func() {
		container.mailChanNotifier = internal.NewMailChanNotifier(container.GetMailChan())
	})
	return container.mailChanNotifier
}

var mailStorageOnce sync.Once

func (container *Container) GetMailStorage() *internal.MailStorage {
	mailStorageOnce.Do(func() {
		container.mailStorage = &internal.MailStorage{}
	})

	return container.mailStorage
}
