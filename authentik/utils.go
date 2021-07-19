package authentik

func stringToPointer(in string) *string {
	return &in
}

func intToPointer(in int) *int32 {
	i := int32(in)
	return &i
}

func int32ToPointer(in int32) *int32 {
	return &in
}

func boolToPointer(in bool) *bool {
	return &in
}
