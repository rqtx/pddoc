package kms

import (
	"fmt"
	"os"

	"github.com/rqtx/pdoc/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

type JobKMS struct {
	service *kms.KMS
}

func New(region string) *JobKMS {
	sess := session.Must(session.NewSession())
	return &JobKMS{
		service: kms.New(sess, aws.NewConfig().WithRegion(region)),
	}
}

func (job *JobKMS) CreateSection() utils.Section {
	subs := make([]utils.Subsection, 0)

	subs = append(subs, utils.NewSubsection("Keypairs", job.keys()))
	return utils.NewSection("KMS", subs)
}

func (worker *JobKMS) keys() utils.Table {
	keys, err := worker.service.ListKeys(&kms.ListKeysInput{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	data := make([][]string, 0)
	for _, key := range keys.Keys {
		keyid := key.KeyId
		tmp := make([]string, 0)
		tmp = append(tmp, *keyid)
		alias := ""
		aliases, err := worker.service.ListAliases(&kms.ListAliasesInput{KeyId: keyid})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			os.Exit(1)
		}
		for _, al := range aliases.Aliases {
			alias += *(al.AliasName) + "\n"
		}
		tmp = append(tmp, alias[:len(alias)-1])
		data = append(data, [][]string{tmp}...)
	}
	return utils.NewTable([]string{"Key ID", "Aliases"}, data, "Keypairs")
}
