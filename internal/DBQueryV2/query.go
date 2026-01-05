package dbqueryv2

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/Rtarun3606k/TakaTime/internal/types"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// generateOutput connects to Mongo, gathers all data, and returns the formatted string.
func GenerateOutput(client *mongo.Client) string {
	var sb strings.Builder

	// 1. Fetch Data
	ctx := context.TODO()
	coll := client.Database("takatime").Collection("logs") // Ensure this matches your DB

	var wg sync.WaitGroup
	var projects, languages []types.Stat
	var totalDetailed float64

	// A. Breakdown (Last 7 Days)
	wg.Add(1)
	go func() {
		defer wg.Done()
		projects, languages, totalDetailed = getBreakdown(ctx, coll, 7)
	}()

	// B. History
	history := make(map[string]time.Duration)
	ranges := []struct {
		Label string
		Days  int
	}{
		{"Yesterday", 1},
		{"Last 7 Days", 7},
		{"Last 30 Days", 30},
		{"All Time", 365},
	}
	var mu sync.Mutex

	for _, r := range ranges {
		wg.Add(1)
		go func(label string, d int) {
			defer wg.Done()
			dur := getTotalDuration(ctx, coll, d)
			mu.Lock()
			history[label] = dur
			mu.Unlock()
		}(r.Label, r.Days)
	}
	wg.Wait()

	// 2. Build GitHub "Dashboard" Markdown 🎨

	// --- Header (Using GitHub Alerts) ---
	start := time.Now().AddDate(0, 0, -7)
	totalDur := time.Duration(totalDetailed) * time.Second
	h := int(totalDur.Hours())
	m := int(totalDur.Minutes()) % 60

	// Blue Box for Title & Date
	fmt.Fprintf(&sb, "> [!NOTE]\n> **TakaTime Dashboard**\n> _%s_ to _%s_\n\n", start.Format("Jan 02"), time.Now().Format("Jan 02"))

	// Green Box for Total Time (Highlights the most important stat)
	fmt.Fprintf(&sb, "> [!TIP]\n> **Total Coding Time (7d):** %dh %dm\n\n", h, m)

	// --- Trends (Formatted Table) ---
	sb.WriteString("#### 📈 Trends\n")
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)

	fmt.Fprintln(w, "| Period\t| Duration\t| Period\t| Duration\t|")
	fmt.Fprintln(w, "| :---\t| :---\t| :---\t| :---\t|")

	fmt.Fprintf(w, "| %s\t| **%s**\t| %s\t| **%s**\t|\n", "Yesterday", formatDuration(history["Yesterday"]), "Last 7 Days", formatDuration(history["Last 7 Days"]))
	fmt.Fprintf(w, "| %s\t| **%s**\t| %s\t| **%s**\t|\n", "Last 30 Days", formatDuration(history["Last 30 Days"]), "All Time", formatDuration(history["All Time"]))

	w.Flush()
	sb.WriteString(buf.String())
	sb.WriteString("\n")

	// --- Languages (Blue Emoji Bars) ---
	sb.WriteString("#### 💻 Languages\n")
	sb.WriteString("| Language | Time | Percentage |\n")
	sb.WriteString("| :--- | :--- | :--- |\n")
	for _, s := range languages {
		if s.Duration > 0 {
			dur := formatDuration(time.Duration(s.Duration) * time.Second)
			bar := generateBar(s.Percent, "🟦")
			fmt.Fprintf(&sb, "| **%s** | %s | %s %.1f%% |\n", s.Name, dur, bar, s.Percent)
		}
	}
	sb.WriteString("\n")

	// --- Projects (Green Emoji Bars) ---
	sb.WriteString("#### 🔥 Projects\n")
	sb.WriteString("| Project | Time | Percentage |\n")
	sb.WriteString("| :--- | :--- | :--- |\n")
	for _, s := range projects {
		if s.Duration > 0 {
			dur := formatDuration(time.Duration(s.Duration) * time.Second)
			bar := generateBar(s.Percent, "🟩")
			fmt.Fprintf(&sb, "| **%s** | %s | %s %.1f%% |\n", s.Name, dur, bar, s.Percent)
		}
	}

	return sb.String()
}

