package instagram

import "botvideosaver/internal/utils/common"

type clientImplWithRetry struct {
	*clientImpl
	retryCfg common.RetryConfig
}

func NewClientWithRetry(retryCfg common.RetryConfig) Client {
	return &clientImplWithRetry{
		clientImpl: &clientImpl{},
		retryCfg:   retryCfg,
	}
}

func (c *clientImplWithRetry) GetVideoURL(url string) (string, error) {
	return common.DoWithRetryAndReturn(
		c.retryCfg,
		func() (string, error) {
			return c.GetVideoURL(url)
		},
	)
}
