package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

type CType string

const (
	TypeDailyAt4P  CType = "daily_at_4P"
	TypeWeeklyAt9A CType = "weekly_at_9A"
)

const (
	triggerFrequency = 15
)

func main() {
	lambda.Start(runCron)
}

func runCron(ctx context.Context, event interface{}) error {
	fmt.Printf("received event: %+v\n", event)
	return runDigest()
}

func runDigest() error {
	t := time.Now().UTC()
	for _, v := range timezones {
		fmt.Println("checking timzone:", v)
		triggered := false
		for _, tz := range v {
			loc, err := time.LoadLocation(tz)
			if err != nil {
				fmt.Println("golang's timezone pkg unknown location --> ", tz)
				continue
			}
			tLoc := t.In(loc)
			if tLoc.Hour() == 16 && tLoc.Minute() >= 0 && tLoc.Minute() <= triggerFrequency {
				triggered = true
				sendDailyDigest(tz)
			} else if tLoc.Hour() == 9 && tLoc.Minute() >= 0 && tLoc.Minute() <= triggerFrequency && tLoc.Weekday() == time.Monday {
				triggered = true
				sendWeeklyDigest(tz)
			} else {
				//fmt.Println("no match. not triggered")
			}
		}
		if !triggered {
			fmt.Println("no match. not triggered")
		}
	}
	return nil
}

func sendDailyDigest(zone string) {
	fmt.Println("triggered: ", fmt.Sprintf("daily, %v", zone))
	post(zone, TypeDailyAt4P)
}

func sendWeeklyDigest(zone string) {
	fmt.Println("triggered: ", fmt.Sprintf("weekly, %v", zone))
	post(zone, TypeWeeklyAt9A)
}

func post(zone string, cType CType) {
	postBody, _ := json.Marshal(map[string]string{
		"zone":  zone,
		"type":  string(cType),
		"token": "<TOKEN>",
	})
	reqBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("<YOUR APP ENDPOINT>", "application/json", reqBody)
	//Handle Error
	if err != nil {
		fmt.Printf("cannot post request to the rewind server please check %v\n", err)
		//bad practice: ignoring error here....
		return
	}
	defer resp.Body.Close()
}

