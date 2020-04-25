package backend

import (
	"fmt"
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
	Outch chan string
	sem chan int
	Stop chan int
}

func NewGrep() Grep {
	return Grep{
		allFiles: true,
		maxdepth: -1,
		sem: make(chan int, 1024),
	}
}

func (g *Grep) Find(data frontend.Data) error {
	re, err := regexp.Compile(data.Pattern)
	if err != nil {
		return err
	}

	g.Outch = make(chan string)
	g.Stop = make(chan int, 1)
	g.pattern = re
	g.glob = data.Glob

	go func(path string, g *Grep) {
		g.walkDir(path, 0)
		g.wg.Wait()
		close(g.Outch)
		close(g.Stop)
	}(data.Path, g)

	return nil
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
		fmt.Println(err)
	}

	return ok
}

func (g *Grep) checkMatch(fpath string) {
	defer g.wg.Done()

	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		fmt.Println(err)
		return
	}

	if g.pattern.Match(b) {
		g.Outch <- fpath
	}
	<-g.sem
}

// FIXME: make the goroutine exit properly.
// Recursively walks in a directory tree.
func (g *Grep) walkDir(root string, depth int) {
	select {
	case <-g.Stop:
		runtime.Goexit()

	default:
		if depth != g.maxdepth {
			files, err := g.readDir(root)
			if err != nil {
				g.Outch <- err.Error()
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
							g.sem <- 1
							go g.checkMatch(fpath)
						}
					}
				}
			}
		}
	}
}
