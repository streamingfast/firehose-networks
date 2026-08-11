# Change log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## v0.2.3

### Added

* Added StreamingFast Firehose and Substreams endpoints (`hoodi.eth.streamingfast.io:443`) to the `hoodi` network, which the upstream registry lists with Pinax endpoints only.

* Added service overrides, a local mechanism to add Firehose and Substreams endpoints to networks that already exist in the upstream registry. Overridden endpoints take precedence over the registry ones and duplicates are dropped.

## v0.2.2

* Updated fallback registry to latest version 0.7.34.

* Added `GetRegistry` function to access the full network registry without filtering.

* Added `FindAll` method to `NetworkRegistry` for finding multiple networks matching a key (also available as global function).

* Added `Search` method to `NetworkRegistry` for finding networks using regular expression patterns (also available as global function).

* Improved performance: `GetSubstreamsRegistry` and `GetFirehoseRegistry` now return cached filtered registries instead of filtering on each call.

## v0.2.1

* Added internal dummy blockchain override.

* Updated fallback registry to latest version 0.7.16.

## v0.2.0

* Updated to fallback registry to latest version 0.7.6.

* Updated to Network Registry version 0.7.6 which comes with breaking change, read about them at https://github.com/graphprotocol/networks-registry/releases/tag/v0.7.0.

## v0.1.0

* Initial version extracted from [substreams](https://github.com/streamingfast/substreams) commit [432e588](https://github.com/streamingfast/substreams/commit/432e58897ceb873d164703a87c0d46f859191669).