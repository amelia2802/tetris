package tetris

import "time"

var speeds = map[int]time.Duration{
	1: 1000 * time.Millisecond,
	2: 900 * time.Millisecond,
	3: 800 * time.Millisecond,
	4: 700 * time.Millisecond,
	5: 600 * time.Millisecond,
	6: 500 * time.Millisecond,
}

func Speed(level int) time.Duration {
	if speed, ok := speeds[level]; ok {
		return speed
	}
	return speeds[6]
}
