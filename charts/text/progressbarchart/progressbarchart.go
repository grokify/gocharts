package progressbarchart

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/grokify/gocharts/v2/data/histogram"
)

type Tasks []Task

func (tasks Tasks) CurrentCountMax() int {
	max := 0
	for i, task := range tasks {
		if i == 0 {
			max = task.CurrentCount
		} else if task.CurrentCount > max {
			max = task.CurrentCount
		}
	}
	return max
}

// SetMaxCountsMax sets each Task's TotalCount to the sum of all CurrentCounts
func (tasks Tasks) SetMaxCountsMax() {
	maxCount := 0

	// Compute sum of all CurrentCounts
	for i, task := range tasks {
		if i == 0 {
			maxCount = task.CurrentCount
		} else if task.CurrentCount > maxCount {
			maxCount = task.CurrentCount
		}
	}

	// Set each Task's TotalCount
	for i := range tasks {
		tasks[i].MaxCount = maxCount
	}
}

// SetTotalCountsSum sets each Task's TotalCount to the sum of all CurrentCounts
func (tasks Tasks) SetMaxCountsSum() {
	sumCount := 0

	// Compute sum of all CurrentCounts
	for _, task := range tasks {
		sumCount += task.CurrentCount
	}

	// Set each Task's TotalCount
	for i := range tasks {
		tasks[i].MaxCount = sumCount
	}
}

// Task represents a single progress item
type Task struct {
	Label        string
	CurrentCount int
	MaxCount     int
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
	tasks.SetMaxCountsSum()
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
			Label:    task.Label,
			MaxCount: tasksProgress.CurrentCountMax(),
		}
		for j := i; j < len(tasksProgress); j++ {
			taskFunnel.CurrentCount += tasksProgress[j].CurrentCount
		}
		tasksFunnel = append(tasksFunnel, taskFunnel)
	}
	tasksFunnel.SetMaxCountsMax()
	return tasksFunnel
}

// ProgressLine: same as before
func ProgressLine(label string, current, max, maxLabelLength int) string {
	const barWidth = 15

	// Ellipsize if too long
	if len(label) > maxLabelLength {
		if maxLabelLength > 3 {
			label = label[:maxLabelLength-3] + "..."
		} else {
			label = label[:maxLabelLength]
		}
	}

	if max <= 0 {
		return fmt.Sprintf("%-*s |%-*s  N/A (0/0)", maxLabelLength, label, barWidth, "")
	}

	ratio := float64(current) / float64(max)
	if ratio > 1 {
		ratio = 1
	}

	filled := int(ratio * barWidth)
	empty := barWidth - filled

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	percent := int(ratio * 100)

	return fmt.Sprintf("%-*s |%s  %3d%% (%d/%d)", maxLabelLength, label, bar, percent, current, max)
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
		line := ProgressLine(t.Label, t.CurrentCount, t.MaxCount, maxLabelLength)
		sb.WriteString(line)
		if i < len(tasks)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func ChartsTextFromHistogram(h *histogram.Histogram, inclHeader, inclProgress, inclFunnel bool, startNum *int) (string, error) {
	if h == nil {
		return "", errors.New("histogram cannot be nil")
	}

	var sb strings.Builder
	var useNums bool
	curNum := 0
	if startNum != nil {
		useNums = true
		curNum = *startNum
	}
	if inclProgress {
		tasks := NewTasksFromHistogram(h)
		if inclHeader {
			var headerParts []string
			if useNums {
				headerParts = append(headerParts, strconv.Itoa(curNum)+".")
				curNum++
			}
			if name := strings.TrimSpace(h.Name); name != "" {
				headerParts = append(headerParts, name)
			}
			headerParts = append(headerParts, "Progress\n\n")
			if _, err := sb.WriteString(strings.Join(headerParts, " ")); err != nil {
				return "", err
			}
		}
		if _, err := sb.WriteString(tasks.ProgressBarChartText()); err != nil {
			return "", err
		}
	}
	if inclFunnel {
		tasks := NewTasksFunnelFromHistogram(h)
		if inclProgress {
			if _, err := sb.WriteString("\n\n"); err != nil {
				return "", err
			}
		}
		if inclHeader {
			var headerParts []string
			if useNums {
				headerParts = append(headerParts, strconv.Itoa(curNum)+".")
			}
			if name := strings.TrimSpace(h.Name); name != "" {
				headerParts = append(headerParts, name)
			}
			headerParts = append(headerParts, "Funnel\n\n")
			if _, err := sb.WriteString(strings.Join(headerParts, " ")); err != nil {
				return "", err
			}
		}
		if _, err := sb.WriteString(tasks.ProgressBarChartText()); err != nil {
			return "", err
		}
	}

	return sb.String(), nil
}