// --- Helpers ---

func generateBar(percent float64, fillIcon string) string {
	const width = 10
	blocks := int((percent / 100) * width)
	if blocks > width {
		blocks = width
	}
	// Use ⬜ for empty space
	return fmt.Sprintf("%s%s", strings.Repeat(fillIcon, blocks), strings.Repeat("⬜", width-blocks))
}

func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}

// --- Helpers (Modified to write to Builder) ---

func writeSection(sb *strings.Builder, title string, stats []types.Stat, color string) {
	fmt.Fprintf(sb, "%s%s%s\n", types.Bold, title, types.Reset)
	for _, s := range stats {
		t := formatDuration(time.Duration(s.Duration) * time.Second)
		bar := generateBar(s.Percent, color)
		fmt.Fprintf(sb, "   %-12s %s%7s%s  %s\n", s.Name, types.Bold, t, types.Reset, bar)
	}
}

func getBreakdown(ctx context.Context, coll *mongo.Collection, days int) ([]types.Stat, []types.Stat, float64) {
	start := time.Now().AddDate(0, 0, -days)
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "timestamp", Value: bson.D{{Key: "$gte", Value: start}}}}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "project", Value: "$project"}, {Key: "language", Value: "$language"}}},
			{Key: "total", Value: bson.D{{Key: "$sum", Value: "$duration"}}},
		}}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, nil, 0 // Handle error gracefully
	}

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, nil, 0
	}

	pMap := make(map[string]float64)
	lMap := make(map[string]float64)
	var total float64

	for _, res := range results {
		dur := res["total"].(float64)
		var proj, lang string

		// --- FIX STARTS HERE ---
		// Safely extract project and language regardless of whether Mongo returns M or D
		switch v := res["_id"].(type) {
		case bson.M:
			proj, _ = v["project"].(string)
			lang, _ = v["language"].(string)
		case bson.D:
			for _, elem := range v {
				if elem.Key == "project" {
					proj, _ = elem.Value.(string)
				}
				if elem.Key == "language" {
					lang, _ = elem.Value.(string)
				}
			}
		case nil:
			continue // Skip if _id is null
		}
		// --- FIX ENDS HERE ---

		if proj == "" {
			proj = "Unknown"
		}
		if lang == "" {
			lang = "Plain Text"
		}

		pMap[proj] += dur
		lMap[strings.ToLower(lang)] += dur
		total += dur
	}
	return processTopN(pMap, total), processTopN(lMap, total), total
}

func getTotalDuration(ctx context.Context, coll *mongo.Collection, days int) time.Duration {
	start := time.Now().AddDate(0, 0, -days)
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "timestamp", Value: bson.D{{Key: "$gte", Value: start}}}}}},
		{{Key: "$group", Value: bson.D{{Key: "_id", Value: nil}, {Key: "total", Value: bson.D{{Key: "$sum", Value: "$duration"}}}}}},
	}
	cursor, _ := coll.Aggregate(ctx, pipeline)
	var results []bson.M
	if cursor.All(ctx, &results) == nil && len(results) > 0 {
		return time.Duration(results[0]["total"].(float64)) * time.Second
	}
	return 0
}

func processTopN(m map[string]float64, total float64) []types.Stat {
	var stats []types.Stat
	for k, v := range m {
		stats = append(stats, types.Stat{Name: k, Duration: v})
	}
	sort.Slice(stats, func(i, j int) bool { return stats[i].Duration > stats[j].Duration })

	if len(stats) > types.TopN {
		var otherDur float64
		for i := types.TopN; i < len(stats); i++ {
			otherDur += stats[i].Duration
		}
		stats = stats[:types.TopN]
		if otherDur > 0 {
			stats = append(stats, types.Stat{Name: "Other", Duration: otherDur})
		}
	}
	for i := range stats {
		if total > 0 {
			stats[i].Percent = (stats[i].Duration / total) * 100
		}
	}
	return stats
}
