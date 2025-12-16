package gorthanc

import (
	"fmt"

	"github.com/proencaj/gorthanc/types"
)

// Modality represents a DICOM modality configuration
type Modality struct {
	// Orthanc Name or ID for this Modality
	Name string `json:"Name,omitempty"`

	// Application Entity Title (AET) of the remote modality
	AET string `json:"AET"`

	// Host/IP address of the remote modality
	Host string `json:"Host"`

	// Port number of the remote modality
	Port int `json:"Port"`

	// Manufacturer name (optional)
	Manufacturer string `json:"Manufacturer,omitempty"`

	// Whether to allow echo requests
	AllowEcho *bool `json:"AllowEcho,omitempty"`

	// Whether to allow C-FIND requests
	AllowFind *bool `json:"AllowFind,omitempty"`

	// Whether to allow C-GET requests
	AllowGet *bool `json:"AllowGet,omitempty"`

	// Whether to allow C-MOVE requests
	AllowMove *bool `json:"AllowMove,omitempty"`

	// Whether to allow C-STORE requests
	AllowStore *bool `json:"AllowStore,omitempty"`

	// Timeout for DICOM operations (in seconds)
	Timeout int `json:"Timeout,omitempty"`

	// Internal client for making requests to modality
	client *Client
}

func (m *Modality) String() string {
	return fmt.Sprintf("<Modality | Name: '%s', AET: %s (%s:%d)>", m.Name, m.AET, m.Host, m.Port)
}

func (c *Client) GetModalityNames() ([]string, error) {
	var modalities []string
	if err := c.get("modalities", &modalities); err != nil {
		return nil, err
	}
	return modalities, nil
}

func (c *Client) GetModalities() ([]*Modality, error) {
	var modalities map[string]Modality
	if err := c.get("modalities?expand", &modalities); err != nil {
		return nil, err
	}
	var modalityList []*Modality

	for name, details := range modalities {
		details.Name = name
		details.client = c
		modalityList = append(modalityList, &details)
	}
	return modalityList, nil
}

func (c *Client) GetModality(modalityName string) (*Modality, error) {
	var modality Modality
	path := fmt.Sprintf("modalities/%s/configuration", modalityName)

	if err := c.get(path, &modality); err != nil {
		return nil, err
	}
	modality.Name = modalityName
	modality.client = c

	return &modality, nil
}

func (c *Client) NewModality(modalityName string, request *types.ModalityCreateRequest) (*Modality, error) {
	if err := c.CreateOrUpdateModality(modalityName, request); err != nil {
		return nil, err
	}

	return c.GetModality(modalityName)
}

func (m *Modality) Update(request *types.ModalityCreateRequest) error {
	return m.client.CreateOrUpdateModality(m.Name, request)
}

func (c *Client) CreateOrUpdateModality(modalityName string, request *types.ModalityCreateRequest) error {
	path := fmt.Sprintf("modalities/%s", modalityName)

	// Convert the request to the format expected by Orthanc
	// Orthanc expects an array: [AET, Host, Port, Manufacturer]
	modalityArray := []interface{}{
		request.AET,
		request.Host,
		request.Port,
	}

	// Add manufacturer if provided
	if request.Manufacturer != "" {
		modalityArray = append(modalityArray, request.Manufacturer)
	}

	if err := c.put(path, modalityArray, nil); err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteModality(modalityName string) error {
	path := fmt.Sprintf("modalities/%s", modalityName)

	if err := c.delete(path, nil); err != nil {
		return err
	}

	return nil
}

func (m *Modality) Echo() error {
	return m.client.EchoModality(m.Name)
}

func (c *Client) EchoModality(modalityName string) error {
	path := fmt.Sprintf("modalities/%s/echo", modalityName)

	if err := c.post(path, nil, nil); err != nil {
		return err
	}

	return nil
}

func (m *Modality) Store(resourceID string) error {
	return m.client.StoreToModality(m.Name, resourceID)
}

func (c *Client) StoreToModality(modalityName, resourceID string) error {
	path := fmt.Sprintf("modalities/%s/store", modalityName)

	if err := c.post(path, resourceID, nil); err != nil {
		return err
	}

	return nil
}

func (m *Modality) StoreWithOptions(request *types.ModalityStoreRequest) (*types.ModalityStoreResult, error) {
	return m.client.StoreToModalityWithOptions(m.Name, request)
}

func (c *Client) StoreToModalityWithOptions(modalityName string, request *types.ModalityStoreRequest) (*types.ModalityStoreResult, error) {
	path := fmt.Sprintf("modalities/%s/store", modalityName)

	var result types.ModalityStoreResult
	if err := c.post(path, request, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (m *Modality) Query(request *types.ModalityFindRequest) (string, error) {
	return m.client.QueryModality(m.Name, request)
}

func (c *Client) QueryModality(modalityName string, request *types.ModalityFindRequest) (string, error) {
	path := fmt.Sprintf("modalities/%s/query", modalityName)

	var queryResponse map[string]interface{}
	if err := c.post(path, request, &queryResponse); err != nil {
		return "", err
	}

	// Extract the query ID from the response
	queryID, ok := queryResponse["ID"].(string)
	if !ok {
		return "", fmt.Errorf("failed to get query ID from response")
	}
	return queryID, nil
}

func (m *Modality) Find(request *types.ModalityFindRequest, answers interface{}) error {
	return m.client.FindInModality(m.Name, request, answers)
}

func (c *Client) FindInModality(modalityName string, request *types.ModalityFindRequest, answers interface{}) error {
	queryID, err := c.QueryModality(modalityName, request)
	if err != nil {
		return err
	}

	// TODO: Turn 'expand' and 'simplify' into new parameters or options, so the user can change based on what is needs
	answersPath := fmt.Sprintf("queries/%s/answers?expand=true&simplify=true", queryID)

	// var answers []map[string]interface{}
	if err := c.get(answersPath, &answers); err != nil {
		return err
	}

	return nil
}

func (m *Modality) Move(request *types.ModalityMoveRequest) (*types.ModalityMoveResult, error) {
	return m.client.MoveFromModality(m.Name, request)
}

func (c *Client) MoveFromModality(modalityName string, request *types.ModalityMoveRequest) (*types.ModalityMoveResult, error) {
	path := fmt.Sprintf("modalities/%s/move", modalityName)

	var result types.ModalityMoveResult
	if err := c.post(path, request, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (m *Modality) Get(request *types.ModalityGetRequest) error {
	return m.client.GetFromModality(m.Name, request)
}

func (c *Client) GetFromModality(modalityName string, request *types.ModalityGetRequest) error {
	path := fmt.Sprintf("modalities/%s/get", modalityName)

	if err := c.post(path, request, nil); err != nil {
		return err
	}

	return nil
}
