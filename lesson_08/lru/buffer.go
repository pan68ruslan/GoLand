package lru

import "log/slog"

type CircularBuffer struct {
	data []string
	head int
}

func NewCircularBuffer(n int) *CircularBuffer {
	return &CircularBuffer{
		data: make([]string, 0, n),
	}
}

func (cb *CircularBuffer) Add(key string) string {
	removed := ""
	for i, v := range cb.data {
		if v == key {
			cb.data = append(cb.data[:i], cb.data[i+1:]...)
			if cb.head > i {
				cb.head--
			}
			break
		}
	}
	if len(cb.data) < cap(cb.data) {
		cb.data = append(cb.data, key)
		Logger.Info("Added key", slog.String("key", key))
	} else {
		removed = cb.data[cb.head]
		cb.data[cb.head] = key
		cb.head = (cb.head + 1) % cap(cb.data)
		Logger.Info("Removed and added key",
			slog.String("removed", removed),
			slog.String("key", key))
	}
	return removed
}
