package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// This command will generate a link similar to
// https://www.youngliving.com/vo/#/signup/new-start?sponsorid=24437300&enrollerid=24437300&isocountrycode=AE&culture=en-US&type=member
var (
	sponsorIDMsg  = "The direct uplink of the new member"
	enrollerIDMsg = "The id of the enroller of the new member"
	countryMsg    = "The iso country code where the new member will be registered"
	memberTypeMsg = "The member type"
	langMsg       = "The default language or know as culture in YoungLiving"

	sponsorID  = flag.Int("sponsor-id", 0, sponsorIDMsg)
	enrollerID = flag.Int("enroller-id", 0, enrollerIDMsg)
	country    = flag.String("country", "AE", countryMsg)
	memberType = flag.String("type", "member", memberTypeMsg)
	lang       = flag.String("lang", "en-US", langMsg)
)

const (
	BitlyApiURL      = "https://api-ssl.bitly.com"
	BitlyAccessToken = "59c888f1adc6b428a2520d8c193b671011784f8d"
	BaseURL          = "https://www.youngliving.com"
	URLPath          = "/vo/#/signup/new-start"
)

func main() {
	flag.Parse()
	u, err := url.Parse(BaseURL + URLPath)
	if err != nil {
		panic(err)
	}

	if *sponsorID == 0 {
		//var tmp string
		fmt.Print(sponsorIDMsg + ": ")
		_, err = fmt.Scanf("%d\n", sponsorID)
		if err != nil {
			panic(fmt.Errorf("invalid sponsor-id: %v", err))
		}
	}

	if *enrollerID == 0 {
		fmt.Print(enrollerIDMsg + ": ")
		_, err = fmt.Scanf("%d\n", enrollerID)
		if err != nil {
			panic(fmt.Errorf("invalid enroller-id: %v", err))
		}
	}

	q := u.Query()
	q.Add("sponsorid", strconv.Itoa(*sponsorID))
	q.Add("enrollerid", strconv.Itoa(*enrollerID))
	q.Add("country", *country)
	q.Add("country", *country)
	q.Add("memberType", *memberType)
	q.Add("lang", *lang)
	u.RawQuery = q.Encode()

	longURL := u.String()
	req, err := http.NewRequest(
		http.MethodPost,
		BitlyApiURL+"/v4/shorten",
		strings.NewReader(fmt.Sprintf(`{"long_url": "%s"}`, longURL)),
	)
	if err != nil {
		panic(fmt.Errorf("error sending request to bitly api: %v", err))
	}

	req.Header.Add("Authorization", "Bearer "+BitlyAccessToken)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(fmt.Errorf("error shortening url: %v", err))
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode > 201 {

		panic(fmt.Errorf("error received from bitly api: %s", body))
	}

	sr := ShortenReponse{}
	err = json.Unmarshal(body, &sr)
	if err != nil {
		panic(fmt.Errorf("error parsing response from bitly api: %v", err))
	}

	fmt.Printf("\nFull generated url: %s\n", longURL)

	fmt.Printf("\nShortened url: %s\n", sr.Link)
}

type ShortenReponse struct {
	CreatedAt      string        `json:"created_at"`
	ID             string        `json:"id"`
	Link           string        `json:"link"`
	CustomBitlinks []interface{} `json:"custom_bitlinks"`
	LongURL        string        `json:"long_url"`
	Archived       bool          `json:"archived"`
	Tags           []interface{} `json:"tags"`
	Deeplinks      []interface{} `json:"deeplinks"`
	References     struct {
		Group string `json:"group"`
	} `json:"references"`
}
