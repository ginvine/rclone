package api

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

const (
	// default time format for almost all request and responses
	timeFormat = "2006-01-02-T15:04:05Z0700"
	// the API server seems to use a different format
	apiTimeFormat = "2006-01-02T15:04:05Z07:00"
)

// Time represents time values in the Jottacloud API. It uses a custom RFC3339 like format.
type Time time.Time

// UnmarshalXML turns XML into a Time
func (t *Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}
	if v == "" {
		*t = Time(time.Time{})
		return nil
	}
	newTime, err := time.Parse(timeFormat, v)
	if err == nil {
		*t = Time(newTime)
	}
	return err
}

// MarshalXML turns a Time into XML
func (t *Time) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(t.String(), start)
}

// Return Time string in Jottacloud format
func (t Time) String() string { return time.Time(t).Format(timeFormat) }

// APIString returns Time string in Jottacloud API format
func (t Time) APIString() string { return time.Time(t).Format(apiTimeFormat) }

// Flag is a hacky type for checking if an attribute is present
type Flag bool

// UnmarshalXMLAttr sets Flag to true if the attribute is present
func (f *Flag) UnmarshalXMLAttr(attr xml.Attr) error {
	*f = true
	return nil
}

// MarshalXMLAttr : Do not use
func (f *Flag) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	attr := xml.Attr{
		Name:  name,
		Value: "false",
	}
	return attr, errors.New("unimplemented")
}

// TokenJSON is the struct representing the HTTP response from OAuth2
// providers returning a token in JSON form.
type TokenJSON struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int32  `json:"expires_in"` // at least PayPal returns string, while most return number
}

/*
GET http://www.jottacloud.com/JFS/<account>

<user time="2018-07-18-T21:39:10Z" host="dn-132">
	<username>12qh1wsht8cssxdtwl15rqh9</username>
	<account-type>free</account-type>
	<locked>false</locked>
	<capacity>5368709120</capacity>
	<max-devices>-1</max-devices>
	<max-mobile-devices>-1</max-mobile-devices>
	<usage>0</usage>
	<read-locked>false</read-locked>
	<write-locked>false</write-locked>
	<quota-write-locked>false</quota-write-locked>
	<enable-sync>true</enable-sync>
	<enable-foldershare>true</enable-foldershare>
	<devices>
		<device>
			<name xml:space="preserve">Jotta</name>
			<display_name xml:space="preserve">Jotta</display_name>
			<type>JOTTA</type>
			<sid>5c458d01-9eaf-4f23-8d3c-2486fd9704d8</sid>
			<size>0</size>
			<modified>2018-07-15-T22:04:59Z</modified>
		</device>
	</devices>
</user>
*/

// DriveInfo represents a Jottacloud account
type DriveInfo struct {
	Username          string        `xml:"username"`
	AccountType       string        `xml:"account-type"`
	Locked            bool          `xml:"locked"`
	Capacity          int64         `xml:"capacity"`
	MaxDevices        int           `xml:"max-devices"`
	MaxMobileDevices  int           `xml:"max-mobile-devices"`
	Usage             int64         `xml:"usage"`
	ReadLocked        bool          `xml:"read-locked"`
	WriteLocked       bool          `xml:"write-locked"`
	QuotaWriteLocked  bool          `xml:"quota-write-locked"`
	EnableSync        bool          `xml:"enable-sync"`
	EnableFolderShare bool          `xml:"enable-foldershare"`
	Devices           []JottaDevice `xml:"devices>device"`
}

