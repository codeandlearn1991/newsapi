version_settings(constraint='>=0.33.21')

load('ext://helm_resource', 'helm_resource', 'helm_repo')
load('ext://secret', 'secret_from_dict')

helm_repo('bitnami', 'https://charts.bitnami.com/bitnami', resource_name='helm-repo-bitnami')
helm_resource('postgres', 'bitnami/postgresql', resource_deps=['helm-repo-bitnami'], flags=[
  '--set', '--image.tag=16.0.0',
  '--set', 'auth.postgresPassword=password'
])
k8s_resource('postgres', port_forwards=[
  port_forward(15432, 5432, name='postgres'),
])

k8s_yaml(secret_from_dict(name='database-secret', namespace='news-service', inputs={
  'host': 'postgres-postgresql.default.svc.cluster.local',
  'dbname': 'postgres',
  'password': 'password',
  'port': '5432',
  'user': 'postgres'
}))

docker_build('news-api-server', '.', dockerfile='Dockerfile', build_args={"APP": "api-server"})
docker_build('news-migrate', '.', dockerfile='Dockerfile', build_args={"APP": "migrate"})

k8s_yaml([
  'deployment/namespace.yaml', 
  'deployment/deployment.yaml', 
  'deployment/service.yaml',
  'deployment/migrate.yaml'])

k8s_resource(workload='news-api-server', port_forwards=[
  port_forward(8080, 8080, name='news-api-server')
])