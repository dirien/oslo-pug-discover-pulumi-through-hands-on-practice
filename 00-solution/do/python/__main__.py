import pulumi
import pulumi_digitalocean as digitalocean

cluster_region = "fra1"
node_pool_name = "default"
node_count = 1
version = "1.31.1-do.5"
do_cluster = digitalocean.KubernetesCluster("do-cluster",
    region=cluster_region,
    version=version,
    destroy_all_associated_resources=True,
    node_pool={
        "name": node_pool_name,
        "size": "s-2vcpu-2gb",
        "node_count": node_count,
    })
pulumi.export("name", do_cluster.name)
pulumi.export("kubeconfig", do_cluster.kube_configs[0].raw_config)
