package servicex

import "context"

type Service interface {
	Start(context.Context) error
	Stop(context.Context) error
}

type Compose struct{ services []Service }

func NewCompose(services ...Service) *Compose {
	return &Compose{services: append([]Service(nil), services...)}
}

func (c *Compose) Start(ctx context.Context) error {
	for _, service := range c.services {
		if err := service.Start(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (c *Compose) Stop(ctx context.Context) error {
	for i := len(c.services) - 1; i >= 0; i-- {
		if err := c.services[i].Stop(ctx); err != nil {
			return err
		}
	}
	return nil
}
