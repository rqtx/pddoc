package waf

import (
	"fmt"
	"os"

	"github.com/rqtx/pddoc/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/waf"
)

type JobWAF struct {
	service *waf.WAF
}

func New(region string) *JobWAF {
	sess := session.Must(session.NewSession())
	return &JobWAF{
		service: waf.New(sess, aws.NewConfig().WithRegion(region)),
	}
}

func (job *JobWAF) CreateSection() utils.Section {
	subs := make([]utils.Subsection, 0)

	subs = append(subs, utils.NewSubsection("WebACLs", job.webACLs()))
	subs = append(subs, utils.NewSubsection("IPSets", job.ipSets()))
	return utils.NewSection("WAF", subs)
}

func (job *JobWAF) webACLs() utils.Table {
	webacls, err := job.service.ListWebACLs(&waf.ListWebACLsInput{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	data := make([][]string, 0)
	for _, web := range webacls.WebACLs {
		tmp := make([]string, 0)
		tmp = append(tmp, *(web.Name))
		tmp = append(tmp, *(web.WebACLId))
		data = append(data, [][]string{tmp}...)
	}
	return utils.NewTable([]string{"Name", "ID"}, data, "WebACLs")
}

func (job *JobWAF) ipSets() utils.Table {
	ipsets, err := job.service.ListIPSets(&waf.ListIPSetsInput{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	data := make([][]string, 0)
	for _, set := range ipsets.IPSets {
		tmp := make([]string, 0)
		tmp = append(tmp, *(set.Name))
		tmp = append(tmp, *(set.IPSetId))
		data = append(data, [][]string{tmp}...)
	}
	return utils.NewTable([]string{"Name", "ID"}, data, "IPsets")
}
