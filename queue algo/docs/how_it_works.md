# How It Works: The Complete Series (Algo 1 to Algo 7)

This guide explains the inner workings of each scheduler version in simple terms.

## 🔴 Algo 1: The Priority Queue (Score-Based)
*   **Concept**: A simple To-Do list sorted by "Score".
*   **How it works**:
    1.  Every time a job arrives, we calculate `Score = WaitTime * Weight + JobSize`.
    2.  We put the job in a **Heap** (Priority Queue). The highest score floats to the top.
    3.  When a worker is free, we pop the top job.
*   **Why we moved on**: Efficient for small scales, but heavy jobs starved (low score compared to fast small jobs).

## 🟠 Algo 2: Bucket Reservation
*   **Concept**: Separate lists for separate sizes.
*   **How it works**:
    1.  We have buckets: `[0-50MB], [50-100MB] ... [4GB+]`.
    2.  **Reservation**: If the `4GB+` bucket has too many jobs (Avalanche), we signal big workers to **Stop** taking small jobs. They sit idle until a big job arrives.
*   **Why we moved on**: Searching linear buckets ($O(N)$) is slow when you have 1000 buckets.

## 🟡 Algo 3: Segment Tree
*   **Concept**: Organize buckets into a Binary Tree for fast searching.
*   **How it works**:
    1.  Instead of checking Bucket 1, then 2, then 3... we ask the Tree Root: "Where is the largest job?"
    2.  The Tree guides us Left or Right to find the best job in **Logarithmic Time** ($O(\log N)$).
*   **Why we moved on**: Fast matching, but managing the tree structure (pointers) is memory-heavy and slow for millions of items.

## 🟢 Algo 4: O(1) Bitmask
*   **Concept**: Use CPU bits to represent empty/full buckets.
*   **How it works**:
    1.  We map memory chunks to bits in a 64-bit integer. `1` = Has Jobs, `0` = Empty.
    2.  To find a job, we use a CPU instruction (`TrailingZeros`) to find the first `1` bit instantly ($O(1)$).
*   **Why we moved on**: 64 bits only cover 3200MB. We needed to support 16GB phones.

## 🔵 Algo 5: Tiered Bitmask (16GB)
*   **Concept**: Multi-Level Maps (like Google Maps Zoom Levels).
*   **How it works**:
    1.  **Tier 0**: Low Zoom (0-4GB) with high detail.
    2.  **Tier 1**: High Zoom (4-16GB) with lower detail.
    3.  **Multi-Dimension**: We can filter by `Region` AND `Battery` instantly using bitwise AND.
*   **Why we moved on**: If 1 Million identical phones joined, the fixed-size buffers overflowed (Density Limit).

## 🟣 Algo 6: Infinite Density (Linked Queues)
*   **Concept**: A Queue that grows forever.
*   **How it works**:
    1.  Instead of a fixed array, we use a **Chain of Arrays**.
    2.  If Bucket A fills up, we instantly attach partial Bucket B.
    3.  **Result**: You can have 100 or 10 Million phones in the same bucket without crashing.
*   **Why we moved on**: The Server still has to track every single phone. If 1M phones disconnect, the server panics dealing with timeouts.

## ⚪ Algo 7: Yggdrasil (Pull-Based)
*   **Concept**: Inversion of Control. The Server is a Bulletin Board.
*   **How it works**:
    1.  **Server**: Tracks 0 Phones. Only tracks Jobs.
    2.  **Phone**: Connects and asks "Do you have work for a [US, 6GB] phone?"
    3.  **Smart Matching**: The Server forces the 6GB phone to check the **Heavy Job Queue** first. If empty, it checks lighter queues.
*   **Final Result**: Massive Stability. Dead phones are ignored. Live phones work at their own pace. Throughput: **250k jobs/sec**.
