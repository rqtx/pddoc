package lambda

import (
	"fmt"
	"os"

	"github.com/rqtx/pdoc/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

type JobLambda struct {
	service *lambda.Lambda
}

func New(region string) *JobLambda {
	sess := session.Must(session.NewSession())
	return &JobLambda{
		service: lambda.New(sess, aws.NewConfig().WithRegion(region)),
	}
}

func (job *JobLambda) CreateSection() utils.Section {
	subs := make([]utils.Subsection, 0)

	subs = append(subs, utils.NewSubsection("Functions", job.lambda()))
	return utils.NewSection("Lambda", subs)
}

func (job *JobLambda) lambda() utils.Table {
	lambdas, err := job.service.ListFunctions(&lambda.ListFunctionsInput{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	data := make([][]string, 0)
	for _, lambda := range lambdas.Functions {
		tmp := make([]string, 0)
		tmp = append(tmp, *(lambda.FunctionName))
		tmp = append(tmp, *(lambda.FunctionArn))
		data = append(data, [][]string{tmp}...)
	}
	return utils.NewTable([]string{"Name", "ARN"}, data, "Functions")
}
