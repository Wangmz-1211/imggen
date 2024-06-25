package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

var supported_models = []string{"dall-e-2", "dall-e-3"}
var supported_sizes = map[string][]string{
    "dall-e-2": {"256x256", "512x512", "1024x1024"},
    "dall-e-3": {"1024x1024", "1792x1024", "1024x1792"},
}
var supported_styles = map[string][]string{
    "dall-e-2": {"vivid"},
    "dall-e-3": {"vivid", "natural"},
}
var supported_qualities = map[string][]string{
    "dall-e-2": {"standard"},
    "dall-e-3": {"standard", "hd"},
}
var supported_outputs = []string{"list", "json"}

type ImageGenerationResponse struct {
    Created int `json:"created"`
    Data []ImageGenerationData `json:"data"`
}

type ImageGenerationData struct {
    Url string `json:"url"`
    Revised_prompt string `json:"revised_prompt"`
    B64_json string `json:"b64_json"`
}

type ImageGenerationErrorResponse struct {
    Error ImageGenerationErrorData `json:"error"`
}

type ImageGenerationErrorData struct {
    Code string `json:"code"`
    Message string `json:"message"`
    Param string `json:"param"`
    Type string `json:"type"`
}

func main() {
    // flags
    model := flag.String("model", "dall-e-2", "The model to use for image generation.\n Options: 'dall-e-2', 'dall-e-3'\n")
    size := flag.String("size", "1024x1024", "Size of the image to generate.\n Options for dall-e-2: '256x256', '512x512', '1024x1024'\n Options for dall-e-3: '1024x1024', '1792x1024', '1024x1792'\n")
    style := flag.String("style", "vivid", "Style of the image to generate. This flag is only supported for model 'dall-e-3'.\n Options: 'vivid', 'natural'\n")
    quality := flag.String("quality", "standard", "Quality of the image to generate. This flag is only supported for model 'dall-e-3'.\n Options: 'standard', 'hd'\n")
    output := flag.String("output", "list", "Output format.\n Options: 'list', 'json'\n")
    // set usage information
    flag.Usage = func() {
        w := flag.CommandLine.Output()
        fmt.Fprintf(w, "Usage of imggen:\n\n\t%s [options] [prompt]\n\nOptions:\n", os.Args[0])
        flag.PrintDefaults()
        fmt.Fprintf(w, "\nEnvironment Variables:\n\n  IMGGEN_API_KEY: The API key to use for image generation.\n\n  IMGGEN_API_ENDPOINT: (Optional) The url to send request. Default to https://api.openai.com/v1\n")
    }
    flag.Parse()
    prompt := flag.Arg(0)
    // check prompt
    if prompt == "" {
        error_information("Prompt is required.")
    }
    // check Options
    check_option(*model, *size, *style, *quality, output)
    // environment variables
    api_key, ok := os.LookupEnv("IMGGEN_API_KEY")
    if !ok {
        error_information("IMGGEN_API_KEY was not set.")
    }
    api_endpoint, ok := os.LookupEnv("IMGGEN_API_ENDPOINT")
    if !ok {
        api_endpoint = "https://api.openai.com/v1"
    }

    // set request body
    requestBody := build_request_body(*model, prompt, *size, *style, *quality)
    req, err := http.NewRequest("POST", api_endpoint + "/images/generations", requestBody)
    if err != nil {
        error_information("Failed to create a request.")
    }
    // set request headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer " + api_key)

    // send request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        error_information("Failed to send a request. Maybe the IMGGEN_API_ENDPOINT is invalid.")
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        error_information("Failed to read response.")
    }

    // parse response
    if resp.StatusCode != 200 {
        var imageGenerationErrorResponse ImageGenerationErrorResponse
        err = json.Unmarshal(body, &imageGenerationErrorResponse)
        if err != nil {
            error_information("Failed to parse error response, body is as follows:\n" + string(body))
        }
        if *output == "json" {
            fmt.Println(string(body))
            os.Exit(1)
        }
        error_information(imageGenerationErrorResponse.Error.Message)
    }

    var imageGenerationResponse ImageGenerationResponse
    err = json.Unmarshal(body, &imageGenerationResponse)
    if err != nil {
        error_information("Failed to parse response, body is as follows:\n" + string(body))
    }

    // output 
    switch *output {
    case "list":
        fmt.Println("Image URL:\n " + imageGenerationResponse.Data[0].Url)
        fmt.Println("\nRevised Prompt:\n " + imageGenerationResponse.Data[0].Revised_prompt)
    case "json":
        fmt.Println(string(body))
    default:
        fmt.Println("Image URL:\n " + imageGenerationResponse.Data[0].Url)
        fmt.Println("\nRevised Prompt:\n " + imageGenerationResponse.Data[0].Revised_prompt)
    }
}

func build_request_body(model, prompt, size, style, quality string) *bytes.Buffer {
    if model == "dall-e-3" {
        return bytes.NewBufferString(
            fmt.Sprintf(`{"model": "%s", "prompt": "%s", "size": "%s", "style": "%s", "quality": "%s"}`, model, prompt, size, style, quality),
        )
    } else {
        return bytes.NewBufferString(
            fmt.Sprintf(`{"model": "%s", "prompt": "%s", "size": "%s"}`, model, prompt, size),
        )
    }
}

/**
 * Check if the options are supported.
 */
func check_option(model string, size string, style string, quality string, output *string) {
    if !contains(supported_models, model) {
        error_information(fmt.Sprintf("Model '%s' is not supported.", model))
    }
    if !contains(supported_sizes[model], size) {
        error_information(fmt.Sprintf("Size '%s' is not supported for model '%s'.", size, model))
    }
    if !contains(supported_styles[model], style) {
        error_information(fmt.Sprintf("Style '%s' is not supported for model '%s'.", style, model))
    }
    if !contains(supported_qualities[model], quality) {
        error_information(fmt.Sprintf("Quality '%s' is not supported for model '%s'.", quality, model))
    }
    if !contains(supported_outputs, *output) {
        fmt.Printf("Output '%s' is not supported. Default to 'list'.\n", *output)
    }
}


/**
 * Check if the string is in the array.
 */
func contains(arr []string, s string) bool {
    for _, a := range arr {
        if a == s {
            return true
        }
    }
    return false
}

func error_information(s string) {
    flag.Usage()
    fmt.Printf("\nError:\n\n  %s\n", s)
    os.Exit(1)
}
