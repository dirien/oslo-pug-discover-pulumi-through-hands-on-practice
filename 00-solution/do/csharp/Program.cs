using System.Collections.Generic;
using System.Linq;
using Pulumi;
using DigitalOcean = Pulumi.DigitalOcean;

return await Deployment.RunAsync(() => 
{
    var clusterRegion = "fra1";
    var nodePoolName = "default";
    var nodeCount = 1;
    var version = "1.32.2-do.0";

    var doCluster = new DigitalOcean.KubernetesCluster("do-cluster", new()
    {
        Region = clusterRegion,
        Version = version,
        DestroyAllAssociatedResources = true,
        NodePool = new DigitalOcean.Inputs.KubernetesClusterNodePoolArgs
        {
            Name = nodePoolName,
            Size = "s-2vcpu-2gb",
            NodeCount = nodeCount,
        },
    });

    return new Dictionary<string, object?>
    {
        ["name"] = doCluster.Name,
        ["kubeconfig"] = doCluster.KubeConfigs.Apply(kubeConfigs => kubeConfigs[0].RawConfig),
    };
});

