package main

import "sync"

// partialResult is what one worker hands back to the master: its chunk,
// transposed in isolation.
type partialResult struct {
	data [][]int
}

// Transpose is the master: it splits input into row-chunks, hands each
// chunk to a worker goroutine, waits for all of them to finish, then
// reassembles the partial results into the final transposed matrix.
func Transpose(input [][]int, numWorkers int) [][]int {
	if len(input) == 0 {
		return nil
	}
	chunks := splitRows(input, numWorkers)
	results := make([]partialResult, len(chunks))

	var wg sync.WaitGroup
	for i, chunk := range chunks {
		wg.Add(1)
		go func(i int, chunk [][]int) {
			defer wg.Done()
			results[i] = partialResult{data: transposeChunk(chunk)} // the worker
		}(i, chunk)
	}
	wg.Wait() // barrier: block until every worker has reported back

	return mergeColumns(results)
}

// splitRows divides input into up to numWorkers contiguous row-chunks.
func splitRows(input [][]int, numWorkers int) [][][]int {
	total := len(input)
	chunkSize := (total + numWorkers - 1) / numWorkers
	var chunks [][][]int
	for start := 0; start < total; start += chunkSize {
		end := start + chunkSize
		if end > total {
			end = total
		}
		chunks = append(chunks, input[start:end])
	}
	return chunks
}

// transposeChunk is the work a single worker performs: transpose its own
// row-chunk, with no knowledge of the other workers or the full matrix.
func transposeChunk(rows [][]int) [][]int {
	if len(rows) == 0 {
		return nil
	}
	cols := len(rows[0])
	out := make([][]int, cols)
	for c := range out {
		out[c] = make([]int, len(rows))
	}
	for r, row := range rows {
		for c, v := range row {
			out[c][r] = v
		}
	}
	return out
}

// mergeColumns is the master's aggregation step: each worker's transposed
// chunk becomes one column-block of the final matrix, in original row order.
func mergeColumns(results []partialResult) [][]int {
	if len(results) == 0 {
		return nil
	}
	rows := len(results[0].data)
	cols := 0
	for _, r := range results {
		cols += len(r.data[0])
	}
	out := make([][]int, rows)
	for r := range out {
		out[r] = make([]int, cols)
	}
	colOffset := 0
	for _, res := range results {
		width := len(res.data[0])
		for r := 0; r < rows; r++ {
			copy(out[r][colOffset:colOffset+width], res.data[r])
		}
		colOffset += width
	}
	return out
}
