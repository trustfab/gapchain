package infrastructure

import (
	"log"
	"time"
)

// RetryTask đại diện cho một tác vụ cross-channel cần thử lại nếu thất bại
type RetryTask struct {
	ID        string
	Action    func() error
	MaxRetry  int
	RetryCh   chan struct{} // channel báo hiệu hoàn thành (cho test)
}

var queue = make(chan *RetryTask, 100)

func init() {
	// Khởi động worker nhẹ ngầm bắt các task retry (Saga basic)
	go startRetryWorker()
}

// EnqueueRetryTask đẩy một task cross-channel (e.g. cập nhật inventory) vào hàng đợi
// Nếu chaincode lỗi (timeout), sẽ thử lại với exponential backoff.
func EnqueueRetryTask(id string, maxRetry int, action func() error) {
	log.Printf("[Saga Queue] Đã đưa task %s vào danh sách chờ xử lý nền", id)
	queue <- &RetryTask{
		ID:       id,
		Action:   action,
		MaxRetry: maxRetry,
	}
}

func startRetryWorker() {
	for task := range queue {
		go processTask(task)
	}
}

func processTask(task *RetryTask) {
	retries := 0
	delay := 1 * time.Second

	for retries < task.MaxRetry {
		err := task.Action()
		if err == nil {
			log.Printf("[Saga Queue] Task %s THÀNH CÔNG (sau %d lần thử)", task.ID, retries)
			if task.RetryCh != nil {
				close(task.RetryCh)
			}
			return
		}

		retries++
		log.Printf("[Saga Queue] Task %s LỖI (lần %d/%d): %v. Thử lại sau %v", task.ID, retries, task.MaxRetry, err, delay)
		time.Sleep(delay)
		delay *= 2 // Exponential backoff
	}

	// Nếu tới đây là failed (Max retries reached), có thể cần gửi cảnh báo cho admin
	log.Printf("[Saga Queue] ❌ CẢNH BÁO: Task %s THẤT BẠI hoàn toàn sau %d lần thử. Cần can thiệp tay!", task.ID, task.MaxRetry)
	if task.RetryCh != nil {
		close(task.RetryCh)
	}
}
