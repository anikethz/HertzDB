package utils

type ConstantInteger [20]uint8

func (ci ConstantInteger) CtoI() int64 {
	var res int64 = 0
	var t int64 = 1

	for _, v := range ci {
		if v == 0 {
			t *= 10
			continue
		}
		res += int64(v) * t
		t *= 10
	}

	return res
}

func ItoC(input int64) ConstantInteger {
	result := ConstantInteger{}
	for i := 0; input > 0; {
		result[i] = uint8(input % 10)
		input /= 10
		i++
	}

	return result
}
