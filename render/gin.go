package render

import (
	"github.com/ezaurum/cthulthu/render/boongeoppang"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"html/template"
)

//check implementation
var _ render.HTMLRender = Render{}

func Default() Render {
	return New(boongeoppang.DefaultTemplateDir)
}

// New instance
func New(templateDir string) Render {

	var b *boongeoppang.TemplateContainer
	if gin.IsDebugging() {
		b = boongeoppang.LoadDebug(templateDir)
	} else {
		b = boongeoppang.Load(templateDir)
	}

	return Render{
		templateContainer: b,
	}
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
