package crypto

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/google/uuid"
)

const publicKey = "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUFzYUtYT1grMm9GeDlWTjVDNVNvQgpMY29BU0Y2U2lSN3dFQjI1VUNhM3Nod0JZMHFreWc4RGphcVF4REhlT2ZrcGhkUnllTWQrMzVMMURSQ0Z5Y0NpCnV6ZURINXJHQVBPRFJDNGgvUndlbmFPSFhtSnZ4Nzd1UzRVeFE3bS84SVRQMldWc25HMlRpTURLbFc3NDhQSHEKR1ZzRk5OdjZ5TkdxZ1ZoZkw5UDJyWjJUUE05MWxJMEpZYmxSNGlidXZTNlowUzRrZUtTSE1oNU45cFl1dEk2bAp4NWlVNXIzcG40OFhNVW1YSlNIZjYwd1FRVm9lVmJ0dk1yWmlwelJienNCNlg4ZGk1WUs0WWhKQisxZmV3cGgxClZLRkRZNmtibmR6RG1qNGhpdExEeE12bzlDNGVCUFZxUEFhcEFJZFZKS1czb0NwbHVjWGcveSttS0t5RVQvRmgKeFFJREFRQUIKLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tCg=="

func TestName(t *testing.T) {
	generate_jazoest(uuid.MustParse("f2423d0f-35d5-4147-9390-03a2c0137f72"))

}

func generate_jazoest(phoneId uuid.UUID) string {
	var sum int32
	for _, s := range phoneId.String() {
		sum += s
	}

	return strconv.FormatInt(int64(sum), 10)
}

func BenchmarkName1(b *testing.B) {
	var d = int32(24324)
	for i := 0; i < b.N; i++ {
		strconv.FormatInt(int64(d), 10)
	}
}

func BenchmarkName2(b *testing.B) {
	var d = int32(24324)
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("%d", d)
	}
}
