package client

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type User struct {
	Country      string        `json:"country"`
	EmailAddress string        `json:"emailAddress"`
	FirstName    string        `json:"firstName"`
	ID           string        `json:"id"`
	LastName     string        `json:"lastName"`
	Roles        []Role        `json:"roles"`
	TeamIds      []interface{} `json:"teamIds"`
	Username     string        `json:"username"`
}

type NewUser struct {
	Roles    []string `json:"roles"`
	Username string   `json:"username"`
}

type NewReturnUser struct {
	InviterUsername string        `json:"inviterUsername"`
	OrgID           string        `json:"orgId"`
	OrgName         string        `json:"orgName"`
	TeamIds         []interface{} `json:"teamIds"`
	Username        string        `json:"username"`
}
type UpdateUser struct {
	Roles []Role `json:"roles"`
}
type Role struct {
	GroupID  string `json:"groupId,omitempty"`
	RoleName string `json:"roleName"`
	OrgID    string `json:"orgId,omitempty"`
}

var (
	Errors = make(map[int]string)
)

type ErrorStruct struct {
	ErrorCode string `json:"errorCode,omitempty"`
}

func init() {
	Errors[400] = "Bad Request, StatusCode = 400"
	Errors[404] = "User Does Not Exist , StatusCode = 404"
	Errors[409] = "User Already Exist, StatusCode = 409"
	Errors[401] = "Unauthorized Access, StatusCode = 401"
	Errors[429] = "User Has Sent Too Many Request, StatusCode = 429"
}

type Client struct {
	publickey  string
	privateKey string
	orgid      string
	httpClient *http.Client
}

func NewClient(publickey string, privateKey string, orgid string) *Client {
	return &Client{
		publickey:  publickey,
		privateKey: privateKey,
		orgid:      orgid,
		httpClient: &http.Client{},
	}
}

func (c *Client) GetUser(username string) (*User, error) {
	url := "https://cloud.mongodb.com/api/atlas/v1.0/users/byName/" + username
	method := "GET"
	req, err := http.NewRequest(method, url, nil)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	digestParts := digestParts(resp)
	digestParts["uri"] = url
	digestParts["method"] = method
	digestParts["username"] = c.publickey
	digestParts["password"] = c.privateKey
	req, err = http.NewRequest(method, url, nil)
	req.Header.Set("Authorization", getDigestAuthrization(digestParts))
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {

		fmt.Print(err)
		return nil, fmt.Errorf(err.Error())
	}
	statuscode := (int)(resp.StatusCode)
	if statuscode >= 200 && statuscode <= 400 {
		newbody := resp.Body
		userInfo := &User{}
		err = json.NewDecoder(newbody).Decode(userInfo)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		return userInfo, nil
	}
	fmt.Print(statuscode)

	newbody := resp.Body
	errorData := &ErrorStruct{}
	err = json.NewDecoder(newbody).Decode(errorData)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	if errorData.ErrorCode != "" {
		return nil, fmt.Errorf("Error : User Does Not Exist")
	}
	return nil, fmt.Errorf("Error : Unauthorized Access")
}

func (c *Client) AddNewUser(item *NewUser) (*NewReturnUser, error) {
	fmt.Println("New user")
	userjson := NewUser{
		Roles:    item.Roles,
		Username: item.Username,
	}
	reqjson, _ := json.Marshal(userjson)
	payload := strings.NewReader(string(reqjson))
	url := "https://cloud.mongodb.com/api/atlas/v1.0/orgs/" + c.orgid + "/invites?pretty=true"
	method := "POST"
	req, err := http.NewRequest(method, url, nil)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	digestParts := digestParts(resp)
	digestParts["uri"] = url
	digestParts["method"] = method
	digestParts["username"] = c.publickey
	digestParts["password"] = c.privateKey
	req, err = http.NewRequest(method, url, payload)
	req.Header.Set("Authorization", getDigestAuthrization(digestParts))
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	statuscode := (int)(resp.StatusCode)
	if statuscode >= 200 && statuscode <= 400 {
		newbody := resp.Body
		newItemUser := &NewReturnUser{}
		err = json.NewDecoder(newbody).Decode(newItemUser)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		return newItemUser, nil
	}
	return nil, fmt.Errorf("Error : %v ", Errors[statuscode])
}

