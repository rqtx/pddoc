package aws

import (
	"github.com/rqtx/pddoc/utils"
	"github.com/rqtx/pddoc/workers"
	"github.com/rqtx/pddoc/workers/aws/jobs/account"
	"github.com/rqtx/pddoc/workers/aws/jobs/cloudwatch"
	"github.com/rqtx/pddoc/workers/aws/jobs/ec2"
	"github.com/rqtx/pddoc/workers/aws/jobs/ecr"
	"github.com/rqtx/pddoc/workers/aws/jobs/ecs"
	"github.com/rqtx/pddoc/workers/aws/jobs/iam"
	"github.com/rqtx/pddoc/workers/aws/jobs/kms"
	"github.com/rqtx/pddoc/workers/aws/jobs/lambda"
	"github.com/rqtx/pddoc/workers/aws/jobs/storage"
	"github.com/rqtx/pddoc/workers/aws/jobs/vpc"
	"github.com/rqtx/pddoc/workers/aws/jobs/waf"
)

type WorkerAWS struct {
	jobs []workers.Job
}

func NewWorker(region string) *WorkerAWS {
	jobs := make([]workers.Job, 0)
	jobs = append(jobs, account.New(region))
	jobs = append(jobs, iam.New(region))
	jobs = append(jobs, vpc.New(region))
	jobs = append(jobs, kms.New(region))
	jobs = append(jobs, ecs.New(region))
	jobs = append(jobs, ec2.New(region))
	jobs = append(jobs, ecr.New(region))
	jobs = append(jobs, storage.New(region))
	jobs = append(jobs, cloudwatch.New(region))
	jobs = append(jobs, lambda.New(region))
	jobs = append(jobs, waf.New(region))
	return &WorkerAWS{
		jobs: jobs,
	}
}

func (worker *WorkerAWS) GetSetctions() []utils.Section {
	secs := make([]utils.Section, 0)
	for _, job := range worker.jobs {
		secs = append(secs, job.CreateSection())
	}
	return secs
}
