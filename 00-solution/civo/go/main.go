package main

import (
	"github.com/pulumi/pulumi-civo/sdk/v2/go/civo"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		firewall, err := civo.NewFirewall(ctx, "firewall", &civo.FirewallArgs{
			CreateDefaultRules: pulumi.Bool(true),
			Region:             pulumi.String("LON1"),
		})
		if err != nil {
			return err
		}
		cluster, err := civo.NewKubernetesCluster(ctx, "cluster", &civo.KubernetesClusterArgs{
			FirewallId:        firewall.ID(),
			Region:            pulumi.String("LON1"),
			Cni:               pulumi.String("cilium"),
			Applications:      pulumi.String("-traefik2-nodeport,civo-cluster-autoscaler,traefik2-loadbalancer"),
			KubernetesVersion: pulumi.String("1.28.7-k3s1"),
			WriteKubeconfig:   pulumi.Bool(true),
			Pools: &civo.KubernetesClusterPoolsArgs{
				Size:      pulumi.String("g4s.kube.medium"),
				NodeCount: pulumi.Int(1),
			},
		})
		if err != nil {
			return err
		}
		ctx.Export("name", cluster.Name)
		ctx.Export("kubeconfig", cluster.Kubeconfig)
		return nil
	})
}
