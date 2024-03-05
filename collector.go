package pgo

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"time"
)

type Collector interface {
	Start() error
}

type collector struct {
	client          *ApiClient
	endpoint        string
	cronTime        string
	cron            *cron.Cron
	uploader        Uploader
	applicationName string
}

func NewCollector(endpoint string, cronTime string, uploader Uploader, applicationName string) Collector {
	apiClient := NewHttpClient()

	return &collector{client: apiClient, endpoint: endpoint, cronTime: cronTime, uploader: uploader, applicationName: applicationName, cron: cron.New()}
}

func (c *collector) Start() error {
	_, err := c.cron.AddFunc(c.cronTime, func() {
		objectName := fmt.Sprintf("%s-%s", c.applicationName, time.Now().Format("2006-01-02T15:04:05"))
		pgoFile := c.getPgo() // write to minio or write to folder
		c.uploader.Upload(context.Background(), c.applicationName, objectName, pgoFile)
	})
	if err != nil {
		return err
	}
	c.cron.Start()
	return nil
}

func (c *collector) getPgo() []byte {
	response, _ := c.client.Get(c.endpoint)
	return response.Body
}
