# Examples

## How to use

Apply everything using curl

### Creating a new project

You can change the variables of the `new_project.json` file before applying

```bash
curl -XPOST -H 'Content-Type: application/json' \
	localhost:8080/projects \
	-d@examples/new_project.json
```

### Creating a new subscription

Change the `<id>` field in the URI to match the project you created.

```bash
PROJECT_ID=<id>
curl -XPOST -H 'Content-Type: application/json' \
	localhost:8080/projects/${PROJECT_ID}/subscriptions \
	-d@examples/new_subscription.json
```

If you don't record your project's id, run

```bash
curl -XGET localhost:8080/projects
```
