package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// універсальна функція для виклику будь-якого ендпоінта
func callAPI(endpoint string, payload interface{}) (map[string]interface{}, error) {
	url := "http://localhost:8080" + endpoint
	body, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}

func callAPIArray(endpoint string, payload interface{}) ([]map[string]interface{}, error) {
	url := "http://localhost:8080" + endpoint
	body, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func callAPIStringArray(endpoint string, payload interface{}) ([]string, error) {
	url := "http://localhost:8080" + endpoint
	body, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result []string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// ---------- Клієнтські методи ----------

// /create_collection
func CreateCollection(collection string) {
	res, err := callAPI("/create_collection", map[string]interface{}{"collection_name": collection})
	fmt.Println("CREATE_COLLECTION:", res, "ERR:", err)
}

// /put_document
func PutDocument(collection string, doc map[string]interface{}) {
	res, err := callAPI("/put_document", map[string]interface{}{"collection_name": collection, "document": doc})
	fmt.Println("PUT:", res, "ERR:", err)
}

// /get_document
func GetDocument(collection, userID string) {
	res, err := callAPI("/get_document", map[string]interface{}{"collection_name": collection, "user_id": userID})
	fmt.Println("GET:", res, "ERR:", err)
}

// /list_documents
func ListDocuments(collection string) {
	//url := "http://localhost:8080
	payload := map[string]interface{}{"collection_name": collection}
	res, err := callAPIArray("/list_documents", payload)
	fmt.Println("LIST:", res, "ERR:", err)
}

// /delete_document
func DeleteDocument(collection, userID string) {
	res, err := callAPI("/delete_document", map[string]interface{}{"collection_name": collection, "user_id": userID})
	fmt.Println("DELETE:", res, "ERR:", err)
}

// /delete_collection
func DeleteCollection(collection string) {
	res, err := callAPI("/delete_collection", map[string]interface{}{"collection_name": collection})
	fmt.Println("DELETE_COLLECTION:", res, "ERR:", err)
}

// /create_index
func CreateIndex(collection, field string) {
	res, err := callAPI("/create_index", map[string]interface{}{"collection_name": collection, "field": field})
	fmt.Println("CREATE_INDEX:", res, "ERR:", err)
}

// /delete_index
func DeleteIndex(collection, indexName string) {
	res, err := callAPI("/delete_index", map[string]interface{}{"collection_name": collection, "index_name": indexName})
	fmt.Println("DELETE_INDEX:", res, "ERR:", err)
}

// /list_collections
func ListCollections() {
	res, err := callAPIStringArray("/list_collections", map[string]interface{}{})
	fmt.Println("LIST_COLLECTIONS:", res, "ERR:", err)
}
