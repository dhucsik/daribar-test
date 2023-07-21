package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"

	"github.com/dhucsik/daribar-test/config"
	"github.com/dhucsik/daribar-test/model"
	"github.com/dhucsik/daribar-test/utils"
)

func ListenForDataUpdates(conf *config.Config) (*utils.Fanout, error) {
	f := utils.NewFanout()
	log.Println("fanout created")
	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	l := pq.NewListener(conf.DbUrl, 50*time.Millisecond, 10*time.Second, reportProblem)
	log.Println("listener created")
	err := l.Listen("data_updates")
	if err != nil {
		return nil, err
	}

	if err = l.Ping(); err != nil {
		return nil, err
	}
	log.Println("ping")
	go func() {
		pingTicker := time.Tick(10 * time.Second)
		for {
			var n *pq.Notification
			select {
			case <-pingTicker:
				log.Println("l")
				f.Publish(nil)
				continue
			case n = <-l.Notify:
				if n == nil {
					log.Print("data_updates: reconnected")
					continue
				}
			}
			// Let's assume that we have some formatted notification like "Operation: INSERT, OrderID: .., "
			// And then we formatted it
			order := model.Order{
				OrderID:   45,
				Phone:     "4552",
				CreatedAt: time.Now(),
				IsOpen:    true,
			}
			update := &model.Update{
				Order: order,
				Inc:   true,
			}

			f.Publish(update)
		}
	}()

	return f, nil

}
