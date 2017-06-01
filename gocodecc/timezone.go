package gocodecc

// default time zone
var siteTimezone = "UTC"

// setTimezone set the timezone of the site
func setTimezone(tz string) {
	siteTimezone = tz
}
