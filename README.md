# Statictask Newsletter

The simplest newsletter service for static sites.

This project aims to cover the necessity of a simple newsletter service
for small and medium blogs or any service that doesn't need all the
fency resources of another expensive SaaS.

The design is simple enough to be deployed by any programmer in a any
infrastructure. The unique dependency is Postgres, so two containers
might be enough depending on your usage.

We're also going to deploy this system in a form of a simple SaaS for
people that doesn't know how to deal with all cloud stuff. This is a
futher step since what we have today doesn't support multi-tenant
workloads.

The idea is to make it free for most users that have a very simple
use-cases and sell honest plans just for heavy users. All the features
might be always available for everyone regardless of the plan being
paid or not.

Wait for good news at https://statictask.io

## Development

### Environment configuration

In order to deploy deps locally run

    make fix

It will install development and system deps.

### Building and running with Docker and Docker Compose

In order to make it easy to test locally new features, you can use docker
to build the current code and run it with compose. The `docker-compose.yml`
file has all the necessary environment variables and configurations you
need by default. The only additional thing you'll have to do is to apply
the migrations to the brand new database.

    # terminal 1
    make run

    # terminal 2
    make migrate

The service is accessible under `http://localhost:8080/`

## Production

### Building production-ready Docker images

Note: you'll need permissions to submit new images to our registry.

The default docker registry used by the organization is `statictask/newsletter`.
Docker tags reflect the version of the code (git tags). The `latest` tag
is always pointing to the newer git tag.

To generate a new image, please make sure the code is correct and tested,
we don't have a CI/CD setup yet to avoid more infrastructure costs. Then
create a new git tag annotated with a brief description about the new release.

    git tag -a -m "This is an example" 1.2.3

We use semmantic versioning, create new tags according to this rule.

Now all you have to do is to build and release a new docker image to our
registry.

    make docker-release

### Deploy to K8s

The system is prepared to be deployed on a K8s cluster using Helm. So make
sure you have `helm` installed.

You need to define in your `.env` some variables to fill necessary information
to create new helm releases.

```
KUBERNETES_IMAGE_VERSION="0.0.4"
NEWSLETTER_POSTGRES_HOST="db.example.com"
NEWSLETTER_POSTGRES_PASSWORD="example"
```

Then let make do the job for you

    make kubernetes-deploy

### Migrating production databases

Newsletter's docker image has a `/migrate` command that you can use to migrate
database from your workloads on Kubernetes.

    kubectl exec -it -n newsletter deploy/newsletter -- /migrate up