/*
GET http://www.jottacloud.com/JFS/<account>/<device>

<device time="2018-07-23-T20:21:50Z" host="dn-158">
	<name xml:space="preserve">Jotta</name>
	<display_name xml:space="preserve">Jotta</display_name>
	<type>JOTTA</type>
	<sid>5c458d01-9eaf-4f23-8d3c-2486fd9704d8</sid>
	<size>0</size>
	<modified>2018-07-15-T22:04:59Z</modified>
	<user>12qh1wsht8cssxdtwl15rqh9</user>
	<mountPoints>
		<mountPoint>
			<name xml:space="preserve">Archive</name>
			<size>0</size>
		<modified>2018-07-15-T22:04:59Z</modified>
		</mountPoint>
		<mountPoint>
			<name xml:space="preserve">Shared</name>
			<size>0</size>
			<modified></modified>
		</mountPoint>
		<mountPoint>
			<name xml:space="preserve">Sync</name>
			<size>0</size>
			<modified></modified>
		</mountPoint>
	</mountPoints>
	<metadata first="" max="" total="3" num_mountpoints="3"/>
</device>
*/

// JottaDevice represents a Jottacloud Device
type JottaDevice struct {
	Name        string            `xml:"name"`
	DisplayName string            `xml:"display_name"`
	Type        string            `xml:"type"`
	Sid         string            `xml:"sid"`
	Size        int64             `xml:"size"`
	User        string            `xml:"user"`
	MountPoints []JottaMountPoint `xml:"mountPoints>mountPoint"`
}

/*
GET http://www.jottacloud.com/JFS/<account>/<device>/<mountpoint>

<mountPoint time="2018-07-24-T20:35:02Z" host="dn-157">
	<name xml:space="preserve">Sync</name>
	<path xml:space="preserve">/12qh1wsht8cssxdtwl15rqh9/Jotta</path>
	<abspath xml:space="preserve">/12qh1wsht8cssxdtwl15rqh9/Jotta</abspath>
	<size>0</size>
	<modified></modified>
	<device>Jotta</device>
	<user>12qh1wsht8cssxdtwl15rqh9</user>
	<folders>
		<folder name="test"/>
	</folders>
	<metadata first="" max="" total="1" num_folders="1" num_files="0"/>
</mountPoint>
*/

// JottaMountPoint represents a Jottacloud mountpoint
type JottaMountPoint struct {
	Name    string        `xml:"name"`
	Size    int64         `xml:"size"`
	Device  string        `xml:"device"`
	Folders []JottaFolder `xml:"folders>folder"`
	Files   []JottaFile   `xml:"files>file"`
}

/*
GET http://www.jottacloud.com/JFS/<account>/<device>/<mountpoint>/<folder>

<folder name="test" time="2018-07-24-T20:41:37Z" host="dn-158">
	<path xml:space="preserve">/12qh1wsht8cssxdtwl15rqh9/Jotta/Sync</path>
	<abspath xml:space="preserve">/12qh1wsht8cssxdtwl15rqh9/Jotta/Sync</abspath>
	<folders>
		<folder name="t2"/>c
	</folders>
	<files>
		<file name="block.csv" uuid="f6553cd4-1135-48fe-8e6a-bb9565c50ef2">
			<currentRevision>
				<number>1</number>
				<state>COMPLETED</state>
				<created>2018-07-05-T15:08:02Z</created>
				<modified>2018-07-05-T15:08:02Z</modified>
				<mime>application/octet-stream</mime>
				<size>30827730</size>
				<md5>1e8a7b728ab678048df00075c9507158</md5>
				<updated>2018-07-24-T20:41:10Z</updated>
			</currentRevision>
		</file>
	</files>
	<metadata first="" max="" total="2" num_folders="1" num_files="1"/>
</folder>
*/

// JottaFolder represents a JottacloudFolder
type JottaFolder struct {
	XMLName    xml.Name
	Name       string        `xml:"name,attr"`
	Deleted    Flag          `xml:"deleted,attr"`
	Path       string        `xml:"path"`
	CreatedAt  Time          `xml:"created"`
	ModifiedAt Time          `xml:"modified"`
	Updated    Time          `xml:"updated"`
	Folders    []JottaFolder `xml:"folders>folder"`
	Files      []JottaFile   `xml:"files>file"`
}

