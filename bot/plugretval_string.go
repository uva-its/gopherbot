// Code generated by "stringer -type=PlugRetVal"; DO NOT EDIT.

package bot

import "strconv"

const _PlugRetVal_name = "NormalFailMechanismFailConfigurationErrorUntrustedPlugin"

var _PlugRetVal_index = [...]uint8{0, 6, 10, 23, 41, 56}

func (i PlugRetVal) String() string {
	if i < 0 || i >= PlugRetVal(len(_PlugRetVal_index)-1) {
		return "PlugRetVal(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _PlugRetVal_name[_PlugRetVal_index[i]:_PlugRetVal_index[i+1]]
}