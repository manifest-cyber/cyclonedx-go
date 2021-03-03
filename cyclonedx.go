package cyclonedx

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

const (
	BOMFormat      = "CycloneDX"
	defaultVersion = 1
	SpecVersion    = "1.2"
	XMLNamespace   = "http://cyclonedx.org/schema/bom/1.2"
)

type AttachedText struct {
	Content     string `json:"content" xml:",innerxml"`
	ContentType string `json:"contentType,omitempty" xml:"content-type,attr,omitempty"`
	Encoding    string `json:"encoding,omitempty" xml:"encoding,attr,omitempty"`
}

type BOM struct {
	// XML specific fields
	XMLName xml.Name `json:"-" xml:"bom"`
	XMLNS   string   `json:"-" xml:"xmlns,attr"`

	// JSON specific fields
	BOMFormat   string `json:"bomFormat" xml:"-"`
	SpecVersion string `json:"specVersion" xml:"-"`

	SerialNumber       string               `json:"serialNumber,omitempty" xml:"serialNumber,attr,omitempty"`
	Version            int                  `json:"version" xml:"version,attr"`
	Metadata           *Metadata            `json:"metadata,omitempty" xml:"metadata,omitempty"`
	Components         *[]Component         `json:"components,omitempty" xml:"components>component,omitempty"`
	Services           *[]Service           `json:"services,omitempty" xml:"services>service,omitempty"`
	ExternalReferences *[]ExternalReference `json:"externalReferences,omitempty" xml:"externalReferences>reference,omitempty"`
	Dependencies       *[]Dependency        `json:"dependencies,omitempty" xml:"dependencies>dependency,omitempty"`
}

func NewBOM() *BOM {
	return &BOM{
		XMLNS:       XMLNamespace,
		BOMFormat:   BOMFormat,
		SpecVersion: SpecVersion,
		Version:     defaultVersion,
	}
}

type BOMFileFormat int

const (
	BOMFileFormatXML BOMFileFormat = iota
	BOMFileFormatJSON
)

// Bool is a convenience function to transform a value of the primitive type bool to a pointer of bool
func Bool(value bool) *bool {
	return &value
}

type ComponentType string

const (
	ComponentTypeApplication ComponentType = "application"
	ComponentTypeContainer   ComponentType = "container"
	ComponentTypeDevice      ComponentType = "device"
	ComponentTypeFile        ComponentType = "file"
	ComponentTypeFirmware    ComponentType = "firmware"
	ComponentTypeFramework   ComponentType = "framework"
	ComponentTypeLibrary     ComponentType = "library"
	ComponentTypeOS          ComponentType = "operating-system"
)

type Commit struct {
	UID       string              `json:"uid,omitempty" xml:"uid,omitempty"`
	URL       string              `json:"url,omitempty" xml:"url,omitempty"`
	Author    *IdentifiableAction `json:"author,omitempty" xml:"author,omitempty"`
	Committer *IdentifiableAction `json:"committer,omitempty" xml:"committer,omitempty"`
	Message   string              `json:"message,omitempty" xml:"message,omitempty"`
}

type Component struct {
	Supplier           *OrganizationalEntity `json:"supplier,omitempty" xml:"supplier,omitempty"`
	Author             string                `json:"author,omitempty" xml:"author,omitempty"`
	Publisher          string                `json:"publisher,omitempty" xml:"publisher,omitempty"`
	Group              string                `json:"group,omitempty" xml:"group,omitempty"`
	Name               string                `json:"name" xml:"name"`
	Version            string                `json:"version" xml:"version"`
	Description        string                `json:"description,omitempty" xml:"description,omitempty"`
	Scope              Scope                 `json:"scope,omitempty" xml:"scope,omitempty"`
	Hashes             *[]Hash               `json:"hashes,omitempty" xml:"hashes>hash,omitempty"`
	Licenses           *[]LicenseChoice      `json:"licenses,omitempty" xml:"licenses>license,omitempty"`
	Copyright          string                `json:"copyright,omitempty" xml:"copyright,omitempty"`
	CPE                string                `json:"cpe,omitempty" xml:"cpe,omitempty"`
	PackageURL         string                `json:"purl,omitempty" xml:"purl,omitempty"`
	SWID               *SWID                 `json:"swid,omitempty" xml:"swid,omitempty"`
	Modified           *bool                 `json:"modified,omitempty" xml:"modified,omitempty"`
	Pedigree           *Pedigree             `json:"pedigree,omitempty" xml:"pedigree,omitempty"`
	ExternalReferences *[]ExternalReference  `json:"externalReferences,omitempty" xml:"externalReferences>reference,omitempty"`
	Components         *[]Component          `json:"components,omitempty" xml:"components>component,omitempty"`
	BOMRef             string                `json:"bom-ref,omitempty" xml:"bom-ref,attr,omitempty"`
	MIMEType           string                `json:"mime-type,omitempty" xml:"mime-type,attr,omitempty"`
	Type               ComponentType         `json:"type" xml:"type,attr"`
}

