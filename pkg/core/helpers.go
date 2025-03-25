package core

/**
 * Read messages from a given channel but do NOT block if empty
 */
func ReadCh[T any](ch chan T, cb func(T) error) error {
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return nil
			}
			return cb(msg)
		default:
			return nil
		}
	}
}
