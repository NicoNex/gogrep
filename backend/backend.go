package backend

import (
	// "flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sync"

	"github.com/NicoNex/gogrep/ui"
)

type Grep struct {
	root string
	pattern *regexp.Regexp
	wg sync.WaitGroup
	glob string
	allFiles bool
	maxdepth int
	outch chan string
}

func NewGrep() Grep {
	return Grep{
		allFiles: true,
		maxdepth: -1,
	}
}

func (g *Grep) Find(data ui.Data) (chan string, error) {
	fmt.Println(data)
	re, err := regexp.Compile(data.Pattern)
	if err != nil {
		return nil, err
	}

	g.outch = make(chan string)
	g.pattern = re
	g.glob = data.Glob
	g.root = data.Path
	go func(g *Grep) {
		g.walkDir(0)
		close(g.outch)
	}(g)

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
		g.outch <- fpath
	}
}

// Recursively walks in a directory tree.
func (g *Grep) walkDir(depth int) {
	if depth != g.maxdepth {
		files, err := g.readDir(g.root)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, finfo := range files {
			fname := finfo.Name()
			fpath := filepath.Join(g.root, fname)

			if fname[0] != '.' || g.allFiles {
				if finfo.IsDir() {
					g.walkDir(depth+1)
				} else {
					if g.glob == "" || g.matchGlob(fpath) {
						g.wg.Add(1)
						go g.checkMatch(fpath)
					}
				}
			}
		}
	}
}

// func main() {
// 	var pattern string
// 	var files []string

// 	flag.BoolVar(&prnt, "p", false, "Print to stdout.")
// 	flag.BoolVar(&verbose, "v", false, "Verbose, explain what is being done.")
// 	flag.StringVar(&glob, "g", "", "Add a pattern the file names must match to be edited.")
// 	flag.BoolVar(&allFiles, "a", false, "Includes hidden files (starting with a dot).")
// 	flag.IntVar(&maxdepth, "l", -1, "Max depth.")
// 	flag.Usage = usage
// 	flag.Parse()

// 	if flag.NArg() >= 3 {
// 		pattern = flag.Arg(0)
// 		repl = flag.Arg(1)
// 		files = flag.Args()[2:]
// 	} else {
// 		usage()
// 		return
// 	}

// 	regex, err := regexp.Compile(pattern)
// 	if err != nil {
// 		die(err)
// 	}
// 	re = regex

// 	for _, f := range files {
// 		finfo, err := os.Stat(f)
// 		if err != nil {
// 			die(err)
// 		}

// 		if finfo.IsDir() {
// 			walkDir(f, 0)
// 		} else {
// 			wg.Add(1)
// 			go edit(f)
// 		}
// 		wg.Wait()
// 	}
// }