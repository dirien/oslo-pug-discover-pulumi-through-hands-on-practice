import * as pulumi from "@pulumi/pulumi";
import * as digitalocean from "@pulumi/digitalocean";

const clusterRegion = "fra1";
const nodePoolName = "default";
const nodeCount = 1;
const version = "1.31.1-do.5";
const doCluster = new digitalocean.KubernetesCluster("do-cluster", {
    region: clusterRegion,
    version: version,
    destroyAllAssociatedResources: true,
    nodePool: {
        name: nodePoolName,
        size: "s-2vcpu-2gb",
        nodeCount: nodeCount,
    },
});
export const name = doCluster.name;
export const kubeconfig = doCluster.kubeConfigs.apply(kubeConfigs => kubeConfigs[0].rawConfig);
