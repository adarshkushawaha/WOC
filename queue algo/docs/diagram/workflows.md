# Algorithm Workflows

## Algo 1: Priority Queue (Score Based)
Push Model. Uses `container/heap` with a scoring formula to prioritize jobs.

```mermaid
graph TD
    A[New Job] --> B{Calculate Score}
    B -->|Score = WaitTime + Size| C[Global Priority Queue]
    D[Worker Available] --> E[Check Queue]
    C -->|Highest Score| E
    E --> F[Match & Execute]
```

## Algo 2: Bucket Reservation
Push Model. Uses fixed buckets (e.g., 0-50MB) to organize jobs.

```mermaid
graph TD
    A[New Job] --> B{Check Size}
    B --> C[Bucket 0-50MB]
    B --> D[Bucket 50-100MB]
    B --> E[Bucket 4GB+]
    F[Worker Available] --> G{Linear Search}
    C --> G
    D --> G
    E --> G
    G -->|First Fit| H[Match & Execute]
```

## Algo 3: Segment Tree (Logarithmic)
Push Model. Uses a Segment Tree to find the best fit in $O(\log N)$.

```mermaid
graph TD
    A[New Job] --> B[Tree Insert]
    C[Worker Available] --> D[Tree Query Root]
    D -->|Navigate Left/Right| E[Find Max Capacity Node]
    E --> F[Match & Bin Pack]
```

## Algo 4: O(1) Bitmask
Push Model. Uses CPU bitwise operations for instant matching.

```mermaid
graph TD
    A[New Job] --> B[Map to Bit Index]
    B --> C[Set Bit = 1]
    D[Worker Available] --> E[TrailingZeros Search]
    C -->|Lowest Set Bit| E
    E --> F[Match 1]
```

## Algo 5: Tiered Bitmask
Push Model. Multi-dimensional filtering with tiered memory ranges.

```mermaid
graph TD
    A[New Job] --> B{Get Tier}
    B --> C[Tier 0: 0-4GB]
    B --> D[Tier 1: 4-16GB]
    E[Worker Available] --> F{Get Masks}
    F -->|Region & Battery| G[Calculate Intersection]
    C & D --> G
    G --> H[Match]
```

## Algo 6: Linked Queue (Infinite Density)
Push Model. Linked-List of Arrays to handle millions of identical items.

```mermaid
graph TD
    A[New Job] --> B{Bucket Full?}
    B -->|No| C[Add to Current Chunk]
    B -->|Yes| D[Link New Chunk]
    E[Worker Available] --> F[Pop Head Chunk]
    F --> G[Match]
```

## Algo 7: Yggdrasil (Pull Based)
Pull Model. Inversion of control where workers request work.

```mermaid
sequenceDiagram
    participant W as Worker (Phone)
    participant S as Server (Dispatcher)
    participant Q as Sharded Queue
    
    loop Every N seconds
        W->>S: PullJob(RAM, Region, Battery)
        S->>Q: Check Heavy Tier
        alt Has Heavy Job
            Q-->>S: Return Big Job
        else Empty
            S->>Q: Check Lower Tier (Backfill)
            Q-->>S: Return Small Job
        end
        S-->>W: Job Payload
        W->>W: Execute Lambda
    end
```
