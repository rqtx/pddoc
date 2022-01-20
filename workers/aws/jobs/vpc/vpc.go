package vpc

import (
	"fmt"
	"os"

	"github.com/rqtx/pdoc/utils"
	"github.com/rqtx/pdoc/workers/aws/helpers"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type JobVPC struct {
	service *ec2.EC2
}

func New(region string) *JobVPC {
	sess := session.Must(session.NewSession())
	return &JobVPC{
		service: ec2.New(sess, aws.NewConfig().WithRegion(region)),
	}
}

func (job *JobVPC) CreateSection() utils.Section {
	subs := make([]utils.Subsection, 0)

	subs = append(subs, utils.NewSubsection("VPCs", job.vpc()))
	subs = append(subs, utils.NewSubsection("Subnets", job.subnets()))
	subs = append(subs, utils.NewSubsection("Route Tables", job.routeTables()))
	subs = append(subs, utils.NewSubsection("Security Groups", job.securityGroups()))
	return utils.NewSection("VPC", subs)
}

func (job *JobVPC) vpc() utils.Table {
	vpcs, err := job.service.DescribeVpcs(&ec2.DescribeVpcsInput{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	data := make([][]string, 0)
	for _, vpc := range vpcs.Vpcs {
		tmp := make([]string, 0)
		if value := helpers.GetTagValue(vpc.Tags, "Name"); value != nil {
			tmp = append(tmp, *value)
		} else {
			tmp = append(tmp, "")
		}
		tmp = append(tmp, *(vpc.VpcId))
		tmp = append(tmp, *(vpc.OwnerId))
		tmp = append(tmp, *(vpc.CidrBlock))
		data = append(data, [][]string{tmp}...)
	}
	return utils.NewTable([]string{"Name", "ID", "OwnerId", "CIDR"}, data, "VPCs")
}

func (job *JobVPC) subnets() utils.Table {
	subnets, err := job.service.DescribeSubnets(&ec2.DescribeSubnetsInput{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	data := make([][]string, 0)
	for _, subnet := range subnets.Subnets {
		tmp := make([]string, 0)
		if value := helpers.GetTagValue(subnet.Tags, "Name"); value != nil {
			tmp = append(tmp, *value)
		} else {
			tmp = append(tmp, "")
		}
		tmp = append(tmp, *(subnet.SubnetId))
		tmp = append(tmp, *(subnet.AvailabilityZone))
		tmp = append(tmp, *(subnet.CidrBlock))
		data = append(data, [][]string{tmp}...)
	}
	return utils.NewTable([]string{"Name", "ID", "Availability Zone", "CIDR"}, data, "Subnets")
}

func (job *JobVPC) routeTables() utils.Table {
	rts, err := job.service.DescribeRouteTables(&ec2.DescribeRouteTablesInput{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	data := make([][]string, 0)
	for _, rt := range rts.RouteTables {
		tmp := make([]string, 0)
		if value := helpers.GetTagValue(rt.Tags, "Name"); value != nil {
			tmp = append(tmp, *value)
		} else {
			tmp = append(tmp, "")
		}
		tmp = append(tmp, *(rt.RouteTableId))
		data = append(data, [][]string{tmp}...)
	}
	return utils.NewTable([]string{"Name", "ID"}, data, "Route Tables")
}

func (job *JobVPC) securityGroups() utils.Table {
	sgs, err := job.service.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	data := make([][]string, 0)
	for _, sg := range sgs.SecurityGroups {
		tmp := make([]string, 0)
		tmp = append(tmp, *(sg.GroupName))
		tmp = append(tmp, *(sg.GroupId))
		tmp = append(tmp, *(sg.VpcId))
		data = append(data, [][]string{tmp}...)
	}
	return utils.NewTable([]string{"Name", "ID", "VPC ID"}, data, "Security Groups")
}
