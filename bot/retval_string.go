// Code generated by "stringer -type=RetVal"; DO NOT EDIT.

package bot

import "strconv"

const _RetVal_name = "OkUserNotFoundChannelNotFoundAttributeNotFoundFailedUserDMFailedChannelJoinDatumNotFoundDatumLockExpiredDataFormatErrorBrainFailedInvalidDatumKeyInvalidDblPtrInvalidCfgStructNoConfigFoundRetryPromptReplyNotMatchedUseDefaultValueTimeoutExpiredInterruptedMatcherNotFoundNoUserEmailNoBotEmailMailErrorTaskNotFoundMissingArguments"

var _RetVal_index = [...]uint16{0, 2, 14, 29, 46, 58, 75, 88, 104, 119, 130, 145, 158, 174, 187, 198, 213, 228, 242, 253, 268, 279, 289, 298, 310, 326}

func (i RetVal) String() string {
	if i < 0 || i >= RetVal(len(_RetVal_index)-1) {
		return "RetVal(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _RetVal_name[_RetVal_index[i]:_RetVal_index[i+1]]
}
