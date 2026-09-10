### K8s Config ###

# Uncomment to use secrets
# k8s_yaml('./infra/development/k8s/secrets.yaml')

k8s_yaml('./infra/development/k8s/app-config.yaml')

### End of K8s Config ###

# Local restart wrapper so we do not pull tiltdev/restart-helper from Docker Hub.
restart_wrapper_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/tilt-restart-wrapper ./infra/development/restart'
local_resource(
  'restart-wrapper-compile',
  restart_wrapper_compile_cmd,
  deps=['./infra/development/restart'], labels="compiles")

### API Gateway ###

gateway_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/api-gateway ./services/api-gateway'
if os.name == 'nt':
  gateway_compile_cmd = './infra/development/docker/api-gateway-build.bat'

local_resource(
  'api-gateway-compile',
  gateway_compile_cmd,
  deps=['./services/api-gateway', './shared'], labels="compiles")


docker_build(
  'ride-sharing/api-gateway',
  '.',
  entrypoint=['/app/build/tilt-restart-wrapper', '--watch_file=/tmp/.restart-proc', '/app/build/api-gateway'],
  dockerfile='./infra/development/docker/api-gateway.Dockerfile',
  only=[
    './build/api-gateway',
    './build/tilt-restart-wrapper',
    './shared',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
    run('date > /tmp/.restart-proc'),
  ],
)

k8s_yaml('./infra/development/k8s/api-gateway-deployment.yaml')
k8s_resource('api-gateway', port_forwards=8081,
             resource_deps=['api-gateway-compile', 'restart-wrapper-compile'], labels="services")
### End of API Gateway ###
### Trip Service ###

trip_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/trip-service ./services/trip-service/cmd/main.go'
if os.name == 'nt':
  trip_compile_cmd = './infra/development/docker/trip-build.bat'

local_resource(
  'trip-service-compile',
  trip_compile_cmd,
  deps=['./services/trip-service', './shared'], labels="compiles")

docker_build(
  'ride-sharing/trip-service',
  '.',
  entrypoint=['/app/build/tilt-restart-wrapper', '--watch_file=/tmp/.restart-proc', '/app/build/trip-service'],
  dockerfile='./infra/development/docker/trip-service.Dockerfile',
  only=[
    './build/trip-service',
    './build/tilt-restart-wrapper',
    './shared',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
    run('date > /tmp/.restart-proc'),
  ],
)

k8s_yaml('./infra/development/k8s/trip-service-deployment.yaml')
k8s_resource('trip-service', resource_deps=['trip-service-compile', 'restart-wrapper-compile'], labels="services")

### End of Trip Service ###
### Web Frontend ###

docker_build(
  'ride-sharing/web',
  '.',
  dockerfile='./infra/development/docker/web.Dockerfile',
)

k8s_yaml('./infra/development/k8s/web-deployment.yaml')
k8s_resource('web', port_forwards=3000, labels="frontend")

### End of Web Frontend ###
