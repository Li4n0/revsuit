package qqwry

// Area return IpArea according to ip
func Area(ip string) string {
	if GetQQWry() == nil {
		return ""
	}
	ipData := GetQQWry().SearchByIPv4(ip)
	if ipData.Area == " CZ88.NET" {
		return ipData.Country
	}
	return ipData.Country + " " + ipData.Area
}
