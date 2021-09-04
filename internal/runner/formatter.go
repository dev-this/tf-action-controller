package runner

import (
	"fmt"
	"time"
)

func workOutPrefixEmoji(isCompleted bool, isSuccessful bool) string {
	if isCompleted && isSuccessful {
		return ":heavy_check_mark:"
	}

	if isCompleted && !isSuccessful {
		return ":x:"
	}

	// in progress
	return ":thinking:"
}

func FormatSections(sections ...*Execution) string {
	details := ""
	shouldCollapse := len(sections) > 1

	if shouldCollapse {
		for i, section := range sections {
			prefix := workOutPrefixEmoji(section.Completed, section.Successful)
			if i == 0 {
				details += fmt.Sprintf("## Terraform init\n\nFinished at `%s` after %d seconds\n\n<details><summary> ", section.StartedAt.Format(time.RFC850), section.GetSecondsRanFor()) + prefix + " View detailed logs</summary>\n\n```hcl\n\n" + section.Details + "\n\n```\n\n</details>\n"
				continue
			}

			details += fmt.Sprintf("\n---\n\n## Terraform plan\n\nFinished at `%s` after %d seconds\n\n<details open><summary> ", section.StartedAt.Format(time.RFC850), section.GetSecondsRanFor()) + prefix + " View detailed logs</summary>\n\n```hcl\n\n" + section.Details + "\n\n```\n\n</details>\n"
		}
	}

	if !shouldCollapse && len(sections) == 1 {
		prefix := workOutPrefixEmoji(sections[0].Completed, sections[0].Successful)
		details += "## Terraform init\n\nCompleted in X, finishing at Y\n\n<details open><summary> " + prefix + " View detailed logs</summary>\n\n```hcl\n\n" + sections[0].Details + "\n\n```\n\n</details>"
	}

	return details
}
