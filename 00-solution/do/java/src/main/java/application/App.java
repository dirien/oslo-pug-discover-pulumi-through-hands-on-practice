package application;

import com.pulumi.Context;
import com.pulumi.Pulumi;
import com.pulumi.digitalocean.KubernetesCluster;
import com.pulumi.digitalocean.KubernetesClusterArgs;
import com.pulumi.digitalocean.inputs.KubernetesClusterNodePoolArgs;

public class App {
    public static void main(String[] args) {
        Pulumi.run(App::stack);
    }

    public static void stack(Context ctx) {
        final var clusterRegion = "fra1";
        final var nodePoolName = "default";
        final var nodeCount = 1;
        final var version = "1.31.1-do.5";

        var doCluster = new KubernetesCluster("do-cluster", KubernetesClusterArgs.builder()
            .region(clusterRegion)
            .version(version)
            .destroyAllAssociatedResources(true)
            .nodePool(KubernetesClusterNodePoolArgs.builder()
                .name(nodePoolName)
                .size("s-2vcpu-2gb")
                .nodeCount(nodeCount)
                .build())
            .build());

        ctx.export("name", doCluster.name());
        ctx.export("kubeconfig", doCluster.kubeConfigs().applyValue(kubeConfigs -> kubeConfigs.get(0).rawConfig()));
    }
}
