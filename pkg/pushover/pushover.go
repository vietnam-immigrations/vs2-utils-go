package pushover

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/samber/lo"

	"github.com/nam-truong-le/lambda-utils-go/pkg/logger"
	"github.com/vietnam-immigrations/vs2-utils-go/pkg/db"
)

const (
	pushoverURL = "https://api.pushover.net/1/messages.json"
)

func Send(ctx context.Context, title, message string) error {
	log := logger.FromContext(ctx)
	log.Infof("Pushover message [%s]", title)

	cfg, err := db.GetConfig(ctx)
	if err != nil {
		log.Errorf("Failed to get global config: %s", err)
		return err
	}
	users := lo.Map(strings.Split(cfg.PushoverUsers, ","), func(item string, _ int) string {
		return strings.TrimSpace(item)
	})
	log.Infof("Pushover to users: %s", strings.Join(users, ", "))

	httpClient := http.Client{
		Timeout: 30 * time.Second,
	}
	errs := make([]error, 0)
	for _, user := range users {
		form := url.Values{
			"token":   []string{cfg.PushoverToken},
			"user":    []string{user},
			"message": []string{message},
			"title":   []string{title},
			"html":    []string{"1"},
		}
		res, err := httpClient.PostForm(pushoverURL, form)
		if err != nil {
			log.Errorf("Failed to send http request to pushover [%s] to user [%s]: %s", title, user, err)
			errs = append(errs, err)
			continue
		}
		if res.StatusCode != 200 {
			log.Errorf("Pushover [%s] to user [%s] returns [%d]", title, user, res.StatusCode)
			errs = append(errs, err)
			continue
		}
		log.Infof("Pushover succeeded [%s]: %s", user, title)
	}
	if len(errs) > 0 {
		log.Errorf("Pushover error")
		errMessages := lo.Map(errs, func(e error, _ int) string { return e.Error() })
		return fmt.Errorf("pushover error: %s", strings.Join(errMessages, " | "))
	}
	return nil
}
