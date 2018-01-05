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
