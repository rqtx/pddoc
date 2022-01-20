package helpers

import (
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go/service/ec2"
)

func GetTagValue(tags []*ec2.Tag, key string) *string {
	for _, tag := range tags {
		if *tag.Key == key {
			return tag.Value
		}
	}
	return nil
}

func GetNameFromARN(arn string) (string, error) {
	re := regexp.MustCompile("arn:aws:ecs:\\S+:\\d+:\\w+\\/\\S+")
	if match := re.FindAllString(arn, -1); match == nil {
		return "", fmt.Errorf("Invalide ARN")
	}
	re = regexp.MustCompile("[^/]*$")
	return re.FindStringSubmatch(arn)[0], nil
}
