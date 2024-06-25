### Description

`imggen` is a CLI tool for generating images using OpenAI API. You should generate the API key from OpenAI or other platforms and configure it to the environment variable before you use. I don’t provide API key.

### Usage


```sh
imggen [options] [prompt]
```

### Options
- `model`  The model to use for image generation. (default `dall-e-2`)

    - `dall-e-2`
    - `dall-e-3`

- `quality`  Quality of the image to generate. This flag is only supported for model `dall-e-3`. (default `standard`)

    - `standard`
    - `hd`

- `size`  Size of the image to generate. (default `1024x1024`)

    Options for `dall-e-2`:

    - `256x256`
    - `512x512`
    - `1024x1024`

    Options for `dall-e-3`: 

    - `1024x1024`
    - `1792x1024`
    - `1024x1792`    	 

- `style`  Style of the image to generate. This flag is only supported for model `dall-e-3`. (default `vivid`)

    - `vivid`
    - `natural`

- `output`  Output format. (default `list`)

    - `list`
    - `json`


### Environment Variables

  - `IMGGEN_API_KEY`: The API key to use for image generation.
  - `IMGGEN_API_ENDPOINT`: (Optional) The url to send request. Default https://api.openai.com/v1

### Examples

```shell
# This command will show the response in formatted JSON format
imggen --model dall-e-3 --output json "A cute baby sea otter wearing a beret" | jq '.'

# This command will download the image to ./img.png
imggen --model dall-e-3 --output json "A cute baby sea otter wearing a beret" | jq '.data[0].url' | xargs wget -q -O img.png
```



