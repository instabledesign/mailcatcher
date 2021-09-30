package service

import (
	"encoding/json"
	"fmt"
	"github.com/instabledesign/mailcatcher/internal/server/http/handler"
	"log"
	"sync"
)

var brokerOnce sync.Once

func (container *Container) GetBroker() *handler.Broker {
	brokerOnce.Do(func() {
		container.broker = handler.NewBroker()
		go func() {
			for mail := range container.GetMailChan() {
				fmt.Println("NEW MESSAGE", mail)
				data, err := json.Marshal(mail)
				if err != nil {
					log.Println("error during marshaling email", mail, err)
				}
				container.broker.Notifier <- data
			}
		}()
	})

	return container.broker
}
