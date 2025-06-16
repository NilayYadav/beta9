shadeform@shadecloud:~/beta9$ make k3d-down
bash bin/k3d.sh down
INFO[0000] Using config file hack/k3d.yaml (k3d.io/v1alpha5#simple) 
INFO[0000] Deleting cluster 'beta9'                     
INFO[0007] Deleting cluster network 'k3d-beta9'         
INFO[0007] Deleting 1 attached volumes...               
INFO[0007] Removing cluster details from default kubeconfig... 
INFO[0007] Removing standalone kubeconfig file (if there is one)... 
INFO[0007] Successfully deleted cluster beta9!          
shadeform@shadecloud:~/beta9$ make stop
cd hack && okteto down --file okteto.yaml
 x  Invalid okteto context 'k3d-beta9'
    Please run 'okteto context' to select one context
make: *** [Makefile:59: stop] Error 1
shadeform@shadecloud:~/beta9$ make k3d-up
bash bin/k3d.sh up
[+] Building 0.2s (11/11) FINISHED                                                                                                                   docker:default
 => [internal] load build definition from Dockerfile.k3d                                                                                                       0.0s
 => => transferring dockerfile: 1.38kB                                                                                                                         0.0s
 => resolve image config for docker-image://docker.io/docker/dockerfile:1.6                                                                                    0.1s
 => CACHED docker-image://docker.io/docker/dockerfile:1.6@sha256:ac85f380a63b13dfcefa89046420e1781752bab202122f8f50032edf31be0021                              0.0s
 => [internal] load metadata for docker.io/nvidia/cuda:12.8.0-base-ubuntu22.04                                                                                 0.1s
 => [internal] load metadata for docker.io/rancher/k3s:v1.31.5-k3s1                                                                                            0.1s
 => [internal] load .dockerignore                                                                                                                              0.0s
 => => transferring context: 145B                                                                                                                              0.0s
 => [stage-1 1/3] FROM docker.io/nvidia/cuda:12.8.0-base-ubuntu22.04@sha256:12242992c121f6cab0ca11bccbaaf757db893b3065d7db74b933e59f321b2cf4                   0.0s
 => [k3s 1/1] FROM docker.io/rancher/k3s:v1.31.5-k3s1@sha256:53cf744fe2fabf140cee240d2db70d13a4f2d98f1a13c98f58f457f423096917                                  0.0s
 => CACHED [stage-1 2/3] RUN <<EOT (set -eu...)                                                                                                                0.0s
 => CACHED [stage-1 3/3] COPY --from=k3s /bin /bin                                                                                                             0.0s
 => exporting to image                                                                                                                                         0.0s
 => => exporting layers                                                                                                                                        0.0s
 => => writing image sha256:15c9fb2a0905e1ef51311de9d5c3820f98eae9069d08d437d77847a991781701                                                                   0.0s
 => => naming to localhost:5001/rancher/k3s:latest                                                                                                             0.0s
INFO[0000] Using config file hack/k3d.yaml (k3d.io/v1alpha5#simple) 
INFO[0000] portmapping '1993:1993' targets the loadbalancer: defaulting to [servers:*:proxy agents:*:proxy] 
INFO[0000] portmapping '1994:1994' targets the loadbalancer: defaulting to [servers:*:proxy agents:*:proxy] 
INFO[0000] portmapping '9000:9000' targets the loadbalancer: defaulting to [servers:*:proxy agents:*:proxy] 
INFO[0000] portmapping '8008:8008' targets the loadbalancer: defaulting to [servers:*:proxy agents:*:proxy] 
INFO[0000] portmapping '9900:9900' targets the loadbalancer: defaulting to [servers:*:proxy agents:*:proxy] 
INFO[0000] Prep: Network                                
INFO[0000] Created network 'k3d-beta9'                  
INFO[0000] Created image volume k3d-beta9-images        
INFO[0000] Creating node 'registry.localhost'           
INFO[0000] Successfully created registry 'registry.localhost' 
INFO[0000] Starting new tools node...                   
INFO[0000] Starting node 'k3d-beta9-tools'              
INFO[0001] Creating node 'k3d-beta9-server-0'           
INFO[0001] Creating LoadBalancer 'k3d-beta9-serverlb'   
INFO[0001] Using the k3d-tools node to gather environment information 
INFO[0001] HostIP: using network gateway 172.18.0.1 address 
INFO[0001] Starting cluster 'beta9'                     
INFO[0001] Starting servers...                          
INFO[0001] Starting node 'k3d-beta9-server-0'           
INFO[0003] All agents already running.                  
INFO[0003] Starting helpers...                          
INFO[0003] Starting node 'registry.localhost'           
INFO[0003] Starting node 'k3d-beta9-serverlb'           
INFO[0010] Injecting records for hostAliases (incl. host.k3d.internal) and for 3 network members into CoreDNS configmap... 
INFO[0013] Cluster 'beta9' created successfully!        
INFO[0013] You can now use it like this:                
kubectl cluster-info
namespace/beta9 created
Property "contexts.k3d-beta9.namespace" set.
 ✓  Using beta9 @ k3d-beta9
shadeform@shadecloud:~/beta9$ make start
 i  Using beta9 @ k3d-beta9 as context
 x  Application 'beta9-gateway' not found in namespace 'beta9'
    Verify that your application is running and your okteto context is pointing to the right namespace
    Or set the 'autocreate' field in your okteto manifest if you want to create a standalone development container
    More information is available here: https://okteto.com/docs/reference/okteto-cli/#up
make: *** [Makefile:48: start] Error 1
shadeform@shadecloud:~/beta9$ kustomize build --enable-helm manifests/kustomize/overlays/cluster-dev | kubectl apply -f-
Warning: resource namespaces/beta9 is missing the kubectl.kubernetes.io/last-applied-configuration annotation which is required by kubectl apply. kubectl apply should only be used on resources created declaratively by either kubectl create --save-config or kubectl apply. The missing annotation will be patched automatically.
namespace/beta9 configured
namespace/monitoring created
namespace/tailscale created
serviceaccount/gateway created
serviceaccount/juicefs-redis-master created
serviceaccount/localstack created
serviceaccount/postgresql created
serviceaccount/redis-master created
serviceaccount/elasticsearch-kibana created
serviceaccount/elasticsearch-master created
serviceaccount/fluent-bit created
serviceaccount/grafana created
serviceaccount/victoria-metrics-agent created
serviceaccount/victoria-metrics-single created
role.rbac.authorization.k8s.io/localstack created
role.rbac.authorization.k8s.io/grafana created
clusterrole.rbac.authorization.k8s.io/gateway-role created
clusterrole.rbac.authorization.k8s.io/victoria-metrics-agent-clusterrole created
clusterrole.rbac.authorization.k8s.io/fluent-bit created
clusterrole.rbac.authorization.k8s.io/grafana-clusterrole created
clusterrole.rbac.authorization.k8s.io/victoria-metrics-single-clusterrole created
rolebinding.rbac.authorization.k8s.io/localstack created
rolebinding.rbac.authorization.k8s.io/grafana created
clusterrolebinding.rbac.authorization.k8s.io/gateway-role-binding created
clusterrolebinding.rbac.authorization.k8s.io/victoria-metrics-agent-clusterrolebinding created
clusterrolebinding.rbac.authorization.k8s.io/fluent-bit created
clusterrolebinding.rbac.authorization.k8s.io/grafana-clusterrolebinding created
clusterrolebinding.rbac.authorization.k8s.io/victoria-metrics-single-clusterrolebinding created
configmap/juicefs-redis-configuration created
configmap/juicefs-redis-health created
configmap/juicefs-redis-scripts created
configmap/localstack-init-scripts-config created
configmap/postgresql-extended-configuration created
configmap/postgresql-init-scripts created
configmap/redis-configuration created
configmap/redis-health created
configmap/redis-scripts created
configmap/elasticsearch-kibana-conf created
configmap/fluent-bit created
configmap/grafana created
configmap/victoria-metrics-agent-config created
secret/beta9-config created
secret/gateway-sa-token created
secret/juicefs-secret created
secret/postgresql created
secret/fluent-bit-g444h6dhhm created
secret/grafana created
service/beta9-gateway created
service/juicefs-redis-headless created
service/juicefs-redis-master created
service/juicefs-s3-gateway created
service/localstack created
service/postgresql created
service/postgresql-hl created
service/redis-headless created
service/redis-master created
service/elasticsearch created
service/elasticsearch-kibana created
service/elasticsearch-master-hl created
service/fluent-bit created
service/grafana created
service/victoria-logs-single created
service/victoria-metrics-single created
persistentvolumeclaim/beta9-images created
persistentvolumeclaim/localstack created
persistentvolumeclaim/elasticsearch-kibana created
deployment.apps/beta9-gateway created
deployment.apps/juicefs-s3-gateway created
deployment.apps/localstack created
deployment.apps/elasticsearch-kibana created
deployment.apps/grafana created
deployment.apps/victoria-metrics-agent created
statefulset.apps/juicefs-redis-master created
statefulset.apps/postgresql created
statefulset.apps/redis-master created
statefulset.apps/elasticsearch-master created
statefulset.apps/victoria-logs-single created
statefulset.apps/victoria-metrics-single created
poddisruptionbudget.policy/elasticsearch-kibana created
poddisruptionbudget.policy/elasticsearch-master created
daemonset.apps/nvidia-device-plugin-daemonset created
daemonset.apps/fluent-bit created
networkpolicy.networking.k8s.io/elasticsearch-kibana created
networkpolicy.networking.k8s.io/elasticsearch-master created
shadeform@shadecloud:~/beta9$ make start
 i  Using beta9 @ k3d-beta9 as context
 ✓  Images successfully pulled
 x  Couldn't connect to your development container: local port 1994 is already in-use in your local machine: port is already allocated
    Find additional logs at: /home/shadeform/.okteto/beta9/beta9-gateway/okteto.log
make: *** [Makefile:48: start] Error 1
shadeform@shadecloud:~/beta9$ make stop
cd hack && okteto down --file okteto.yaml
 i  Using beta9 @ k3d-beta9 as context
 ✓  Development container 'beta9-gateway' deactivated
shadeform@shadecloud:~/beta9$ ls
CODE_OF_CONDUCT.md  LICENSE   README.md  cmd     docker  e2e     go.sum  manifests  proto  setupgo.sh
CONTRIBUTING.md     Makefile  bin        deploy  docs    go.mod  hack    pkg        sdk    static
shadeform@shadecloud:~/beta9$ kustomize build --enable-helm manifests/kustomize/overlays/cluster-dev | kubectl apply -f-
namespace/beta9 unchanged
namespace/monitoring unchanged
namespace/tailscale unchanged
serviceaccount/gateway unchanged
serviceaccount/juicefs-redis-master unchanged
serviceaccount/localstack unchanged
serviceaccount/postgresql unchanged
serviceaccount/redis-master unchanged
serviceaccount/elasticsearch-kibana unchanged
serviceaccount/elasticsearch-master unchanged
serviceaccount/fluent-bit unchanged
serviceaccount/grafana unchanged
serviceaccount/victoria-metrics-agent unchanged
serviceaccount/victoria-metrics-single unchanged
role.rbac.authorization.k8s.io/localstack unchanged
role.rbac.authorization.k8s.io/grafana configured
clusterrole.rbac.authorization.k8s.io/gateway-role unchanged
clusterrole.rbac.authorization.k8s.io/victoria-metrics-agent-clusterrole unchanged
clusterrole.rbac.authorization.k8s.io/fluent-bit unchanged
clusterrole.rbac.authorization.k8s.io/grafana-clusterrole configured
clusterrole.rbac.authorization.k8s.io/victoria-metrics-single-clusterrole unchanged
rolebinding.rbac.authorization.k8s.io/localstack unchanged
rolebinding.rbac.authorization.k8s.io/grafana unchanged
clusterrolebinding.rbac.authorization.k8s.io/gateway-role-binding unchanged
clusterrolebinding.rbac.authorization.k8s.io/victoria-metrics-agent-clusterrolebinding unchanged
clusterrolebinding.rbac.authorization.k8s.io/fluent-bit unchanged
clusterrolebinding.rbac.authorization.k8s.io/grafana-clusterrolebinding unchanged
clusterrolebinding.rbac.authorization.k8s.io/victoria-metrics-single-clusterrolebinding unchanged
configmap/juicefs-redis-configuration unchanged
configmap/juicefs-redis-health unchanged
configmap/juicefs-redis-scripts unchanged
configmap/localstack-init-scripts-config unchanged
configmap/postgresql-extended-configuration unchanged
configmap/postgresql-init-scripts unchanged
configmap/redis-configuration unchanged
configmap/redis-health unchanged
configmap/redis-scripts unchanged
configmap/elasticsearch-kibana-conf unchanged
configmap/fluent-bit unchanged
configmap/grafana unchanged
configmap/victoria-metrics-agent-config unchanged
secret/beta9-config configured
secret/gateway-sa-token unchanged
secret/juicefs-secret unchanged
secret/postgresql unchanged
secret/fluent-bit-g444h6dhhm unchanged
secret/grafana configured
service/beta9-gateway unchanged
service/juicefs-redis-headless unchanged
service/juicefs-redis-master configured
service/juicefs-s3-gateway unchanged
service/localstack configured
service/postgresql configured
service/postgresql-hl unchanged
service/redis-headless unchanged
service/redis-master configured
service/elasticsearch configured
service/elasticsearch-kibana configured
service/elasticsearch-master-hl unchanged
service/fluent-bit unchanged
service/grafana unchanged
service/victoria-logs-single unchanged
service/victoria-metrics-single unchanged
persistentvolumeclaim/beta9-images unchanged
persistentvolumeclaim/localstack unchanged
persistentvolumeclaim/elasticsearch-kibana unchanged
deployment.apps/beta9-gateway configured
deployment.apps/juicefs-s3-gateway configured
deployment.apps/localstack configured
deployment.apps/elasticsearch-kibana configured
deployment.apps/grafana configured
deployment.apps/victoria-metrics-agent unchanged
statefulset.apps/juicefs-redis-master configured
statefulset.apps/postgresql configured
statefulset.apps/redis-master configured
statefulset.apps/elasticsearch-master configured
statefulset.apps/victoria-logs-single configured
statefulset.apps/victoria-metrics-single configured
poddisruptionbudget.policy/elasticsearch-kibana configured
poddisruptionbudget.policy/elasticsearch-master configured
daemonset.apps/nvidia-device-plugin-daemonset unchanged
daemonset.apps/fluent-bit configured
networkpolicy.networking.k8s.io/elasticsearch-kibana configured
networkpolicy.networking.k8s.io/elasticsearch-master configured
shadeform@shadecloud:~/beta9$ make start
 i  Using beta9 @ k3d-beta9 as context
 ✓  Images successfully pulled
 x  Couldn't connect to your development container: local port 1994 is already in-use in your local machine: port is already allocated
    Find additional logs at: /home/shadeform/.okteto/beta9/beta9-gateway/okteto.log
make: *** [Makefile:48: start] Error 1
shadeform@shadecloud:~/beta9$ make k3d-down
bash bin/k3d.sh down
INFO[0000] Using config file hack/k3d.yaml (k3d.io/v1alpha5#simple) 
INFO[0000] Deleting cluster 'beta9'                     
INFO[0009] Deleting cluster network 'k3d-beta9'         
INFO[0009] Deleting 1 attached volumes...               
INFO[0009] Removing cluster details from default kubeconfig... 
INFO[0009] Removing standalone kubeconfig file (if there is one)... 
INFO[0009] Successfully deleted cluster beta9!          
shadeform@shadecloud:~/beta9$ k3d-up
k3d-up: command not found
shadeform@shadecloud:~/beta9$ make k3d-up
bash bin/k3d.sh up
[+] Building 0.3s (11/11) FINISHED                                                                        docker:default
 => [internal] load build definition from Dockerfile.k3d                                                            0.0s
 => => transferring dockerfile: 1.38kB                                                                              0.0s
 => resolve image config for docker-image://docker.io/docker/dockerfile:1.6                                         0.1s
 => CACHED docker-image://docker.io/docker/dockerfile:1.6@sha256:ac85f380a63b13dfcefa89046420e1781752bab202122f8f5  0.0s
 => [internal] load metadata for docker.io/nvidia/cuda:12.8.0-base-ubuntu22.04                                      0.1s
 => [internal] load metadata for docker.io/rancher/k3s:v1.31.5-k3s1                                                 0.1s
 => [internal] load .dockerignore                                                                                   0.0s
 => => transferring context: 145B                                                                                   0.0s
 => [stage-1 1/3] FROM docker.io/nvidia/cuda:12.8.0-base-ubuntu22.04@sha256:12242992c121f6cab0ca11bccbaaf757db893b  0.0s
 => [k3s 1/1] FROM docker.io/rancher/k3s:v1.31.5-k3s1@sha256:53cf744fe2fabf140cee240d2db70d13a4f2d98f1a13c98f58f45  0.0s
 => CACHED [stage-1 2/3] RUN <<EOT (set -eu...)                                                                     0.0s
 => CACHED [stage-1 3/3] COPY --from=k3s /bin /bin                                                                  0.0s
 => exporting to image                                                                                              0.0s
 => => exporting layers                                                                                             0.0s
 => => writing image sha256:15c9fb2a0905e1ef51311de9d5c3820f98eae9069d08d437d77847a991781701                        0.0s
 => => naming to localhost:5001/rancher/k3s:latest                                                                  0.0s
INFO[0000] Using config file hack/k3d.yaml (k3d.io/v1alpha5#simple) 
INFO[0000] portmapping '1993:1993' targets the loadbalancer: defaulting to [servers:*:proxy agents:*:proxy] 
INFO[0000] portmapping '1994:1994' targets the loadbalancer: defaulting to [servers:*:proxy agents:*:proxy] 
INFO[0000] portmapping '9000:9000' targets the loadbalancer: defaulting to [servers:*:proxy agents:*:proxy] 
INFO[0000] portmapping '8008:8008' targets the loadbalancer: defaulting to [servers:*:proxy agents:*:proxy] 
INFO[0000] portmapping '9900:9900' targets the loadbalancer: defaulting to [servers:*:proxy agents:*:proxy] 
INFO[0000] Prep: Network                                
INFO[0000] Created network 'k3d-beta9'                  
INFO[0000] Created image volume k3d-beta9-images        
INFO[0000] Creating node 'registry.localhost'           
INFO[0000] Successfully created registry 'registry.localhost' 
INFO[0000] Starting new tools node...                   
INFO[0000] Starting node 'k3d-beta9-tools'              
INFO[0001] Creating node 'k3d-beta9-server-0'           
INFO[0001] Creating LoadBalancer 'k3d-beta9-serverlb'   
INFO[0001] Using the k3d-tools node to gather environment information 
INFO[0001] HostIP: using network gateway 172.18.0.1 address 
INFO[0001] Starting cluster 'beta9'                     
INFO[0001] Starting servers...                          
INFO[0001] Starting node 'k3d-beta9-server-0'           
INFO[0004] All agents already running.                  
INFO[0004] Starting helpers...                          
INFO[0004] Starting node 'registry.localhost'           
INFO[0004] Starting node 'k3d-beta9-serverlb'           
INFO[0011] Injecting records for hostAliases (incl. host.k3d.internal) and for 3 network members into CoreDNS configmap... 
INFO[0014] Cluster 'beta9' created successfully!        
INFO[0014] You can now use it like this:                
kubectl cluster-info
namespace/beta9 created
Property "contexts.k3d-beta9.namespace" set.
 ✓  Using beta9 @ k3d-beta9
shadeform@shadecloud:~/beta9$ kustomize build --enable-helm manifests/kustomize/overlays/cluster-dev | kubectl apply -f-
Warning: resource namespaces/beta9 is missing the kubectl.kubernetes.io/last-applied-configuration annotation which is required by kubectl apply. kubectl apply should only be used on resources created declaratively by either kubectl create --save-config or kubectl apply. The missing annotation will be patched automatically.
namespace/beta9 configured
namespace/monitoring created
namespace/tailscale created
serviceaccount/gateway created
serviceaccount/juicefs-redis-master created
serviceaccount/localstack created
serviceaccount/postgresql created
serviceaccount/redis-master created
serviceaccount/elasticsearch-kibana created
serviceaccount/elasticsearch-master created
serviceaccount/fluent-bit created
serviceaccount/grafana created
serviceaccount/victoria-metrics-agent created
serviceaccount/victoria-metrics-single created
role.rbac.authorization.k8s.io/localstack created
role.rbac.authorization.k8s.io/grafana created
clusterrole.rbac.authorization.k8s.io/gateway-role created
clusterrole.rbac.authorization.k8s.io/victoria-metrics-agent-clusterrole created
clusterrole.rbac.authorization.k8s.io/fluent-bit created
clusterrole.rbac.authorization.k8s.io/grafana-clusterrole created
clusterrole.rbac.authorization.k8s.io/victoria-metrics-single-clusterrole created
rolebinding.rbac.authorization.k8s.io/localstack created
rolebinding.rbac.authorization.k8s.io/grafana created
clusterrolebinding.rbac.authorization.k8s.io/gateway-role-binding created
clusterrolebinding.rbac.authorization.k8s.io/victoria-metrics-agent-clusterrolebinding created
clusterrolebinding.rbac.authorization.k8s.io/fluent-bit created
clusterrolebinding.rbac.authorization.k8s.io/grafana-clusterrolebinding created
clusterrolebinding.rbac.authorization.k8s.io/victoria-metrics-single-clusterrolebinding created
configmap/juicefs-redis-configuration created
configmap/juicefs-redis-health created
configmap/juicefs-redis-scripts created
configmap/localstack-init-scripts-config created
configmap/postgresql-extended-configuration created
configmap/postgresql-init-scripts created
configmap/redis-configuration created
configmap/redis-health created
configmap/redis-scripts created
configmap/elasticsearch-kibana-conf created
configmap/fluent-bit created
configmap/grafana created
configmap/victoria-metrics-agent-config created
secret/beta9-config created
secret/gateway-sa-token created
secret/juicefs-secret created
secret/postgresql created
secret/fluent-bit-g444h6dhhm created
secret/grafana created
service/beta9-gateway created
service/juicefs-redis-headless created
service/juicefs-redis-master created
service/juicefs-s3-gateway created
service/localstack created
service/postgresql created
service/postgresql-hl created
service/redis-headless created
service/redis-master created
service/elasticsearch created
service/elasticsearch-kibana created
service/elasticsearch-master-hl created
service/fluent-bit created
service/grafana created
service/victoria-logs-single created
service/victoria-metrics-single created
persistentvolumeclaim/beta9-images created
persistentvolumeclaim/localstack created
persistentvolumeclaim/elasticsearch-kibana created
deployment.apps/beta9-gateway created
deployment.apps/juicefs-s3-gateway created
deployment.apps/localstack created
deployment.apps/elasticsearch-kibana created
deployment.apps/grafana created
deployment.apps/victoria-metrics-agent created
statefulset.apps/juicefs-redis-master created
statefulset.apps/postgresql created
statefulset.apps/redis-master created
statefulset.apps/elasticsearch-master created
statefulset.apps/victoria-logs-single created
statefulset.apps/victoria-metrics-single created
poddisruptionbudget.policy/elasticsearch-kibana created
poddisruptionbudget.policy/elasticsearch-master created
daemonset.apps/nvidia-device-plugin-daemonset created
daemonset.apps/fluent-bit created
networkpolicy.networking.k8s.io/elasticsearch-kibana created
networkpolicy.networking.k8s.io/elasticsearch-master created
shadeform@shadecloud:~/beta9$ make start
 i  Using beta9 @ k3d-beta9 as context
 ✓  Images successfully pulled
 x  Couldn't connect to your development container: local port 1994 is already in-use in your local machine: port is already allocated
    Find additional logs at: /home/shadeform/.okteto/beta9/beta9-gateway/okteto.log
make: *** [Makefile:48: start] Error 1
shadeform@shadecloud:~/beta9$ 