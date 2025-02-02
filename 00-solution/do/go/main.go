package main

import (
	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		clusterRegion := "fra1"
		nodePoolName := "default"
		nodeCount := 1
		version := "1.31.1-do.5"
		doCluster, err := digitalocean.NewKubernetesCluster(ctx, "do-cluster", &digitalocean.KubernetesClusterArgs{
			Region:                        pulumi.String(clusterRegion),
			Version:                       pulumi.String(version),
			DestroyAllAssociatedResources: pulumi.Bool(true),
			NodePool: &digitalocean.KubernetesClusterNodePoolArgs{
				Name:      pulumi.String(nodePoolName),
				Size:      pulumi.String("s-2vcpu-2gb"),
				NodeCount: pulumi.IntPtr(nodeCount),
			},
		})
		if err != nil {
			return err
		}
		ctx.Export("name", doCluster.Name)
		ctx.Export("kubeconfig", doCluster.KubeConfigs.ApplyT(func(kubeConfigs []digitalocean.KubernetesClusterKubeConfig) (*string, error) {
			return kubeConfigs[0].RawConfig, nil
		}).(pulumi.StringPtrOutput))
		return nil
	})
}
