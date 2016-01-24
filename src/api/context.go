package api

import (
	"github.com/golang/glog"
	"models"
)

var EmptyFlags models.StringArray
var EmptyLicenses models.StringArray

func NewContext() *Context {
	_model := &Context{}
	_model.Context = models.NewContext()

	return _model
}

type Context struct {
	*models.Context

	Session   *Session
	RouteName string
}

func (c *Context) LogError(err *models.AppError) {
	// glog.Errorf("", ...)
	glog.Errorf("[%v]\tcode=%v\tmsg=%v\tip=%v\trid=%v\tcid=%v\tdev_msg=%v",
		err.T.Path,
		err.StatusCode,
		err.Message,
		err.T.Ip,
		err.T.RequestId,
		err.T.ClientId,
		err.DevMessage,
	)
}
