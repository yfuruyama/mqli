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

Multi-line query

```
mql> fetch ... \
  -> | within 1m \
  -> | top 3
```

## Disclaimer
This is not an officially supported Google Cloud product.
