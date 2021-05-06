package main

/*
func TestGetEntries(t *testing.T) {
	type result struct {
		ksats []ksat
		err   error
	}
	cases := []struct {
		desc     string
		expected result
	}{
		{
			"gets all cases",
			result{
				ksats: []ksat{
					ksat{id: 1, name: "first", usage: "usage here"},
					ksat{id: 2, name: "second", usage: "more usage here"},
					ksat{id: 3, name: "hasP", usage: "usage"},
				},
				err: nil,
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			ksats, err := getKsats()
			if err != tc.expected.err {
				t.Fatalf("errors don't match: %v, %v", err, tc.expected.err)
			}
			for i, item := range tc.expected.ksats {
				if i > len(ksats)-1 {
					t.Fatalf("ksats returned have less items than expected: items missing at ith '%v': %v", i, ksats[i:])
				}
				if item.id != tc.expected.ksats[i].id {
					t.Fatalf("ids don't match: %v, %v", item.id, tc.expected.ksats[i].id)
				}
				if item.name != tc.expected.ksats[i].name {
					t.Fatalf("name don't match: %v, %v", item.name, tc.expected.ksats[i].name)
				}
				if item.usage != tc.expected.ksats[i].usage {
					t.Fatalf("usages don't match: %v, %v", item.usage, tc.expected.ksats[i].usage)
				}
			}
		})
	}
}
*/