type DataClassification struct {
	Flow           DataFlow `json:"flow" xml:"flow,attr"`
	Classification string   `json:"classification" xml:",innerxml"`
}

type DataFlow string

const (
	DataFlowBidirectional DataFlow = "bi-directional"
	DataFlowInbound       DataFlow = "inbound"
	DataFlowOutbound      DataFlow = "outbound"
	DataFlowUnknown       DataFlow = "unknown"
)

type Dependency struct {
	Ref          string        `json:"-" xml:"ref,attr"`
	Dependencies *[]Dependency `json:"-" xml:"dependency,omitempty"`
}

// dependencyJSON is temporarily used for marshalling and unmarshalling Dependency instances to and from JSON
type dependencyJSON struct {
	Ref       string   `json:"ref" xml:"-"`
	DependsOn []string `json:"dependsOn,omitempty" xml:"-"`
}

func (d Dependency) MarshalJSON() ([]byte, error) {
	if d.Dependencies == nil || len(*d.Dependencies) == 0 {
		return json.Marshal(&dependencyJSON{
			Ref: d.Ref,
		})
	}

	dependencyRefs := make([]string, len(*d.Dependencies))
	for i, dependency := range *d.Dependencies {
		dependencyRefs[i] = dependency.Ref
	}

	return json.Marshal(&dependencyJSON{
		Ref:       d.Ref,
		DependsOn: dependencyRefs,
	})
}

func (d *Dependency) UnmarshalJSON(bytes []byte) error {
	dependency := new(dependencyJSON)
	if err := json.Unmarshal(bytes, dependency); err != nil {
		return err
	}
	d.Ref = dependency.Ref

	if len(dependency.DependsOn) == 0 {
		return nil
	}

	dependencies := make([]Dependency, len(dependency.DependsOn))
	for i, dep := range dependency.DependsOn {
		dependencies[i] = Dependency{
			Ref: dep,
		}
	}
	d.Dependencies = &dependencies

	return nil
}

type Diff struct {
	Text *AttachedText `json:"text,omitempty" xml:"text,omitempty"`
	URL  string        `json:"url,omitempty" xml:"url,omitempty"`
}

type ExternalReference struct {
	URL     string                `json:"url" xml:"url"`
	Comment string                `json:"comment,omitempty" xml:"comment,omitempty"`
	Type    ExternalReferenceType `json:"type" xml:"type,attr"`
}

type ExternalReferenceType string

const (
	ERTypeAdvisories    ExternalReferenceType = "advisories"
	ERTypeBOM           ExternalReferenceType = "bom"
	ERTypeBuildMeta     ExternalReferenceType = "build-meta"
	ERTypeBuildSystem   ExternalReferenceType = "build-system"
	ERTypeChat          ExternalReferenceType = "chat"
	ERTypeDistribution  ExternalReferenceType = "distribution"
	ERTypeDocumentation ExternalReferenceType = "documentation"
	ERTypeLicense       ExternalReferenceType = "license"
	ERTypeMailingList   ExternalReferenceType = "mailing-list"
	ERTypeOther         ExternalReferenceType = "other"
	ERTypeIssueTracker  ExternalReferenceType = "issue-tracker"
	ERTypeSocial        ExternalReferenceType = "social"
	ERTypeSupport       ExternalReferenceType = "support"
	ERTypeVCS           ExternalReferenceType = "vcs"
	ERTypeWebsite       ExternalReferenceType = "website"
)

