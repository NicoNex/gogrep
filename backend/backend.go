package backend

import (
	// "fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sync"

	"github.com/NicoNex/gogrep/frontend"
)

type Grep struct {
	pattern *regexp.Regexp
	wg sync.WaitGroup
	glob string
	allFiles bool
	maxdepth int
	outch chan string
	sem chan bool
	Stop chan bool
}

func NewGrep() Grep {
	return Grep{
		allFiles: true,
		maxdepth: -1,
		sem: make(chan bool, 1024),
		Stop: make(chan bool, 1),
	}
}

func (g *Grep) Find(data frontend.Data) (chan string, error) {
	re, err := regexp.Compile(data.Pattern)
	if err != nil {
		return nil, err
	}

	g.outch = make(chan string)
	g.pattern = re
	g.glob = data.Glob

	go func(path string, g *Grep) {
		g.walkDir(path, 0)
		g.wg.Wait()
		close(g.outch)
	}(data.Path, g)

	return g.outch, nil
}

func (g *Grep) readDir(filename string) ([]os.FileInfo, error) {
	file, err := os.Open(filename)
	if err != nil {
		return []os.FileInfo{}, err
	}
	defer file.Close()
	return file.Readdir(-1)
}

func (g *Grep) matchGlob(fname string) bool {
	if runtime.GOOS == "windows" {
		re := regexp.MustCompile(`[^\\]+$`)
		fname = re.FindString(fname)
	}

	ok, err := filepath.Match(g.glob, fname)
	if err != nil {
		g.outch <- err.Error()
	}

	return ok
}

func (g *Grep) checkMatch(fpath string) {
	defer g.wg.Done()

	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		g.outch <- err.Error()
		return
	}

	if g.pattern.Match(b) {
		g.outch <- fpath
	}
	<-g.sem
}

// Recursively walks in a directory tree.
func (g *Grep) walkDir(root string, depth int) {
	select {
	case <-g.Stop:
		return

	default:
		if depth != g.maxdepth {
			files, err := g.readDir(root)
			if err != nil {
				g.outch <- err.Error()
				return
			}

			for _, finfo := range files {
				fname := finfo.Name()
				fpath := filepath.Join(root, fname)

				if fname[0] != '.' || g.allFiles {
					if finfo.IsDir() {
						g.walkDir(fpath, depth+1)
					} else {
						if g.glob == "" || g.matchGlob(fname) {
							g.wg.Add(1)
							g.sem <- true
							go g.checkMatch(fpath)
						}
					}
				}
			}
		}
	}
}
