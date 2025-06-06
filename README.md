# Datacite Dump Tool

[![DOI](https://zenodo.org/badge/270958623.svg)](https://zenodo.org/badge/latestdoi/270958623)

----

Ifficial datacite data exports are available on a yearly basis at
[https://datafiles.datacite.org/datafiles](https://datafiles.datacite.org/datafiles).

This tool allows to assemble a snapshot of datacite at any time, in the rare
case you the above data snapshots do not cover your needs.

A recent version of datacite metadata generated with dcdump can be found here:
[https://archive.org/details/datacite-2024-07-31](https://archive.org/details/datacite-2024-07-31).

----

As of Fall 2019 the [datacite API](https://support.datacite.org/docs/api) is
a bit flaky: [#237](https://github.com/datacite/lupo/issues/237),
[#851](https://github.com/datacite/datacite/issues/851),
[#188](https://github.com/datacite/datacite/issues/188),
[#709](https://github.com/datacite/datacite/issues/709)
[#897](https://github.com/datacite/datacite/issues/897),
[#898](https://github.com/datacite/datacite/issues/898),
[#1805](https://github.com/datacite/datacite/issues/1805).

This tool tries to get a data dump from the API, until a [full](https://github.com/datacite/datacite/issues/1805)
[dump](https://github.com/datacite/datacite/issues/709) [might be
available](https://github.com/datacite/datacite/issues/851#issuecomment-538718411).

This data has been ingested into [fatcat](https://scholar.archive.org/fatcat/), via
[fatcat_import.py](https://github.com/internetarchive/fatcat/blob/master/python/fatcat_import.py)
in 01/2020.

Built at the [Internet Archive](https://archive.org).

## Install and Build

You'll need the [go](https://golang.org/cmd/go/) tool installed (i.e. [installed go](https://golang.org/doc/install)).

```
$ git clone https://git.archive.org/webgroup/dcdump.git
$ cd dcdump
$ make
```

Or install with the Go tool:

```
$ go install github.com/miku/dcdump/cmd/dcdump@latest
```

## Usage

The basic idea is to request small enough chunks (intervals) from the API to
eventuall capture most records. Some hand-holding may be required; e.g. request
most data via "daily" or "hourly" slices and if gaps remain (e.g. because the
number of updates in a given time slice exceeds the maximum number of records
sent by the api), use "every minute" slices for the rest.

```
$ dcdump -h
Usage of dcdump:
  -A    do not include affiliation information
  -d string
        directory, where to put harvested files (default ".")
  -debug
        only print intervals then exit
  -e value
        end date for harvest (default 2022-07-04)
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

## Affiliations

Affiliations are requested by default (turn if off with `-A`).

* [Can I see more detailed affiliation information in the REST API?](https://support.datacite.org/docs/can-i-see-more-detailed-affiliation-information-in-the-rest-api)

Example:

```json
{
  "data": [
    {
      "id": "10.3886/e100985v1",
      "type": "dois",
      "attributes": {
        "doi": "10.3886/e100985v1",
        "identifiers": [
          {
            "identifier": "https://doi.org/10.3886/e100985v1",
            "identifierType": "DOI"
          }
        ],
        "creators": [
          {
            "name": "Porter, Joshua J.",
            "nameType": "Personal",
            "givenName": "Joshua J.",
            "familyName": "Porter",
            "affiliation": [
              {
                "name": "George Washington University"
              }
            ],
            "nameIdentifiers": []
          }
        ],
      ...
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

Or, more modern:

```shell
$ fd 'dcdump-.*ndjson' -x cat | jq -rc '.data[]' > datacite.ndjson # may contain dups
```

Again, this is ugly, but should all be obsolete as soon as [a public data
dump](https://github.com/datacite/datacite/issues/709) is available.

Another way:

```shell
$ fd | \
    parallel --block 20M -j 32 cat | \
    parallel --block 30M --pipe 'jq -rc .data[]' | \
    pv -l | \
    zstd -c -T0 > /tmp/dcdump-2024-04-17.json.zst
```

## Duration

A duration data point, about 80h.

```
$ dcdump -version
dcdump 5ae0556 2020-01-21T16:25:10Z

$ dcdump -i e
...
INFO[294683] 1075178 date slices succeeded

real    4911m23.343s
user    930m54.034s
sys     173m7.383s
```

After 80h, the total size amounts to about 78G.

## Archive Items

* [https://archive.org/details/datacite-2024-04-17](https://archive.org/details/datacite-2024-04-17)
* [https://archive.org/details/datacite-2024-01-26](https://archive.org/details/datacite-2024-01-26)
* [https://archive.org/details/datacite_dump_20230713](https://archive.org/details/datacite_dump_20230713)
* [https://archive.org/details/datacite_dump_20221118](https://archive.org/details/datacite_dump_20221118)
* [https://archive.org/details/datacite_dump_20211022](https://archive.org/details/datacite_dump_20211022)
* [https://archive.org/details/datacite_dump_20200824](https://archive.org/details/datacite_dump_20200824)
* [https://archive.org/details/datacite_dump_20191122](https://archive.org/details/datacite_dump_20191122)

### Initial snapshot

A datacite snapshot from 11/2019 is available as part of the [Bulk
Bibliographic Metadata](https://archive.org/details/ia_biblio_metadata)
collection at
[Datacite Dump 20191122](https://archive.org/details/datacite_dump_20191122).

> 18210075 items, 72GB uncompressed.

### Updates

See: [#1805](https://github.com/datacite/datacite/issues/1805),
[#709](https://github.com/datacite/datacite/issues/709) and
[ia_biblio_metadata](https://archive.org/details/ia_biblio_metadata) for
updates.

* [https://archive.org/details/datacite_dump_20211022](https://archive.org/details/datacite_dump_20211022); 25859678 unique (lowercased) DOI

```
$ curl -sL https://archive.org/download/datacite_dump_20211022/datacite_dump_20211022.json.zst | \
    zstdcat -c -T0 | jq -rc '.id'

10.1001/jama.289.8.989
10.1001/jama.293.14.1723-a
10.1001/jamainternmed.2013.9245
10.1001/jamaneurol.2015.4885
10.1002/2014gb004975
10.1002/2014gl061020
10.1002/2014jc009965
10.1002/2014jd022411
10.1002/2015gb005314
10.1002/2015gl065259
...
```

* [https://archive.org/details/datacite_dump_20200824](https://archive.org/details/datacite_dump_20200824); 19606708 [unique DOI](https://archive.org/download/datacite_dump_20200824/datacite_20200824_doi.tsv.xz)

```
$ xz -T0 -cd datacite.ndjson.xz | wc
18210075 2562859030 72664858976

$ xz -T0 -cd datacite.ndjson.xz | sha1sum
6fa3bbb1fe07b42e021be32126617b7924f119fb  -
```

----

JI:KNIEKQ2QKJFEGVCTFUZDIMZQBI
