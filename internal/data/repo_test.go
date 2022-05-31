package data

// func TestRepo(t *testing.T) {

// 	r := *NewRepository()

// 	if r == nil {
// 		t.Fatalf("Faield to initialize reopsitory")
// 	}

// 	// test data
// 	d := Data{
// 		Payload: "TEST VALUE",
// 	}

// 	got, err := r.PostData(NetworkType, d)
// 	if err != nil {
// 		t.Fatalf("PostData with %v failed: %v", d, err)
// 	}
// 	if len(got.Metadata.ID) == 0 {
// 		t.Fatalf("PostData has no ID: %v", got)
// 	}
// 	if len(got.Metadata.LastUpdated) == 0 {
// 		t.Fatalf("PostData has no LastUpdated: %v", got)
// 	}
// 	if len(got.Metadata.Type) == 0 {
// 		t.Fatalf("PostData has no Type: %v", got)
// 	}

// 	// t.Logf("Got1: %v", got)

// 	d2 := Data{
// 		Payload: "TEST VALUE - TWO",
// 	}

// 	got2, err2 := r.PostData(NetworkType, d2)
// 	if err2 != nil {
// 		t.Fatalf("PostData #2 with %v failed: %v", d2, err2)
// 	}

// 	i1, e1 := strconv.Atoi(got.Metadata.ID)
// 	i2, e2 := strconv.Atoi(got2.Metadata.ID)

// 	if e1 != nil || e2 != nil {
// 		t.Fatalf("ID conversion failed")
// 	}

// 	if i2 <= i1 {
// 		t.Fatalf("New if was expected: i1=%d  i2=%d", i1, i2)
// 	}

// 	// t.Logf("Got2: %v", got2)

// }