type Hash struct {
	Algorithm HashAlgorithm `json:"alg" xml:"alg,attr"`
	Value     string        `json:"content" xml:",innerxml"`
}

type HashAlgorithm string

const (
	HashAlgoMD5         HashAlgorithm = "MD5"
	HashAlgoSHA1        HashAlgorithm = "SHA-1"
	HashAlgoSHA256      HashAlgorithm = "SHA-256"
	HashAlgoSHA384      HashAlgorithm = "SHA-384"
	HashAlgoSHA512      HashAlgorithm = "SHA-512"
	HashAlgoSHA3_256    HashAlgorithm = "SHA3-256"
	HashAlgoSHA3_512    HashAlgorithm = "SHA3-512"
	HashAlgoBlake2b_256 HashAlgorithm = "BLAKE2b-256"
	HashAlgoBlake2b_384 HashAlgorithm = "BLAKE2b-384"
	HashAlgoBlake2b_512 HashAlgorithm = "BLAKE2b-512"
	HashAlgoBlake3      HashAlgorithm = "BLAKE3"
)

type IdentifiableAction struct {
	Timestamp string `json:"timestamp,omitempty" xml:"timestamp,omitempty"`
	Name      string `json:"name,omitempty" xml:"name,omitempty"`
	EMail     string `json:"email,omitempty" xml:"email,omitempty"`
}

type Issue struct {
	ID          string    `json:"id" xml:"id"`
	Name        string    `json:"name" xml:"name"`
	Description string    `json:"description" xml:"description"`
	Source      *Source   `json:"source,omitempty" xml:"source,omitempty"`
	References  *[]string `json:"references,omitempty" xml:"references>url,omitempty"`
	Type        IssueType `json:"type" xml:"type,attr"`
}

type IssueType string

const (
	IssueTypeDefect      IssueType = "defect"
	IssueTypeEnhancement IssueType = "enhancement"
	IssueTypeSecurity    IssueType = "security"
)

type License struct {
	ID   string        `json:"id,omitempty" xml:"id,omitempty"`
	Name string        `json:"name,omitempty" xml:"name,omitempty"`
	Text *AttachedText `json:"text,omitempty" xml:"text,omitempty"`
	URL  string        `json:"url,omitempty" xml:"url,omitempty"`
}

type LicenseChoice struct {
	License    *License `json:"license,omitempty" xml:"-"`
	Expression string   `json:"expression,omitempty" xml:"-"`
}

func (l LicenseChoice) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	if l.License != nil && l.Expression != "" {
		return fmt.Errorf("either license or expression must be set, but not both")
	}

	if l.License != nil {
		return e.EncodeElement(l.License, xml.StartElement{Name: xml.Name{Local: "license"}})
	} else if l.Expression != "" {
		expressionElement := xml.StartElement{Name: xml.Name{Local: "expression"}}
		e.EncodeToken(expressionElement)
		e.EncodeToken(xml.CharData(l.Expression))
		return e.EncodeToken(xml.EndElement{Name: expressionElement.Name})
	}

	// Neither license nor expression set - don't write anything
	return nil
}

func (l *LicenseChoice) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if start.Name.Local == "license" {
		license := new(License)
		if err := d.DecodeElement(license, &start); err != nil {
			return err
		}
		l.License = license
		l.Expression = ""
		return nil
	} else if start.Name.Local == "expression" {
		expression := new(string)
		if err := d.DecodeElement(expression, &start); err != nil {
			return err
		}
		l.License = nil
		l.Expression = *expression
		return nil
	}

	return xml.UnmarshalError(fmt.Sprintf("cannot unmarshal element %#v", start))
}

type Metadata struct {
	Timestamp   string                   `json:"timestamp,omitempty" xml:"timestamp,omitempty"`
	Tools       *[]Tool                  `json:"tools,omitempty" xml:"tools>tool,omitempty"`
	Authors     *[]OrganizationalContact `json:"authors,omitempty" xml:"authors>author,omitempty"`
	Component   *Component               `json:"component,omitempty" xml:"component,omitempty"`
	Manufacture *OrganizationalEntity    `json:"manufacture,omitempty" xml:"manufacture,omitempty"`
	Supplier    *OrganizationalEntity    `json:"supplier,omitempty" xml:"supplier,omitempty"`
}

