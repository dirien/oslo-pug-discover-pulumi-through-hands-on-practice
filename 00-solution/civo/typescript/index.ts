import * as pulumi from "@pulumi/pulumi";
import * as civo from "@pulumi/civo";

const firewall = new civo.Firewall("firewall", {
    createDefaultRules: true,
    region: "LON1",
});
const cluster = new civo.KubernetesCluster("cluster", {
    firewallId: firewall.id,
    region: "LON1",
    cni: "cilium",
    applications: "-traefik2-nodeport,civo-cluster-autoscaler,traefik2-loadbalancer",
    kubernetesVersion: "1.28.7-k3s1",
    writeKubeconfig: true,
    pools: {
        size: "g4s.kube.medium",
        nodeCount: 1,
    },
});
export const name = cluster.name;
export const kubeconfig = cluster.kubeconfig;
