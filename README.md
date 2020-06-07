mqli [![CircleCI](https://circleci.com/gh/yfuruyama/mqli.svg?style=svg)](https://circleci.com/gh/yfuruyama/mqli)
===
An interactive client for Google Cloud [Monitoring Query Language](https://cloud.google.com/monitoring/mql).

## Install

```
$ go get -u github.com/yfuruyama/mqli
```

## Usage

```
$ mqli -p ${PROJECT_ID}
```

## Example

```
$ mqli -p ${PROJECT_ID}
mql> fetch gce_instance::compute.googleapis.com/instance/cpu/utilization | within 5m
+------------+-------------------+---------------------+-----------------+----------------------+-----------------------+
| project_id | zone              | instance_id         | instance_name   | time                 | utilization           |
+------------+-------------------+---------------------+-----------------+----------------------+-----------------------+
| example    | asia-northeast1-a | 8558245152915036717 | instance-01     | 2020-05-17T10:55:00Z | 0.008204130810478697  |
| example    | asia-northeast1-a | 8558245152915036717 | instance-01     | 2020-05-17T10:54:00Z | 0.00755872624336007   |
| example    | us-central1-f     | 6250447186313570186 | gke-instance-01 | 2020-05-17T10:55:00Z | 0.05655076492888232   |
| example    | us-central1-f     | 6250447186313570186 | gke-instance-01 | 2020-05-17T10:54:00Z | 0.05569375121267512   |
| example    | us-central1-f     | 1734065169329771604 | instance-02     | 2020-05-17T10:55:00Z | 0.014347280275615049  |
| example    | us-central1-f     | 1734065169329771604 | instance-02     | 2020-05-17T10:54:00Z | 0.012462649474400678  |
+------------+-------------------+---------------------+-----------------+----------------------+-----------------------+
6 points in result
```

You can input multi-line query with a backslash(`\`) in the end of the input.

```
mql> fetch gce_instance::compute.googleapis.com/instance/cpu/utilization \
  -> | within 5m \
  -> | filter zone = 'asia-northeast1-a'
+------------+-------------------+---------------------+-----------------+----------------------+-----------------------+
| project_id | zone              | instance_id         | instance_name   | time                 | utilization           |
+------------+-------------------+---------------------+-----------------+----------------------+-----------------------+
| example    | asia-northeast1-a | 8558245152915036717 | instance-01     | 2020-05-17T10:55:00Z | 0.008204130810478697  |
| example    | asia-northeast1-a | 8558245152915036717 | instance-01     | 2020-05-17T10:54:00Z | 0.00755872624336007   |
+------------+-------------------+---------------------+-----------------+----------------------+-----------------------+
2 points in result
```

## Disclaimer
This is not an officially supported Google Cloud product.
