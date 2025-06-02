### Подключение Prometheus

```bash

helm install monitoring-stack prometheus-community/kube-prometheus-stack \
  -n monitoring \
  -f helm/kube-prometheus-stack.yaml

```

### Просмотр собранных метрик 

```bash

kubectl port-forward svc/monitoring-stack-kube-prom-prometheus 9090:9090 -n monitoring

```

#### После этого Prometheus будет доступен по http://localhost:9090
##### Можно ввести интересующую Вас метрику в поле ввода и Вы получите данные о ней:

1. expense_publisher_requests_total: Инкрементируется при получении каждого gRPC-вызова.
2. expense_publisher_errors_total: Инкрементируется, когда в обработчике gRPC-метода возникает ошибка (например: невалидный JWT, ошибка валидации тела сообщения или внутренняя ошибка при отправке в Kafka).
3. expense_publisher_request_duration_seconds: Гистограмма длительности обработки запроса (в секундах).
4. expensereader_grpc_requests_total{method="GetReport"}: Инкрементируется при каждом вызове gRPC-метода GetReport.
5. expensereader_grpc_errors_total{method="GetReport"}: Инкрементируется, когда GetReport возвращает ошибку (например, JWT недействителен или внутренняя ошибка БД).
6. expensereader_grpc_duration_seconds{method="GetReport"}: Гистограмма длительности одного выполнения GetReport (получения агрегации из БД, формирования ответа).
7. expensereader_grpc_requests_total{method="GetBills"}: Инкрементируется при каждом вызове метода GetBills.
8. expensereader_grpc_errors_total{method="GetBills"}: Инкрементируется, если в GetBills произошла ошибка (например, ошибка работы с БД).
9. expensereader_grpc_duration_seconds{method="GetBills"}: Гистограмма времени выполнения GetBills (поиск строк в таблице по ключу).
10. userservice_grpc_requests_total{method="Register"}: Инкрементируется при каждом вызове метода регистрации нового пользователя.
11. userservice_grpc_errors_total{method="Register"}: Инкрементируется при каждой неуспешой попытокой регистрации (невалидные данные, конфликт username и т. д.).
12. userservice_grpc_duration_seconds{method="Register"}: Гистограмма времени длительности операции регистрации (хеширование пароля, запись в БД, генерация JWT).
13. userservice_grpc_requests_total{method="Login"}: Инкрементируется при каждом вызове метода логина.
14. userservice_grpc_errors_total{method="Login"}: Инкрементируется при каждой неуспешой попытоке входа (неправильный пароль, отсутствующий пользователь и т. д.).
15. userservice_grpc_duration_seconds{method="Login"}: Гистограмма времени процесса аутентификации и генерации JWT.
16. expense_writer_bills_processed_total: Инкрементируется каждый раз, когда батч из Kafka успешно вставился в БД.
17. expense_writer_bills_failed_total: Инкрементируется, если при обработке какого-то батча произошла ошибка (например, связь с БД упала, нарушение ограничений и т. д.).
18. expense_writer_bill_process_latency_seconds: Измеряет время исполнения batchCreateBills(...).
