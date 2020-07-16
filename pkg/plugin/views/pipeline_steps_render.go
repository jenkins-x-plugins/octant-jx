package views

import (
	"fmt"
	"log"
	"strings"
	"time"

	v1 "github.com/jenkins-x/jx-api/pkg/apis/jenkins.io/v1"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/jenkins-x/octant-jx/pkg/plugin"
	"github.com/jenkins-x/octant-jx/pkg/plugin/util"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	indentation = "  "
)

type PipelineStepRenderer struct {
	Writer           strings.Builder
	PipelineActivity *v1.PipelineActivity
	Pod              *unstructured.Unstructured
}

// ToStepsView renders a markdown description of the pipeline
func ToStepsView(pa *v1.PipelineActivity, pod *unstructured.Unstructured) *component.Text {
	r := &PipelineStepRenderer{
		PipelineActivity: pa,
		Pod:              pod,
	}
	w := &r.Writer
	if pa != nil {
		for _, step := range pa.Spec.Steps {
			if addStepRow(r, &step, "") {
				break
			}
		}
	}
	text := w.String()
	return component.NewMarkdownText(text)
}

func addStepRow(w *PipelineStepRenderer, parent *v1.PipelineActivityStep, indent string) bool {
	stage := parent.Stage
	preview := parent.Preview
	promote := parent.Promote
	pending := false
	if stage != nil {
		if addStageRow(w, stage, indent) {
			pending = true
		}
	} else if preview != nil {
		if addPreviewRow(w, stage, preview, indent) {
			pending = true
		}
	} else if promote != nil {
		if addPromoteRow(w, stage, promote, indent) {
			pending = true
		}
	} else {
		log.Printf("Unknown step kind %#v", parent)
	}
	return pending
}

func addStageRow(w *PipelineStepRenderer, stage *v1.StageActivityStep, indent string) bool {
	name := "Stage"
	if stage.Name != "" {
		name = ""
	}
	pending := addStepRowItem(w, stage, &stage.CoreActivityStep, indent, name, "")

	indent += indentation
	for _, step := range stage.Steps {
		if addStepRowItem(w, stage, &step, indent, "", "") {
			pending = true
			break
		}
	}
	return pending
}

func addPreviewRow(w *PipelineStepRenderer, stage *v1.StageActivityStep, parent *v1.PreviewActivityStep, indent string) bool {
	pullRequestURL := parent.PullRequestURL
	if pullRequestURL == "" {
		pullRequestURL = parent.Environment
	}
	pending := addStepRowItem(w, stage, &parent.CoreActivityStep, indent, "Preview", colorInfo(pullRequestURL))
	indent += indentation

	appURL := parent.ApplicationURL
	if appURL != "" {
		if addStepRowItem(w, stage, &parent.CoreActivityStep, indent, "Preview Application", colorInfo(appURL)) {
			pending = true
		}
	}
	return pending
}

func addPromoteRow(w *PipelineStepRenderer, stage *v1.StageActivityStep, parent *v1.PromoteActivityStep, indent string) bool {
	pending := addStepRowItem(w, stage, &parent.CoreActivityStep, indent, "Promote: "+parent.Environment, "")
	indent += indentation

	pullRequest := parent.PullRequest
	update := parent.Update
	if pullRequest != nil {
		if addStepRowItem(w, stage, &pullRequest.CoreActivityStep, indent, "PullRequest", describePromotePullRequest(pullRequest)) {
			pending = true
		}
	}
	if update != nil {
		if addStepRowItem(w, stage, &update.CoreActivityStep, indent, "Update", describePromoteUpdate(update)) {
			pending = true
		}
	}
	appURL := parent.ApplicationURL
	if appURL != "" {
		if addStepRowItem(w, stage, &update.CoreActivityStep, indent, "Promoted", " Application is at: "+colorInfo(appURL)) {
			pending = true
		}
	}
	return pending
}

