package iam

import (
	"fmt"
	"os"

	"github.com/rqtx/pddoc/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

type JobIAM struct {
	service *iam.IAM
}

func New(region string) *JobIAM {
	sess := session.Must(session.NewSession())
	return &JobIAM{
		service: iam.New(sess, aws.NewConfig().WithRegion(region)),
	}
}

func (worker *JobIAM) CreateSection() utils.Section {
	subs := make([]utils.Subsection, 0)

	subs = append(subs, utils.NewSubsection("Groups", worker.groups()))
	subs = append(subs, utils.NewSubsection("Users", worker.users()))
	subs = append(subs, utils.NewSubsection("Policies", worker.policies()))
	subs = append(subs, utils.NewSubsection("Roles", worker.roles()))
	return utils.NewSection("IAM", subs)
}

func (job *JobIAM) groups() utils.Table {
	groups, err := job.service.ListGroups(&iam.ListGroupsInput{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	data := make([][]string, 0)
	for _, group := range groups.Groups {
		groupName := group.GroupName
		tmp := make([]string, 0)
		tmp = append(tmp, *groupName)
		listPolicies, err := job.service.ListAttachedGroupPolicies(&iam.ListAttachedGroupPoliciesInput{GroupName: groupName})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			os.Exit(1)
		}
		plc := ""
		for _, policy := range listPolicies.AttachedPolicies {
			plc += *(policy.PolicyName) + "\n"
		}
		tmp = append(tmp, plc[:len(plc)-1])
		data = append(data, [][]string{tmp}...)
	}
	return utils.NewTable([]string{"Group", "Policies"}, data, "Groups")
}

func (job *JobIAM) users() utils.Table {
	users, err := job.service.ListUsers(&iam.ListUsersInput{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	data := make([][]string, 0)
	for _, user := range users.Users {
		userName := user.UserName
		temp := make([]string, 0)
		temp = append(temp, *userName)
		listGroups, err := job.service.ListGroupsForUser(&iam.ListGroupsForUserInput{UserName: userName})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			os.Exit(1)
		}
		plc := ""
		for _, group := range listGroups.Groups {
			plc += *(group.GroupName) + "\n"
		}
		if len(plc) > 0 {
			temp = append(temp, plc[:len(plc)-1])
		} else {
			temp = append(temp, "")
		}
		data = append(data, [][]string{temp}...)
	}
	return utils.NewTable([]string{"Users", "Groups"}, data, "Users")
}

func (job *JobIAM) policies() utils.Table {
	scope := "Local" //List only customer managed policies https://docs.aws.amazon.com/sdk-for-go/api/service/iam/#IAM.ListGroupsForUser
	listPolicies, err := job.service.ListPolicies(&iam.ListPoliciesInput{Scope: &scope})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	data := make([][]string, 0)
	for _, plc := range listPolicies.Policies {
		temp := make([]string, 0)
		temp = append(temp, *(plc.PolicyName))
		data = append(data, [][]string{temp}...)
	}
	return utils.NewTable([]string{"Policy"}, data, "Policies")
}

func (job *JobIAM) roles() utils.Table {
	listRoles, err := job.service.ListRoles(&iam.ListRolesInput{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	data := make([][]string, 0)
	for _, role := range listRoles.Roles {
		temp := make([]string, 0)
		temp = append(temp, *(role.RoleName))
		data = append(data, [][]string{temp}...)
	}
	return utils.NewTable([]string{"Role"}, data, "Roles")
}
