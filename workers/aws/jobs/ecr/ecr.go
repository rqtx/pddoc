package ecr

import (
	"fmt"
	"os"

	"github.com/rqtx/pdoc/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

type JobECR struct {
	service *ecr.ECR
}

func New(region string) *JobECR {
	sess := session.Must(session.NewSession())
	return &JobECR{
		service: ecr.New(sess, aws.NewConfig().WithRegion(region)),
	}
}

func (job *JobECR) CreateSection() utils.Section {
	subs := make([]utils.Subsection, 0)

	subs = append(subs, utils.NewSubsection("Repositories", job.repositories()))
	return utils.NewSection("ECR", subs)
}

func (job *JobECR) repositories() utils.Table {
	repos, err := job.service.DescribeRepositories(&ecr.DescribeRepositoriesInput{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	data := make([][]string, 0)
	for _, repo := range repos.Repositories {
		tmp := make([]string, 0)
		tmp = append(tmp, *repo.RepositoryName)
		data = append(data, [][]string{tmp}...)
	}
	return utils.NewTable([]string{"Name"}, data, "Repositories")
}
