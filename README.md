# go-musthave-metrics-tpl

Шаблон репозитория для трека «Сервер сбора метрик и алертинга».

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m v2 template https://github.com/Yandex-Practicum/go-musthave-metrics-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/v2 .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).

## Структура проекта

Приведённая в этом репозитории структура проекта является рекомендуемой, но не обязательной.

Это лишь пример организации кода, который поможет вам в реализации сервиса.

При необходимости можно вносить изменения в структуру проекта, использовать любые библиотеки и предпочитаемые структурные паттерны организации кода приложения, например:
- **DDD** (Domain-Driven Design)
- **Clean Architecture**
- **Hexagonal Architecture**
- **Layered Architecture**

## PPROF


File: server
Build ID: a6c6ea25e169705060e11a1d8c8dc245a32d1338
Type: alloc_space
Time: 2026-07-04 15:33:07 MSK
Showing nodes accounting for -15232.43MB, 12.41% of 122699.25MB total
Dropped 385 nodes (cum <= 613.50MB)
      flat  flat%   sum%        cum   cum%
-21957.54MB 17.90% 17.90% -19785.30MB 16.13%  github.com/WorstOfAny/go-musthave-metrics-tpl/internal/handler.(*metricsController).listAll-range1
 2437.16MB  1.99% 15.91%  2926.05MB  2.38%  compress/flate.NewWriter (inline)
 2436.24MB  1.99% 13.92%  2436.24MB  1.99%  bytes.growSlice
