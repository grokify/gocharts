package progressbarchart

import (
	"fmt"
	"strings"

	"github.com/grokify/gocharts/v2/data/histogram"
)

type Tasks []Task

// SetTotalCountsSum sets each Task's TotalCount to the sum of all CurrentCounts
func (tasks Tasks) SetTotalCountsSum() {
	totalCount := 0

	// Compute sum of all CurrentCounts
	for _, task := range tasks {
		totalCount += task.CurrentCount
	}

	// Set each Task's TotalCount
	for i := range tasks {
		tasks[i].TotalCount = totalCount
	}
}

// Task represents a single progress item
type Task struct {
	Label        string
	CurrentCount int
	TotalCount   int
}

func NewTasksFromHistogram(h *histogram.Histogram) Tasks {
	tasks := Tasks{}
	if h == nil {
		return tasks
	}
	itemNames := h.ItemNamesOrderOrDefault()
	for _, itemName := range itemNames {
		count, ok := h.Items[itemName]
		if !ok {
			count = 0
		}
		tasks = append(tasks, Task{Label: itemName, CurrentCount: count})
	}
	tasks.SetTotalCountsSum()
	return tasks
}

func NewTasksFunnelFromHistogram(h *histogram.Histogram) Tasks {
	tasksFunnel := Tasks{}
	if h == nil {
		return tasksFunnel
	}
	tasksProgress := NewTasksFromHistogram(h)
	for i, task := range tasksProgress {
		taskFunnel := Task{
			Label:      task.Label,
			TotalCount: task.TotalCount,
		}
		for j := i; j < len(tasksProgress); j++ {
			taskFunnel.CurrentCount += tasksProgress[j].CurrentCount
		}
		tasksFunnel = append(tasksFunnel, taskFunnel)
	}
	return tasksFunnel
}

// ProgressLine: same as before
func ProgressLine(label string, current, total, maxLabelLength int) string {
	const barWidth = 15

	// Ellipsize if too long
	if len(label) > maxLabelLength {
		if maxLabelLength > 3 {
			label = label[:maxLabelLength-3] + "..."
		} else {
			label = label[:maxLabelLength]
		}
	}

	if total <= 0 {
		return fmt.Sprintf("%-*s |%-*s  N/A (0/0)", maxLabelLength, label, barWidth, "")
	}

	ratio := float64(current) / float64(total)
	if ratio > 1 {
		ratio = 1
	}

	filled := int(ratio * barWidth)
	empty := barWidth - filled

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	percent := int(ratio * 100)

	return fmt.Sprintf("%-*s |%s  %3d%% (%d/%d)", maxLabelLength, label, bar, percent, current, total)
}

func (tasks Tasks) ProgressBarChartText() string {
	// Determine maxLabelLength
	maxLabelLength := 0
	for _, t := range tasks {
		if len(t.Label) > maxLabelLength {
			maxLabelLength = len(t.Label)
		}
	}

	var sb strings.Builder
	for i, t := range tasks {
		line := ProgressLine(t.Label, t.CurrentCount, t.TotalCount, maxLabelLength)
		sb.WriteString(line)
		if i < len(tasks)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
