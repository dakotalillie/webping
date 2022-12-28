package main

func HasTransitionedIntoErrorState(prevRecords []PingRecord) bool {
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
