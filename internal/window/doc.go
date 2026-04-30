// Package window provides a sliding-window counter for tracking event
// frequency over a rolling time period.
//
// The window is divided into N equal-sized buckets. Each call to Add records
// events in the current bucket. As time passes, old buckets are rotated out
// so that Count always reflects only events within the configured duration.
//
// Typical usage:
//
//	w := window.New(time.Minute, 60) // 60 one-second buckets
//	w.Add(1)
//	fmt.Println(w.Count())
package window
