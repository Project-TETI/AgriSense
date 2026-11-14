package inspection

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var codePattern = regexp.MustCompile(`^[a-z0-9]+(_[a-z0-9]+)*$`)

type Prediction struct {
	Plant      string   `json:"plant"`
	ClassName  string   `json:"class_name"`
	Confidence *float64 `json:"confidence_score"`
	Status     string   `json:"status"`
}

type MLError struct {
	Status int
	Code   string
}

func (e *MLError) Error() string { return e.Code }

type MLClient struct {
	url    string
	client *http.Client
}

func NewMLClient(base string) (*MLClient, error) {
	u, err := url.Parse(base)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Hostname() == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || (u.Path != "" && u.Path != "/") {
		return nil, errors.New("invalid ML service URL")
	}
	return &MLClient{url: strings.TrimRight(base, "/") + "/api/v1/ml/predict", client: &http.Client{Timeout: 15 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}}, nil
}

func (m *MLClient) Predict(ctx context.Context, plant, contentType string, image []byte) (Prediction, error) {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", `form-data; name="file"; filename="leaf"`)
	header.Set("Content-Type", contentType)
	part, err := w.CreatePart(header)
	if err != nil {
		return Prediction{}, err
	}
	if _, err := part.Write(image); err != nil {
		return Prediction{}, err
	}
	if err := w.WriteField("plant", plant); err != nil {
		return Prediction{}, err
	}
	if err := w.Close(); err != nil {
		return Prediction{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, m.url, &body)
	if err != nil {
		return Prediction{}, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := m.client.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return Prediction{}, &MLError{http.StatusGatewayTimeout, "ML_TIMEOUT"}
		}
		return Prediction{}, &MLError{http.StatusServiceUnavailable, "ML_UNAVAILABLE"}
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 65537))
	if err != nil || len(data) > 65536 {
		return Prediction{}, &MLError{http.StatusBadGateway, "ML_INVALID_RESPONSE"}
	}
	if resp.StatusCode != http.StatusOK {
		var failure struct {
			Detail struct {
				Code string `json:"error_code"`
			} `json:"detail"`
		}
		if resp.StatusCode == http.StatusUnprocessableEntity && json.Unmarshal(data, &failure) == nil && failure.Detail.Code == "BLURRY_IMAGE" {
			return Prediction{}, &MLError{http.StatusUnprocessableEntity, "BLURRY_IMAGE"}
		}
		if resp.StatusCode == http.StatusBadRequest {
			return Prediction{}, &MLError{http.StatusBadRequest, "ML_INPUT_REJECTED"}
		}
		return Prediction{}, &MLError{http.StatusServiceUnavailable, "ML_UNAVAILABLE"}
	}
	var prediction Prediction
	if err := json.Unmarshal(data, &prediction); err != nil || prediction.Plant != plant || prediction.Confidence == nil || *prediction.Confidence < 0 || *prediction.Confidence > 1 ||
		len(prediction.ClassName) > 150 || !codePattern.MatchString(prediction.ClassName) ||
		(prediction.Status != "SUCCESS" && prediction.Status != "UNCERTAIN") ||
		(prediction.Status == "UNCERTAIN" && prediction.ClassName != "unknown") || (prediction.Status == "SUCCESS" && prediction.ClassName == "unknown") {
		return Prediction{}, &MLError{http.StatusBadGateway, "ML_INVALID_RESPONSE"}
	}
	return prediction, nil
}
