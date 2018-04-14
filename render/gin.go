package render

import (
	"github.com/ezaurum/cthulthu/render/boongeoppang"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"log"
	"path/filepath"
)

//check implementation
var _ render.HTMLRender = Render{}

func Default(engine *gin.Engine) Render {
	return New(boongeoppang.DefaultTemplateDir, engine)
}

// New instance
func New(templateDir string, engine *gin.Engine) Render {

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
					if ev.Op&fsnotify.Create == fsnotify.Create &&
						".tmpl" == filepath.Ext(ev.Name) {
						i.templateContainer = boongeoppang.Load(templateDir)
						//TODO 나중에 처널로 어떻게...
						engine.HTMLRender = i
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
		Template: layout.Layout,
		Data:     data,
	}
}
