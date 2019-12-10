# Datacite Dump Tool

As of Fall 2019 the [datacite API](https://support.datacite.org/docs/api) is
a bit flaky: [#237](https://github.com/datacite/lupo/issues/237),
[#851](https://github.com/datacite/datacite/issues/851),
[#188](https://github.com/datacite/datacite/issues/188),
[#709](https://github.com/datacite/datacite/issues/709)
[#897](https://github.com/datacite/datacite/issues/897),
[#898](https://github.com/datacite/datacite/issues/898).

This tool tries to get a data dump from the API, until a [full
dump](https://github.com/datacite/datacite/issues/709) [might be
available](https://github.com/datacite/datacite/issues/851#issuecomment-538718411).

This data will be ingested into [fatcat](https://fatcat.wiki/).

## Install and Build

You'll need the [go](https://golang.org/cmd/go/) tool installed (i.e. [installed go](https://golang.org/doc/install)).

```
$ git clone https://git.archive.org/webgroup/dcdump.git
$ cd dcdump
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
	end date for harvest (default 2019-12-10)
  -i string
	[w]eekly, [d]daily, [h]ourly, [e]very minute (default "d")
  -l int
	upper limit for number of requests (default 16777216)
  -p string
	file prefix for harvested files (default "dcdump-")
  -s value
	start date for harvest (default 2018-01-01)
  -sleep duration
	backoff after HTTP error (default 5m0s)
  -version
	show version
  -w int
	parallel workers (approximate) (default 4)
```

## Examples

The dcdump tool uses [datacite API version
2](https://support.datacite.org/docs/api). We query for intervals and via
cursor to circumvent the Index Deep Paging Problem (limit as of 12/2019 is
[10000 records for a query](https://support.datacite.org/docs/pagination), 400
pages x 25 records per page).

To just list the intervals (depending on the -i flag), use the `-debug` flag:

```
$ dcdump -i h -s 2019-10-01 -e 2019-10-02 -debug
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


So create some temporary dir (to not pollute the current directory with the
harvested files).

```
$ mkdir tmp
```

Start harvesting (minute intervals, into `tmp`, with 2 workers).

```
$ dcdump -i e -d tmp -w 2
```

The time windows are not adjusted dynamically. Worse, it seems that even with
a low profile harvest (two workers, backoffs, retries) and minute
intervals, the harvest still can stall (maybe with a 403 or 500).

If a specific time window fails repeatedly, you can manually touch the file, e.g.

```
$ touch tmp/dcdump-20190801114700-20190801114759.ndjson
```

The dcdump tool checks for the existence of the file, before harvesting; this
way it's possible to skip unfetchable slices.

After successful runs, concatenate the data to get a newline delimited single file dump of datacite.

```
$ cat tmp/*ndjson | sort -u > datacite.ndjson
```

Again, this is ugly, but should all be obsolete as soon as [a public data
dump](https://github.com/datacite/datacite/issues/709) is available.

## Archive Item

A datacite snapshot from 11/2019 is available as part of the [Bulk
Bibliographic Metadata](https://archive.org/details/ia_biblio_metadata)
collection at
[Datacite Dump 20191122](https://archive.org/details/datacite_dump_20191122).

> 18210075 items, 72GB uncompressed.

```
$ xz -T0 -cd datacite.ndjson.xz | wc
18210075 2562859030 72664858976

$ xz -T0 -cd datacite.ndjson.xz | sha1sum
6fa3bbb1fe07b42e021be32126617b7924f119fb  -
```

----

Refs SPECPRJCTS-2430.
