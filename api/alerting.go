package api

import (
	"fmt"
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/shieldwall/database"
	"time"
)

func (api *API) alertingLoop() {

	for {
		log.Debug("checking alerts ...")

		if agents, err := database.FindAgentsForAlert(); err != nil {
			log.Error("error querying database for alerts: %v", err)
		} else if num := len(agents); num > 0 {
			log.Debug("found %d agents for alerts", num)

			for _, agent := range agents {
				doAlert := false
				secsSinceLastAlert := uint(time.Since(agent.AlertAt).Seconds())

				// alert just once
				if agent.AlertPeriod == 0 && secsSinceLastAlert > database.MaxAlertPeriod {
					log.Debug("period 0, delta s from alert is %d, limit %d", secsSinceLastAlert, database.MaxAlertPeriod)
					doAlert = true
				} else if agent.AlertPeriod != 0 && secsSinceLastAlert >= agent.AlertPeriod {
					log.Debug("period %d, delta s from alert is %d", agent.AlertPeriod, secsSinceLastAlert)
					doAlert = true
				} else {
					log.Debug("already sent at %s for agent %d", agent.AlertAt, agent.ID)
				}

				if doAlert {
					log.Info("sending alert for agent %s (%d)", agent.Name, agent.ID)

					if user, err := database.FindUserByID(int(agent.UserID)); err != nil {
						log.Error("error searching for agent's user: %v", err)
					} else {
						emailSubject := fmt.Sprintf("shieldwall.me alert for %s", agent.Name)
						emailBody := fmt.Sprintf("Your agent '%s' has not been active since %s.",
							agent.Name,
							time.Since(agent.SeenAt))

						if err = api.sendmail.Send(api.mail.From, user.Email, emailSubject, emailBody); err != nil {
							log.Error("error sending alert email to %s: %v", user.Email, err)
						} else {
							agent.AlertAt = time.Now()
							if err = database.Save(&agent); err != nil {
								log.Error("error updating agent alert_at: %v", err)
							} else {
								log.Debug("agent alert_at updated")
							}
						}
					}
				}

			}
		}

		time.Sleep(time.Second * 30)
	}
}
