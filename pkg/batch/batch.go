// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package batch

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"os"

	"github.com/moov-io/fincen"
	"github.com/moov-io/fincen/pkg/cash_payments"
	"github.com/moov-io/fincen/pkg/currency_transaction"
	"github.com/moov-io/fincen/pkg/exempt_designation"
	"github.com/moov-io/fincen/pkg/financial_accounts"
	"github.com/moov-io/fincen/pkg/suspicious_activity"
)

const (
	ReportTypeSubmission = "SUBMISSION"
)

func NewReport(args ...string) *EFilingBatchXML {

	reportXml := EFilingBatchXML{}

	rType := "UNKNOWN"
	if len(args) > 0 {
		rType = args[0]
	}
	if rType == ReportTypeSubmission {
		reportXml.StatusCode = "A"
	} else {
		if !fincen.CheckInvolved(rType, "CTRX", "SARX", "DOEPX", "FBARX", "8300X") {
			reportXml.FormTypeCode = rType
		}
	}

	return &reportXml
}

func CreateReportWithBuffer(buf []byte) (*EFilingBatchXML, error) {

	reportXml := EFilingBatchXML{}

	err := xml.Unmarshal(buf, &reportXml)
	if err == nil {
		return &reportXml, nil
	}

	err = json.Unmarshal(buf, &reportXml)
	if err == nil {
		return &reportXml, nil
	}

	return nil, errors.New("unable to create batch, invalid input data")
}

func CreateReportWithFile(path string) (*EFilingBatchXML, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("opening file %s: %w", path, err)
	}

	r, err := CreateReportWithBuffer(raw)
	if err != nil {
		return nil, fmt.Errorf("unable to parse file: %w", err)
	}

	return r, nil
}

type EFilingBatchXML struct {
	XMLName                 xml.Name                 `xml:"EFilingBatchXML"`
	SeqNum                  fincen.SeqNumber         `xml:"SeqNum,attr"`
	StatusCode              string                   `xml:"StatusCode,attr,omitempty" json:",omitempty"`
	TotalAmount             float64                  `xml:"TotalAmount,attr,omitempty" json:",omitempty"`
	PartyCount              int64                    `xml:"PartyCount,attr,omitempty" json:",omitempty"`
	ActivityCount           int64                    `xml:"ActivityCount,attr,omitempty" json:",omitempty"`
	AccountCount            int64                    `xml:"AccountCount,attr,omitempty" json:",omitempty"`
	ActivityAttachmentCount int64                    `xml:"ActivityAttachmentCount,attr,omitempty" json:",omitempty"`
	AttachmentCount         int64                    `xml:"AttachmentCount,attr,omitempty" json:",omitempty"`
	JointlyOwnedOwnerCount  int64                    `xml:"JointlyOwnedOwnerCount,attr,omitempty" json:",omitempty"`
	NoFIOwnerCount          int64                    `xml:"NoFIOwnerCount,attr,omitempty" json:",omitempty"`
	ConsolidatedOwnerCount  int64                    `xml:"ConsolidatedOwnerCount,attr,omitempty" json:",omitempty"`
	Attrs                   []xml.Attr               `xml:",any,attr"`
	FormTypeCode            string                   `xml:"FormTypeCode,omitempty" json:",omitempty"`
	Activity                []fincen.ElementActivity `xml:"Activity,omitempty" json:",omitempty"`
	EFilingSubmissionXML    *EFilingSubmissionXML    `xml:"EFilingSubmissionXML,omitempty" json:",omitempty"`
}

type dummyXML struct {
	XMLName    xml.Name
	Attrs      []xml.Attr       `xml:",any,attr"`
	SeqNum     fincen.SeqNumber `xml:"SeqNum,attr"`
	StatusCode string           `xml:"StatusCode,attr,omitempty" json:",omitempty"`
	Content    []byte           `xml:",innerxml"`
}

