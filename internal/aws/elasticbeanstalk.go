package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk"
	"github.com/aws/aws-sdk-go-v2/config"
)

type ElasticBeanstalkEnvironment struct {
	ID          string
	Name        string
	AppName     string
	Status      string
	Health      string
	CNAME       string
}

type ElasticBeanstalkFetcher struct {
	client *elasticbeanstalk.Client
}

func NewElasticBeanstalkFetcher(ctx context.Context, region string) (*ElasticBeanstalkFetcher, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	return &ElasticBeanstalkFetcher{client: elasticbeanstalk.NewFromConfig(cfg)}, nil
}

func (f *ElasticBeanstalkFetcher) FetchAll(ctx context.Context) ([]ElasticBeanstalkEnvironment, error) {
	out, err := f.client.DescribeEnvironments(ctx, &elasticbeanstalk.DescribeEnvironmentsInput{
		IncludeDeleted: aws.Bool(false),
	})
	if err != nil {
		return nil, err
	}
	var envs []ElasticBeanstalkEnvironment
	for _, e := range out.Environments {
		envs = append(envs, mapBeanstalkEnvironment(e))
	}
	return envs, nil
}

func mapBeanstalkEnvironment(e elasticbeanstalk_types.EnvironmentDescription) ElasticBeanstalkEnvironment {
	env := ElasticBeanstalkEnvironment{
		Status: string(e.Status),
		Health: string(e.Health),
	}
	if e.EnvironmentId != nil {
		env.ID = aws.ToString(e.EnvironmentId)
	}
	if e.EnvironmentName != nil {
		env.Name = aws.ToString(e.EnvironmentName)
	}
	if e.ApplicationName != nil {
		env.AppName = aws.ToString(e.ApplicationName)
	}
	if e.CNAME != nil {
		env.CNAME = aws.ToString(e.CNAME)
	}
	return env
}
