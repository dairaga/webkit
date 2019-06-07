package webkit

import (
	"github.com/dairaga/config"
	"github.com/dairaga/log"
)

var _maxFileSize int64 = 1024 * 1024

func init() {
	tmp := config.GetInt64("file.max", 1024)
	_maxFileSize = tmp * 1024
	log.Debugf("upload file max size: %d", _maxFileSize)
}
