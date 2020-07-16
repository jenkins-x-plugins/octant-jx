package views

import (
	"fmt"
	"sort"
	"strings"
	"time"

	v1 "github.com/jenkins-x/jx-api/pkg/apis/jenkins.io/v1"
	"github.com/jenkins-x/octant-jx/pkg/common/viewhelpers"
	"github.com/jenkins-x/octant-jx/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func ToStatus(pa *v1.PipelineActivity) *component.Text {
	s := &pa.Spec
	return component.NewText(string(s.Status))
}

func ToPipelineOwner(pa *v1.PipelineActivity) component.Component {
	return component.NewText(pa.Spec.GitOwner)
}

func ToPipelineRepository(pa *v1.PipelineActivity) component.Component {
	/*
			owner := s.GitOwner
		repository := s.GitRepository

		if s.GitURL != "" {
			repository = fmt.Sprintf("[%s](%s)", repository, s.GitURL)
		}

	*/
	return component.NewText(pa.Spec.GitRepository)
}

func ToPipelineName(pa *v1.PipelineActivity) component.Component {
	s := &pa.Spec
	name := s.Build
	if s.Context != "" {
		name += " " + s.Context
	}
	md := fmt.Sprintf("[#%s](%s)", name, plugin.GetPipelineLink(pa.Name))
	return component.NewMarkdownText(md)
}

func ToRepository(pa *v1.PipelineActivity) component.Component {
	return component.NewMarkdownText(ToRepositoryMarkdown(pa))
}

func ToRepositoryMarkdown(pa *v1.PipelineActivity) string {
	s := &pa.Spec
	owner := s.GitOwner
	repository := s.GitRepository
	repoLink := s.GitURL

	answer := viewhelpers.ToOwnerRepositoryLinkMarkdown(owner, repository, repoLink)
	branch := strings.ToLower(s.GitBranch)
	if repoLink != "" && strings.HasPrefix(branch, "pr-") {
		prNumber := strings.TrimPrefix(branch, "pr-")
		if prNumber != "" {
			prLink := strings.TrimSuffix(repoLink, ".git")
			prLink = strings.TrimSuffix(prLink, "/")
			prLink += "/pull/" + prNumber
			answer += fmt.Sprintf(` / PR <a href="%s" title="Pull Request #%s">#%s</a>`, prLink, prNumber, prNumber)
		}
	}
	return answer
}

// ToNameMarkdown returns the markdown of the name
func ToNameMarkdown(pa *v1.PipelineActivity) string {
	s := &pa.Spec
	branchAndBuild := fmt.Sprintf("%s #%s", s.GitBranch, s.Build)
	buildLabel := branchAndBuild
	if s.Context != "" {
		buildLabel += " " + s.Context
	}
	return buildLabel
}

func ToPipelineLastStepStatus(pa *v1.PipelineActivity, addContext bool, addTimestamp bool) component.Component {
	status := ""
	if pa != nil && pa.Spec.Status != v1.ActivityStatusTypeNone {
		status = ToPipelineStatusMarkup(pa.Spec.Status)
	}
	lastStep := ToLastStepMarkdown(pa)
	if status != "" {
		lastStep = status + " " + lastStep
	}
	context := pa.Spec.Context
	if addContext && context != "" {
		lastStep += " : " + context
	}
	if addTimestamp {
		durationText := pipelineDuration(pa)
		if durationText != "" {
			title := fmt.Sprintf("Started at %s", viewhelpers.ToTimeMarkdown(pa.Spec.StartedTimestamp))
			lastStep += fmt.Sprintf(`&nbsp;&nbsp;<span class="badge" title="%s">%s</span>`, title, durationText)
		}
	}
	return component.NewMarkdownText(lastStep)
}

func pipelineDuration(pa *v1.PipelineActivity) string {
	start := pa.Spec.StartedTimestamp
	if start == nil {
		return ""
	}
	if pa.Spec.CompletedTimestamp != nil {
		return durationMarkup(start, pa.Spec.CompletedTimestamp)
	}
	u := start.Time
	return viewhelpers.ToDurationString(u)
}

// ToPipelineStartCompleteTimeMarkdown returns the time the pipeline started and completed
func ToPipelineStartCompleteTimeMarkdown(pa *v1.PipelineActivity) string {
	s := &pa.Spec
	if s.StartedTimestamp == nil {
		return ""
	}
	from := viewhelpers.ToDurationMarkdown(s.StartedTimestamp.Time, "started at ")
	if s.CompletedTimestamp == nil {
		return from
	}
	to := viewhelpers.ToDurationMarkdown(s.CompletedTimestamp.Time, "completed at ")
	return from + `&nbsp;<clr-icon shape="minus"></clr-icon>&nbsp;` + to
}

// ToLastStepMarkdown returns the  running step
func ToLastStepMarkdown(pa *v1.PipelineActivity) string {
	s := &pa.Spec
	steps := s.Steps
	if len(steps) > 0 {
		step := steps[len(steps)-1]
		st := step.Stage
		if st != nil {
			ssteps := st.Steps
			for i := len(ssteps) - 1; i >= 0; i-- {
				ss := ssteps[i]
				if ss.Status == v1.ActivityStatusTypePending && i > 0 {
					continue
				}
				return ss.Name
			}
			return st.Name
		}
		promote := step.Promote
		if promote != nil {
			pr := promote.PullRequest
			prURL := pr.PullRequestURL
			if pr != nil && prURL != "" {
				prName := "PR"
				i := strings.LastIndex(prURL, "/")
				if i > 0 && i < len(prURL) {
					prName = prURL[i+1:]
				}
				title := pr.Name
				if promote.Environment != "" {
					title = fmt.Sprintf("Promote to %s", strings.Title(promote.Environment))
				}
				return fmt.Sprintf(`<span>%s <a href="%s">#%s</a></span>`, title, prURL, prName)
			}
			return promote.Name
		}
		preview := step.Preview
		if preview != nil {
			if preview.ApplicationURL != "" {
				title := preview.Name
				if title == "" {
					title = "Preview"
				}
				return fmt.Sprintf(`<span>Promote <a href="%s">%s</a></span>`, preview.ApplicationURL, title)
			}
			return preview.Name
		}
		return st.Name
	}
	return ""
}

func ToDuration(pa *v1.PipelineActivity) component.Component {
	s := &pa.Spec
	start := s.StartedTimestamp
	if start == nil {
		return component.NewText("")
	}
	t := start.Time

	complete := s.CompletedTimestamp
	if complete != nil {
		d := complete.Time.Sub(t)
		t = time.Now().Add(d * -1)
	}
	return component.NewTimestamp(t)
}

// SortPipelines sorts pipelines in name order then with newest build first
func SortPipelines(resources []*v1.PipelineActivity) {
	sort.Slice(resources, func(i, j int) bool {
		r1 := resources[i]
		r2 := resources[j]
		if strings.Compare(r1.Name, r2.Name) < 0 {
			return true
		}

		b1 := viewhelpers.PipelineBuildNumber(r1)
		b2 := viewhelpers.PipelineBuildNumber(r2)
		if b1 < b2 {
			return true
		}
		return viewhelpers.ResourceTimeLessThan(r1.Spec.StartedTimestamp, r2.Spec.StartedTimestamp)
	})
}