func (c *Client) UpdateUser(roles []string, userId string) (*User, error) {
	fmt.Println("Update user")
	ois := []Role{}
	for _, role := range roles {
		oi := Role{
			OrgID:    c.orgid,
			RoleName: role,
		}
		ois = append(ois, oi)
	}
	updatevalue := UpdateUser{
		Roles: ois,
	}
	reqjson, _ := json.Marshal(updatevalue)
	payload := strings.NewReader(string(reqjson))
	url := "https://cloud.mongodb.com/api/public/v1.0/users/" + userId
	method := "PATCH"
	req, err := http.NewRequest(method, url, nil)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error client do ,", err)
	}
	defer resp.Body.Close()
	digestParts := digestParts(resp)
	digestParts["uri"] = url
	digestParts["method"] = method
	digestParts["username"] = c.publickey
	digestParts["password"] = c.privateKey
	req, err = http.NewRequest(method, url, payload)
	req.Header.Set("Authorization", getDigestAuthrization(digestParts))
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("eror :", err)
	}
	statuscode := (int)(resp.StatusCode)
	if statuscode >= 200 && statuscode < 400 {
		newbody := resp.Body
		newItemUser := &User{}
		err = json.NewDecoder(newbody).Decode(newItemUser)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		return newItemUser, nil
	}

	return nil, fmt.Errorf("Error : %v ", Errors[statuscode])
}

func (c *Client) DeleteUser(UserId string) error {
	url := "https://cloud.mongodb.com/api/atlas/v1.0/orgs/" + c.orgid + "/users/" + UserId
	method := "DELETE"
	req, err := http.NewRequest(method, url, nil)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	digestParts := digestParts(resp)
	digestParts["uri"] = url
	digestParts["method"] = method
	digestParts["username"] = "bnaouyco"
	digestParts["password"] = "0d9f4ebf-2153-48e6-a579-3eeb5a9758e8"
	req, err = http.NewRequest(method, url, nil)
	req.Header.Set("Authorization", getDigestAuthrization(digestParts))
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		log.Println("[DELETE ERROR]: ", err)
		return err
	}
	statuscode := (int)(resp.StatusCode)
	if statuscode >= 200 && statuscode <= 400 {
		return nil
	}

	log.Println(Errors[statuscode])
	return fmt.Errorf("Error : %v \n", Errors[statuscode])

}

func digestParts(resp *http.Response) map[string]string {
	result := map[string]string{}
	if len(resp.Header["Www-Authenticate"]) > 0 {
		wantedHeaders := []string{"nonce", "realm", "qop"}
		responseHeaders := strings.Split(resp.Header["Www-Authenticate"][0], ",")
		for _, r := range responseHeaders {
			for _, w := range wantedHeaders {
				if strings.Contains(r, w) {
					result[w] = strings.Split(r, `"`)[1]
				}
			}
		}
	}
	return result
}

func getMD5(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func getCnonce() string {
	b := make([]byte, 8)
	io.ReadFull(rand.Reader, b)
	return fmt.Sprintf("%x", b)[:16]
}

func getDigestAuthrization(digestParts map[string]string) string {
	d := digestParts
	ha1 := getMD5(d["username"] + ":" + d["realm"] + ":" + d["password"])
	ha2 := getMD5(d["method"] + ":" + d["uri"])
	nonceCount := 00000001
	cnonce := getCnonce()
	response := getMD5(fmt.Sprintf("%s:%s:%v:%s:%s:%s", ha1, d["nonce"], nonceCount, cnonce, d["qop"], ha2))
	authorization := fmt.Sprintf(`Digest username="%s", realm="%s", nonce="%s", uri="%s", cnonce="%s", nc="%v", qop="%s", response="%s"`,
		d["username"], d["realm"], d["nonce"], d["uri"], cnonce, nonceCount, d["qop"], response)
	return authorization
}

func (c *Client) IsRetry(err error) bool {
	if err != nil {
		if strings.Contains(err.Error(), "429") == true {
			return true
		}
	}
	return false
}
