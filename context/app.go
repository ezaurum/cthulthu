package context

import (
	"errors"
	"github.com/ezaurum/cthulthu/database"
	"github.com/ezaurum/cthulthu/generators"
	"github.com/ezaurum/owlbear"
	"github.com/labstack/echo/v4"
	"reflect"
	"strings"
	"sync"
)

// 어플리케이션 레벨 콘텍스트, 싱글톤
type Application interface {
	AddPersistedResource(interface{}) Resource
	AddAllPersistedResource(...interface{}) []Resource
	SetRepository(repository database.Repository)
	Repository() database.Repository
	SetNodeNumber(number int64)
	NodeNumber() int64
	SetDebug(b bool)
	Debug() bool
	ResourceInterfaces() ([]interface{}, error)
	SetIDGenerators(idGenerators generators.IDGenerators)
	Generators() *generators.IDGenerators
	Router
	InitRoute(*echo.Echo, AssignFunc) error

	// event notifier
	SetEventNotifier(notifierMap *owlbear.NotifierMap)
	Notifier
}

var _ Application = &app{}

type app struct {
	nodeNumber   int64
	IDGenerators generators.IDGenerators
	debug        bool

	repository         database.Repository
	persistedResources []Resource
	router
	eventNotifier *owlbear.NotifierMap
}

func (a *app) SetEventNotifier(notifierMap *owlbear.NotifierMap) {
	a.eventNotifier = notifierMap
}

func (a *app) InitRoute(e *echo.Echo, handlerFunc AssignFunc) error {
	return a.Assign(e, a, handlerFunc)
}

func (a *app) SetNodeNumber(number int64) {
	a.nodeNumber = number
}

func (a *app) NodeNumber() int64 {
	return a.nodeNumber
}

func (a *app) SetDebug(b bool) {
	a.debug = b
}

func (a *app) Debug() bool {
	return a.debug
}

func (a *app) SetIDGenerators(idGenerators generators.IDGenerators) {
	a.IDGenerators = idGenerators
}

func (a *app) Generators() *generators.IDGenerators {
	return &a.IDGenerators
}

var _app Application
var once sync.Once

func App() Application {
	if nil == _app {
		once.Do(func() {
			_app = &app{
				eventNotifier: owlbear.New(),
			}
		})
	}
	return _app
}

func (a *app) SetRepository(repository database.Repository) {
	a.repository = repository
}

func (a *app) Repository() database.Repository {
	return a.repository
}

func (a *app) ResourceInterfaces() ([]interface{}, error) {
	var rr []interface{}
	for _, r := range a.persistedResources {
		if r.Type == nil {
			return nil, ErrResourceInvalid
		}
		rr = append(rr, r.Type)
	}
	return rr, nil
}

func (a *app) AddAllPersistedResource(resourceType ...interface{}) []Resource {
	var rr []Resource
	for _, res := range resourceType {
		rr = append(rr, a.AddPersistedResource(res))
	}
	return rr
}

func (a *app) AddPersistedResource(resourceType interface{}) Resource {
	s := reflect.TypeOf(resourceType).String()
	r := Resource{
		Name:         strings.Replace(s, "*", "", -1),
		Type:         resourceType,
		ResourceType: PersistedResourceType,
	}

	a.persistedResources = append(a.persistedResources, r)
	return r
}

var ErrResourceInvalid = errors.New("resource invalid")

const (
	PersistedResourceType        = "resource.type.persisted"
	HandlerFuncResourceType      = "resource.type.handlerFunc"
	GroupHandlerFuncResourceType = "resource.type.groupHandlerFunc"
)
