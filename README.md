mqli
===
An interactive client for Google Cloud [Monitoring Query Language](https://cloud.google.com/monitoring/mql).

## Usage

```
$ mqli -p ${PROJECT_ID}
```

## Install

```
$ go get -u github.com/yfuruyama/mqli
```

## Example

```
$ mqli -p ${PROJECT_ID}
mql> 
```

## Disclaimer
This is not an officially supported Google Cloud product.

## TODO
- both start_date and end_date?
- value key name
- "resource.xxx and metric.xxx" remove