package pipelines

import (
	"context"

	v1 "github.com/jenkins-x/jx-api/v4/pkg/apis/jenkins.io/v1"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	jxpipeline "github.com/jenkins-x/jx-pipeline/pkg/pipelines"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/pkg/errors"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/store"
)

var lookupByNameWhichDoesntExistBreaksOctant = true

func GetPipelines(ctx context.Context, client service.Dashboard, ns string) ([]v1.PipelineActivity, error) {
	dl, err := client.List(ctx, store.Key{
		APIVersion: "jenkins.io/v1",
		Kind:       "PipelineActivity",
		Namespace:  ns,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to ")
	}

	paList := []v1.PipelineActivity{}
	if dl != nil {
		for k, v := range dl.Items {
			var pa *v1.PipelineActivity
			pa, err = viewhelpers.ToPipelineActivity(&dl.Items[k])
			if err != nil {
				log.Logger().Infof("failed to convert to PipelineActivity for %s: %s", v.GetName(), err.Error())
				continue
			}
			if pa != nil {
				paList = append(paList, *pa)
			}
		}
	}

	dl, err = client.List(ctx, store.Key{
		APIVersion: "tekton.dev/v1beta1",
		Kind:       "PipelineRun",
		Namespace:  ns,
	})
	if dl != nil {
		for k, v := range dl.Items {
			pr := &v1beta1.PipelineRun{}
			err = viewhelpers.ToStructured(&dl.Items[k], pr)
			if err != nil {
				log.Logger().Infof("failed to convert to PipelineActivity for %s: %s", v.GetName(), err.Error())
				continue
			}
			pr.Name = jxpipeline.ToPipelineActivityName(pr, paList)
			if pr.Name == "" {
				continue
			}

			var pa *v1.PipelineActivity
			for i := range paList {
				r := &paList[i]
				if r.Name == pr.Name {
					pa = &paList[i]
					break
				}
			}
			if pa == nil {
				paList = append(paList, v1.PipelineActivity{})
				pa = &paList[len(paList)-1]
			}
			jxpipeline.ToPipelineActivity(pr, pa, false)
		}
	}
	return paList, err
}

func GetPipeline(ctx context.Context, client service.Dashboard, ns, name string) (*v1.PipelineActivity, error) {
	if lookupByNameWhichDoesntExistBreaksOctant {
		paList, err := GetPipelines(ctx, client, ns)
		if err != nil {
			log.Logger().Infof("failed to list PipelineActivity in namespace %s: %s", ns, err.Error())
		}
		for i := range paList {
			pa := &paList[i]
			if pa.Name == name {
				return pa, nil
			}
		}
		return nil, nil
	}
	u, err := viewhelpers.GetResourceByName(ctx, client, "jenkins.io/v1", "PipelineActivity", name, ns)
	if err != nil {
		paList, err2 := GetPipelines(ctx, client, ns)
		if err2 == nil {
			for i := range paList {
				pa := &paList[i]
				if pa.Name == name {
					return pa, nil
				}
			}
		}
		return nil, nil
	}
	pa, err := viewhelpers.ToPipelineActivity(u)
	if err != nil {
		return pa, errors.Wrapf(err, "failed to load PipelineActivity")
	}

	dl, err := client.List(ctx, store.Key{
		APIVersion: "tekton.dev/v1beta1",
		Kind:       "PipelineRun",
		Namespace:  ns,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load PipelineRuns")
	}
	paList := []v1.PipelineActivity{*pa}
	if dl != nil {
		for k, v := range dl.Items {
			pr := &v1beta1.PipelineRun{}
			err := viewhelpers.ToStructured(&dl.Items[k], pr)
			if err != nil {
				log.Logger().Infof("failed to convert to PipelineActivity for %s: %s", v.GetName(), err.Error())
				continue
			}

			pr.Name = jxpipeline.ToPipelineActivityName(pr, paList)
			if pr.Name == name {
				jxpipeline.ToPipelineActivity(pr, pa, false)
				return pa, nil
			}
		}
	}
	return pa, nil
}