var timezones = map[string][]string{
	"GMT": {
		"Africa/Abidjan",
		"Africa/Accra",
		"Africa/Bamako",
		"Africa/Banjul",
		"Africa/Bissau",
		"Africa/Conakry",
		"Africa/Dakar",
		"Africa/Freetown",
		"Africa/Lome",
		"Africa/Monrovia",
		"Africa/Nouakchott",
		"Africa/Ouagadougou",
		"Africa/Sao_Tome",
		"Africa/Timbuktu",
		"America/Danmarkshavn",
		"Antarctica/Troll",
		"Atlantic/Reykjavik",
		"Atlantic/St_Helena",
		"Eire",
		"Etc/GMT",
		"Etc/GMT+0",
		"Etc/GMT-0",
		"Etc/GMT0",
		"Etc/Greenwich",
		"Europe/Belfast",
		"Europe/Dublin",
		"Europe/Guernsey",
		"Europe/Isle_of_Man",
		"Europe/Jersey",
		"Europe/London",
		"GB",
		"GB-Eire",
		"GMT",
		"GMT+0",
		"GMT-0",
		"GMT0",
		"Greenwich",
		"Iceland",
	},
	"GHST": []string{
		"Africa/Accra",
	},
	"EAT": []string{
		"Africa/Addis_Ababa",
		"Africa/Asmara",
		"Africa/Asmera",
		"Africa/Dar_es_Salaam",
		"Africa/Djibouti",
		"Africa/Juba",
		"Africa/Kampala",
		"Africa/Mogadishu",
		"Africa/Nairobi",
		"Indian/Antananarivo",
		"Indian/Comoro",
		"Indian/Mayotte",
	},
	"CET": []string{
		"Africa/Algiers",
		"Africa/Ceuta",
		"Africa/Tunis",
		"Arctic/Longyearbyen",
		"Atlantic/Jan_Mayen",
		"CET",
		"Europe/Amsterdam",
		"Europe/Andorra",
		"Europe/Belgrade",
		"Europe/Berlin",
		"Europe/Bratislava",
		"Europe/Brussels",
		"Europe/Budapest",
		"Europe/Busingen",
		"Europe/Copenhagen",
		"Europe/Gibraltar",
		"Europe/Ljubljana",
		"Europe/Luxembourg",
		"Europe/Madrid",
		"Europe/Malta",
		"Europe/Monaco",
		"Europe/Oslo",
		"Europe/Paris",
		"Europe/Podgorica",
		"Europe/Prague",
		"Europe/Rome",
		"Europe/San_Marino",
		"Europe/Sarajevo",
		"Europe/Skopje",
		"Europe/Stockholm",
		"Europe/Tirane",
		"Europe/Vaduz",
		"Europe/Vatican",
		"Europe/Vienna",
		"Europe/Warsaw",
		"Europe/Zagreb",
		"Europe/Zurich",
		"Poland",
	},
	"WAT": []string{
		"Africa/Bangui",
		"Africa/Brazzaville",
		"Africa/Douala",
		"Africa/Kinshasa",
		"Africa/Lagos",
		"Africa/Libreville",
		"Africa/Luanda",
		"Africa/Malabo",
		"Africa/Ndjamena",
		"Africa/Niamey",
		"Africa/Porto-Novo",
	},
	"CAT": {
		"Africa/Blantyre",
		"Africa/Bujumbura",
		"Africa/Gaborone",
		"Africa/Harare",
		"Africa/Khartoum",
		"Africa/Kigali",
		"Africa/Lubumbashi",
		"Africa/Lusaka",
		"Africa/Maputo",
		"Africa/Windhoek",
	},
	"EET": {
		"Africa/Cairo",
		"Africa/Tripoli",
		"Asia/Amman",
		"Asia/Beirut",
		"Asia/Damascus",
		"Asia/Famagusta",
		"Asia/Gaza",
		"Asia/Hebron",
		"Asia/Nicosia",
		"EET",
		"Egypt",
		"Europe/Athens",
		"Europe/Bucharest",
		"Europe/Chisinau",
		"Europe/Helsinki",
		"Europe/Kaliningrad",
		"Europe/Kiev",
		"Europe/Mariehamn",
		"Europe/Nicosia",
		"Europe/Riga",
		"Europe/Sofia",
		"Europe/Tallinn",
		"Europe/Tiraspol",
		"Europe/Uzhgorod",
		"Europe/Vilnius",
		"Europe/Zaporozhye",
		"Libya",
	},
	"EEST": {
		"Africa/Cairo",
		"Asia/Amman",
		"Asia/Beirut",
		"Asia/Damascus",
		"Asia/Famagusta",
		"Asia/Gaza",
		"Asia/Hebron",
		"Asia/Nicosia",
		"EET",
		"Egypt",
		"Europe/Athens",
		"Europe/Bucharest",
		"Europe/Chisinau",
		"Europe/Helsinki",
		"Europe/Kaliningrad",
		"Europe/Kiev",
		"Europe/Mariehamn",
		"Europe/Nicosia",
		"Europe/Riga",
		"Europe/Sofia",
		"Europe/Tallinn",
		"Europe/Tiraspol",
		"Europe/Uzhgorod",
		"Europe/Vilnius",
		"Europe/Zaporozhye",
	},
	"WET": {
		"Africa/Casablanca",
		"Africa/El_Aaiun",
		"Atlantic/Canary",
		"Atlantic/Faeroe",
		"Atlantic/Faroe",
		"Atlantic/Madeira",
		"Europe/Lisbon",
		"Portugal",
		"WET",
	},
	"WEST": {
		"Africa/Casablanca",
		"Africa/El_Aaiun",
		"Atlantic/Canary",
		"Atlantic/Faeroe",
		"Atlantic/Faroe",
		"Atlantic/Madeira",
		"Europe/Lisbon",
		"Portugal",
		"WET",
	},
	"CEST": {
		"Africa/Ceuta",
		"Africa/Tunis",
		"Antarctica/Troll",
		"Arctic/Longyearbyen",
		"Atlantic/Jan_Mayen",
		"CET",
		"Europe/Amsterdam",
		"Europe/Andorra",
		"Europe/Belgrade",
		"Europe/Berlin",
		"Europe/Bratislava",
		"Europe/Brussels",
		"Europe/Budapest",
		"Europe/Busingen",
		"Europe/Copenhagen",
		"Europe/Gibraltar",
		"Europe/Ljubljana",
		"Europe/Luxembourg",
		"Europe/Madrid",
		"Europe/Malta",
		"Europe/Monaco",
		"Europe/Oslo",
		"Europe/Paris",
		"Europe/Podgorica",
		"Europe/Prague",
		"Europe/Rome",
		"Europe/San_Marino",
		"Europe/Sarajevo",
		"Europe/Skopje",
		"Europe/Stockholm",
		"Europe/Tirane",
		"Europe/Vaduz",
		"Europe/Vatican",
		"Europe/Vienna",
		"Europe/Warsaw",
		"Europe/Zagreb",
		"Europe/Zurich",
		"Poland",
	},
	"SAST": {
		"Africa/Johannesburg",
		"Africa/Maseru",
		"Africa/Mbabane",
	},
	"CAST": {
		"Africa/Khartoum",
	},
	"HAT": {
		"America/Adak",
		"America/Atka",
		"HST",
		"Pacific/Honolulu",
		"Pacific/Johnston",
		"US/Aleutian",
		"US/Hawaii",
	},
	"HAST": {
		"America/Adak",
		"America/Atka",
		"HST",
		"Pacific/Honolulu",
		"Pacific/Johnston",
		"US/Aleutian",
		"US/Hawaii",
	},
	"HADT": {
		"America/Adak",
		"America/Atka",
		"Pacific/Honolulu",
		"Pacific/Johnston",
		"US/Aleutian",
		"US/Hawaii",
	},
	"AKT": {
		"America/Anchorage",
		"America/Juneau",
		"America/Metlakatla",
		"America/Nome",
		"America/Sitka",
		"America/Yakutat",
		"US/Alaska",
	},
	"AKST": {
		"America/Anchorage",
		"America/Juneau",
		"America/Metlakatla",
		"America/Nome",
		"America/Sitka",
		"America/Yakutat",
		"US/Alaska",
	},
	"AKDT": {
		"America/Anchorage",
		"America/Juneau",
		"America/Metlakatla",
		"America/Nome",
		"America/Sitka",
		"America/Yakutat",
		"US/Alaska",
	},
	"AT": {
		"America/Anguilla",
		"America/Antigua",
		"America/Aruba",
		"America/Barbados",
		"America/Blanc-Sablon",
		"America/Curacao",
		"America/Dominica",
		"America/Glace_Bay",
		"America/Goose_Bay",
		"America/Grenada",
		"America/Guadeloupe",
		"America/Halifax",
		"America/Kralendijk",
		"America/Lower_Princes",
		"America/Marigot",
		"America/Martinique",
		"America/Moncton",
		"America/Montserrat",
		"America/Port_of_Spain",
		"America/Puerto_Rico",
		"America/Santo_Domingo",
		"America/St_Barthelemy",
		"America/St_Kitts",
		"America/St_Lucia",
		"America/St_Thomas",
		"America/St_Vincent",
		"America/Thule",
		"America/Tortola",
		"America/Virgin",
		"Atlantic/Bermuda",
		"Canada/Atlantic",
	},
	"AST": {
		"America/Anguilla",
		"America/Antigua",
		"America/Aruba",
		"America/Barbados",
		"America/Blanc-Sablon",
		"America/Curacao",
		"America/Dominica",
		"America/Glace_Bay",
		"America/Goose_Bay",
		"America/Grenada",
		"America/Guadeloupe",
		"America/Halifax",
		"America/Kralendijk",
		"America/Lower_Princes",
		"America/Marigot",
		"America/Martinique",
		"America/Moncton",
		"America/Montserrat",
		"America/Port_of_Spain",
		"America/Puerto_Rico",
		"America/Santo_Domingo",
		"America/St_Barthelemy",
		"America/St_Kitts",
		"America/St_Lucia",
		"America/St_Thomas",
		"America/St_Vincent",
		"America/Thule",
		"America/Tortola",
		"America/Virgin",
		"Asia/Aden",
		"Asia/Baghdad",
		"Asia/Bahrain",
		"Asia/Kuwait",
		"Asia/Qatar",
		"Asia/Riyadh",
		"Atlantic/Bermuda",
		"Canada/Atlantic",
	},
	"BRT": {
		"America/Araguaina",
		"America/Bahia",
		"America/Belem",
		"America/Fortaleza",
		"America/Maceio",
		"America/Recife",
		"America/Santarem",
		"America/Sao_Paulo",
		"Brazil/East",
	},
	"BRST": {
		"America/Araguaina",
		"America/Bahia",
		"America/Belem",
		"America/Fortaleza",
		"America/Maceio",
		"America/Recife",
		"America/Sao_Paulo",
		"Brazil/East",
	},
	"ART": {
		"America/Argentina/Buenos_Aires",
		"America/Argentina/Catamarca",
		"America/Argentina/ComodRivadavia",
		"America/Argentina/Cordoba",
		"America/Argentina/Jujuy",
		"America/Argentina/La_Rioja",
		"America/Argentina/Mendoza",
		"America/Argentina/Rio_Gallegos",
		"America/Argentina/Salta",
		"America/Argentina/San_Juan",
		"America/Argentina/San_Luis",
		"America/Argentina/Tucuman",
		"America/Argentina/Ushuaia",
		"America/Buenos_Aires",
		"America/Catamarca",
		"America/Cordoba",
		"America/Jujuy",
		"America/Mendoza",
		"America/Rosario",
	},
	"ARST": {
		"America/Argentina/Buenos_Aires",
		"America/Argentina/Catamarca",
		"America/Argentina/ComodRivadavia",
		"America/Argentina/Cordoba",
		"America/Argentina/Jujuy",
		"America/Argentina/La_Rioja",
		"America/Argentina/Mendoza",
		"America/Argentina/Rio_Gallegos",
		"America/Argentina/Salta",
		"America/Argentina/San_Juan",
		"America/Argentina/Tucuman",
		"America/Argentina/Ushuaia",
		"America/Buenos_Aires",
		"America/Catamarca",
		"America/Cordoba",
		"America/Jujuy",
		"America/Mendoza",
		"America/Rosario",
	},
	"PYT": {
		"America/Asuncion",
	},
	"PYST": {
		"America/Asuncion",
	},
	"ET": {
		"America/Atikokan",
		"America/Cancun",
		"America/Cayman",
		"America/Coral_Harbour",
		"America/Detroit",
		"America/Fort_Wayne",
		"America/Grand_Turk",
		"America/Indiana/Indianapolis",
		"America/Indiana/Marengo",
		"America/Indiana/Petersburg",
		"America/Indiana/Vevay",
		"America/Indiana/Vincennes",
		"America/Indiana/Winamac",
		"America/Indianapolis",
		"America/Iqaluit",
		"America/Jamaica",
		"America/Kentucky/Louisville",
		"America/Kentucky/Monticello",
		"America/Louisville",
		"America/Montreal",
		"America/Nassau",
		"America/New_York",
		"America/Nipigon",
		"America/Panama",
		"America/Pangnirtung",
		"America/Port-au-Prince",
		"America/Thunder_Bay",
		"America/Toronto",
		"Canada/Eastern",
		"EST",
		"EST5EDT",
		"Jamaica",
		"US/East-Indiana",
		"US/Eastern",
		"US/Michigan",
	},
	"EST": {
		"America/Atikokan",
		"America/Cancun",
		"America/Cayman",
		"America/Coral_Harbour",
		"America/Detroit",
		"America/Fort_Wayne",
		"America/Grand_Turk",
		"America/Indiana/Indianapolis",
		"America/Indiana/Marengo",
		"America/Indiana/Petersburg",
		"America/Indiana/Vevay",
		"America/Indiana/Vincennes",
		"America/Indiana/Winamac",
		"America/Indianapolis",
		"America/Iqaluit",
		"America/Jamaica",
		"America/Kentucky/Louisville",
		"America/Kentucky/Monticello",
		"America/Louisville",
		"America/Montreal",
		"America/Nassau",
		"America/New_York",
		"America/Nipigon",
		"America/Panama",
		"America/Pangnirtung",
		"America/Port-au-Prince",
		"America/Thunder_Bay",
		"America/Toronto",
		"Canada/Eastern",
		"EST",
		"EST5EDT",
		"Jamaica",
		"US/East-Indiana",
		"US/Eastern",
		"US/Michigan",
	},
	"CT": {
		"America/Bahia_Banderas",
		"America/Belize",
		"America/Chicago",
		"America/Costa_Rica",
		"America/El_Salvador",
		"America/Guatemala",
		"America/Havana",
		"America/Indiana/Knox",
		"America/Indiana/Tell_City",
		"America/Knox_IN",
		"America/Managua",
		"America/Matamoros",
		"America/Menominee",
		"America/Merida",
		"America/Mexico_City",
		"America/Monterrey",
		"America/North_Dakota/Beulah",
		"America/North_Dakota/Center",
		"America/North_Dakota/New_Salem",
		"America/Rainy_River",
		"America/Rankin_Inlet",
		"America/Regina",
		"America/Resolute",
		"America/Swift_Current",
		"America/Tegucigalpa",
		"America/Winnipeg",
		"Asia/Chongqing",
		"Asia/Chungking",
		"Asia/Harbin",
		"Asia/Kashgar",
		"Asia/Macao",
		"Asia/Macau",
		"Asia/Shanghai",
		"Asia/Taipei",
		"CST6CDT",
		"Canada/Central",
		"Canada/Saskatchewan",
		"Cuba",
		"Mexico/General",
		"PRC",
		"ROC",
		"US/Central",
		"US/Indiana-Starke",
	},
	"CST": {
		"America/Bahia_Banderas",
		"America/Belize",
		"America/Chicago",
		"America/Costa_Rica",
		"America/El_Salvador",
		"America/Guatemala",
		"America/Havana",
		"America/Indiana/Knox",
		"America/Indiana/Tell_City",
		"America/Knox_IN",
		"America/Managua",
		"America/Matamoros",
		"America/Menominee",
		"America/Merida",
		"America/Mexico_City",
		"America/Monterrey",
		"America/North_Dakota/Beulah",
		"America/North_Dakota/Center",
		"America/North_Dakota/New_Salem",
		"America/Rainy_River",
		"America/Rankin_Inlet",
		"America/Regina",
		"America/Resolute",
		"America/Swift_Current",
		"America/Tegucigalpa",
		"America/Winnipeg",
		"Asia/Chongqing",
		"Asia/Chungking",
		"Asia/Harbin",
		"Asia/Kashgar",
		"Asia/Macao",
		"Asia/Macau",
		"Asia/Shanghai",
		"Asia/Taipei",
		"Asia/Urumqi",
		"CST6CDT",
		"Canada/Central",
		"Canada/Saskatchewan",
		"Cuba",
		"Mexico/General",
		"PRC",
		"ROC",
		"US/Central",
		"US/Indiana-Starke",
	},
	"CDT": {
		"America/Bahia_Banderas",
		"America/Belize",
		"America/Chicago",
		"America/Costa_Rica",
		"America/El_Salvador",
		"America/Guatemala",
		"America/Havana",
		"America/Indiana/Knox",
		"America/Indiana/Tell_City",
		"America/Knox_IN",
		"America/Managua",
		"America/Matamoros",
		"America/Menominee",
		"America/Merida",
		"America/Mexico_City",
		"America/Monterrey",
		"America/North_Dakota/Beulah",
		"America/North_Dakota/Center",
		"America/North_Dakota/New_Salem",
		"America/Rainy_River",
		"America/Rankin_Inlet",
		"America/Resolute",
		"America/Tegucigalpa",
		"America/Winnipeg",
		"Asia/Chongqing",
		"Asia/Chungking",
		"Asia/Macao",
		"Asia/Macau",
		"Asia/Shanghai",
		"Asia/Taipei",
		"CST6CDT",
		"Canada/Central",
		"Cuba",
		"Mexico/General",
		"PRC",
		"ROC",
		"US/Central",
		"US/Indiana-Starke",
	},
	"ADT": {
		"America/Barbados",
		"America/Blanc-Sablon",
		"America/Glace_Bay",
		"America/Goose_Bay",
		"America/Halifax",
		"America/Martinique",
		"America/Moncton",
		"America/Puerto_Rico",
		"America/Thule",
		"Asia/Baghdad",
		"Atlantic/Bermuda",
		"Canada/Atlantic",
	},
	"AMT": {
		"America/Boa_Vista",
		"America/Campo_Grande",
		"America/Cuiaba",
		"America/Manaus",
		"America/Porto_Velho",
		"Asia/Yerevan",
		"Brazil/West",
	},
	"AMST": {
		"America/Boa_Vista",
		"America/Campo_Grande",
		"America/Cuiaba",
		"America/Manaus",
		"America/Porto_Velho",
		"Asia/Yerevan",
		"Brazil/West",
	},
	"COT": {
		"America/Bogota",
	},
	"COST": {
		"America/Bogota",
	},
	"MT": {
		"America/Boise",
		"America/Cambridge_Bay",
		"America/Chihuahua",
		"America/Creston",
		"America/Dawson_Creek",
		"America/Denver",
		"America/Edmonton",
		"America/Fort_Nelson",
		"America/Hermosillo",
		"America/Inuvik",
		"America/Mazatlan",
		"America/Ojinaga",
		"America/Phoenix",
		"America/Shiprock",
		"America/Whitehorse",
		"America/Yellowknife",
		"Canada/Mountain",
		"Canada/Yukon",
		"MST",
		"MST7MDT",
		"Mexico/BajaSur",
		"Navajo",
		"US/Arizona",
		"US/Mountain",
	},
	"MST": {
		"America/Boise",
		"America/Cambridge_Bay",
		"America/Chihuahua",
		"America/Creston",
		"America/Dawson_Creek",
		"America/Denver",
		"America/Edmonton",
		"America/Fort_Nelson",
		"America/Hermosillo",
		"America/Inuvik",
		"America/Mazatlan",
		"America/Ojinaga",
		"America/Phoenix",
		"America/Shiprock",
		"America/Whitehorse",
		"America/Yellowknife",
		"Canada/Mountain",
		"Canada/Yukon",
		"MST",
		"MST7MDT",
		"Mexico/BajaSur",
		"Navajo",
		"US/Arizona",
		"US/Mountain",
	},
	"MDT": {
		"America/Boise",
		"America/Cambridge_Bay",
		"America/Chihuahua",
		"America/Denver",
		"America/Edmonton",
		"America/Hermosillo",
		"America/Inuvik",
		"America/Mazatlan",
		"America/Ojinaga",
		"America/Phoenix",
		"America/Shiprock",
		"America/Yellowknife",
		"Canada/Mountain",
		"MST7MDT",
		"Mexico/BajaSur",
		"Navajo",
		"US/Arizona",
		"US/Mountain",
	},
	"VET": {
		"America/Caracas",
	},
	"GFT": {
		"America/Cayenne",
	},
	"PT": {
		"America/Dawson",
		"America/Ensenada",
		"America/Los_Angeles",
		"America/Santa_Isabel",
		"America/Tijuana",
		"America/Vancouver",
		"Canada/Pacific",
		"Mexico/BajaNorte",
		"PST8PDT",
		"US/Pacific",
	},
	"PST": {
		"America/Dawson",
		"America/Ensenada",
		"America/Los_Angeles",
		"America/Santa_Isabel",
		"America/Tijuana",
		"America/Vancouver",
		"Canada/Pacific",
		"Mexico/BajaNorte",
		"PST8PDT",
		"Pacific/Pitcairn",
		"US/Pacific",
	},
	"EDT": {
		"America/Detroit",
		"America/Fort_Wayne",
		"America/Grand_Turk",
		"America/Guayaquil",
		"America/Indiana/Indianapolis",
		"America/Indiana/Marengo",
		"America/Indiana/Petersburg",
		"America/Indiana/Vevay",
		"America/Indiana/Vincennes",
		"America/Indiana/Winamac",
		"America/Indianapolis",
		"America/Iqaluit",
		"America/Jamaica",
		"America/Kentucky/Louisville",
		"America/Kentucky/Monticello",
		"America/Louisville",
		"America/Montreal",
		"America/Nassau",
		"America/New_York",
		"America/Nipigon",
		"America/Pangnirtung",
		"America/Port-au-Prince",
		"America/Thunder_Bay",
		"America/Toronto",
		"Canada/Eastern",
		"EST5EDT",
		"Jamaica",
		"US/East-Indiana",
		"US/Eastern",
		"US/Michigan",
	},
	"ACT": {
		"America/Eirunepe",
		"America/Porto_Acre",
		"America/Rio_Branco",
		"Australia/Adelaide",
		"Australia/Broken_Hill",
		"Australia/Darwin",
		"Australia/North",
		"Australia/South",
		"Australia/Yancowinna",
		"Brazil/Acre",
	},
	"ACST": {
		"America/Eirunepe",
		"America/Porto_Acre",
		"America/Rio_Branco",
		"Australia/Adelaide",
		"Australia/Broken_Hill",
		"Australia/Darwin",
		"Australia/North",
		"Australia/South",
		"Australia/Yancowinna",
		"Brazil/Acre",
	},
	"PDT": {
		"America/Ensenada",
		"America/Los_Angeles",
		"America/Santa_Isabel",
		"America/Tijuana",
		"America/Vancouver",
		"Canada/Pacific",
		"Mexico/BajaNorte",
		"PST8PDT",
		"US/Pacific",
	},
	"WGT": {
		"America/Godthab",
		"America/Nuuk",
	},
	"WGST": {
		"America/Godthab",
		"America/Nuuk",
	},
	"ECT": {
		"America/Guayaquil",
	},
	"GYT": {
		"America/Guyana",
	},
	"BOT": {
		"America/La_Paz",
	},
	"BST": {
		"America/La_Paz",
		"Europe/Belfast",
		"Europe/Guernsey",
		"Europe/Isle_of_Man",
		"Europe/Jersey",
		"Europe/London",
		"GB",
		"GB-Eire",
		"Pacific/Bougainville",
	},
	"PET": {
		"America/Lima",
	},
	"PEST": {
		"America/Lima",
	},
	"PMST": {
		"America/Miquelon",
	},
	"PMDT": {
		"America/Miquelon",
	},
	"UYT": {
		"America/Montevideo",
	},
	"UYST": {
		"America/Montevideo",
	},
	"FNT": {
		"America/Noronha",
		"Brazil/DeNoronha",
	},
	"FNST": {
		"America/Noronha",
		"Brazil/DeNoronha",
	},
	"SRT": {
		"America/Paramaribo",
	},
	"CLT": {
		"America/Punta_Arenas",
		"America/Santiago",
		"Chile/Continental",
	},
	"CLST": {
		"America/Santiago",
		"Chile/Continental",
	},
	"EHDT": {
		"America/Santo_Domingo",
	},
	"EGT": {
		"America/Scoresbysund",
	},
	"EGST": {
		"America/Scoresbysund",
	},
	"NT": {
		"America/St_Johns",
		"Canada/Newfoundland",
	},
	"NST": {
		"America/St_Johns",
		"Canada/Newfoundland",
	},
	"NDT": {
		"America/St_Johns",
		"Canada/Newfoundland",
	},
	"AWT": {
		"Antarctica/Casey",
		"Australia/Perth",
		"Australia/West",
	},
	"AWST": {
		"Antarctica/Casey",
		"Australia/Perth",
		"Australia/West",
	},
	"DAVT": {
		"Antarctica/Davis",
	},
	"DDUT": {
		"Antarctica/DumontDUrville",
	},
	"MIST": {
		"Antarctica/Macquarie",
	},
	"MAWT": {
		"Antarctica/Mawson",
	},
	"NZT": {
		"Antarctica/McMurdo",
		"Antarctica/South_Pole",
		"NZ",
		"Pacific/Auckland",
	},
	"NZST": {
		"Antarctica/McMurdo",
		"Antarctica/South_Pole",
		"NZ",
		"Pacific/Auckland",
	},
	"NZDT": {
		"Antarctica/McMurdo",
		"Antarctica/South_Pole",
		"NZ",
		"Pacific/Auckland",
	},
	"ROTT": {
		"Antarctica/Palmer",
		"Antarctica/Rothera",
	},
	"SYOT": {
		"Antarctica/Syowa",
	},
	"VOST": {
		"Antarctica/Vostok",
	},
	"ALMT": {
		"Asia/Almaty",
		"Asia/Qostanay",
	},
	"ALMST": {
		"Asia/Almaty",
	},
	"ANAT": {
		"Asia/Anadyr",
	},
	"AQTT": {
		"Asia/Aqtau",
		"Asia/Aqtobe",
		"Asia/Atyrau",
	},
	"AQTST": {
		"Asia/Aqtobe",
	},
	"TMT": {
		"Asia/Ashgabat",
		"Asia/Ashkhabad",
	},
	"AZT": {
		"Asia/Baku",
	},
	"AZST": {
		"Asia/Baku",
	},
	"ICT": {
		"Asia/Bangkok",
		"Asia/Ho_Chi_Minh",
		"Asia/Phnom_Penh",
		"Asia/Saigon",
		"Asia/Vientiane",
	},
	"KRAT": {
		"Asia/Barnaul",
		"Asia/Krasnoyarsk",
		"Asia/Novokuznetsk",
	},
	"KGT": {
		"Asia/Bishkek",
	},
	"BNT": {
		"Asia/Brunei",
	},
	"IST": {
		"Asia/Calcutta",
		"Asia/Colombo",
		"Asia/Jerusalem",
		"Asia/Kolkata",
		"Asia/Tel_Aviv",
		"Eire",
		"Europe/Dublin",
		"Israel",
	},
	"YAKT": {
		"Asia/Chita",
		"Asia/Khandyga",
		"Asia/Yakutsk",
	},
	"YAKST": {
		"Asia/Chita",
		"Asia/Khandyga",
		"Asia/Yakutsk",
	},
	"CHOT": {
		"Asia/Choibalsan",
	},
	"CHOST": {
		"Asia/Choibalsan",
	},
	"BDT": {
		"Asia/Dacca",
		"Asia/Dhaka",
	},
	"BDST": {
		"Asia/Dacca",
		"Asia/Dhaka",
	},
	"TLT": {
		"Asia/Dili",
	},
	"GST": {
		"Asia/Dubai",
		"Asia/Muscat",
		"Atlantic/South_Georgia",
	},
	"TJT": {
		"Asia/Dushanbe",
	},
	"TSD": {
		"Asia/Dushanbe",
	},
	"HKT": {
		"Asia/Hong_Kong",
		"Hongkong",
	},
	"HKST": {
		"Asia/Hong_Kong",
		"Hongkong",
	},
	"HOVT": {
		"Asia/Hovd",
	},
	"HOVST": {
		"Asia/Hovd",
	},
	"IRKT": {
		"Asia/Irkutsk",
	},
	"IRKST": {
		"Asia/Irkutsk",
	},
	"TRT": {
		"Asia/Istanbul",
		"Europe/Istanbul",
		"Turkey",
	},
	"WIB": {
		"Asia/Jakarta",
		"Asia/Pontianak",
	},
	"WIT": {
		"Asia/Jayapura",
	},
	"IDT": {
		"Asia/Jerusalem",
		"Asia/Tel_Aviv",
		"Israel",
	},
	"AFT": {
		"Asia/Kabul",
	},
	"PETT": {
		"Asia/Kamchatka",
	},
	"PKT": {
		"Asia/Karachi",
	},
	"PKST": {
		"Asia/Karachi",
	},
	"NPT": {
		"Asia/Kathmandu",
		"Asia/Katmandu",
	},
	"KRAST": {
		"Asia/Krasnoyarsk",
	},
	"MYT": {
		"Asia/Kuala_Lumpur",
		"Asia/Kuching",
	},
	"MLAST": {
		"Asia/Kuala_Lumpur",
	},
	"BORTST": {
		"Asia/Kuching",
	},
	"MAGT": {
		"Asia/Magadan",
	},
	"MAGST": {
		"Asia/Magadan",
		"Asia/Srednekolymsk",
	},
	"WITA": {
		"Asia/Makassar",
		"Asia/Ujung_Pandang",
	},
	"PHT": {
		"Asia/Manila",
	},
	"PHST": {
		"Asia/Manila",
	},
	"NOVT": {
		"Asia/Novosibirsk",
		"Asia/Tomsk",
	},
	"OMST": {
		"Asia/Omsk",
	},
	"OMSST": {
		"Asia/Omsk",
	},
	"ORAT": {
		"Asia/Oral",
	},
	"KT": {
		"Asia/Pyongyang",
		"Asia/Seoul",
		"ROK",
	},
	"KST": {
		"Asia/Pyongyang",
		"Asia/Seoul",
		"ROK",
	},
	"QYZT": {
		"Asia/Qyzylorda",
	},
	"QYZST": {
		"Asia/Qyzylorda",
	},
	"MMT": {
		"Asia/Rangoon",
		"Asia/Yangon",
	},
	"SAKT": {
		"Asia/Sakhalin",
	},
	"UZT": {
		"Asia/Samarkand",
		"Asia/Tashkent",
	},
	"UZST": {
		"Asia/Samarkand",
		"Asia/Tashkent",
	},
	"KDT": {
		"Asia/Seoul",
		"ROK",
	},
	"SGT": {
		"Asia/Singapore",
		"Singapore",
	},
	"MALST": {
		"Asia/Singapore",
		"Singapore",
	},
	"SRET": {
		"Asia/Srednekolymsk",
	},
	"GET": {
		"Asia/Tbilisi",
	},
	"IRST": {
		"Asia/Tehran",
		"Iran",
	},
	"IRDT": {
		"Asia/Tehran",
		"Iran",
	},
	"BTT": {
		"Asia/Thimbu",
		"Asia/Thimphu",
	},
	"JST": {
		"Asia/Tokyo",
		"Japan",
	},
	"JDT": {
		"Asia/Tokyo",
		"Japan",
	},
	"ULAT": {
		"Asia/Ulaanbaatar",
		"Asia/Ulan_Bator",
	},
	"ULAST": {
		"Asia/Ulaanbaatar",
		"Asia/Ulan_Bator",
	},
	"VLAT": {
		"Asia/Ust-Nera",
		"Asia/Vladivostok",
	},
	"VLAST": {
		"Asia/Ust-Nera",
		"Asia/Vladivostok",
	},
	"YEKT": {
		"Asia/Yekaterinburg",
	},
	"YEKST": {
		"Asia/Yekaterinburg",
	},
	"AZOT": {
		"Atlantic/Azores",
	},
	"AZOST": {
		"Atlantic/Azores",
	},
	"CVT": {
		"Atlantic/Cape_Verde",
	},
	"FKT": {
		"Atlantic/Stanley",
	},
	"AET": {
		"Australia/ACT",
		"Australia/Brisbane",
		"Australia/Canberra",
		"Australia/Currie",
		"Australia/Hobart",
		"Australia/Lindeman",
		"Australia/Melbourne",
		"Australia/NSW",
		"Australia/Queensland",
		"Australia/Sydney",
		"Australia/Tasmania",
		"Australia/Victoria",
	},
	"AEST": {
		"Australia/ACT",
		"Australia/Brisbane",
		"Australia/Canberra",
		"Australia/Currie",
		"Australia/Hobart",
		"Australia/Lindeman",
		"Australia/Melbourne",
		"Australia/NSW",
		"Australia/Queensland",
		"Australia/Sydney",
		"Australia/Tasmania",
		"Australia/Victoria",
	},
	"AEDT": {
		"Australia/ACT",
		"Australia/Brisbane",
		"Australia/Canberra",
		"Australia/Currie",
		"Australia/Hobart",
		"Australia/Lindeman",
		"Australia/Melbourne",
		"Australia/NSW",
		"Australia/Queensland",
		"Australia/Sydney",
		"Australia/Tasmania",
		"Australia/Victoria",
	},
	"ACDT": {
		"Australia/Adelaide",
		"Australia/Broken_Hill",
		"Australia/Darwin",
		"Australia/South",
		"Australia/Yancowinna",
	},
	"ACWT": {
		"Australia/Eucla",
	},
	"ACWST": {
		"Australia/Eucla",
	},
	"ACWDT": {
		"Australia/Eucla",
	},
	"LHT": {
		"Australia/LHI",
		"Australia/Lord_Howe",
	},
	"LHST": {
		"Australia/LHI",
		"Australia/Lord_Howe",
	},
	"LHDT": {
		"Australia/LHI",
		"Australia/Lord_Howe",
	},
	"AWDT": {
		"Australia/Perth",
		"Australia/West",
	},
	"EAST": {
		"Chile/EasterIsland",
		"Pacific/Easter",
	},
	"EASST": {
		"Chile/EasterIsland",
		"Pacific/Easter",
		"Pacific/Galapagos",
	},
	"GMT-1": {
		"Etc/GMT+1",
	},
	"GMT-10": {
		"Etc/GMT+10",
	},
	"GMT-11": {
		"Etc/GMT+11",
	},
	"GMT-12": {
		"Etc/GMT+12",
	},
	"GMT-2": {
		"Etc/GMT+2",
	},
	"GMT-3": {
		"Etc/GMT+3",
	},
	"GMT-4": {
		"Etc/GMT+4",
	},
	"GMT-5": {
		"Etc/GMT+5",
	},
	"GMT-6": {
		"Etc/GMT+6",
	},
	"GMT-7": {
		"Etc/GMT+7",
	},
	"GMT-8": {
		"Etc/GMT+8",
	},
	"GMT-9": {
		"Etc/GMT+9",
	},
	"GMT+1": {
		"Etc/GMT-1",
	},
	"GMT+10": {
		"Etc/GMT-10",
	},
	"GMT+11": {
		"Etc/GMT-11",
	},
	"GMT+12": {
		"Etc/GMT-12",
	},
	"GMT+13": {
		"Etc/GMT-13",
	},
	"GMT+14": {
		"Etc/GMT-14",
	},
	"GMT+2": {
		"Etc/GMT-2",
	},
	"GMT+3": {
		"Etc/GMT-3",
	},
	"GMT+4": {
		"Etc/GMT-4",
	},
	"GMT+5": {
		"Etc/GMT-5",
	},
	"GMT+6": {
		"Etc/GMT-6",
	},
	"GMT+7": {
		"Etc/GMT-7",
	},
	"GMT+8": {
		"Etc/GMT-8",
	},
	"GMT+9": {
		"Etc/GMT-9",
	},
	"UTC": {
		"Etc/UCT",
		"Etc/UTC",
		"Etc/Universal",
		"Etc/Zulu",
		"UCT",
		"UTC",
		"Universal",
		"Zulu",
	},
	"SAMT": {
		"Europe/Astrakhan",
		"Europe/Samara",
		"Europe/Ulyanovsk",
	},
	"MSK": {
		"Europe/Kirov",
		"Europe/Minsk",
		"Europe/Moscow",
		"Europe/Simferopol",
		"W-SU",
	},
	"MSD": {
		"Europe/Kirov",
		"Europe/Moscow",
		"W-SU",
	},
	"GMT+04:00": {
		"Europe/Saratov",
	},
	"VOLT": {
		"Europe/Volgograd",
	},
	"IOT": {
		"Indian/Chagos",
	},
	"CXT": {
		"Indian/Christmas",
	},
	"CCT": {
		"Indian/Cocos",
	},
	"TFT": {
		"Indian/Kerguelen",
	},
	"SCT": {
		"Indian/Mahe",
	},
	"MVT": {
		"Indian/Maldives",
	},
	"MUT": {
		"Indian/Mauritius",
	},
	"MUST": {
		"Indian/Mauritius",
	},
	"RET": {
		"Indian/Reunion",
	},
	"IRT": {
		"Iran",
	},
	"MHT": {
		"Kwajalein",
		"Pacific/Kwajalein",
		"Pacific/Majuro",
	},
	"MET": {
		"MET",
	},
	"MEST": {
		"MET",
	},
	"CHAT": {
		"NZ-CHAT",
		"Pacific/Chatham",
	},
	"CHAST": {
		"NZ-CHAT",
		"Pacific/Chatham",
	},
	"CHADT": {
		"NZ-CHAT",
		"Pacific/Chatham",
	},
	"WST": {
		"Pacific/Apia",
	},
	"WSDT": {
		"Pacific/Apia",
	},
	"CHUT": {
		"Pacific/Chuuk",
		"Pacific/Truk",
		"Pacific/Yap",
	},
	"VUT": {
		"Pacific/Efate",
	},
	"VUST": {
		"Pacific/Efate",
	},
	"PHOT": {
		"Pacific/Enderbury",
	},
	"TKT": {
		"Pacific/Fakaofo",
	},
	"FJT": {
		"Pacific/Fiji",
	},
	"FJST": {
		"Pacific/Fiji",
	},
	"TVT": {
		"Pacific/Funafuti",
	},
	"GALT": {
		"Pacific/Galapagos",
	},
	"GAMT": {
		"Pacific/Gambier",
	},
	"SBT": {
		"Pacific/Guadalcanal",
	},
	"ChST": {
		"Pacific/Guam",
		"Pacific/Saipan",
	},
	"GDT": {
		"Pacific/Guam",
		"Pacific/Saipan",
	},
	"LINT": {
		"Pacific/Kiritimati",
	},
	"KOST": {
		"Pacific/Kosrae",
	},
	"MART": {
		"Pacific/Marquesas",
	},
	"SST": {
		"Pacific/Midway",
		"Pacific/Pago_Pago",
		"Pacific/Samoa",
		"US/Samoa",
	},
	"NRT": {
		"Pacific/Nauru",
	},
	"NUT": {
		"Pacific/Niue",
	},
	"NFT": {
		"Pacific/Norfolk",
	},
	"NFDT": {
		"Pacific/Norfolk",
	},
	"NCT": {
		"Pacific/Noumea",
	},
	"NCST": {
		"Pacific/Noumea",
	},
	"PWT": {
		"Pacific/Palau",
	},
	"PONT": {
		"Pacific/Pohnpei",
		"Pacific/Ponape",
	},
	"PGT": {
		"Pacific/Port_Moresby",
	},
	"CKT": {
		"Pacific/Rarotonga",
	},
	"CKHST": {
		"Pacific/Rarotonga",
	},
	"TAHT": {
		"Pacific/Tahiti",
	},
	"GILT": {
		"Pacific/Tarawa",
	},
	"TOT": {
		"Pacific/Tongatapu",
	},
	"TOST": {
		"Pacific/Tongatapu",
	},
	"WAKT": {
		"Pacific/Wake",
	},
	"WFT": {
		"Pacific/Wallis",
	},
}
