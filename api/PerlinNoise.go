package api

import (
	"github.com/aquilax/go-perlin"
	"github.com/swaggest/openapi-go"
	"hash/fnv"
	"httpServer/validation"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"net/http"
	"strconv"
)

type PerlinNoiseRequest struct {
	Width     int     `query:"width" description:"Width of perlin image" example:"512" default:"512" required:"false"`
	Height    int     `query:"height" description:"Height of perlin image" example:"512" default:"512" required:"false"`
	Alpha     float64 `query:"alpha" description:"Alpha of perlin image" example:"2" default:"2" required:"false"`
	Beta      float64 `query:"beta" description:"Beta of perlin image" example:"2" default:"2" required:"false"`
	ScaleX    float64 `query:"scalex" description:"The x scale of perlin image" example:"5" default:"5" required:"false"`
	ScaleY    float64 `query:"scaley" description:"The y scale of perlin image" example:"5" default:"5" required:"false"`
	Iteration int     `query:"n" description:"Iteration of perlin image" example:"5" default:"5" required:"false"`
	Seed      string  `query:"seed" description:"Seed of perlin image. Empty for random" example:"abc123" required:"false"`
}

func RoutePerlinNoise(path string, builder *RouteBuilder) {
	builder.Mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		widthStr := r.URL.Query().Get("width")
		heightStr := r.URL.Query().Get("height")
		alphaStr := r.URL.Query().Get("alpha")
		betaStr := r.URL.Query().Get("beta")
		scalexStr := r.URL.Query().Get("scalex")
		scaleyStr := r.URL.Query().Get("scaley")
		iterationStr := r.URL.Query().Get("n")
		seedStr := r.URL.Query().Get("seed")

		width := 512
		height := 512
		alpha := 2.0
		beta := 2.0
		scalex := 5.0
		scaley := 5.0
		iteration := int32(5)
		var seed int64

		if widthStr != "" {
			var err error
			width, err = strconv.Atoi(widthStr)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("Invalid width parameter"))
				return
			}
		}
		if heightStr != "" {
			var err error
			height, err = strconv.Atoi(heightStr)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("Invalid height parameter"))
				return
			}
		}
		if alphaStr != "" {
			var err error
			alpha, err = strconv.ParseFloat(alphaStr, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("Invalid alpha parameter"))
				return
			}
		}
		if betaStr != "" {
			var err error
			beta, err = strconv.ParseFloat(betaStr, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("Invalid beta parameter"))
				return
			}
		}
		if scalexStr != "" {
			var err error
			scalex, err = strconv.ParseFloat(scalexStr, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("Invalid scalex parameter"))
				return
			}
		}
		if scaleyStr != "" {
			var err error
			scaley, err = strconv.ParseFloat(scaleyStr, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("Invalid scaley parameter"))
				return
			}
		}
		if iterationStr != "" {
			var err error
			iteration32, err := strconv.Atoi(iterationStr)
			iteration = int32(iteration32)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("Invalid iteration parameter"))
				return
			}
		}

		var errorsAggregate = InvalidArgumentBadRequestResponse{}
		errorsAggregate.Errors = make(map[string][]*validation.ValidateError)
		ok, errors := validation.Validate(int64(height), validation.DefaultValidateOptions,
			validation.Integer.NotLessThan(1),
			validation.Integer.NotGreaterThan(4096),
		)
		if !ok {
			errorsAggregate.Errors["height"] = errors
		}
		ok, errors = validation.Validate(int64(width), validation.DefaultValidateOptions,
			validation.Integer.NotLessThan(1),
			validation.Integer.NotGreaterThan(4096),
		)
		if !ok {
			errorsAggregate.Errors["width"] = errors
		}
		ok, errors = validation.Validate(int64(iteration), validation.DefaultValidateOptions,
			validation.Integer.NotLessThan(0),
			validation.Integer.NotGreaterThan(50),
		)
		if !ok {
			errorsAggregate.Errors["iteration"] = errors
		}

		if len(errorsAggregate.Errors) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			jsonByte, _ := errorsAggregate.ToJson()
			_, _ = w.Write(jsonByte)
			return
		}

		if seedStr == "" {
			seed = rand.Int63()
		} else {
			seed = stringToInt64(seedStr)
		}
		random := rand.NewSource(seed)
		img := image.NewGray(image.Rect(0, 0, width, height))
		noise := perlin.NewPerlinRandSource(alpha, beta, iteration, random)
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				n := noise.Noise2D(float64(x)/float64(width)*scalex, float64(y)/float64(height)*scaley)
				grayValue := uint8((n + 1) * 127.5)
				img.SetGray(x, y, color.Gray{Y: grayValue})
			}
		}

		w.Header().Set("Content-Type", "image/png")
		err := png.Encode(w, img)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Failed to encode image: " + err.Error()))
			return
		}
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
		return
	})
}

func ConfigurePerlinNoise(path string, builder *OpenApiBuilder) error {
	context, err := builder.OpenApiReflector.NewOperationContext(http.MethodPost, path)
	if err != nil {
		return err
	}
	context.SetDescription("Generates a Perlin noise image using GEGL https://gitlab.gnome.org/GNOME/gegl")
	context.SetTags("image")
	context.AddRespStructure(new(string), func(cu *openapi.ContentUnit) {
		cu.HTTPStatus = http.StatusOK
		cu.Description = "The Perlin noise image"
		cu.ContentType = "image/png"
		cu.IsDefault = true
	})
	context.AddRespStructure(new(string), func(cu *openapi.ContentUnit) {
		cu.HTTPStatus = http.StatusBadRequest
		cu.Description = "Invalid request parameter"
		cu.ContentType = "text/plain"
	})
	context.AddReqStructure(new(PerlinNoiseRequest), func(cu *openapi.ContentUnit) {
		cu.IsDefault = true
	})
	err = builder.OpenApiReflector.AddOperation(context)
	if err != nil {
		return err
	}
	format := "binary"
	// quite dirty but works
	builder.OpenApiReflector.Spec.Paths.MapOfPathItemValues[path].MapOfOperationValues["post"].Responses.Default.Response.Content["image/png"].Schema.Schema.Format = &format
	return nil
}
func stringToInt64(str string) int64 {
	h := fnv.New64a()
	h.Write([]byte(str))
	hashValue := h.Sum64()

	// 可选：将哈希值转为负数范围内的值
	return int64(hashValue & math.MaxInt64)
}
