package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	exporter "github.com/Rtarun3606k/TakaTime/internal/Exporter"
	"github.com/Rtarun3606k/TakaTime/internal/db"
	"github.com/Rtarun3606k/TakaTime/internal/types"
)

// dateLayout is the format accepted by --from and --to flags.
const dateLayout = "2006-01-02"

func main() {
	// ── Flags ──────────────────────────────────────────────────────────────────
	uri     := flag.String("uri", "", "MongoDB connection URI (falls back to $MONGO_URI)")
	format  := flag.String("format", "csv", "Output format: csv | json")
	fromStr := flag.String("from", "", "Start date inclusive, YYYY-MM-DD (optional)")
	toStr   := flag.String("to", "", "End date inclusive, YYYY-MM-DD (optional)")
	output  := flag.String("output", "", "Output file path (default: stdout)")
	version := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *version {
		fmt.Println(types.Version)
		return
	}

	// ── Validate --format ──────────────────────────────────────────────────────
	if *format != "csv" && *format != "json" {
		fmt.Fprintf(os.Stderr, "error: --format must be 'csv' or 'json', got %q\n", *format)
		os.Exit(1)
	}

	// ── Resolve MongoDB URI ────────────────────────────────────────────────────
	mongoURI := *uri
	if mongoURI == "" {
		mongoURI = os.Getenv("MONGO_URI")
	}
	if mongoURI == "" {
		fmt.Fprintln(os.Stderr, "error: MongoDB URI required — pass --uri or set $MONGO_URI")
		os.Exit(1)
	}

	// ── Parse date flags ───────────────────────────────────────────────────────
	var from, to time.Time
	if *fromStr != "" {
		t, err := time.Parse(dateLayout, *fromStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: --from %q is not a valid date (expected YYYY-MM-DD)\n", *fromStr)
			os.Exit(1)
		}
		from = t
	}
	if *toStr != "" {
		t, err := time.Parse(dateLayout, *toStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: --to %q is not a valid date (expected YYYY-MM-DD)\n", *toStr)
			os.Exit(1)
		}
		to = t
	}
	if !from.IsZero() && !to.IsZero() && to.Before(from) {
		fmt.Fprintln(os.Stderr, "error: --to must not be before --from")
		os.Exit(1)
	}

	// ── Connect to MongoDB ─────────────────────────────────────────────────────
	client, err := db.ConnectToDataBase(mongoURI)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: could not connect to MongoDB: %v\n", err)
		os.Exit(1)
	}
	defer client.Disconnect(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// ── Fetch logs ─────────────────────────────────────────────────────────────
	rows, err := exporter.FetchAllLogs(ctx, client, exporter.FilterOptions{
		From: from,
		To:   to,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to fetch logs: %v\n", err)
		os.Exit(1)
	}

	// ── Open output destination ────────────────────────────────────────────────
	dest := os.Stdout
	if *output != "" {
		dir := filepath.Dir(*output)
		if dir != "." && dir != "" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "error: cannot create output directory %q: %v\n", dir, err)
				os.Exit(1)
			}
		}
		f, err := os.Create(*output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: cannot create output file %q: %v\n", *output, err)
			os.Exit(1)
		}
		defer f.Close()
		dest = f
	}

	// ── Write ──────────────────────────────────────────────────────────────────
	switch *format {
	case "csv":
		if err := exporter.WriteCSV(dest, rows); err != nil {
			fmt.Fprintf(os.Stderr, "error: CSV write failed: %v\n", err)
			os.Exit(1)
		}
	case "json":
		if err := exporter.WriteJSON(dest, rows); err != nil {
			fmt.Fprintf(os.Stderr, "error: JSON write failed: %v\n", err)
			os.Exit(1)
		}
	}

	// ── Done ───────────────────────────────────────────────────────────────────
	if *output != "" {
		fmt.Fprintf(os.Stderr, "✅  Exported %d records → %s\n", len(rows), *output)
	} else {
		log.Printf("exported %d records", len(rows))
	}
}