package generated_program;

import com.pulumi.Context;
import com.pulumi.Pulumi;
import com.pulumi.core.Output;
import com.pulumi.civo.Firewall;
import com.pulumi.civo.FirewallArgs;
import com.pulumi.civo.KubernetesCluster;
import com.pulumi.civo.KubernetesClusterArgs;
import com.pulumi.civo.inputs.KubernetesClusterPoolsArgs;
import java.util.List;
import java.util.ArrayList;
import java.util.Map;
import java.io.File;
import java.nio.file.Files;
import java.nio.file.Paths;

public class App {
    public static void main(String[] args) {
        Pulumi.run(App::stack);
    }

    public static void stack(Context ctx) {
        var firewall = new Firewall("firewall", FirewallArgs.builder()
            .createDefaultRules(true)
            .region("LON1")
            .build());

        var cluster = new KubernetesCluster("cluster", KubernetesClusterArgs.builder()
            .firewallId(firewall.id())
            .region("LON1")
            .cni("cilium")
            .applications("-traefik2-nodeport,civo-cluster-autoscaler,traefik2-loadbalancer")
            .kubernetesVersion("1.28.7-k3s1")
            .writeKubeconfig(true)
            .pools(KubernetesClusterPoolsArgs.builder()
                .size("g4s.kube.medium")
                .nodeCount(1)
                .build())
            .build());

        ctx.export("name", cluster.name());
        ctx.export("kubeconfig", cluster.kubeconfig());
    }
}
