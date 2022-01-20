package cloudwatch

import (
	"fmt"
	"os"

	"github.com/rqtx/pdoc/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

type JobCW struct {
	service *cloudwatch.CloudWatch
}

func New(region string) *JobCW {
	sess := session.Must(session.NewSession())
	return &JobCW{
		service: cloudwatch.New(sess, aws.NewConfig().WithRegion(region)),
	}
}

func (job *JobCW) CreateSection() utils.Section {
	subs := make([]utils.Subsection, 0)
	subs = append(subs, utils.NewSubsection("Dashbords", job.dashboards()))
	return utils.NewSection("CloudWatch", subs)
}

func (job *JobCW) dashboards() utils.Table {
	cw, err := job.service.ListDashboards(&cloudwatch.ListDashboardsInput{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	data := make([][]string, 0)
	for _, dbord := range cw.DashboardEntries {
		tmp := make([]string, 0)
		tmp = append(tmp, *(dbord.DashboardName))
		tmp = append(tmp, *(dbord.DashboardName))
		data = append(data, [][]string{tmp}...)
	}
	return utils.NewTable([]string{"NAME", "ARN"}, data, "Dashbords")
}
