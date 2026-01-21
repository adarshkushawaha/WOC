package model

type JobBucket struct {
	MinSize int
	MaxSize int
	Jobs    []*Job
}

func (b *JobBucket) Push(j *Job) {
	b.Jobs = append(b.Jobs, j)
}

func (b *JobBucket) Pop() *Job {
	if len(b.Jobs) == 0 {
		return nil
	}
	j := b.Jobs[0]
	b.Jobs = b.Jobs[1:]
	return j
}

func (b *JobBucket) IsEmpty() bool {
	return len(b.Jobs) == 0
}
