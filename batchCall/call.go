package batchCall

type Operation func(interface{}) (interface{}, error)

type Result struct {
	Id       int
	RetryCnt int
	Params   interface{}
	Result   interface{}
	Errors   []error
}

func BathCall(params []interface{}, op Operation, nWorkers int, maxRetry int) []Result {
	results := make([]Result, len(params))
	resultChan := make(chan Result, len(params))
	jobs := make(chan Result, len(params))

	for i, p := range params {
		jobs <- Result{
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
