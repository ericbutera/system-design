# Extract Transform Load (ETL) Data Pipeline

This is an example of how to build a generic data pipeline to ingest disparate data sources, transform into a unified format, and load into a target data store.

```mermaid
flowchart LR
  Extract --> |SourceData| Transform --> |UnifiedFormat| Load
```

1. Pull data from source system
2. Transform source format to unified format
3. Store unified format in target system

## Tenets of ETL

- data is immutable
- transforms yield new data
- process is idempotent
- prefer passing data location over actual data
- ETL is implemented a Directed Acyclic Graph (DAG)
- all steps should be able to be replayed (example: fixing a bug in transform and being able to fix the data)

## Overview

Extract:

- should handle any source data format
- shouldn't pull the same data continuously (store or emit events)
- should tolerate new unknown fields
- use resilency patterns to be a good api consumer (backoff, retry, circuit breaker)
- handle resuming from a failure
- support full run (full data extraction) and incremental run (only new data since offset)

Transform:

- clear mapping from source to target
- transforms yield new data (never update existing)

Load:

- in this implementation, load calls the SaaS SDK to save assets into the system
- store high quality data for downstream consumers
- the destination can be anything but usually is either a data warehouse or a real-time DB

## Rate Limits

Rate limiting ETL can be a challenge. Vendor systems can have global rate limits, user rate limits, and per service rate limits. These all need to be accounted for in the ETL system.

A common approach is to have a shared state rate limiter (Redis) that can keep track of rate limits at each level across concurrent runs.

## Deduplication

Another common issue is deduplication. A lot of times it isn't that big of a deal to update an asset multiple times, but there are other times where two different processes might try to update the same asset.

One strategy is to have a monotonic counter for each ETL process. As assets are updated ensure that they are persisted with the counter of the current ETL process. If the counter is less than the current ETL process, then the asset has already been updated and can be skipped. If this happens, be sure to emit an event so this descision can be audited.

## Future Enhancements

- concurrency control (multiple data fetching, uploads)
- run type
  - incremental run since last known offset
  - full run for initial run or backfill reconciliation
