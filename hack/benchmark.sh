#!/usr/bin/env bash

#docker run -i loadimpact/k6 run -e FUNC_NAME=withHashProcessWithoutPool --vus 10 --duration 10s - < bench/script.js
#docker run -i loadimpact/k6 run -e FUNC_NAME=withHashLibWithoutPool --vus 10 --duration 10s - < bench/script.js
docker run -i loadimpact/k6 run -e FUNC_NAME=withHashLibWithPool --vus 10 --duration 10s - < bench/script.js