type batchDummy struct {
	XMLName                 xml.Name              `xml:"EFilingBatchXML"`
	SeqNum                  fincen.SeqNumber      `xml:"SeqNum,attr"`
	StatusCode              string                `xml:"StatusCode,attr,omitempty" json:",omitempty"`
	TotalAmount             float64               `xml:"TotalAmount,attr,omitempty" json:",omitempty"`
	PartyCount              int64                 `xml:"PartyCount,attr,omitempty" json:",omitempty"`
	ActivityCount           int64                 `xml:"ActivityCount,attr,omitempty" json:",omitempty"`
	AccountCount            int64                 `xml:"AccountCount,attr,omitempty" json:",omitempty"`
	ActivityAttachmentCount int64                 `xml:"ActivityAttachmentCount,attr,omitempty" json:",omitempty"`
	AttachmentCount         int64                 `xml:"AttachmentCount,attr,omitempty" json:",omitempty"`
	JointlyOwnedOwnerCount  int64                 `xml:"JointlyOwnedOwnerCount,attr,omitempty" json:",omitempty"`
	NoFIOwnerCount          int64                 `xml:"NoFIOwnerCount,attr,omitempty" json:",omitempty"`
	ConsolidatedOwnerCount  int64                 `xml:"ConsolidatedOwnerCount,attr,omitempty" json:",omitempty"`
	Attrs                   []xml.Attr            `xml:",any,attr"`
	FormTypeCode            string                `xml:"FormTypeCode,omitempty" json:",omitempty"`
	Activity                []dummyXML            `xml:"Activity,omitempty" json:",omitempty"`
	EFilingSubmissionXML    *EFilingSubmissionXML `xml:"EFilingSubmissionXML,omitempty" json:",omitempty"`
}

type batchAttr struct {
	TotalAmount             float64
	PartyCount              int64
	ActivityCount           int64
	AccountCount            int64
	ActivityAttachmentCount int64
	AttachmentCount         int64
	JointlyOwnedOwnerCount  int64
	NoFIOwnerCount          int64
	ConsolidatedOwnerCount  int64
}

func (r *EFilingBatchXML) copy(org batchDummy) {
	// copy object
	r.XMLName = org.XMLName
	r.Attrs = org.Attrs
	r.SeqNum = org.SeqNum
	r.StatusCode = org.StatusCode
	r.TotalAmount = org.TotalAmount
	r.PartyCount = org.PartyCount
	r.ActivityCount = org.ActivityCount
	r.AccountCount = org.AccountCount
	r.ActivityAttachmentCount = org.ActivityAttachmentCount
	r.AttachmentCount = org.AttachmentCount
	r.JointlyOwnedOwnerCount = org.JointlyOwnedOwnerCount
	r.NoFIOwnerCount = org.NoFIOwnerCount
	r.ConsolidatedOwnerCount = org.ConsolidatedOwnerCount
	r.FormTypeCode = org.FormTypeCode
	r.EFilingSubmissionXML = org.EFilingSubmissionXML
}

func (r *EFilingBatchXML) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var dummy batchDummy
	if err := d.DecodeElement(&dummy, &start); err != nil {
		return err
	}

	r.copy(dummy)

	for i := range dummy.Activity {
		act := dummy.Activity[i]

		buf, err := xml.Marshal(&act)
		if err != nil {
			return fincen.NewErrValueInvalid("Activity")
		}

		constructor := activityConstructor[r.FormTypeCode]
		if constructor == nil {
			return fincen.NewErrValueInvalid("FormTypeCode")
		}

		elm := constructor()
		if err = xml.Unmarshal(buf, elm); err != nil {
			return fincen.NewErrValueInvalid("Activity")
		}

		r.Activity = append(r.Activity, elm)
	}

	return nil
}

func (r EFilingBatchXML) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	a := struct {
		XMLName                 xml.Name                 `xml:"EFilingBatchXML"`
		SeqNum                  fincen.SeqNumber         `xml:"SeqNum,attr"`
		StatusCode              string                   `xml:"StatusCode,attr,omitempty" json:",omitempty"`
		TotalAmount             float64                  `xml:"TotalAmount,attr,omitempty" json:",omitempty"`
		PartyCount              int64                    `xml:"PartyCount,attr,omitempty" json:",omitempty"`
		ActivityCount           int64                    `xml:"ActivityCount,attr,omitempty" json:",omitempty"`
		AccountCount            int64                    `xml:"AccountCount,attr,omitempty" json:",omitempty"`
		ActivityAttachmentCount int64                    `xml:"ActivityAttachmentCount,attr,omitempty" json:",omitempty"`
		AttachmentCount         int64                    `xml:"AttachmentCount,attr,omitempty" json:",omitempty"`
		JointlyOwnedOwnerCount  int64                    `xml:"JointlyOwnedOwnerCount,attr,omitempty" json:",omitempty"`
		NoFIOwnerCount          int64                    `xml:"NoFIOwnerCount,attr,omitempty" json:",omitempty"`
		ConsolidatedOwnerCount  int64                    `xml:"ConsolidatedOwnerCount,attr,omitempty" json:",omitempty"`
		Attrs                   []xml.Attr               `xml:",any,attr"`
		FormTypeCode            string                   `xml:"FormTypeCode,omitempty" json:",omitempty"`
		Activity                []fincen.ElementActivity `xml:"Activity,omitempty" json:",omitempty"`
		EFilingSubmissionXML    *EFilingSubmissionXML    `xml:"EFilingSubmissionXML,omitempty" json:",omitempty"`
	}(r)

	for index := 0; index < len(a.Attrs); index++ {
		switch a.Attrs[index].Name.Local {
		case "schemaLocation", "xsi", "fc2":
			a.Attrs = append(a.Attrs[:index], a.Attrs[index+1:]...)
			index--
		}
	}

	a.Attrs = append(a.Attrs, xml.Attr{
		Name: xml.Name{
			Local: "xsi:schemaLocation",
		},
		Value: "www.server.gov/base https://www.fincen.gov/base https://www.fincen.gov/base/EFL_8300XBatchSchema.xsd",
	})

	a.Attrs = append(a.Attrs, xml.Attr{
		Name: xml.Name{
			Local: "xmlns:xsi",
		},
		Value: "http://www.w3.org/2001/XMLSchema-instance",
	})

	a.Attrs = append(a.Attrs, xml.Attr{
		Name: xml.Name{
			Local: "xsi:fc2",
		},
		Value: "www.server.gov/base",
	})

	return e.EncodeElement(&a, start)
}

