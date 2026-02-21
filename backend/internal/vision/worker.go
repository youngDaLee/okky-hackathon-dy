package vision

import "context"

var jobCh = make(chan string, 100)

func StartWorker(ctx context.Context, repo VisionRepository) {
}
