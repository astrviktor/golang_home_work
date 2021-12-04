//go:build !bench
// +build !bench

package hw10programoptimization

import (
	"archive/zip"
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDomainStat(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

	t.Run("slow find 'com'", func(t *testing.T) {
		result, err := GetDomainStatSlow(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("slow find 'gov'", func(t *testing.T) {
		result, err := GetDomainStatSlow(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("slow find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStatSlow(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})
}

func TestCountDomains(t *testing.T) {
	testUsers := users{
		User{1, "Howard Mendoza", "0Oliver", "aliquid_qui_ea@Browsedrive.gov", "6-866-899-36-79", "InAQJvsq", "Blackbird Place 25"},
		User{2, "Jesse Vasquez", "qRichardson", "mLynch@broWsecat.com", "9-373-949-64-00", "SiZLeNSGn", "Fulton Hill 80"},
		User{3, "Clarence Olson", "RachelAdams", "RoseSmith@Browsecat.com", "988-48-97", "71kuz3gA5w", "Monterey Park 39"},
		User{4, "Gregory Reid", "tButler", "5Moore@Teklist.net", "520-04-16", "r639qLNu", "Sunfield Park 20"},
		User{5, "Janice Rose", "KeithHart", "nulla@Linktype.com", "146-91-01", "acSBF5", "Russell Trail 61"},
	}

	t.Run("count 'com'", func(t *testing.T) {
		result, err := countDomains(testUsers, "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("count 'net'", func(t *testing.T) {
		result, err := countDomains(testUsers, "net")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"teklist.net": 1}, result)
	})

	t.Run("count 'unknown'", func(t *testing.T) {
		result, err := countDomains(testUsers, "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})
}

// go test -bench=BenchmarkGetDomainStat -benchmem -benchtime 20s -count 10 | tee benchmark/old
// go test -bench=BenchmarkGetDomainStat -benchmem -benchtime 20s -count 10 | tee benchmark/new
// go test -bench=BenchmarkGetDomainStat -cpuprofile=benchmark/old-cpu.out -memprofile=benchmark/old-mem.out .
// go test -bench=BenchmarkGetDomainStat -benchmem -cpuprofile=benchmark/old-cpu.out -memprofile=benchmark/old-mem.out -x .
// go tool pprof -http=":8090" hw10_program_optimization.test benchmark/old-cpu.out
// go tool pprof hw10_program_optimization.test benchmark/old-mem.out
func BenchmarkGetDomainStat(b *testing.B) {
	r, _ := zip.OpenReader("testdata/users.dat.zip")
	defer r.Close()
	data, _ := r.File[0].Open()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GetDomainStat(data, "biz")
	}
	b.StopTimer()
}