func (r *EFilingBatchXML) AppendActivity(act fincen.ElementActivity) error {
	if act == nil {
		return errors.New("invalid activity")
	}

	if !fincen.CheckInvolved(r.FormTypeCode, "CTRX", "SARX", "DOEPX", "FBARX", "8300X") ||
		r.FormTypeCode != act.FormTypeCode() {
		return errors.New("invalid form type")
	}

	r.Activity = append(r.Activity, act)

	return nil
}

func (r *EFilingBatchXML) fieldInclusionReport() error {

	if len(r.Activity) < 1 {
		return fincen.NewErrValueInvalid("Activity")
	}

	return nil
}

func (r *EFilingBatchXML) fieldInclusionSubmission() error {
	if r.EFilingSubmissionXML == nil {
		return fincen.NewErrValueInvalid("EFilingSubmissionXML")
	}

	return nil
}

func (r EFilingBatchXML) generateAttrs() batchAttr {
	s := batchAttr{}
	// The count of all <Activity> elements in the batch
	s.ActivityCount = int64(len(r.Activity))

	for _, activity := range r.Activity {
		s.TotalAmount += activity.TotalAmount()

		switch r.FormTypeCode {
		case "8300X":
			// The count of all <Party> elements in the batch where the
			// <ActivityPartyTypeCode> is equal to 16, 23, 4, 3, and 8 (combined)
			s.PartyCount += activity.PartyCount("16", "23", "4", "3", "8")
		case "DOEPX":
			// The count of all <Party> elements in the batch where the
			//<ActivityPartyTypeCode> is equal to 3, 11, 12, and 45 (combined)
			s.PartyCount += activity.PartyCount("3", "11", "12", "45")
		case "CTRX":
			// The total count of all <Party> elements recorded in the batch file.
			s.PartyCount += activity.PartyCount()
		case "SARX":
			// The count of all <Party> elements in the batch where the
			// <ActivityPartyTypeCode> is equal to “33” (Subject)
			s.PartyCount += activity.PartyCount("33")

		case "FBARX":
			// The total count of <Party> elements where the <ActivityPartyTypeCode> element is equal to “41”
			s.PartyCount += activity.PartyCount("41")

			// AccountCount. The total count of all <Account> elements recorded in the batch file.
			if elm, ok := activity.(*financial_accounts.ActivityType); ok {
				s.AccountCount = int64(len(elm.Account))
			}

			// The total count of <Party> elements where the <ActivityPartyTypeCode>
			// element is equal to “42”
			s.JointlyOwnedOwnerCount += activity.PartyCount("42")

			// The total count of <Party> elements where the <ActivityPartyTypeCode> element is
			// equal to “43”
			s.NoFIOwnerCount += activity.PartyCount("43")

			// The total count of <Party> elements where the <ActivityPartyTypeCode>
			// equal to “44”
			s.NoFIOwnerCount += activity.PartyCount("44")
		}

	}

	return s
}

func (r *EFilingBatchXML) GenerateAttrs() error {

	s := r.generateAttrs()

	r.ActivityCount = s.ActivityCount
	r.TotalAmount = s.TotalAmount
	r.PartyCount = s.PartyCount
	r.ActivityAttachmentCount = s.ActivityAttachmentCount
	r.AttachmentCount = s.AttachmentCount
	r.JointlyOwnedOwnerCount = s.JointlyOwnedOwnerCount
	r.NoFIOwnerCount = s.NoFIOwnerCount
	r.ConsolidatedOwnerCount = s.ConsolidatedOwnerCount

	return nil
}

