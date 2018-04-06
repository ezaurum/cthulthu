package render

import (
	"github.com/ezaurum/cthulthu/render/boongeoppang"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"html/template"
	"log"
	"fmt"
)

//check implementation
var _ render.HTMLRender = Render{}

func Default() Render {
	return New(boongeoppang.DefaultTemplateDir)
}

// New instance
func New(templateDir string) Render {

	var b *boongeoppang.TemplateContainer
	b = boongeoppang.Load(templateDir)
	i := Render{
		templateContainer: b,
	}

	if gin.IsDebugging() {

		boongeoppang.WatchDir(templateDir, func(watcher *fsnotify.Watcher) {
			for {
				select {
				case ev := <-watcher.Events:
					if ev.Op != 0 {
						i.templateContainer = boongeoppang.Load(templateDir)
						fmt.Println("reload remplate")
						fmt.Println("reload remplate")
						fmt.Println("reload remplate")
					}
				case err := <-watcher.Errors:
					log.Fatal("error:", err)
				}
			}
		})
	}

	return i
}

type Render struct {
	templateContainer *boongeoppang.TemplateContainer
}

// Instance find by name
func (r Render) Instance(name string, data interface{}) render.Render {
	layout, isExist := r.templateContainer.Get(name)
	if !isExist {
		panic("not exist template " + name)
	}
	return render.HTML{
		Template: layout.Layout.(*template.Template),
		Data:     data,
	}
}
