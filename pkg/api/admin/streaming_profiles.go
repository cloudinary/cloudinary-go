package admin

import (
	"cloudinary-labs/cloudinary-go/pkg/api"
	"cloudinary-labs/cloudinary-go/pkg/transformation"
	"context"
)

const (
	StreamingProfiles api.EndPoint = "streaming_profiles"
)

func (a *Api) ListStreamingProfiles(ctx context.Context) (*ListStreamingProfilesResult, error) {
	res := &ListStreamingProfilesResult{}
	_, err := a.get(ctx, StreamingProfiles, nil, res)

	return res, err
}

type ListStreamingProfilesResult struct {
	Data  []StreamingProfile `json:"data"`
	Error api.ErrorResp      `json:"error,omitempty"`
}

type StreamingProfile struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Predefined  bool   `json:"predefined"`
}

type GetStreamingProfileParams struct {
	Name string `json:"-"`
}

func (a *Api) GetStreamingProfile(ctx context.Context, params GetStreamingProfileParams) (*GetStreamingProfileResult, error) {
	res := &GetStreamingProfileResult{}
	_, err := a.get(ctx, api.BuildPath(StreamingProfiles, params.Name), params, res)

	return res, err
}

type GetStreamingProfileResult struct {
	Data  StreamingProfileDetails `json:"data"`
	Error api.ErrorResp           `json:"error,omitempty"`
}

type StreamingProfileDetails struct {
	Name            string                           `json:"name"`
	DisplayName     string                           `json:"display_name,omitempty"`
	Predefined      bool                             `json:"predefined"`
	Representations []StreamingProfileRepresentation `json:"representations,omitempty"`
}

type StreamingProfileRepresentation struct {
	Transformation transformation.Transformation `json:"transformation"`
}

type RawStreamingProfileRepresentation struct {
	Transformation transformation.RawTransformation `json:"transformation"`
}

type CreateStreamingProfileParams struct {
	Name            string                              `json:"name"`
	DisplayName     string                              `json:"display_name,omitempty"`
	Representations []RawStreamingProfileRepresentation `json:"representations"`
}

func (a *Api) CreateStreamingProfile(ctx context.Context, params CreateStreamingProfileParams) (*GetStreamingProfileResult, error) {
	res := &GetStreamingProfileResult{}
	_, err := a.post(ctx, api.BuildPath(StreamingProfiles), params, res)

	return res, err
}

type UpdateStreamingProfileParams struct {
	Name            string                              `json:"-"`
	DisplayName     string                              `json:"display_name,omitempty"`
	Representations []RawStreamingProfileRepresentation `json:"representations"`
}

func (a *Api) UpdateStreamingProfile(ctx context.Context, params UpdateStreamingProfileParams) (*GetStreamingProfileResult, error) {
	res := &GetStreamingProfileResult{}
	_, err := a.put(ctx, api.BuildPath(StreamingProfiles, params.Name), params, res)

	return res, err
}

type DeleteStreamingProfileParams struct {
	Name string `json:"-"`
}

func (a *Api) DeleteStreamingProfile(ctx context.Context, params DeleteStreamingProfileParams) (*DeleteStreamingProfileResult, error) {
	res := &DeleteStreamingProfileResult{}
	_, err := a.delete(ctx, api.BuildPath(StreamingProfiles, params.Name), params, res)

	return res, err
}

type DeleteStreamingProfileResult struct {
	Message string        `json:"message"`
	Error   api.ErrorResp `json:"error,omitempty"`
}
