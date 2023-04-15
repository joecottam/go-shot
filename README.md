# go-shot

One shot is enough!

## What is this?

This repo is meant to be an exercise to practice some concurrency concepts in Golang.

## The Problem

Build a worker that consumes messages from a queue. Each message has an `appId` referring to which app it belongs to. Messages need to eventually be sent as notifications to devices. Messages need to be collected in batches. Each batch should only contain messages for the same app. A batch is sent as a notification if it has 10 messages (`MaxBatchSize`) or 10 seconds (`MaxBatchInterval`) have passed since the first message was added to the batch.
- The message queue is in AWS SQS
- Notifications are sent by publishing them on an AWS SNS topic with the `appId`.

## How to use?

- Fork this repo
- Build your worker inside `./worker`
- Run tests using `docker compose up --build` and see whether your worker passes all tests.
