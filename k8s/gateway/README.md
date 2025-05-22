Как запустить?

1. quickstart.sh
2. kubectl create namespace envoy
3. kubectl apply -f gateway.yaml
4. kubectl create namespace website # если еще не установили helm вебсайта
5. kubectl apply -f frontend-gateway.yaml