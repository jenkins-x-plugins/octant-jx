package pipelines

import (
	"context"
	"strings"

	v1 "github.com/jenkins-x/jx-api/pkg/apis/jenkins.io/v1"
	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/store"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetPipelines(ctx context.Context, client service.Dashboard, ns string) ([]*v1.PipelineActivity, error) {
	dl, err := client.List(ctx, store.Key{
		APIVersion: "jenkins.io/v1",
		Kind:       "PipelineActivity",
		Namespace:  ns,
	})

	paList := []*v1.PipelineActivity{}
	if dl != nil {
		for k, v := range dl.Items {
			pa, err := viewhelpers.ToPipelineActivity(&dl.Items[k])
			if err != nil {
				log.Logger().Infof("failed to convert to PipelineActivity for %s: %s", v.GetName(), err.Error())
				continue
			}
			if pa != nil {
				paList = append(paList, pa)
			}
		}
	}

	if len(paList) == 0 {
		dl, err = client.List(ctx, store.Key{
			APIVersion: "tekton.dev/v1beta1",
			Kind:       "PipelineRun",
			Namespace:  ns,
		})
		if dl != nil {
			for k, v := range dl.Items {
				pr := &v1beta1.PipelineRun{}
				err := viewhelpers.ToStructured(&dl.Items[k], pr)
				if err != nil {
					log.Logger().Infof("failed to convert to PipelineActivity for %s: %s", v.GetName(), err.Error())
					continue
				}

				pa := ToPipelineActivity(pr)
				if pa != nil {
					paList = append(paList, pa)
				}
			}
		}
	}
	return paList, err
}

func GetPipeline(ctx context.Context, client service.Dashboard, ns, name string) (*v1.PipelineActivity, error) {
	u, err := viewhelpers.GetResourceByName(ctx, client, "jenkins.io/v1", "PipelineActivity", name, ns)
	if err != nil {
		u2, err2 := viewhelpers.GetResourceByName(ctx, client, "tekton.dev/v1beta1", "PipelineRun", name, ns)
		if err2 == nil {
			pr := &v1beta1.PipelineRun{}
			err2 = viewhelpers.ToStructured(u2, pr)
			if err2 == nil {
				pa := ToPipelineActivity(pr)
				return pa, nil
			}
		}
		return nil, nil
	}
	return viewhelpers.ToPipelineActivity(u)
}

func ToPipelineActivity(pr *v1beta1.PipelineRun) *v1.PipelineActivity {
	annotations := pr.Annotations
	labels := pr.Labels
	pa := &v1.PipelineActivity{}
	pa.Name = pr.Name
	pa.Namespace = pr.Namespace
	pa.Annotations = annotations
	pa.Labels = labels

	ps := &pa.Spec
	if labels != nil {
		if ps.GitOwner == "" {
			ps.GitOwner = labels["lighthouse.jenkins-x.io/refs.org"]
		}
		if ps.GitRepository == "" {
			ps.GitRepository = labels["lighthouse.jenkins-x.io/refs.repo"]
		}
		if ps.GitBranch == "" {
			ps.GitBranch = labels["lighthouse.jenkins-x.io/branch"]
		}
		if ps.Build == "" {
			ps.Build = labels["lighthouse.jenkins-x.io/buildNum"]
		}
		if ps.Context == "" {
			ps.Context = labels["lighthouse.jenkins-x.io/context"]
		}
		if ps.BaseSHA == "" {
			ps.BaseSHA = labels["lighthouse.jenkins-x.io/baseSHA"]
		}
		if ps.LastCommitSHA == "" {
			ps.LastCommitSHA = labels["lighthouse.jenkins-x.io/lastCommitSHA"]
		}
	}
	if annotations != nil {
		if ps.GitURL == "" {
			ps.GitURL = annotations["lighthouse.jenkins-x.io/cloneURI"]
		}
	}

	podName := ""
	if pr.Status.TaskRuns != nil {
		for _, v := range pr.Status.TaskRuns {
			if v.Status == nil {
				continue
			}
			if podName == "" {
				podName = v.Status.PodName
			}
			for _, step := range v.Status.Steps {
				name := step.Name
				var started *metav1.Time
				var completed *metav1.Time
				status := v1.ActivityStatusTypePending

				terminated := step.Terminated
				if terminated != nil {
					if terminated.ExitCode == 0 {
						status = v1.ActivityStatusTypeSucceeded
					} else if !terminated.FinishedAt.IsZero() {
						status = v1.ActivityStatusTypeFailed
					}
					started = &terminated.StartedAt
					completed = &terminated.FinishedAt
				}

				paStep := v1.PipelineActivityStep{
					Kind: v1.ActivityStepKindTypeStage,
					Stage: &v1.StageActivityStep{
						CoreActivityStep: v1.CoreActivityStep{
							Name:               Humanize(name),
							Description:        "",
							Status:             status,
							StartedTimestamp:   started,
							CompletedTimestamp: completed,
						},
					},
				}
				ps.Steps = append(ps.Steps, paStep)
			}
		}
	}
	if podName != "" {
		pa.Labels["podName"] = podName
	}
	return pa
}

// Humanize splits into words and capitalises
func Humanize(text string) string {
	wordsText := strings.ReplaceAll(strings.ReplaceAll(text, "-", " "), "_", " ")
	words := strings.Split(wordsText, " ")
	for i := range words {
		words[i] = strings.Title(words[i])
	}
	return strings.Join(words, " ")
}
