// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package utils

import (
	"errors"
	"fmt"
	"github.com/iwind/TeaGo/types"
	"regexp"
)

// RangeTimes 计算时间点
func RangeTimes(timeFrom string, timeTo string, everyMinutes int32) (result []string, err error) {
	if everyMinutes <= 0 {
		return nil, errors.New("invalid 'everyMinutes'")
	}

	var reg = regexp.MustCompile(`^\d{4}$`)
	if !reg.MatchString(timeFrom) {
		return nil, errors.New("invalid timeFrom '" + timeFrom + "'")
	}
	if !reg.MatchString(timeTo) {
		return nil, errors.New("invalid timeTo '" + timeTo + "'")
	}

	if timeFrom > timeTo {
		// swap
		timeFrom, timeTo = timeTo, timeFrom
	}

	var everyMinutesInt = int(everyMinutes)

	var fromHour = types.Int(timeFrom[:2])
	var fromMinute = types.Int(timeFrom[2:])
	var toHour = types.Int(timeTo[:2])
	var toMinute = types.Int(timeTo[2:])

	if fromMinute%everyMinutesInt == 0 {
		result = append(result, timeFrom)
	}

	for {
		fromMinute += everyMinutesInt
		if fromMinute > 59 {
			fromHour += fromMinute / 60
			fromMinute = fromMinute % 60
		}
		if fromHour > toHour || (fromHour == toHour && fromMinute > toMinute) {
			break
		}
		result = append(result, fmt.Sprintf("%02d%02d", fromHour, fromMinute))
	}

	return
}
