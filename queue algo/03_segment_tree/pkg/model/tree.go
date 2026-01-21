package model

const BucketInterval = 50

type Node struct {
	Count       int // Number of jobs in this range
	MinSize     int
	MaxSize     int
	Left, Right *Node
	Bucket      *JobBucket // Only for leaves
}

type SegmentTree struct {
	Root    *Node
	MaxSize int
}

func NewSegmentTree(maxSizeMB int) *SegmentTree {
	st := &SegmentTree{MaxSize: maxSizeMB}
	st.Root = buildTree(0, maxSizeMB)
	return st
}

func buildTree(min, max int) *Node {
	node := &Node{MinSize: min, MaxSize: max}
	if max-min <= BucketInterval {
		node.Bucket = &JobBucket{MinSize: min, MaxSize: max}
		return node
	}
	mid := min + (max-min)/2
	node.Left = buildTree(min, mid)
	node.Right = buildTree(mid+1, max)
	return node
}

// AddJob updates the tree in O(log B)
func (st *SegmentTree) AddJob(j *Job) {
	update(st.Root, j, 1)
}

func update(n *Node, j *Job, delta int) {
	n.Count += delta
	if n.Bucket != nil {
		if delta > 0 {
			n.Bucket.Push(j)
		}
		return
	}
	mid := n.MinSize + (n.MaxSize-n.MinSize)/2
	if j.SizeMB <= mid {
		update(n.Left, j, delta)
	} else {
		update(n.Right, j, delta)
	}
}

// FindHeaviest finds the best job <= capacity in O(log B)
func (st *SegmentTree) FindHeaviest(capacity, minSize int) *Job {
	return query(st.Root, capacity, minSize)
}

func query(n *Node, cap, min int) *Job {
	// Optimization: If this branch has no jobs, or the range is completely outside our search, skip it.
	if n.Count == 0 || n.MinSize > cap || n.MaxSize < min {
		return nil
	}

	// If leaf, try to pop a job
	if n.Bucket != nil {
		j := n.Bucket.Pop()
		if j != nil {
			n.Count--
		}
		return j
	}

	// Try Right (Heavier) branch first to find the best fit
	job := query(n.Right, cap, min)
	if job == nil {
		// Fallback to Left (Lighter) branch
		job = query(n.Left, cap, min)
	}

	if job != nil {
		n.Count--
	}
	return job
}

func (st *SegmentTree) TotalJobsInRange(min, max int) int {
	return countRange(st.Root, min, max)
}

func countRange(n *Node, min, max int) int {
	if n == nil || n.MinSize > max || n.MaxSize < min {
		return 0
	}
	if n.MinSize >= min && n.MaxSize <= max {
		return n.Count
	}
	return countRange(n.Left, min, max) + countRange(n.Right, min, max)
}
