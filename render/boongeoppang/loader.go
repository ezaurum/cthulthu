package boongeoppang

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"github.com/fsnotify/fsnotify"
)

const (
	baseOf             = "baseof"
	defaultDir         = "_default"
	partialsDir        = "_partials"
	DefaultTemplateDir = "templates"
)

var (
	EmptyLayoutHolder = LayoutHolder{}
)

type LayoutHolder struct {
	Path   string
	Layout interface{}
	Name   string
}

type TemplateContainer struct {
	M        map[string]*LayoutHolder
	Partials map[string]string
	Defaults map[string]string
	Populate func(files ...string) interface{}
	debug    bool
}

func (t TemplateContainer) Set(name string, layout interface{}) {
	get, _ := t.Get(name)
	get.Layout = layout
}

func (t TemplateContainer) Get(name string) (*LayoutHolder, bool) {
	if r, b := t.M[name]; b {
		return r, true
	}

	baseName := path.Base(name)

	if mm, b := t.Defaults[baseName]; b && baseName != baseOf {
		t.M[name] = &LayoutHolder{
			Name: name,
			Path: mm,
		}
		return t.M[name], true
	}

	return nil, false
}

func Default() *TemplateContainer {
	partials := make(map[string]string)
	defaults := make(map[string]string)
	return &TemplateContainer{
		Partials: partials,
		Defaults: defaults,
		M:        make(map[string]*LayoutHolder),
		Populate: populateHtmlTemplate,
	}
}

func LoadDebug(rootDir string) *TemplateContainer {
	d := Default()
	d.debug = true

	load := d.Load(rootDir)

	WatchDir(rootDir, func(watcher *fsnotify.Watcher) {
		fmt.Printf("watch dirs\n")
		for {
			select {
			case ev := <-watcher.Events:
				fmt.Printf("event %v\n", ev)
				if ev.Op != 0 {
					fmt.Println("reload remplate")
					fmt.Println("reload remplate")
					fmt.Println("reload remplate")
					load = d.Load(rootDir)
				}
			case err := <-watcher.Errors:
				log.Fatal("error:", err)
			}
		}
	})

	return load
}

func DefaultLoad() *TemplateContainer {
	return Default().Load(DefaultTemplateDir)
}

func Load(rootDir string) *TemplateContainer {
	return Default().Load(rootDir)
}

func (t *TemplateContainer) Load(rootDir string) *TemplateContainer {
	Defaults := t.Defaults
	holders := t.M

	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if nil != err {
			log.Printf("err before %v, %v", path, err)
		}

		// 디렉토리는 패스
		if info.IsDir() {
			return nil
		}

		filename := info.Name()
		layoutName := strings.TrimSuffix(filename, filepath.Ext(filename))
		if layoutName == "" {
			return fmt.Errorf("file name is empty %v, %v", path, info)
		}

		contentName := filepath.Base(filepath.Dir(path))
		templateKey := contentName + "/" + layoutName
		switch contentName {
		case "":
			return fmt.Errorf("file name is empty %v, %v", path, info)
		case partialsDir:
			t.Partials[layoutName] = path
			break
		case defaultDir:
			Defaults[layoutName] = path
			templateKey = layoutName

			fallthrough
		default:
			holders[templateKey] = &LayoutHolder{
				Name: layoutName,
				Path: path,
			}
			break
		}

		if t.debug {
			log.Printf("key:%v, name:%v, path:%v", templateKey, layoutName, path)
		}

		return err
	})

	t.initiateTemplates()

	return t
}

// initiate html/template
func (t *TemplateContainer) initiateTemplates() {
	var partialsFileNames []string
	for _, v := range t.Partials {
		partialsFileNames = append(partialsFileNames, v)
	}
	// _default/baseof 먼제 체크
	base, isBase := t.Defaults["baseof"]
	for key, value := range t.M {
		layoutName := value.Name

		// 목록의 제일 처음이 기본 템플릿이 된다.
		var files []string

		// 1. baseof
		if isBase {
			files = append(files, base)
		}

		// 2. _default/layout
		if layoutName != key {
			ln, e := t.Defaults[layoutName]
			if e && len(ln) > 0 {
				files = append(files, ln)
			}
		}

		// 3. domain/layout - path
		files = append(files, value.Path)

		// partials added after first object
		if len(files) > 1 {
			files = append(files[:1], append(partialsFileNames, files[1:]...)...)
		} else {
			files = append(files[:1], partialsFileNames...)
		}

		value.Layout = t.Populate(files...)
	}
}

func populateHtmlTemplate(files ...string) interface{} {
	return template.Must(template.ParseFiles(files...))
}
func (t *TemplateContainer) IsDebug() bool {
	return t.debug
}
