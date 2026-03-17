package eleven

import (
	"github.com/11notes/go-eleven/container"
	"github.com/11notes/go-eleven/http"
	"github.com/11notes/go-eleven/util"
)

var(
	Container = &container.Container{}
	HTTP = &http.HTTP{}
	Util = &util.Util{}
)

// wrapper for Util.Log
func Log(t string, m string, args ...interface{}){
	if(Util.IfIsNil(args)){
		Util.Log(t, m)
	}else{
		Util.Log(t, m, args...)
	}
}

// wrapper for Util.LogFatal
func LogFatal(m string, args ...interface{}){
	if(Util.IfIsNil(args)){
		Util.LogFatal(m)
	}else{
		Util.LogFatal(m, args...)
	}
}