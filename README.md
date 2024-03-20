# Introduction

Most modern HTTP API systems implement HTTP callback when they have to execute
some time consuming tasks. This usually is done by calling an API to start the
task. This API registers the request and responds saying the task was
registered. The caller is provided an option to sepcify a callback endpoint
where they will received the status update/details once the task is finished.

![api-service-with-callbacks](https://github.com/josnidhin/go-async-cb-sync-rr/assets/670464/c5225b63-f270-42a0-b9ab-011ab7f35b67)

But sometime you encounter some odd clients which can't handle such systems and
required the API to be Request-Reply without any callbacks. This is a POC which
demos a Redis PubSub based mechasim to convert callback based API to
Request-Reply without changing existing callback based API.
