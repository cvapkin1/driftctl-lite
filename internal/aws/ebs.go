package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type EBSVolume struct {
	ID         string
	State      string
	Size       int32
	VolumeType string
	AZ         string
	Encrypted  bool
}

type ebsClient interface {
	DescribeVolumes(ctx context.Context, params *ec2.DescribeVolumesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeVolumesOutput, error)
}

type EBSFetcher struct {
	client ebsClient
}

func NewEBSFetcher(client ebsClient) *EBSFetcher {
	return &EBSFetcher{client: client}
}

func (f *EBSFetcher) FetchAll(ctx context.Context) ([]EBSVolume, error) {
	out, err := f.client.DescribeVolumes(ctx, &ec2.DescribeVolumesInput{})
	if err != nil {
		return nil, err
	}
	var volumes []EBSVolume
	for _, v := range out.Volumes {
		volumes = append(volumes, mapEBSVolume(v))
	}
	return volumes, nil
}

func mapEBSVolume(v types.Volume) EBSVolume {
	state := string(v.State)
	volumeType := string(v.VolumeType)
	az := ""
	if v.AvailabilityZone != nil {
		az = aws.ToString(v.AvailabilityZone)
	}
	size := int32(0)
	if v.Size != nil {
		size = aws.ToInt32(v.Size)
	}
	return EBSVolume{
		ID:         aws.ToString(v.VolumeId),
		State:      state,
		Size:       size,
		VolumeType: volumeType,
		AZ:         az,
		Encrypted:  aws.ToBool(v.Encrypted),
	}
}
