package storage

import (
	"fmt"
	"os"

	"github.com/rqtx/pddoc/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/efs"
	"github.com/aws/aws-sdk-go/service/s3"
)

type JobStorage struct {
	serviceEFS *efs.EFS
	serviceS3  *s3.S3
}

func New(region string) *JobStorage {
	sess := session.Must(session.NewSession())
	return &JobStorage{
		serviceEFS: efs.New(sess, aws.NewConfig().WithRegion(region)),
		serviceS3:  s3.New(sess, aws.NewConfig().WithRegion(region)),
	}
}

func (job *JobStorage) CreateSection() utils.Section {
	subs := make([]utils.Subsection, 0)

	subs = append(subs, utils.NewSubsection("EFS", job.efs()))
	subs = append(subs, utils.NewSubsection("S3", job.s3()))
	return utils.NewSection("Storage", subs)
}

func (job *JobStorage) efs() utils.Table {
	efs, err := job.serviceEFS.DescribeFileSystems(&efs.DescribeFileSystemsInput{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	data := make([][]string, 0)
	for _, fs := range efs.FileSystems {
		tmp := make([]string, 0)
		tmp = append(tmp, *(fs.FileSystemId))
		tmp = append(tmp, *(fs.FileSystemArn))
		data = append(data, [][]string{tmp}...)
	}
	return utils.NewTable([]string{"ID", "ARN"}, data, "EFS")
}

func (job *JobStorage) s3() utils.Table {
	s3, err := job.serviceS3.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	data := make([][]string, 0)
	for _, bucket := range s3.Buckets {
		tmp := make([]string, 0)
		tmp = append(tmp, *(bucket.Name))
		data = append(data, [][]string{tmp}...)
	}
	return utils.NewTable([]string{"Name"}, data, "Buckets S3")
}
