using System.Collections.Generic;
using System.Linq;
using Pulumi;
using Civo = Pulumi.Civo;

return await Deployment.RunAsync(() => 
{
    var firewall = new Civo.Firewall("firewall", new()
    {
        CreateDefaultRules = true,
        Region = "LON1",
    });

    var cluster = new Civo.KubernetesCluster("cluster", new()
    {
        FirewallId = firewall.Id,
        Region = "LON1",
        Cni = "cilium",
        Applications = "-traefik2-nodeport,civo-cluster-autoscaler,traefik2-loadbalancer",
        KubernetesVersion = "1.28.7-k3s1",
        WriteKubeconfig = true,
        Pools = new Civo.Inputs.KubernetesClusterPoolsArgs
        {
            Size = "g4s.kube.medium",
            NodeCount = 1,
        },
    });

    return new Dictionary<string, object?>
    {
        ["name"] = cluster.Name,
        ["kubeconfig"] = cluster.Kubeconfig,
    };
});

