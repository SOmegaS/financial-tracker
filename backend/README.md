# Бэкенд

Бэкенд часть финансового трекера.

Описания основных сервисов в соответствующих папках


build:

```bash
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t harbor.devops-teta.ru/financial-tracker/expensewriter:0.1.0 \
  expensewriter \
  --push
```
```bash
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t harbor.devops-teta.ru/financial-tracker/userservice:0.1.0 \
  expensepublisher \
  --push
```
```bash
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t harbor.devops-teta.ru/financial-tracker/expensepublisher:0.1.0 \
  expensepublisher \
  --push
```
```bash
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t harbor.devops-teta.ru/financial-tracker/expensereader:0.1.0 \
  expensereader \
  --push
```