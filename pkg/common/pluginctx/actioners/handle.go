package actioners

import (
	"github.com/jenkins-x/jx-logging/v3/pkg/log"

	"github.com/jenkins-x/octant-jx/pkg/common/actions"
	"github.com/jenkins-x/octant-jx/pkg/common/pluginctx"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
)

// CreateHandler creates a handler from the given context
func CreateHandler(c *pluginctx.Context) service.HandlerActionFunc {
	return func(request *service.ActionRequest) error {
		return HandleAction(request, c)
	}
}

// HandleAction handles the action to update the namespace
func HandleAction(request *service.ActionRequest, c *pluginctx.Context) error {
	switch request.ActionName {
	case action.RequestSetNamespace:
		namespace, err := request.Payload.String("namespace")
		if err != nil {
			log.Logger().Infof("failed to handle %s for payload %#v with error: %s", action.RequestSetNamespace, request.Payload, err.Error())
		} else {
			c.Namespace = namespace
			log.Logger().Infof("set the namespace to %s", namespace)
			log.Logger().Infof("payload is %#v", request.Payload)
		}

	case actions.PerformAction:
		action, err := request.Payload.String("action")
		if err != nil {
			namespace, err2 := request.Payload.String("namespace")
			if err2 == nil && namespace != "" {
				c.Namespace = namespace
				log.Logger().Infof("set the namespace to %s", namespace)
			} else {
				log.Logger().Infof("failed to handle %s for payload %#v with error: %s", actions.PerformAction, request.Payload, err.Error())
			}
		} else {
			switch action {
			case actions.TriggerBootJob:
				log.Logger().Infof("Trigger BOOT Job with payload %#v", request.Payload)
				// TODO how to get an alerter?
				return HandleTriggerBootJob(request)

			case actions.TriggerJob:
				log.Logger().Infof("Trigger Job with payload %#v", request.Payload)

				// TODO how to get an alerter?

			default:
				log.Logger().Infof("unknown  PerformAction %#v", request.Payload)
			}
		}

	default:
		log.Logger().Infof("unknown action %s", request.ActionName)
	}
	return nil
}
