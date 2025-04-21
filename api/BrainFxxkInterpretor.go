package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/swaggest/openapi-go"
	"io"
	"net/http"
	"strconv"
)

type BrainFxxkRequest struct {
	Code    string `json:"code" description:"The code of the request" example:"+[.+]"`
	MemSize int    `json:"memSize" description:"The size of the memory in Byte" default:"8"`
	Memory  string `json:"memory" description:"The default memory set in base64. Leave empty for full zero" default:""`
	StdIn   string `json:"stdIn" description:"The input to the program in base64"`
}

type BrainFxxkResponse struct {
	StdOut string `json:"stdOut" description:"The output of the program in base64"`
	Memory string `json:"memory" description:"The memory after the program run in base64"`
}

type Ops rune

const (
	OpsPlus  Ops = '+'
	OpsMinus Ops = '-'
	OpsRight Ops = '>'
	OpsLeft  Ops = '<'
	OpsDot   Ops = '.'
	OpsComma Ops = ','
	OpsOpen  Ops = '['
	OpsClose Ops = ']'
)

func RouteBrainFxxkInterpretor(path string, builder *RouteBuilder) {
	builder.Mux.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodOptions {
			writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			writer.WriteHeader(http.StatusOK)
			return
		}
		if request.Method != http.MethodPost {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			writer.Write([]byte("Method not allowed"))
			return
		}

		body, err := io.ReadAll(request.Body)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				builder.ServiceProvider.Logger.Warning(err.Error())
			}
		}(request.Body)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("Error reading request body"))
			return
		}
		var req BrainFxxkRequest
		err = json.Unmarshal(body, &req)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("Error unmarshalling request body"))
			return
		}

		memory := make([]byte, req.MemSize)
		if req.Memory != "" {
			mem, err := base64.StdEncoding.DecodeString(req.Memory)
			if err != nil {
				writer.WriteHeader(http.StatusBadRequest)
				writer.Write([]byte("Error decoding memory"))
				return
			}
			if len(mem) > len(memory) {
				writer.WriteHeader(http.StatusBadRequest)
				writer.Write([]byte("Memory size is too large"))
				return
			}
			copy(memory, mem)
		}
		stdin, err := base64.StdEncoding.DecodeString(req.StdIn)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("Error decoding stdin"))
			return
		}
		stdout, err := interpret(req.Code, memory, stdin)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("Error interpreting code: " + err.Error()))
			return
		}
		response := BrainFxxkResponse{
			StdOut: base64.StdEncoding.EncodeToString(stdout),
			Memory: base64.StdEncoding.EncodeToString(memory),
		}

		responseBody, err := json.Marshal(response)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("Error marshalling response"))
			return
		}
		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		writer.Header().Set("Content-Length", strconv.Itoa(len(responseBody)))
		writer.Write(responseBody)
		if flusher, ok := writer.(http.Flusher); ok {
			flusher.Flush()
		}
	})
}

func interpret(code string, mem []byte, stdin []byte) (result []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	codeOffset := 0
	memOffset := 0
	stdinOffset := 0
	stdout := make([]byte, 0)
	for codeOffset < len(code) {
		op := code[codeOffset]
		switch Ops(op) {
		case OpsPlus:
			mem[memOffset]++
			break
		case OpsMinus:
			mem[memOffset]--
			break
		case OpsRight:
			memOffset++
			break
		case OpsLeft:
			memOffset--
			break
		case OpsDot:
			stdout = append(stdout, mem[memOffset])
			break
		case OpsComma:
			mem[memOffset] = stdin[stdinOffset]
			stdinOffset++
			break
		case OpsOpen:
			if mem[memOffset] == 0 {
				// skip forward to matching ]
				depth := 1
				for depth > 0 {
					codeOffset++
					switch Ops(code[codeOffset]) {
					case OpsOpen:
						depth++
					case OpsClose:
						depth--
					}
				}
			}
		case OpsClose:
			if mem[memOffset] != 0 {
				// jump back to matching [
				depth := 1
				for depth > 0 {
					codeOffset--
					switch Ops(code[codeOffset]) {
					case OpsOpen:
						depth--
					case OpsClose:
						depth++
					}
				}
			}
		}
		codeOffset++
	}
	return stdout, nil
}

func ConfigureBrainFxxkInterpretor(path string, builder *OpenApiBuilder) error {
	context, err := builder.OpenApiReflector.NewOperationContext(http.MethodPost, path)
	if err != nil {
		return err
	}
	context.AddReqStructure(new(BrainFxxkRequest), func(cu *openapi.ContentUnit) {
		cu.ContentType = "application/json"
		cu.Description = "BrainFxxk request"
	})
	context.AddRespStructure(new(BrainFxxkResponse), func(cu *openapi.ContentUnit) {
		cu.HTTPStatus = http.StatusOK
		cu.ContentType = "application/json"
	})
	context.AddRespStructure(new(string), func(cu *openapi.ContentUnit) {
		cu.HTTPStatus = http.StatusBadRequest
		cu.ContentType = "text/plain"
		cu.Description = "Error"
	})

	err = builder.OpenApiReflector.AddOperation(context)
	if err != nil {
		return err
	}
	return nil
}
