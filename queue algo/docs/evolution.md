# The Evolution of the Dispatcher (Algo 1 to Algo 7)

This document tracks the journey of building a High-Scale Job Scheduler, explaining **Why** we moved from one version to the next.

## Phase 1: The Basics (Algo 1 & 2)

### 🔴 Algo 1: The Priority Queue (Score-Based)
*   **The Idea**: Treat Jobs like a simple checklist. Sort them by "Score" (Wait Time + Size).
*   **The Implementation**: Used a Heap (Priority Queue). The dispatcher popped the "Best" job and gave it to the "Best" worker.
*   **The Problem (The Need for Algo 2)**:
    *   **Starvation**: Heavy jobs (10GB) waited forever because small jobs (100MB) were always "easier" to fit.
    *   **Inefficiency**: A 16GB worker would take a 100MB job, wasting 15.9GB.

### 🟠 Algo 2: Bucket Reservation
*   **The Fix**: We stopped treating all memory as equal. We created **Buckets** (50MB, 100MB... 4GB).
*   **Key Feature**: **Reservation Mode**. If a Heavy Job waits too long (Avalanche), we **Reserve** the big workers. They sit idle and *refuse* small jobs until the big job arrives.
*   **The Problem**: Searching through linear buckets was slow ($O(N)$). We needed speed.

---

## Phase 2: The Optimization (Algo 3 - 6)

### 🟡 Algo 3: Segment Tree
*   **The Fix**: Replaced the Linear List with a **Segment Tree**. Use binary search to find the "Best Fit" bucket.
*   **Result**: Search speed improved from $O(N)$ to Logarithmic $O(\log N)$.

### 🟢 Algo 4: O(1) Bitmask
*   **The Fix**: Replaced the Tree with **CPU Bitmasks**. We map memory classes to Bits (0/1).
*   **Result**: Finding a worker became **Constant Time ($O(1)$)** (~20 nanoseconds).

### 🔵 Algo 5: Tiered Scheduler
*   **The Fix**: Algo 4 couldn't handle 16GB phones (too many bits). We added **Tiers** (High Res vs Low Res maps) and **Multi-Dimensions** (Region + Battery).

### 🟣 Algo 6: Infinite Density
*   **The Fix**: Algo 5 failed if 1 Million phones had the *exact same* specs (Buffer Overflow). We added **Linked Chunk Queues** to allow infinite capacity per bucket.

---

## Phase 3: The Paradigm Shift (Algo 7)

### ⚪ Algo 7: Yggdrasil (Pull-Based)
*   **The Shift**: We stopped "Pushing" jobs.
*   **The Logic**: The Server tracks **Zero Workers**. It only holds Job Queues. Workers (Phones) **Pull** jobs when ready.
*   **The Result**:
    *   **Scale**: Infinite (Server is stateless).
    *   **Reliability**: No "Zombie" phones causing timeouts.
    *   **Throughput**: ~250,000 Jobs/Second.
