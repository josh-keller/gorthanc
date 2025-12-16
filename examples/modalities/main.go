package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/proencaj/gorthanc"
	"github.com/proencaj/gorthanc/types"
)

type Count int

func (c *Count) UnmarshalJSON(data []byte) error {
	if n, err := strconv.Atoi(strings.Trim(string(data), "\"")); err != nil {
		return err
	} else {
		*c = Count(n)
		return nil
	}
}

func main() {
	client, err := gorthanc.NewClient(
		"http://localhost:8042",
		gorthanc.WithBasicAuth("orthanc", "orthanc"),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Example: GetModalities
	modalities, err := client.GetModalities()
	if err != nil {
		log.Fatalf("Failed to get modalities: %v", err)
	}
	fmt.Println(modalities)

	// // Example: GetModalityDetails
	// if len(modalities) > 0 {
	// 	modality, err := client.GetModalityDetails(modalities[0])
	// 	if err != nil {
	// 		log.Fatalf("Failed to get modality details: %v", err)
	// 	}
	//
	// 	jsonData, _ := json.MarshalIndent(modality, "", "  ")
	// 	fmt.Println(string(jsonData))
	// }

	// Example: EchoModality
	for _, mod := range modalities {
		if err = mod.Echo(); err != nil {
			fmt.Printf("Modality '%s' not reachable: %v\n", mod.Name, err)
		} else {
			fmt.Printf("Modality '%s' is reachable!\n", mod.Name)
		}
	}

	// Example: CreateOrUpdateModality
	request := &types.ModalityCreateRequest{
		AET:  "TEST_MODALITY",
		Host: "localhost",
		Port: 4242,
	}
	testModality, err := client.NewModality("TestModality", request)
	// err = client.CreateOrUpdateModality("TEST_MODALITY", request)
	if err != nil {
		log.Fatalf("Failed to create modality: %v", err)
	}
	fmt.Println(testModality)

	if err = client.DeleteModality(testModality.Name); err != nil {
		log.Fatalf("Problem deleting modality: %v", err)
	}
	// Example: FindInModality
	findRequest := &types.ModalityFindRequest{
		Level:     "Study",
		Normalize: gorthanc.BoolPtr(true),
		Query: map[string]string{
			"PatientName": "*",
			"0020,1206":   "", // Series in study
			"0020,1208":   "", // Instances in study
		},
	}
	modalities, err = client.GetModalities()
	if err != nil {
		log.Fatalf("Failed to get modalities: %v", err)
	}

	var findResults []struct {
		PatientName   string
		MRN           string `json:"PatientID"`
		StudyUID      string `json:"StudyInstanceUID"`
		SeriesCount   Count  `json:"NumberOfStudyRelatedSeries"`
		InstanceCount Count  `json:"NumberOfStudyRelatedInstances"`
	}
	// var findResults interface{}
	err = modalities[0].Find(findRequest, &findResults)
	if err != nil {
		log.Fatalf("Failed to query modality: %v", err)
	}
	jsonData, _ := json.MarshalIndent(findResults, "", "  ")
	fmt.Println(string(jsonData))
	fmt.Printf("%+v\n", findResults)
}

//
// 	// Example: MoveFromModality
// 	moveRequest := &types.ModalityMoveRequest{
// 		Level:        "Study",
// 		TargetAet:    "ORTHANC",
// 		Asynchronous: gorthanc.BoolPtr(false),
// 		Permissive:   gorthanc.BoolPtr(false),
// 		Timeout:      30,
// 		Resources: []map[string]interface{}{
// 			{"StudyInstanceUID": "1.2.840.113619.2.55.3.123456789"},
// 		},
// 	}
// 	moveResults, err := client.MoveFromModality("PACS", moveRequest)
// 	if err != nil {
// 		log.Fatalf("Failed to move study: %v", err)
// 	}
// 	jsonData, _ = json.MarshalIndent(moveResults, "", "  ")
// 	fmt.Println(string(jsonData))
//
// 	// Example: GetFromModality
// 	getRequest := &types.ModalityGetRequest{
// 		Level:        "Study",
// 		Asynchronous: gorthanc.BoolPtr(false),
// 		Permissive:   gorthanc.BoolPtr(false),
// 		Timeout:      30,
// 		Resources: []map[string]interface{}{
// 			{"StudyInstanceUID": "1.2.840.113619.2.55.3.123456789"},
// 		},
// 	}
// 	err = client.GetFromModality("PACS", getRequest)
// 	if err != nil {
// 		log.Fatalf("Failed to get study: %v", err)
// 	}
//
// 	// Example: StoreToModality
// 	err = client.StoreToModality("PACS", "study-orthanc-id")
// 	if err != nil {
// 		log.Fatalf("Failed to store study: %v", err)
// 	}
//
// 	// Example: DeleteModality
// 	err = client.DeleteModality("TEST_MODALITY")
// 	if err != nil {
// 		log.Fatalf("Failed to delete modality: %v", err)
// 	}
// }
