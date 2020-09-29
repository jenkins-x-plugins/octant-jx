package pipelines

import (
	"context"
	"strings"

	v1 "github.com/jenkins-x/jx-api/pkg/apis/jenkins.io/v1"
	"github.com/jenkins-x/jx-helpers/pkg/kube/activities"
	"github.com/jenkins-x/jx-helpers/pkg/kube/naming"
	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/pkg/errors"
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
			pr.Name = ToPipelineActivityName(pr, paList)

			var pa *v1.PipelineActivity
			for i, r := range paList {
				if r.Name == pr.Name {
					pa = paList[i]
					break
				}
			}
			if pa == nil {
				pa = &v1.PipelineActivity{}
				paList = append(paList, pa)
			}
			ToPipelineActivity(pr, pa)
		}
	}
	return paList, err
}

func GetPipeline(ctx context.Context, client service.Dashboard, ns, name string) (*v1.PipelineActivity, error) {
	u, err := viewhelpers.GetResourceByName(ctx, client, "jenkins.io/v1", "PipelineActivity", name, ns)
	if err != nil {
		paList, err2 := GetPipelines(ctx, client, ns)
		if err2 == nil {
			for _, pa := range paList {
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
	paList := []*v1.PipelineActivity{pa}
	if dl != nil {
		for k, v := range dl.Items {
			pr := &v1beta1.PipelineRun{}
			err := viewhelpers.ToStructured(&dl.Items[k], pr)
			if err != nil {
				log.Logger().Infof("failed to convert to PipelineActivity for %s: %s", v.GetName(), err.Error())
				continue
			}

			pr.Name = ToPipelineActivityName(pr, paList)
			if pr.Name == name {
				ToPipelineActivity(pr, pa)
				return pa, nil
			}
		}
	}
	return pa, nil
}

func ToPipelineActivityName(pr *v1beta1.PipelineRun, paList []*v1.PipelineActivity) string {
	labels := pr.Labels
	if labels == nil {
		return pr.Name
	}
	build := labels["build"]
	owner := labels["lighthouse.jenkins-x.io/refs.org"]
	if owner == "" {
		owner = labels["owner"]
	}
	repository := labels["lighthouse.jenkins-x.io/refs.repo"]
	if repository == "" {
		repository = labels["repository"]
	}
	branch := labels["lighthouse.jenkins-x.io/branch"]
	if branch == "" {
		branch = labels["branch"]
	}
	if build == "" {
		buildID := labels["lighthouse.jenkins-x.io/buildNum"]
		if buildID == "" {
			return pr.Name
		}
		for _, pa := range paList {
			if pa.Labels == nil {
				continue
			}
			if pa.Labels["buildID"] == buildID {
				if pa.Spec.Build != "" {
					pr.Labels["build"] = pa.Spec.Build
				}
				return pa.Name
			}
		}
		if owner != "" && repository != "" && branch != "" {
			build = "1"
			pr.Labels["build"] = build
		} else {
			return pr.Name
		}
	}
	if owner != "" && repository != "" && branch != "" && build != "" {
		return naming.ToValidName(owner + "-" + repository + "-" + branch + "-" + build)
	}
	return pr.Name
}

func ToPipelineActivity(pr *v1beta1.PipelineRun, pa *v1.PipelineActivity) {
	annotations := pr.Annotations
	labels := pr.Labels
	if pa.APIVersion == "" {
		pa.APIVersion = "jenkins.io/v1"
	}
	if pa.Kind == "" {
		pa.Kind = "PipelineActivity"
	}
	pa.Name = pr.Name
	pa.Namespace = pr.Namespace

	if pa.Annotations == nil {
		pa.Annotations = map[string]string{}
	}
	if pa.Labels == nil {
		pa.Labels = map[string]string{}
	}
	for k, v := range annotations {
		pa.Annotations[k] = v
	}
	for k, v := range labels {
		pa.Labels[k] = v
	}

	ps := &pa.Spec
	if labels != nil {
		if ps.GitOwner == "" {
			ps.GitOwner = labels["lighthouse.jenkins-x.io/refs.org"]
			if ps.GitOwner == "" {
				ps.GitOwner = labels["owner"]
			}
		}
		if ps.GitRepository == "" {
			ps.GitRepository = labels["lighthouse.jenkins-x.io/refs.repo"]
			if ps.GitRepository == "" {
				ps.GitRepository = labels["repository"]
			}
		}
		if ps.GitBranch == "" {
			ps.GitBranch = labels["lighthouse.jenkins-x.io/branch"]
			if ps.GitBranch == "" {
				ps.GitBranch = labels["branch"]
			}
		}
		if ps.Build == "" {
			ps.Build = labels["build"]
			if ps.Build == "" {
				ps.Build = labels["lighthouse.jenkins-x.io/buildNum"]
			}
		}
		if ps.Context == "" {
			ps.Context = labels["lighthouse.jenkins-x.io/context"]
			if ps.Context == "" {
				ps.Context = labels["context"]
			}
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
	var steps []v1.PipelineActivityStep
	if pr.Status.TaskRuns != nil {
		for _, v := range pr.Status.TaskRuns {
			if v.Status == nil {
				continue
			}
			if podName == "" {
				podName = v.Status.PodName
			}

			previousStepTerminated := false
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
					previousStepTerminated = true
				} else if step.Running != nil {
					if previousStepTerminated {
						started = &step.Running.StartedAt
						status = v1.ActivityStatusTypeRunning
					}
					previousStepTerminated = false
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
				steps = append(steps, paStep)
			}
		}
	}

	// if the PipelineActivity has some real steps lets trust it; otherise lets merge any prevew/promote steps
	// with steps from the PipelineRun
	// lets add any missing steps from the PipelineActivity as they may have been created via a `jx promote` step
	hasStep := false
	for _, s := range ps.Steps {
		if s.Kind == v1.ActivityStepKindTypeStage && s.Stage != nil && s.Stage.Name != "Release" {
			hasStep = true
			break
		}
	}
	if !hasStep {
		for _, s := range ps.Steps {
			if s.Kind == v1.ActivityStepKindTypePreview || s.Kind == v1.ActivityStepKindTypePromote {
				steps = append(steps, s)
			}
		}
		ps.Steps = steps
	}

	if len(ps.Steps) == 0 {
		ps.Steps = append(ps.Steps, v1.PipelineActivityStep{
			Kind: v1.ActivityStepKindTypeStage,
			Stage: &v1.StageActivityStep{
				CoreActivityStep: v1.CoreActivityStep{
					Name:   "initialising",
					Status: v1.ActivityStatusTypeRunning,
				},
			},
		})
	}

	if podName != "" {
		pa.Labels["podName"] = podName
	}

	activities.UpdateStatus(pa, false, nil)
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