-1829.13MB  1.49% 15.41% -17752.99MB 14.47%  github.com/WorstOfAny/go-musthave-metrics-tpl/internal/handler.(*metricsController).listAll
 1677.20MB  1.37% 14.05%  3473.30MB  2.83%  github.com/jackc/pgx/v5.AppendRows[go.shape.struct { ID string "json:\"id\" db:\"id\""; MType string "json:\"type\" db:\"mtype\""; Delta *int64 "json:\"delta,omitempty\" db:\"delta\""; Value *float64 "json:\"value,omitempty\" db:\"value\""; Hash string "json:\"hash,omitempty\" db:\"hash\"" },go.shape.[]go.shape.struct { ID string "json:\"id\" db:\"id\""; MType string "json:\"type\" db:\"mtype\""; Delta *int64 "json:\"delta,omitempty\" db:\"delta\""; Value *float64 "json:\"value,omitempty\" db:\"value\""; Hash string "json:\"hash,omitempty\" db:\"hash\"" }]
 -755.04MB  0.62% 14.66%  -852.13MB  0.69%  fmt.Sprintf
  610.55MB   0.5% 14.17%   610.55MB   0.5%  github.com/jackc/pgx/v5.setupStructScanTargets
  558.53MB  0.46% 13.71%  1742.09MB  1.42%  github.com/jackc/pgx/v5.RowToStructByName[go.shape.struct { ID string "json:\"id\" db:\"id\""; MType string "json:\"type\" db:\"mtype\""; Delta *int64 "json:\"delta,omitempty\" db:\"delta\""; Value *float64 "json:\"value,omitempty\" db:\"value\""; Hash string "json:\"hash,omitempty\" db:\"hash\"" }]
  472.01MB  0.38% 13.33%   604.92MB  0.49%  github.com/WorstOfAny/go-musthave-metrics-tpl/internal/model.Metrics.String
  446.36MB  0.36% 12.96%   446.36MB  0.36%  compress/flate.(*compressor).initDeflate (inline)
  302.01MB  0.25% 12.72%   302.01MB  0.25%  internal/bytealg.MakeNoZero
     215MB  0.18% 12.54%      215MB  0.18%  github.com/jackc/pgx/v5/pgtype.scanPlanString.Scan
 -143.31MB  0.12% 12.66%  -143.31MB  0.12%  sync.(*Pool).pinSlow
  124.52MB   0.1% 12.56%   118.02MB 0.096%  context.(*cancelCtx).propagateCancel
   90.52MB 0.074% 12.48%    90.52MB 0.074%  net/textproto.readMIMEHeader
   89.53MB 0.073% 12.41%    89.53MB 0.073%  net/http.Header.Clone (inline)
  -77.51MB 0.063% 12.47%   -77.51MB 0.063%  github.com/jackc/pgx/v5.rawState
   62.52MB 0.051% 12.42%   156.48MB  0.13%  net/http.readRequest
  -61.51MB  0.05% 12.47%  -271.05MB  0.22%  github.com/WorstOfAny/go-musthave-metrics-tpl/internal/repository.(*dbDecorator).BulkSet.func1
   49.51MB  0.04% 12.43%   141.34MB  0.12%  net/http.(*conn).readRequest
   42.50MB 0.035% 12.40%  2478.75MB  2.02%  bytes.(*Buffer).grow
  -18.50MB 0.015% 12.41%  -167.53MB  0.14%  github.com/jackc/pgx/v5.rewriteQuery
   16.49MB 0.013% 12.40% -18290.08MB 14.91%  github.com/WorstOfAny/go-musthave-metrics-tpl/internal/handler.(*metricsController).logRequest-fm.(*metricsController).logRequest.func1
  -15.56MB 0.013% 12.41%  -454.29MB  0.37%  github.com/WorstOfAny/go-musthave-metrics-tpl/internal/handler.(*metricsController).updates
      -7MB 0.0057% 12.42%      280MB  0.23%  github.com/jackc/pgx/v5.(*baseRows).Scan
    2.50MB 0.002% 12.41%  2884.94MB  2.35%  github.com/WorstOfAny/go-musthave-metrics-tpl/internal/client.(*Client).Post
    2.50MB 0.002% 12.41%  3672.39MB  2.99%  github.com/WorstOfAny/go-musthave-metrics-tpl/internal/repository.(*dbDecorator).All
   -1.50MB 0.0012% 12.41% -18103.29MB 14.75%  net/http.(*conn).serve
   -1.50MB 0.0012% 12.41%    91.02MB 0.074%  context.AfterFunc
         0     0% 12.41%  2454.78MB  2.00%  bytes.(*Buffer).Write
         0     0% 12.41%   488.89MB   0.4%  compress/flate.(*compressor).init
         0     0% 12.41%  2921.55MB  2.38%  compress/gzip.(*Writer).Write
         0     0% 12.41%   -77.64MB 0.063%  encoding/json.(*Decoder).Decode
         0     0% 12.41%  2518.84MB  2.05%  fmt.Fprintf
         0     0% 12.41%  -541.43MB  0.44%  github.com/WorstOfAny/go-musthave-metrics-tpl/internal/handler.(*metricsController).bulkJSONCtx-fm.(*metricsController).bulkJSONCtx.func1
         0     0% 12.41% -18286.58MB 14.90%  github.com/WorstOfAny/go-musthave-metrics-tpl/internal/handler.(*metricsController).checkHMAC-fm.(*metricsController).checkHMAC.func1
         0     0% 12.41%    86.53MB 0.071%  github.com/WorstOfAny/go-musthave-metrics-tpl/internal/handler.(*responseWriter).Write
         0     0% 12.41%   -76.64MB 0.062%  github.com/WorstOfAny/go-musthave-metrics-tpl/internal/handler.decodeMetrics
         0     0% 12.41% -18290.08MB 14.91%  github.com/WorstOfAny/go-musthave-metrics-tpl/internal/handler.decodeRequest.func1
         0     0% 12.41% -18179.28MB 14.82%  github.com/WorstOfAny/go-musthave-metrics-tpl/internal/handler.encodeResponse.func1
         0     0% 12.41%  -496.93MB  0.41%  github.com/WorstOfAny/go-musthave-metrics-tpl/internal/handler.jsonTypeCheck.func1
         0     0% 12.41%  -497.93MB  0.41%  github.com/WorstOfAny/go-musthave-metrics-tpl/internal/handler.jsonTypeSet.func1
         0     0% 12.41% -18286.58MB 14.90%  github.com/WorstOfAny/go-musthave-metrics-tpl/internal/handler.recoveryPanic.func1
         0     0% 12.41% -17752.99MB 14.47%  github.com/WorstOfAny/go-musthave-metrics-tpl/internal/handler.textPlainTypeCheck.func1
         0     0% 12.41%  2882.93MB  2.35%  github.com/WorstOfAny/go-musthave-metrics-tpl/internal/observers.(*audit).send
         0     0% 12.41%   191.58MB  0.16%  github.com/WorstOfAny/go-musthave-metrics-tpl/internal/repository.(*dbDecorator).All.func1
         0     0% 12.41% -19785.30MB 16.13%  github.com/WorstOfAny/go-musthave-metrics-tpl/internal/repository.(*dbDecorator).All.func4
         0     0% 12.41%  -271.05MB  0.22%  github.com/WorstOfAny/go-musthave-metrics-tpl/internal/repository.(*dbDecorator).BulkSet
         0     0% 12.41%   -87.06MB 0.071%  github.com/WorstOfAny/go-musthave-metrics-tpl/internal/repository.(*dbDecorator).Get
         0     0% 12.41%  -114.01MB 0.093%  github.com/WorstOfAny/go-musthave-metrics-tpl/internal/repository.retry
         0     0% 12.41% -18250.92MB 14.87%  github.com/go-chi/chi/v5.(*ChainHandler).ServeHTTP
         0     0% 12.41% -18224.59MB 14.85%  github.com/go-chi/chi/v5.(*Mux).ServeHTTP
         0     0% 12.41% -18253.92MB 14.88%  github.com/go-chi/chi/v5.(*Mux).routeHTTP
         0     0% 12.41%  -211.53MB  0.17%  github.com/jackc/pgx/v5.(*Conn).Exec
         0     0% 12.41%   141.03MB  0.11%  github.com/jackc/pgx/v5.(*Conn).Query
         0     0% 12.41%  -211.53MB  0.17%  github.com/jackc/pgx/v5.(*Conn).exec
         0     0% 12.41%  1183.56MB  0.96%  github.com/jackc/pgx/v5.(*namedStructRowScanner).ScanRow
         0     0% 12.41%  3473.30MB  2.83%  github.com/jackc/pgx/v5.CollectRows[go.shape.struct { ID string "json:\"id\" db:\"id\""; MType string "json:\"type\" db:\"mtype\""; Delta *int64 "json:\"delta,omitempty\" db:\"delta\""; Value *float64 "json:\"value,omitempty\" db:\"value\""; Hash string "json:\"hash,omitempty\" db:\"hash\"" }] (inline)
         0     0% 12.41%  -167.53MB  0.14%  github.com/jackc/pgx/v5.NamedArgs.RewriteQuery
         0     0% 12.41%   293.01MB  0.24%  github.com/jackc/pgx/v5.joinFieldNames
         0     0% 12.41%   293.01MB  0.24%  github.com/jackc/pgx/v5.lookupNamedStructFields
         0     0% 12.41%  -208.53MB  0.17%  github.com/jackc/pgx/v5/pgxpool.(*Conn).Exec
         0     0% 12.41%   141.03MB  0.11%  github.com/jackc/pgx/v5/pgxpool.(*Conn).Query
         0     0% 12.41%  -210.03MB  0.17%  github.com/jackc/pgx/v5/pgxpool.(*Pool).Exec
         0     0% 12.41%   164.03MB  0.13%  github.com/jackc/pgx/v5/pgxpool.(*Pool).Query
         0     0% 12.41%      280MB  0.23%  github.com/jackc/pgx/v5/pgxpool.(*poolRows).Scan
         0     0% 12.41%    86.53MB 0.071%  net/http.(*response).Write
         0     0% 12.41%    86.02MB  0.07%  net/http.(*response).WriteHeader
         0     0% 12.41%    86.53MB 0.071%  net/http.(*response).write
         0     0% 12.41% -18286.58MB 14.90%  net/http.HandlerFunc.ServeHTTP
         0     0% 12.41%   -76.15MB 0.062%  net/http.newBufioWriterSize
         0     0% 12.41% -18224.59MB 14.85%  net/http.serverHandler.ServeHTTP
         0     0% 12.41%    90.52MB 0.074%  net/textproto.(*Reader).ReadMIMEHeader (inline)
         0     0% 12.41%   302.01MB  0.25%  strings.(*Builder).Grow
         0     0% 12.41%   302.01MB  0.25%  strings.(*Builder).grow
         0     0% 12.41%  -210.83MB  0.17%  sync.(*Pool).Get
         0     0% 12.41%  -143.31MB  0.12%  sync.(*Pool).pin
