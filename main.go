package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := ec2.NewFromConfig(cfg)
	volumes, err := client.DescribeVolumes(context.Background(), &ec2.DescribeVolumesInput{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("VOLUMES\n\n")

	for _, volume := range volumes.Volumes {
		fmt.Printf("id: %s zone: %s type: %s size: %d state: %s, created: %s\n", *volume.VolumeId, *volume.AvailabilityZone,
			*&volume.VolumeType, *&volume.Size, volume.State, volume.CreateTime)

		fmt.Printf("\tTags: %s\n", joinTags(volume.Tags))

		if len(volume.Attachments) == 0 {
			fmt.Printf("\tOrphan")
			continue
		}

		for _, attachment := range volume.Attachments {
			fmt.Printf("\tAttached to instance: %s as device: %s at: %s\n", *attachment.InstanceId, *attachment.Device,
				attachment.AttachTime)
		}
	}

	fmt.Printf("\nSNAPSHOTS\n\n")

	snapshots, err := client.DescribeSnapshots(context.Background(), &ec2.DescribeSnapshotsInput{OwnerIds: []string{"507980699075"}})
	if err != nil {
		panic(err)
	}

	for _, snapshot := range snapshots.Snapshots {
		fmt.Printf("id: %s volume: %s size: %d state: %s storage tier: %s owner: %s started: %s progress: %s description: %s\n",
			*snapshot.SnapshotId, *snapshot.VolumeId, *snapshot.VolumeSize, snapshot.State, snapshot.StorageTier,
			*snapshot.OwnerId, *snapshot.StartTime, *snapshot.Progress, *snapshot.Description)

		fmt.Printf("\tTags: %s\n", joinTags(snapshot.Tags))
	}
}

func joinTags(tags []types.Tag) string {
	if len(tags) == 0 {
		return ""
	}

	ts := make([]string, 0, len(tags))
	for _, tag := range tags {
		ts = append(ts, fmt.Sprintf("%s: %s", *tag.Key, *tag.Value))
	}

	return strings.Join(ts, ", ")
}
