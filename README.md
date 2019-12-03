# (wip) Datacite Dump Tool

As of Fall 2019 the Datacite API is a bit flaky.

* https://github.com/datacite/lupo/issues/237
* https://github.com/datacite/datacite/issues/851
* https://github.com/datacite/datacite/issues/188
* https://github.com/datacite/datacite/issues/709

This tool tries to get a data dump until from the API, until a [full
dump](https://github.com/datacite/datacite/issues/709) [might be
available](https://github.com/datacite/datacite/issues/851#issuecomment-538718411).

```
$ nslookup api.datacite.org

api.datacite.org        canonical name = lb-813668705.eu-west-1.elb.amazonaws.com.
Name:   lb-813668705.eu-west-1.elb.amazonaws.com
Address: 52.208.253.16
Name:   lb-813668705.eu-west-1.elb.amazonaws.com
Address: 54.76.211.202
```

## Build

```
$ make
```

## Usage

```
$ dcdump -h
Usage of dcdump:
  -d string
        directory, where to put harvested files (default ".")
  -debug
        only print intervals then exit
  -e value
        end date for harvest (default 2019-12-03)
  -i string
        [w]eekly, [d]daily, [h]ourly, [e]very minute (default "d")
  -l int
        upper limit for number of requests (default 16777216)
  -p string
        file prefix for harvested files (default "dcdump-")
  -s value
        start date for harvest (default 2018-01-01)
  -w int
        parallel workers (approximate) (default 4)
```

## Examples

The dcdump tool uses [datacite API version
2](https://support.datacite.org/docs/api). We query for intervals and via
cursor to circumvent the Index Deep Paging Problem (limit as of 12/2019 is
[10000 records for a query](https://support.datacite.org/docs/pagination), 400
pages x 25 records per page).

So create a tempdir.

```
$ mkdir dump
```

To just list the intervals (depending on the -i flag), use the `-debug` flag:

```
$ dcdump -d dump -i h -s 2019-10-01 -e 2019-10-02 -debug
2019-10-01 00:00:00 +0000 UTC -- 2019-10-01 00:59:59.999999999 +0000 UTC
2019-10-01 01:00:00 +0000 UTC -- 2019-10-01 01:59:59.999999999 +0000 UTC
2019-10-01 02:00:00 +0000 UTC -- 2019-10-01 02:59:59.999999999 +0000 UTC
2019-10-01 03:00:00 +0000 UTC -- 2019-10-01 03:59:59.999999999 +0000 UTC
2019-10-01 04:00:00 +0000 UTC -- 2019-10-01 04:59:59.999999999 +0000 UTC
2019-10-01 05:00:00 +0000 UTC -- 2019-10-01 05:59:59.999999999 +0000 UTC
2019-10-01 06:00:00 +0000 UTC -- 2019-10-01 06:59:59.999999999 +0000 UTC
2019-10-01 07:00:00 +0000 UTC -- 2019-10-01 07:59:59.999999999 +0000 UTC
2019-10-01 08:00:00 +0000 UTC -- 2019-10-01 08:59:59.999999999 +0000 UTC
2019-10-01 09:00:00 +0000 UTC -- 2019-10-01 09:59:59.999999999 +0000 UTC
2019-10-01 10:00:00 +0000 UTC -- 2019-10-01 10:59:59.999999999 +0000 UTC
2019-10-01 11:00:00 +0000 UTC -- 2019-10-01 11:59:59.999999999 +0000 UTC
2019-10-01 12:00:00 +0000 UTC -- 2019-10-01 12:59:59.999999999 +0000 UTC
2019-10-01 13:00:00 +0000 UTC -- 2019-10-01 13:59:59.999999999 +0000 UTC
2019-10-01 14:00:00 +0000 UTC -- 2019-10-01 14:59:59.999999999 +0000 UTC
2019-10-01 15:00:00 +0000 UTC -- 2019-10-01 15:59:59.999999999 +0000 UTC
2019-10-01 16:00:00 +0000 UTC -- 2019-10-01 16:59:59.999999999 +0000 UTC
2019-10-01 17:00:00 +0000 UTC -- 2019-10-01 17:59:59.999999999 +0000 UTC
2019-10-01 18:00:00 +0000 UTC -- 2019-10-01 18:59:59.999999999 +0000 UTC
2019-10-01 19:00:00 +0000 UTC -- 2019-10-01 19:59:59.999999999 +0000 UTC
2019-10-01 20:00:00 +0000 UTC -- 2019-10-01 20:59:59.999999999 +0000 UTC
2019-10-01 21:00:00 +0000 UTC -- 2019-10-01 21:59:59.999999999 +0000 UTC
2019-10-01 22:00:00 +0000 UTC -- 2019-10-01 22:59:59.999999999 +0000 UTC
2019-10-01 23:00:00 +0000 UTC -- 2019-10-01 23:59:59.999999999 +0000 UTC
2019-10-02 00:00:00 +0000 UTC -- 2019-10-02 00:59:59.999999999 +0000 UTC
INFO[0000] 25 intervals
```

Start and end date are relatively flexible, for example (minute slices for a single day):

```
$ dcdump -s 2019-05-01 -e '2019-05-01 23:59:59' -i e -debug
2019-05-01 00:00:00 +0000 UTC -- 2019-05-01 00:00:59.999999999 +0000 UTC
...
2019-05-01 23:59:00 +0000 UTC -- 2019-05-01 23:59:59.999999999 +0000 UTC
INFO[0000] 1440 intervals
...
```

The time windows are not adjusted dynamically. So if you know, by accident,
that 2019-08-02 is a critical date (meaning there are [millions of
updates](https://gist.github.com/miku/176edd1222fc42ae3b23234bc9d3cd87#file-freq-tsv-L325)),
you could run three harvests:

* daily, up until 2019-08-01
* every minute, between 2019-08-01 and 2019-08-03
* daily, rest

This is not automated, but can be scripted.

```
$ dcdump -e '2019-07-31 23:59:59' -i daily -d tmp -p 'part-01-'
$ dcdump -s 2019-08-01 -e '2019-08-03 23:59:59' -i e -d tmp -p 'part-02-'
$ dcdump -s 2019-08-04 -i daily -d tmp -p 'part-03-'
```


If a specific time window fails repeatedly, you can manually touch the file, e.g.

```
$ touch tmp/part-02-20190801114700-20190801114759.ndj
```

The dcdump tool checks for the existence of the file, before harvesting; this
way it's possible to skip unfetchable slices.

After successful runs, concatenate the data to get a newline delimited single file dump of datacite.

```
$ cat tmp/*ndj > datacite.ndj
```

Again, this is ugly, but should all be obsolete as soon as [a public data
dump](https://github.com/datacite/datacite/issues/709) is available.
