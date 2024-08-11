package batchCall

type Operation[T any, U any] func(T) (U, error)

type Result[T any, U any] struct {
	Id       int
	RetryCnt int
	Params   T
	Result   U
	Errors   []error
}

type BatchCall[T any, U any] struct {
	Params []T
	Op     Operation[T, U]
}

func (p *BatchCall[T, U]) Call(nWorkers int, maxRetry int) []Result[T, U] {
	params := p.Params
	op := p.Op

	results := make([]Result[T, U], len(params))
	resultChan := make(chan Result[T, U], len(params))
	jobs := make(chan Result[T, U], len(params))

	for i, p := range params {
		jobs <- Result[T, U]{
			Id:       i,
			RetryCnt: 0,
			Params:   p,
		}
	}

	for i := 0; i < nWorkers; i++ {
		go func() {
			for {
				select {
				case job, ok := <-jobs:
					if !ok {
						return
					}
					r, err := op(job.Params)
					if err != nil {
						job.RetryCnt++
						job.Errors = append(job.Errors, err)
						if job.RetryCnt < maxRetry {
							jobs <- job
						} else {
							resultChan <- job
						}
					} else {
						job.Result = r
						resultChan <- job
					}
				}
			}
		}()
	}

	nFinished := 0
	for {
		select {
		case result := <-resultChan:
			results[result.Id] = result
			nFinished++
		}
		if nFinished == len(params) {
			break
		}
	}
	close(jobs)
	close(resultChan)
	return results
}
