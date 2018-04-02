package route

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"path/filepath"
)

func SetStaticFile(r *gin.Engine) func(path string, info os.FileInfo, err error) error {
	return func(path string, info os.FileInfo, err error) error {
		if nil != err {
			log.Printf("err before %v, %v", path, err)
		}

		// 디렉토리는 패스
		if info.IsDir() {
			return nil
		}

		//TODO ignored files?

		base := filepath.Base(path)
		r.StaticFile(base, path)
		return nil
	}
}

//TODO
func InitStaticFiles(r *gin.Engine, staticDir string) {
	// static
	//TODO 디렉토리 여러 군데서 찾도록 하는 것도 필요
	//TODO skin 시스템을 상속가능하도록 하려면...
	r.Static("/images", staticDir+"/images")
	r.Static("/js", staticDir+"/js")
	r.Static("/fonts", staticDir+"/fonts")
	r.Static("/css", staticDir+"/css")
	filepath.Walk(staticDir, SetStaticFile(r))
}
