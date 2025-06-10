package errorutils

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

//These tools should only be used for debugging.

// blocking! wait for wg to finish with debug logging every 5 seconds. If maxCycles is -1, will wait indefinitely otherwise will terminate after maxCycles. Optional identifier for logging.
func MonitorWaitGroup(wg *sync.WaitGroup, maxCycles int, wgName, id string) {
	logrus.Debugf("WG:%s waiting for wg", id)
	idStr := ""
	if id != "" {
		idStr = " " + id
	}
	wgNameStr := ""
	if wgName != "" {
		wgNameStr = " " + wgName
	}
	var cycles int
	done := make(chan bool)
	go func() {
		wg.Wait()
		close(done)
	}()
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
waitg:
	for {
		select {
		case <-ticker.C:
			logrus.Debugf("WG:%s waiting for wg%s", idStr, wgNameStr)
			if maxCycles != -1 && cycles > maxCycles {
				break waitg
			}
		case <-done:
			break waitg
		}
		cycles++
	}
	logrus.Debugf("WG:waitgroup%sfinished%s", wgNameStr, idStr)
}

var centralB_ch chan activityEvent

// tell the world that the block of code is still running
func ActiveBarker(activity string, targetidORpath string, finish_ch chan struct{}) {
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			logrus.Debugf("BARKER:Running %s %s", activity, targetidORpath)
			centralB_ch <- activityEvent{true, activity, targetidORpath}
		case <-finish_ch:
			logrus.Debugf("BARKER:execution ended for %s with %s", activity, targetidORpath)
			ticker.Stop()
			centralB_ch <- activityEvent{false, activity, targetidORpath}
			return
		}
	}
}

type activityEvent struct {
	toAdd   bool //false removes
	name    string
	details string
}

// CentralBarker receives events and manages the tree of activities, printing the state every interval.
func CentralBarker(stateCh chan<- string, reportInterval time.Duration) {
	root := newNode("/")
	ticker := time.NewTicker(reportInterval)
	defer ticker.Stop()

	for {
		select {
		case ev, ok := <-centralB_ch:
			if !ok {
				// Channel closed, do a final report and return
				stateCh <- fmt.Sprintf("=========== FINAL BARKERS TREE ===========\n%s\n===========================================\n",
					root.PrintTree("", 0))
				return
			}
			root.updateActivityWithPrune(ev.details, ev.name, ev.toAdd)
		case <-ticker.C:
			stateCh <- fmt.Sprintf("=========== ACTIVE BARKERS TREE ===========\n%s\n===========================================\n", root.PrintTree("", 0))
		}
	}
}

type node struct {
	name       string
	activities map[string]struct{}
	children   map[string]*node
}

func newNode(name string) *node {
	return &node{
		name:       name,
		activities: make(map[string]struct{}),
		children:   make(map[string]*node),
	}
}

// Handles both absolute and relative paths.
func splitPathElements(path string) []string {
	clean := filepath.Clean(path)
	clean = strings.TrimPrefix(clean, string(os.PathSeparator))
	if clean == "" {
		return nil
	} // nothing to split —
	return strings.Split(clean, string(os.PathSeparator))
}

func (n *node) updateActivityWithPrune(path, activity string, add bool) {
	parts := splitPathElements(path)
	n.recursiveUpdate(parts, activity, add)
}

func (n *node) recursiveUpdate(parts []string, activity string, add bool) bool {
	if len(parts) == 0 {
		if add {
			n.activities[activity] = struct{}{}
		} else {
			delete(n.activities, activity)
		}
	} else {
		childName := parts[0]
		child, ok := n.children[childName]
		if !ok {
			if add {
				child = newNode(childName)
				n.children[childName] = child
			} else {
				// nothing to remove
				return n.shouldPrune()
			}
		}
		if child.recursiveUpdate(parts[1:], activity, add) {
			delete(n.children, childName)
		}
	}
	return n.shouldPrune()
}

// Returns true if this node should be pruned (no activities and no children)
func (n *node) shouldPrune() bool {
	return len(n.activities) == 0 && len(n.children) == 0
}

// Recursive function to generate the text representation of the tree.
func (n *node) PrintTree(prefix string, depth int) string {
	var sb strings.Builder
	indent := strings.Repeat("  ", depth)
	line := indent + n.name
	if len(n.activities) > 0 {
		acts := make([]string, 0, len(n.activities))
		for a := range n.activities {
			acts = append(acts, a)
		}
		sort.Strings(acts)
		line += " [" + strings.Join(acts, ", ") + "]"
	}
	sb.WriteString(line + "\n")
	// Sort children for consistent output
	kids := make([]string, 0, len(n.children))
	for k := range n.children {
		kids = append(kids, k)
	}
	sort.Strings(kids)
	for _, k := range kids {
		sb.WriteString(n.children[k].PrintTree(prefix, depth+1))
	}
	return sb.String()
}
