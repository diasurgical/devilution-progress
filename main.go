// The devilution-progress tool reports the percentage of binary identical
// functions in the Devilution project.
//
// https://github.com/diasurgical/devilution/milestone/4
package main

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	done, total := progress()
	percent := 100 * float64(done) / float64(total)
	fmt.Printf("progress: %.02f%% (%d/%d)\n", percent, done, total)
}

// progress returns the progress of the closed and open issues in the Binary
// identical functions milestone of the Devilution project.
func progress() (done, total int) {
	for page := 1; ; page++ {
		url := fmt.Sprintf("https://github.com/diasurgical/devilution/milestone/4/paginated_issues?closed=1&page=%d", page)
		closedDone, closedTotal := getProgress(url)
		if closedTotal == 0 {
			break
		}
		done += closedDone
		total += closedTotal
	}
	for page := 1; ; page++ {
		url := fmt.Sprintf("https://github.com/diasurgical/devilution/milestone/4/paginated_issues?page=%d", page)
		openDone, openTotal := getProgress(url)
		if openTotal == 0 {
			break
		}
		done += openDone
		total += openTotal
	}
	return done, total
}

// getProgress returns the total progress of all tasks in the given URL.
func getProgress(url string) (done, total int) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatalf("unable to load page %q; %+v", url, err)
	}
	process := func(i int, sel *goquery.Selection) {
		taskDone, taskTotal := parseProgress(sel.Text())
		done += taskDone
		total += taskTotal
	}
	doc.Find(".task-progress-counts").Each(process)
	return done, total
}

// parseProgress parses the progress of a task.
//
// Example input: 10 of 42
func parseProgress(s string) (done, total int) {
	_, err := fmt.Sscanf(s, "%d of %d", &done, &total)
	if err != nil {
		log.Fatalf("unable to parse task progress; expected format `10 of 42`, got `%s`; %v", s, err)
	}
	return done, total
}