/*
GET http://www.jottacloud.com/JFS/<account>/<device>/<mountpoint>/.../<file>

<file name="block.csv" uuid="f6553cd4-1135-48fe-8e6a-bb9565c50ef2">
	<currentRevision>
		<number>1</number>
		<state>COMPLETED</state>
		<created>2018-07-05-T15:08:02Z</created>
		<modified>2018-07-05-T15:08:02Z</modified>
		<mime>application/octet-stream</mime>
		<size>30827730</size>
		<md5>1e8a7b728ab678048df00075c9507158</md5>
		<updated>2018-07-24-T20:41:10Z</updated>
	</currentRevision>
</file>
*/

// JottaFile represents a Jottacloud file
type JottaFile struct {
	XMLName         xml.Name
	Name            string `xml:"name,attr"`
	Deleted         Flag   `xml:"deleted,attr"`
	PublicSharePath string `xml:"publicSharePath"`
	State           string `xml:"currentRevision>state"`
	CreatedAt       Time   `xml:"currentRevision>created"`
	ModifiedAt      Time   `xml:"currentRevision>modified"`
	Updated         Time   `xml:"currentRevision>updated"`
	Size            int64  `xml:"currentRevision>size"`
	MimeType        string `xml:"currentRevision>mime"`
	MD5             string `xml:"currentRevision>md5"`
}

// Error is a custom Error for wrapping Jottacloud error responses
type Error struct {
	StatusCode int    `xml:"code"`
	Message    string `xml:"message"`
	Reason     string `xml:"reason"`
	Cause      string `xml:"cause"`
}

// Error returns a string for the error and statistifes the error interface
func (e *Error) Error() string {
	out := fmt.Sprintf("error %d", e.StatusCode)
	if e.Message != "" {
		out += ": " + e.Message
	}
	if e.Reason != "" {
		out += fmt.Sprintf(" (%+v)", e.Reason)
	}
	return out
}

// AllocateFileRequest to prepare an upload to Jottacloud
type AllocateFileRequest struct {
	Bytes    int64  `json:"bytes"`
	Created  string `json:"created"`
	Md5      string `json:"md5"`
	Modified string `json:"modified"`
	Path     string `json:"path"`
}

// AllocateFileResponse for upload requests
type AllocateFileResponse struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	State     string `json:"state"`
	UploadID  string `json:"upload_id"`
	UploadURL string `json:"upload_url"`
	Bytes     int64  `json:"bytes"`
	ResumePos int64  `json:"resume_pos"`
}

// UploadResponse after an upload
type UploadResponse struct {
	Name      string      `json:"name"`
	Path      string      `json:"path"`
	Kind      string      `json:"kind"`
	ContentID string      `json:"content_id"`
	Bytes     int64       `json:"bytes"`
	Md5       string      `json:"md5"`
	Created   int64       `json:"created"`
	Modified  int64       `json:"modified"`
	Deleted   interface{} `json:"deleted"`
	Mime      string      `json:"mime"`
}

// DeviceRegistrationResponse is the response to registering a device
type DeviceRegistrationResponse struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// CustomerInfo provides general information about the account. Required for finding the correct internal username.
type CustomerInfo struct {
	Username          string      `json:"username"`
	Email             string      `json:"email"`
	Name              string      `json:"name"`
	CountryCode       string      `json:"country_code"`
	LanguageCode      string      `json:"language_code"`
	CustomerGroupCode string      `json:"customer_group_code"`
	BrandCode         string      `json:"brand_code"`
	AccountType       string      `json:"account_type"`
	SubscriptionType  string      `json:"subscription_type"`
	Usage             int64       `json:"usage"`
	Qouta             int64       `json:"quota"`
	BusinessUsage     int64       `json:"business_usage"`
	BusinessQouta     int64       `json:"business_quota"`
	WriteLocked       bool        `json:"write_locked"`
	ReadLocked        bool        `json:"read_locked"`
	LockedCause       interface{} `json:"locked_cause"`
	WebHash           string      `json:"web_hash"`
	AndroidHash       string      `json:"android_hash"`
	IOSHash           string      `json:"ios_hash"`
}
