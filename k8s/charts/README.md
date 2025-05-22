```bash

helm install financial-tracker ./financial-tracker --namespace default

```

```bash

helm install website ./website --namespace website --create-namespace

```

```bash

helm upgrade --install website ./website --namespace website

```
