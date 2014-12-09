package mintpal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func get(url string, pointer interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = json.Unmarshal(b, pointer)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func debug(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}
