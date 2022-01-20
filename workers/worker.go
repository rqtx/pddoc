package workers

import (
	"github.com/rqtx/pddoc/utils"
)

type Job interface {
	CreateSection() utils.Section
}

type Worker interface {
	GetSetctions()
}