func addStepRowItem(w *PipelineStepRenderer, stage *v1.StageActivityStep, step *v1.CoreActivityStep, indent string, name string, description string) bool {
	text := step.Description
	if description != "" {
		if text == "" {
			text = description
		} else {
			text += " " + description
		}
	}
	textName := step.Name
	if textName == "" {
		textName = name
	} else {
		if name != "" {
			textName = name + ":" + textName
		}
	}

	icon := ToPipelineStatusMarkup(step.Status)

	status := ""
	durationText := durationMarkup(step.StartedTimestamp, step.CompletedTimestamp)
	if durationText != "" {
		if status == "" {
			status = " : " + durationText
		} else {
			status += "&nbsp; : " + durationText
		}
	}

	podName := ""
	if w.Pod != nil {
		podName = w.Pod.GetName()
	}

	containerName := FindContainerName(step, w.Pod)
	if containerName != "" {
		paName := w.PipelineActivity.Name
		ns := w.PipelineActivity.Namespace
		if podName != "" {
			status += fmt.Sprintf(`&nbsp;&nbsp;<a href="%s" title="View Step details"><clr-icon shape="details"></clr-icon></a>`, plugin.GetPipelineContainerLink(ns, paName, podName, containerName))
		}
		w.Writer.WriteString(fmt.Sprintf("%s* %s [%s](%s) %s\n", indent, icon, textName, plugin.GetPipelineContainerLogLink(paName, containerName), status))
	} else {
		log.Printf("failed to find container name for step %s and pod %s", step.Name, podName)
		w.Writer.WriteString(fmt.Sprintf("%s* %s %s %s\n", indent, icon, textName, status))
	}
	return step.Status == v1.ActivityStatusTypePending
}

func durationMarkup(start *metav1.Time, end *metav1.Time) string {
	if start == nil || end == nil {
		return ""
	}
	return end.Sub(start.Time).String()
}

func FindContainerName(step *v1.CoreActivityStep, u *unstructured.Unstructured) string {
	if u != nil && step != nil {
		pod, err := viewhelpers.ToPod(u)
		if err != nil {
			log.Printf(fmt.Sprintf("failed to convert to Pod: %s", err.Error()))
			return ""
		}
		names := []string{step.Name}
		return FindContainerNameForStepName(pod, names)
	}
	return ""
}

func FindContainerNameForStepName(pod *corev1.Pod, pipelineActivityStepNames []string) string {
	name := "step-" + strings.ToLower(strings.Join(pipelineActivityStepNames, "-"))
	name = strings.ReplaceAll(name, " ", "-")
	name2 := "step-" + name
	for _, c := range pod.Spec.Containers {
		if c.Name == name || c.Name == name2 {
			return c.Name
		}
	}
	return ""
}

func ToPipelineStatus(pa *v1.PipelineActivity) component.Component {
	if pa == nil || pa.Spec.Status == v1.ActivityStatusTypeNone {
		return component.NewText("")
	}
	return component.NewMarkdownText(ToPipelineStatusMarkup(pa.Spec.Status))
}

func ToPipelineStatusMarkup(statusType v1.ActivityStatusType) string {
	text := statusType.String()
	switch statusType {
	case v1.ActivityStatusTypeFailed, v1.ActivityStatusTypeError:
		return `<clr-icon shape="warning-standard" class="is-solid is-danger" title="Failed"></clr-icon>`
	case v1.ActivityStatusTypeSucceeded:
		return `<clr-icon shape="check-circle" class="is-solid is-success" title="Succeeded"></clr-icon>`
	case v1.ActivityStatusTypePending:
		return `<clr-icon shape="clock" title="Pending"></clr-icon>`
	case v1.ActivityStatusTypeRunning:
		return `<span class="spinner spinner-inline" title="Running"></span>`
	}
	return text
}

func describePromotePullRequest(promote *v1.PromotePullRequestStep) string {
	description := ""
	if promote.PullRequestURL != "" {
		description += " PullRequest: " + colorInfo(promote.PullRequestURL)
	}
	if promote.MergeCommitSHA != "" {
		description += " Merge SHA: " + colorInfo(promote.MergeCommitSHA)
	}
	return description
}

func describePromoteUpdate(promote *v1.PromoteUpdateStep) string {
	description := ""
	for _, status := range promote.Statuses {
		url := status.URL
		state := status.Status

		if url != "" && state != "" {
			description += " Status: " + pullRequestStatusString(state) + " at: " + colorInfo(url)
		}
	}
	return description
}

func colorInfo(text string) string {
	return text
}

func colorError(text string) string {
	return text
}

func pullRequestStatusString(text string) string {
	title := strings.Title(text)
	switch text {
	case "success":
		return colorInfo(title)
	case "error", "failed":
		return colorError(title)
	default:
		return text
	}
}

func timeToString(t *metav1.Time) string {
	if t == nil {
		return ""
	}
	now := &metav1.Time{
		Time: time.Now(),
	}
	return util.DurationString(t, now)
}
