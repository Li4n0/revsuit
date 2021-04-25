package qqwry

// Area returns IpArea according to ip
func Area(ip string) string {
	defer func() {
		_ = recover()
	}()
	if GetQQWry() == nil {
		return ""
	}
	ipData := GetQQWry().SearchByIPv4(ip)
	if ipData.Area == " CZ88.NET" {
		return ipData.Country
	}
	return ipData.Country + " " + ipData.Area
}
