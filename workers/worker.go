package workers

import (
	"github.com/rqtx/pdoc/utils"
)

type Job interface {
	CreateSection() utils.Section
}

type Worker interface {
	GetSetctions()
}
