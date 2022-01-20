package account

import (
	"fmt"
	"os"

	"github.com/rqtx/pdoc/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

type JobAccount struct {
	service *sts.STS
}

func New(region string) *JobAccount {
	sess := session.Must(session.NewSession())
	return &JobAccount{
		service: sts.New(sess, aws.NewConfig().WithRegion(region)),
	}
}

func (job *JobAccount) CreateSection() utils.Section {
	subs := make([]utils.Subsection, 0)
	subs = append(subs, utils.NewSubsection("ID", job.account()))
	return utils.NewSection("Account", subs)
}

func (job *JobAccount) account() utils.Table {
	caller, err := job.service.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	data := make([][]string, 0)
	tmp := make([]string, 0)
	tmp = append(tmp, *(caller.Account))
	data = append(data, [][]string{tmp}...)
	return utils.NewTable([]string{"Account ID"}, data, "Account")
}
