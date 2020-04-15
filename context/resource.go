package context

/* # 어플리케이션 리소스

어플리케이션에서 사용되는 리소스들, 권한 관리라든가 여러 API 사용할 때 필요하다
*/
type Resource struct {
	Name         string
	Type         interface{}
	ResourceType string
}

type HandlerFuncResource struct {
	Resource
	Method      string
	Path        string
	HandlerFunc RequestHandlerFunc
}