func (r *EFilingBatchXML) validateAttrs() error {

	s := r.generateAttrs()

	if r.ActivityCount != s.ActivityCount {
		return fincen.NewErrValueInvalid("ActivityCount")
	}

	if r.AccountCount != s.AccountCount {
		return fincen.NewErrValueInvalid("AccountCount")
	}

	if r.TotalAmount != s.TotalAmount {
		return fincen.NewErrValueInvalid("TotalAmount")
	}

	if r.PartyCount != s.PartyCount {
		return fincen.NewErrValueInvalid("PartyCount")
	}

	if r.ActivityAttachmentCount != s.ActivityAttachmentCount {
		return fincen.NewErrValueInvalid("ActivityAttachmentCount")
	}

	if r.AttachmentCount != s.AttachmentCount {
		return fincen.NewErrValueInvalid("AttachmentCount")
	}

	if r.JointlyOwnedOwnerCount != s.JointlyOwnedOwnerCount {
		return fincen.NewErrValueInvalid("JointlyOwnedOwnerCount")
	}

	if r.NoFIOwnerCount != s.NoFIOwnerCount {
		return fincen.NewErrValueInvalid("NoFIOwnerCount")
	}

	if r.ConsolidatedOwnerCount != s.ConsolidatedOwnerCount {
		return fincen.NewErrValueInvalid("ConsolidatedOwnerCount")
	}

	return nil
}

// Validate args:
//
//	1: disableValidateAttrs
func (r EFilingBatchXML) Validate(args ...string) error {

	if r.StatusCode == "A" {
		// FinCEN XML Acknowledgement
		if err := r.fieldInclusionSubmission(); err != nil {
			return err
		}
	} else {

		// FinCEN XML Batch Reporting
		if !fincen.CheckInvolved(r.FormTypeCode, "CTRX", "SARX", "DOEPX", "FBARX", "8300X") {
			return fincen.NewErrValueInvalid("FormTypeCode")
		}

		if err := r.fieldInclusionReport(); err != nil {
			return err
		}
	}

	if len(args) == 0 {
		if err := r.validateAttrs(); err != nil {
			return err
		}
	}

	return fincen.Validate(&r, args...)
}

func (r EFilingBatchXML) GenerateSeqNumbers() error {
	return fincen.GenerateSeqNumbers(&r)
}

type EFilingSubmissionXML struct {
	XMLName            xml.Name             `xml:"EFilingSubmissionXML"`
	SeqNum             fincen.SeqNumber     `xml:"SeqNum,attr"`
	StatusCode         string               `xml:"StatusCode,attr,omitempty" json:",omitempty"`
	EFilingActivityXML []EFilingActivityXML `xml:"EFilingActivityXML"`
}

func (r EFilingSubmissionXML) Validate(args ...string) error {
	return fincen.Validate(&r, args...)
}

type EFilingActivityXML struct {
	XMLName                 xml.Name                  `xml:"EFilingActivityXML"`
	SeqNum                  fincen.SeqNumber          `xml:"SeqNum,attr"`
	BSAID                   fincen.RestrictNumeric14  `xml:"BSAID"`
	EFilingActivityErrorXML []EFilingActivityErrorXML `xml:"EFilingActivityErrorXML"`
}

func (r EFilingActivityXML) Validate(args ...string) error {
	return fincen.Validate(&r, args...)
}

type EFilingActivityErrorXML struct {
	XMLName              xml.Name                   `xml:"EFilingActivityErrorXML"`
	SeqNum               fincen.SeqNumber           `xml:"SeqNum,attr"`
	ErrorContextText     *fincen.RestrictString4000 `xml:"ErrorContextText,omitempty" json:",omitempty"`
	ErrorElementNameText *fincen.RestrictString512  `xml:"ErrorElementNameText,omitempty" json:",omitempty"`
	ErrorLevelText       *fincen.RestrictString50   `xml:"ErrorLevelText,omitempty" json:",omitempty"`
	ErrorText            *fincen.RestrictString525  `xml:"ErrorText,omitempty" json:",omitempty"`
	ErrorTypeCode        *fincen.RestrictString50   `xml:"ErrorTypeCode,omitempty" json:",omitempty"`
}

func (r EFilingActivityErrorXML) Validate(args ...string) error {
	return fincen.Validate(&r, args...)
}

type constructorFunc func() fincen.ElementActivity

var (
	activityConstructor = map[string]constructorFunc{
		"CTRX":  func() fincen.ElementActivity { return &currency_transaction.ActivityType{} },
		"SARX":  func() fincen.ElementActivity { return &suspicious_activity.ActivityType{} },
		"DOEPX": func() fincen.ElementActivity { return &exempt_designation.ActivityType{} },
		"FBARX": func() fincen.ElementActivity { return &financial_accounts.ActivityType{} },
		"8300X": func() fincen.ElementActivity { return &cash_payments.ActivityType{} },
	}
)
