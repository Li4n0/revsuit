package qqwry

// Area return IpArea according to ip
func Area(ip string) string {
	ipData := GetQQWry().SearchByIPv4(ip)
	if GetQQWry() == nil {
		return ""
	}
	if ipData.Area == " CZ88.NET" {
		return ipData.Country
	} else {
		return ipData.Country + " " + ipData.Area
	}
}
