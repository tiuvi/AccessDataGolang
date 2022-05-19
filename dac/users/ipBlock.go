package http


import(
	"net/http"
	"strconv"
	"encoding/json"
	"io"
	."dac"
)


var IpApp    *SpaceRamSync
var IpBlock  *SpaceRamSync


func InitBlockIp() {

	IpApp = NewSfPermBytes(nil, map[string]int64{"IpApp": 16}, "users", "IpApp").
	InitSync("IpApp")
	
	IpBlock = NewSfPermBytes(nil, map[string]int64{"IpBlock": 16}, "users", "IpBlock").
	InitSync("IpBlock")

}

func NewVisitBlock( ip string){
	
	IpBlock.NewLineString(ip)

}

func IsVisitBlock(ip string)(found bool) {

	if block := IpBlock.GetLineString(ip); block != nil {

		return true
	}

	return false
}



func NewUserBlock( ip string){
	
	IpApp.NewLineString(ip)

}

func IsUserBlock(ip string)(found bool) {

	if block := IpApp.GetLineString(ip); block != nil {

		return true
	}

	return false
}


func GetCountry(ip string)(country string, region string){


	respIpWho, err := http.Get("http://ipwho.is/" + ip + "?fields=country,region")
	if err != nil &&
		IpApp.NRESM(err != nil, err.Error(), "ip", "block") ||
		IpApp.NRESM(respIpWho.StatusCode > 299, "Response failed with status code:" + strconv.Itoa(respIpWho.StatusCode), "ip", "block") {
	}

	body, err := io.ReadAll(respIpWho.Body)
	if err != nil &&
		IpApp.NRESM(err != nil, err.Error(), "ip", "block") {
	}

	//Decodificacion json
	type IpGeoLocation struct {
		Country string `json:"country"`
		Region  string `json:"region"`
	}

	ipGeoLocation := IpGeoLocation{}
	err = json.Unmarshal(body, &ipGeoLocation)
	if err != nil &&
		IpApp.NRESM(err != nil, err.Error(), "ip", "block") {
	}

	err = respIpWho.Body.Close()
	if err != nil &&
		IpApp.NRESM(err != nil, err.Error(), "ip", "block") {
	}

	return ipGeoLocation.Country , ipGeoLocation.Region
}