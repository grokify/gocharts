package slack

import (
	"time"

	"github.com/grokify/gocharts/v2/data/roadmap2"
	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/mogo/time/timeutil"
)

// GetRoadmapExample returns a roamdap copied from public Slack's Platform roadmap
// posted on Trello as of August 20, 2022. The Trello roadmap is available here:
// https://trello.com/b/ZnTQyumQ/slack-platform-roadmap
func GetRoadmapExample() roadmap2.Roadmap {
	qNow := timeutil.QuarterEnd(time.Now().UTC())
	qPrev1 := timeutil.QuarterAdd(qNow, -1)
	qNext1 := timeutil.QuarterAdd(qNow, 1)
	qNext2 := timeutil.QuarterAdd(qNow, 2)

	initiativeAPIAdmin := "API & Administration"
	initiativeAUX := "App User Experience"

	rm := roadmap2.Roadmap{
		Name: "Slack Platform Roadmap",
		Columns: table.Columns{
			"",
			timeutil.FormatQuarterYYQ(qPrev1),
			timeutil.FormatQuarterYYQ(qNow),
			timeutil.FormatQuarterYYQ(qNext1),
			timeutil.FormatQuarterYYQ(qNext2),
		},
		ItemCellFunc: func(i roadmap2.Item) (colIdx int) {
			rmStart := timeutil.QuarterStart(time.Now())
			rmStart = timeutil.QuarterAdd(rmStart, -1)
			dt := i.ReleaseTime
			if dt.Before(rmStart) {
				return -1
			}
			if dt.Before(timeutil.QuarterAdd(rmStart, 1)) {
				return 0
			} else if dt.Before(timeutil.QuarterAdd(rmStart, 2)) {
				return 1
			} else if dt.Before(timeutil.QuarterAdd(rmStart, 3)) {
				return 2
			} else if dt.Before(timeutil.QuarterAdd(rmStart, 5)) {
				return 3
			} else {
				return -1
			}
		},
		StreamNames: []string{initiativeAUX, initiativeAPIAdmin},
		Items: []roadmap2.Item{
			{
				Name:        "Deprecating support for TLS versions 1.1 and earlier",
				Description: `Support for TLS versions 1.2 and later only, resulting in error for Slack services running on TLS versions 1.1 and earlier. This upgrade aligns with industry best practices around security and data integrity.\n\nTLS specific Help Center article coming soon. Refer to "Minimum Requirements for using Slack" in the interim: https://get.slack.help/hc/en-us/articles/115002037526-Minimum-requirements-for-using-Slack`,
				ReleaseTime: qPrev1,
				StreamName:  initiativeAPIAdmin,
			},
			{
				Name:        "Better access to apps in Slack",
				Description: `More ways users can discover your app in Slack through the shortcuts button. Check out the lightning bolt icon to the left of your message composer.\n\nMore details here: https://medium.com/slack-developer-blog/introducing-new-ways-to-interact-with-apps-d66e160b8ae\n\nShortcuts Overview: https://api.slack.com/interactivity/shortcuts`,
				ReleaseTime: qPrev1,
				StreamName:  initiativeAUX,
			},

			{
				Name:        "Persistent link previews",
				Description: "App unfurls available everywhere – private channels, uninstalled apps, etc.",
				ReleaseTime: qNow,
				StreamName:  initiativeAUX,
			},
			{
				Name:        "Sign in with Slack links",
				Description: "When a user clicks on a link from your domain, provision new accounts or map existing accounts for your service using Slack’s verified identity data.",
				ReleaseTime: qNow,
				StreamName:  initiativeAPIAdmin,
			},
			{
				Name:        "App Manifests",
				Description: "A faster way to start building and managing Slack apps. Start using pre-defined manifests, now in beta, with our new guided tutorials. https://api.slack.com/tutorials. To learn more about how to create your own manifests, check out: https://api.slack.com/reference/manifests",
				ReleaseTime: qNow,
				StreamName:  initiativeAPIAdmin,
			},
			{
				Name:        "Message metadata",
				Description: "The ability for apps to provide metadata about system events from their service and send it to Slack along with the messages posted into Slack. Metadata converts unstructured messages into structured events that customers can leverage to automate their work.\n\nNow in closed beta. Learn more here: https://api.slack.com/future/metadata",
				ReleaseTime: qNow,
				StreamName:  initiativeAPIAdmin,
			},

			{
				Name:        "Subscribe in Slack",
				Description: "Users can subscribe a channel to receive custom notifications about a specific resource without needing to invite an app to the channel.",
				ReleaseTime: qNext1,
				StreamName:  initiativeAUX,
			},
			{
				Name:        "Functional building blocks",
				Description: "Components that take actions in your workflow, offered as Slack-native functions and from other tools. Now available for developers to build in closed beta. Learn more here: https://api.slack.com/future/integrations",
				ReleaseTime: qNext1,
				StreamName:  initiativeAPIAdmin,
			},
			{
				Name:        "Slack Command Line Interface (CLI)",
				Description: "A tool for devs to quickly create, install and host an app using standard configurations, including local dev mode. Now in closed beta. Learn more here: https://api.slack.com/future/tools",
				ReleaseTime: qNext1,
				StreamName:  initiativeAPIAdmin,
			},
			{
				Name:        "Slack-first app deployment",
				Description: "Integrated hosting of an app and its data, ensuring Slack-level security and compliance.",
				ReleaseTime: qNext1,
				StreamName:  initiativeAPIAdmin,
			},
			{
				Name:        "Platform administration",
				Description: "Tools for additional integration management APIs, analytics on integration usage, and integration management automation.",
				ReleaseTime: qNext1,
				StreamName:  initiativeAPIAdmin,
			},

			{
				Name:        "Shortcut gallery",
				Description: "A new Slack surface for discovering solution-oriented workflows tailored to your organization, including Slack-native workflows and functions, too.",
				ReleaseTime: qNext2,
				StreamName:  initiativeAPIAdmin,
			},
			{
				Name:        "New Workflow Builder experience",
				Description: "Unlock new use cases with more steps, triggers, conditional logic, and multi-channel workflows.",
				ReleaseTime: qNext2,
				StreamName:  initiativeAPIAdmin,
			},
			{
				Name:        "Easy workflow discoverability",
				Description: "Easily see all your team’s workflows in a channel, pinned to bookmarks, etc.",
				ReleaseTime: qNext2,
				StreamName:  initiativeAPIAdmin,
			},
			{
				Name:        "Easy workflow set up",
				Description: "Customize and publish workflows to a channel without navigating to Workflow Builder.",
				ReleaseTime: qNext2,
				StreamName:  initiativeAPIAdmin,
			},
		},
	}
	return rm
}
