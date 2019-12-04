package dcdump

// DOIResponse is the https://api.datacite.org/dois endpoint response.
// TODO(martin): Sort out the interface{} fields, if necessary.
type DOIResponse struct {
	Data []struct {
		Attributes    interface{} `json:"attributes"`
		Id            string      `json:"id"`
		Relationships struct {
			Client struct {
				Data struct {
					Id   string `json:"id"`
					Type string `json:"type"`
				} `json:"data"`
			} `json:"client"`
		} `json:"relationships"`
		Type string `json:"type"`
	} `json:"data"`
	Included []struct {
		Attributes struct {
			AlternateName interface{}   `json:"alternateName"`
			ClientType    string        `json:"clientType"`
			ContactEmail  string        `json:"contactEmail"`
			Created       string        `json:"created"`
			Description   interface{}   `json:"description"`
			Domains       string        `json:"domains"`
			HasPassword   bool          `json:"hasPassword"`
			IsActive      bool          `json:"isActive"`
			Issn          interface{}   `json:"issn"`
			Language      []interface{} `json:"language"`
			Name          string        `json:"name"`
			Opendoar      interface{}   `json:"opendoar"`
			Re3data       interface{}   `json:"re3data"`
			Symbol        string        `json:"symbol"`
			Updated       string        `json:"updated"`
			Url           interface{}   `json:"url"`
			Year          int64         `json:"year"`
		} `json:"attributes"`
		Id            string `json:"id"`
		Relationships struct {
			Prefixes struct {
				Data []struct {
					Id   string `json:"id"`
					Type string `json:"type"`
				} `json:"data"`
			} `json:"prefixes"`
			Provider struct {
				Data struct {
					Id   string `json:"id"`
					Type string `json:"type"`
				} `json:"data"`
			} `json:"provider"`
		} `json:"relationships"`
		Type string `json:"type"`
	} `json:"included"`
	Links struct {
		Next string `json:"next"`
		Self string `json:"self"`
	} `json:"links"`
	Meta struct {
		Affiliations []struct {
			Count int64  `json:"count"`
			Id    string `json:"id"`
			Title string `json:"title"`
		} `json:"affiliations"`
		Certificates []interface{} `json:"certificates"`
		Clients      []struct {
			Count int64  `json:"count"`
			Id    string `json:"id"`
			Title string `json:"title"`
		} `json:"clients"`
		Created []struct {
			Count int64  `json:"count"`
			Id    string `json:"id"`
			Title string `json:"title"`
		} `json:"created"`
		LinkChecksCitationDoi  int64 `json:"linkChecksCitationDoi"`
		LinkChecksDcIdentifier int64 `json:"linkChecksDcIdentifier"`
		LinkChecksSchemaOrgId  int64 `json:"linkChecksSchemaOrgId"`
		LinkChecksStatus       []struct {
			Count int64  `json:"count"`
			Id    string `json:"id"`
			Title string `json:"title"`
		} `json:"linkChecksStatus"`
		LinksChecked       int64 `json:"linksChecked"`
		LinksWithSchemaOrg []struct {
			Count int64  `json:"count"`
			Id    string `json:"id"`
			Title string `json:"title"`
		} `json:"linksWithSchemaOrg"`
		Page     int64 `json:"page"`
		Prefixes []struct {
			Count int64  `json:"count"`
			Id    string `json:"id"`
			Title string `json:"title"`
		} `json:"prefixes"`
		Providers []struct {
			Count int64  `json:"count"`
			Id    string `json:"id"`
			Title string `json:"title"`
		} `json:"providers"`
		Registered []struct {
			Count int64  `json:"count"`
			Id    string `json:"id"`
			Title string `json:"title"`
		} `json:"registered"`
		ResourceTypes []struct {
			Count int64  `json:"count"`
			Id    string `json:"id"`
			Title string `json:"title"`
		} `json:"resourceTypes"`
		SchemaVersions []struct {
			Count int64  `json:"count"`
			Id    string `json:"id"`
			Title string `json:"title"`
		} `json:"schemaVersions"`
		Sources []struct {
			Count int64  `json:"count"`
			Id    string `json:"id"`
			Title string `json:"title"`
		} `json:"sources"`
		States []struct {
			Count int64  `json:"count"`
			Id    string `json:"id"`
			Title string `json:"title"`
		} `json:"states"`
		Subjects []struct {
			Count int64  `json:"count"`
			Id    string `json:"id"`
			Title string `json:"title"`
		} `json:"subjects"`
		Total      int64 `json:"total"`
		TotalPages int64 `json:"totalPages"`
	} `json:"meta"`
}
