package cli

import (
    "encoding/json"
    "fmt"
)

type JSONResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

func OutputJSONSuccess(data interface{}) error {
    resp := JSONResponse{
        Success: true,
        Data:    data,
    }
    return printJSON(resp)
}

func OutputJSONError(err error) error {
    resp := JSONResponse{
        Success: false,
        Error:   err.Error(),
    }
    return printJSON(resp)
}

func printJSON(v interface{}) error {
    bytes, err := json.MarshalIndent(v, "", "  ")
    if err != nil {
        return err
    }
    fmt.Println(string(bytes))
    return nil
}

// Unified handler for success + error
func Respond(jsonOutput bool, data interface{}, err error) error {
    if err != nil {
        if jsonOutput {
            _ = OutputJSONError(err)
            return nil
        }
        return err
    }

    if jsonOutput {
        return OutputJSONSuccess(data)
    }

    return nil
}