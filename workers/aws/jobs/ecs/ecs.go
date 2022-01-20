package ecs

import (
	"fmt"
	"os"

	"github.com/rqtx/pddoc/utils"
	"github.com/rqtx/pddoc/workers/aws/helpers"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type JobECS struct {
	service *ecs.ECS
}

func New(region string) *JobECS {
	sess := session.Must(session.NewSession())
	return &JobECS{
		service: ecs.New(sess, aws.NewConfig().WithRegion(region)),
	}
}

func (job *JobECS) CreateSection() utils.Section {
	subs := make([]utils.Subsection, 0)

	subs = append(subs, utils.NewSubsection("Clusters", job.clusters()))
	return utils.NewSection("ECS", subs)
}

func (job *JobECS) clusters() utils.Table {
	clusters, err := job.service.ListClusters(&ecs.ListClustersInput{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	data := make([][]string, 0)
	for _, arn := range clusters.ClusterArns {
		tmp := make([]string, 0)
		clusterName, _ := helpers.GetNameFromARN(*arn)
		tmp = append(tmp, clusterName)
		services, _ := job.service.ListServices(&ecs.ListServicesInput{Cluster: arn})
		svcs := ""
		enabled := ""
		for _, svcArn := range services.ServiceArns {
			svcName, err := helpers.GetNameFromARN(*svcArn)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
				os.Exit(1)
			}
			svcs += svcName + "\n"

			service, err := job.service.DescribeServices(&ecs.DescribeServicesInput{Cluster: arn, Services: []*string{&svcName}})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
				os.Exit(1)
			}
			if *(service.Services[0].DesiredCount) > 0 {
				enabled += "YES" + "\n"
			} else {
				enabled += "NO" + "\n"
			}
		}
		tmp = append(tmp, svcs[:len(svcs)-1])
		tmp = append(tmp, enabled[:len(enabled)-1])
		data = append(data, [][]string{tmp}...)
	}
	return utils.NewTable([]string{"Cluster", "Service", "Enabled"}, data, "Clusters")
}
