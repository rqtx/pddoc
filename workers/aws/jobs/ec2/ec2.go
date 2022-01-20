package ec2

import (
	"fmt"
	"os"
	"strconv"

	"github.com/rqtx/pdoc/utils"
	"github.com/rqtx/pdoc/workers/aws/helpers"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/sts"
)

type JobEC2 struct {
	serviceEC2   *ec2.EC2
	serviceASG   *autoscaling.AutoScaling
	serviceELBv2 *elbv2.ELBV2
	serviceSTS   *sts.STS
}

func New(region string) *JobEC2 {
	sess := session.Must(session.NewSession())
	return &JobEC2{
		serviceEC2:   ec2.New(sess, aws.NewConfig().WithRegion(region)),
		serviceASG:   autoscaling.New(sess, aws.NewConfig().WithRegion(region)),
		serviceELBv2: elbv2.New(sess, aws.NewConfig().WithRegion(region)),
		serviceSTS:   sts.New(sess, aws.NewConfig().WithRegion(region)),
	}
}

func (job *JobEC2) CreateSection() utils.Section {
	subs := make([]utils.Subsection, 0)

	subs = append(subs, utils.NewSubsection("AMI", job.ami()))
	subs = append(subs, utils.NewSubsection("Instances", job.instances()))
	subs = append(subs, utils.NewSubsection("Load Balancers", job.loadBalancers()))
	subs = append(subs, utils.NewSubsection("Launch Template", job.launchTemplates()))
	subs = append(subs, utils.NewSubsection("Auto Scaling Groups", job.autoScalingGroups()))
	return utils.NewSection("EC2", subs)
}

func (job *JobEC2) ami() utils.Table {
	caller, err := job.serviceSTS.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	iams, err := job.serviceEC2.DescribeImages(&ec2.DescribeImagesInput{Owners: []*string{caller.Account}})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	data := make([][]string, 0)
	for _, image := range iams.Images {
		tmp := make([]string, 0)
		tmp = append(tmp, *image.ImageId)
		if image.Name != nil {
			tmp = append(tmp, *image.Name)
		} else {
			tmp = append(tmp, "")
		}
		data = append(data, [][]string{tmp}...)
	}
	return utils.NewTable([]string{"ID", "Name"}, data, "AMI")
}

func (job *JobEC2) instances() utils.Table {
	ec2, err := job.serviceEC2.DescribeInstances(&ec2.DescribeInstancesInput{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	data := make([][]string, 0)
	for _, reservations := range ec2.Reservations {
		for _, instance := range reservations.Instances {
			tmp := make([]string, 0)
			tmp = append(tmp, *instance.InstanceId)
			if value := helpers.GetTagValue(instance.Tags, "Name"); value != nil {
				tmp = append(tmp, *value)
			} else {
				tmp = append(tmp, "")
			}
			tmp = append(tmp, *instance.InstanceType)
			tmp = append(tmp, *instance.Placement.AvailabilityZone)
			tmp = append(tmp, *instance.PrivateDnsName)
			tmp = append(tmp, *instance.PrivateIpAddress)
			data = append(data, [][]string{tmp}...)
		}
	}
	return utils.NewTable([]string{"Instance ID", "Name", "Instance Type", "Availability Zone", "Private IPv4 DNS", "Private IPv4"}, data, "Instances")
}

func (job *JobEC2) loadBalancers() utils.Table {
	lbs, err := job.serviceELBv2.DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	data := make([][]string, 0)
	for _, lb := range lbs.LoadBalancers {
		tmp := make([]string, 0)
		tmp = append(tmp, *lb.LoadBalancerName)
		tmp = append(tmp, *lb.LoadBalancerArn)
		tmp = append(tmp, *lb.Type)
		tmp = append(tmp, *lb.Scheme)
		data = append(data, [][]string{tmp}...)
	}
	return utils.NewTable([]string{"Name", "ARN", "Type", "Scheme"}, data, "Launch Template")
}

func (job *JobEC2) launchTemplates() utils.Table {
	templates, err := job.serviceEC2.DescribeLaunchTemplates(&ec2.DescribeLaunchTemplatesInput{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	data := make([][]string, 0)
	for _, template := range templates.LaunchTemplates {
		tmp := make([]string, 0)
		tmp = append(tmp, *template.LaunchTemplateId)
		tmp = append(tmp, *template.LaunchTemplateName)
		data = append(data, [][]string{tmp}...)
	}
	return utils.NewTable([]string{"ID", "Name"}, data, "Launch Template")
}

func (job *JobEC2) autoScalingGroups() utils.Table {
	asgs, err := job.serviceASG.DescribeAutoScalingGroups(&autoscaling.DescribeAutoScalingGroupsInput{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	data := make([][]string, 0)
	for _, asg := range asgs.AutoScalingGroups {
		tmp := make([]string, 0)
		tmp = append(tmp, *asg.AutoScalingGroupName)
		tmp = append(tmp, *asg.LaunchConfigurationName)
		tmp = append(tmp, strconv.FormatInt(*asg.MinSize, 10))
		tmp = append(tmp, strconv.FormatInt(*asg.MaxSize, 10))
		data = append(data, [][]string{tmp}...)
	}
	return utils.NewTable([]string{"Name", "Launch Template", "Min", "Max"}, data, "Auto Scaling Groups")
}
