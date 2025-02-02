import pulumi
import pulumi_civo as civo

firewall = civo.Firewall("firewall",
    create_default_rules=True,
    region="LON1")
cluster = civo.KubernetesCluster("cluster",
    firewall_id=firewall.id,
    region="LON1",
    cni="cilium",
    applications="-traefik2-nodeport,civo-cluster-autoscaler,traefik2-loadbalancer",
    kubernetes_version="1.28.7-k3s1",
    write_kubeconfig=True,
    pools={
        "size": "g4s.kube.medium",
        "node_count": 1,
    })
pulumi.export("name", cluster.name)
pulumi.export("kubeconfig", cluster.kubeconfig)
