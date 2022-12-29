package main

const ErrorStateThreshold = 3

func HasTransitionedIntoErrorState(prevRecords []PingRecord) bool {
	if len(prevRecords) < ErrorStateThreshold-1 {
		// We haven't made enough requests yet
		return false
	}

	if len(prevRecords) == ErrorStateThreshold-1 {
		// We've made enough requests to meet the error threshold (including the request that hasn't yet been added to
		// the database), but there was never a passing result.
		for _, record := range prevRecords {
			if record.Result != "FAIL" {
				return false
			}
		}
		return true
	}

	for i := 0; i < len(prevRecords)-1; i++ {
		if prevRecords[i].Result != "FAIL" {
			return false
		}
	}

	if prevRecords[len(prevRecords)-1].Result != "PASS" {
		return false
	}

	return true
}