type OrganizationalContact struct {
	Name  string `json:"name,omitempty" xml:"name,omitempty"`
	EMail string `json:"email,omitempty" xml:"email,omitempty"`
	Phone string `json:"phone,omitempty" xml:"phone,omitempty"`
}

type OrganizationalEntity struct {
	Name    string                   `json:"name" xml:"name"`
	URL     *[]string                `json:"url,omitempty" xml:"url,omitempty"`
	Contact *[]OrganizationalContact `json:"contact,omitempty" xml:"contact,omitempty"`
}

type Patch struct {
	Diff     *Diff     `json:"diff,omitempty" xml:"diff,omitempty"`
	Resolves *[]Issue  `json:"resolves,omitempty" xml:"resolves>issue,omitempty"`
	Type     PatchType `json:"type" xml:"type,attr"`
}

type PatchType string

const (
	PatchTypeBackport   PatchType = "backport"
	PatchTypeCherryPick PatchType = "cherry-pick"
	PatchTypeMonkey     PatchType = "monkey"
	PatchTypeUnofficial PatchType = "unofficial"
)

type Pedigree struct {
	Ancestors   *[]Component `json:"ancestors,omitempty" xml:"ancestors>component,omitempty"`
	Descendants *[]Component `json:"descendants,omitempty" xml:"descendants>component,omitempty"`
	Variants    *[]Component `json:"variants,omitempty" xml:"variants>component,omitempty"`
	Commits     *[]Commit    `json:"commits,omitempty" xml:"commits>commit,omitempty"`
	Patches     *[]Patch     `json:"patches,omitempty" xml:"patches>patch,omitempty"`
	Notes       string       `json:"notes,omitempty" xml:"notes,omitempty"`
}

type Scope string

const (
	ScopeExcluded Scope = "excluded"
	ScopeOptional Scope = "optional"
	ScopeRequired Scope = "required"
)

type Service struct {
	Provider             *OrganizationalEntity `json:"provider,omitempty" xml:"provider,omitempty"`
	Group                string                `json:"group,omitempty" xml:"group,omitempty"`
	Name                 string                `json:"name" xml:"name"`
	Version              string                `json:"version,omitempty" xml:"version,omitempty"`
	Description          string                `json:"description,omitempty" xml:"description,omitempty"`
	Endpoints            *[]string             `json:"endpoints,omitempty" xml:"endpoints>endpoint,omitempty"`
	Authenticated        *bool                 `json:"authenticated,omitempty" xml:"authenticated,omitempty"`
	CrossesTrustBoundary *bool                 `json:"x-trust-boundary,omitempty" xml:"x-trust-boundary,omitempty"`
	Data                 *[]DataClassification `json:"data,omitempty" xml:"data>classification,omitempty"`
	Licenses             *[]LicenseChoice      `json:"licenses,omitempty" xml:"licenses>license,omitempty"`
	ExternalReferences   *[]ExternalReference  `json:"externalReferences,omitempty" xml:"externalReferences>reference,omitempty"`
	Services             *[]Service            `json:"services,omitempty" xml:"services>service,omitempty"`
	BOMRef               string                `json:"bom-ref,omitempty" xml:"bom-ref,attr,omitempty"`
}

type Source struct {
	Name string `json:"name,omitempty" xml:"name,omitempty"`
	URL  string `json:"url,omitempty" xml:"url,omitempty"`
}

type SWID struct {
	Text       *AttachedText `json:"text,omitempty" xml:"text,omitempty"`
	URL        string        `json:"url,omitempty" xml:"url,attr,omitempty"`
	TagID      string        `json:"tagId" xml:"tagId,attr"`
	Name       string        `json:"name" xml:"name,attr"`
	Version    string        `json:"version,omitempty" xml:"version,attr,omitempty"`
	TagVersion *int          `json:"tagVersion,omitempty" xml:"tagVersion,attr,omitempty"`
	Patch      *bool         `json:"patch,omitempty" xml:"patch,attr,omitempty"`
}

type Tool struct {
	Vendor  string  `json:"vendor,omitempty" xml:"vendor,omitempty"`
	Name    string  `json:"name" xml:"name"`
	Version string  `json:"version,omitempty" xml:"version,omitempty"`
	Hashes  *[]Hash `json:"hashes,omitempty" xml:"hashes>hash,omitempty"`
}
