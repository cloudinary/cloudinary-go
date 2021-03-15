package admin

// Enables you to manage streaming profiles for use with adaptive bitrate streaming.
//
// https://cloudinary.com/documentation/admin_api#adaptive_streaming_profiles

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/cloudinary/cloudinary-go/api"
	"github.com/cloudinary/cloudinary-go/transformation"
)

const (
	StreamingProfiles api.EndPoint = "streaming_profiles"
)

// ListStreamingProfiles lists streaming profiles including built-in and custom profiles.
//
// https://cloudinary.com/documentation/admin_api#get_adaptive_streaming_profiles
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

// GetStreamingProfile gets details of a single streaming profile by name.
//
// https://cloudinary.com/documentation/admin_api#get_details_of_a_single_streaming_profile
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

type StreamingProfileRepresentations []RawStreamingProfileRepresentation

// Server expects to get a string of a json encoded array
func (spRepresentations StreamingProfileRepresentations) MarshalJSON() ([]byte, error) {
	spRepresentationsArray := ([]RawStreamingProfileRepresentation)(spRepresentations)
	paramsJsonObj, _ := json.Marshal(spRepresentationsArray)

	return []byte(strconv.Quote(string(paramsJsonObj))), nil
}

type RawStreamingProfileRepresentation struct {
	Transformation transformation.RawTransformation `json:"transformation"`
}

type CreateStreamingProfileParams struct {
	Name            string                          `json:"name"` // The name to assign to the new streaming profile.
	DisplayName     string                          `json:"display_name,omitempty"`
	Representations StreamingProfileRepresentations `json:"representations"`
}

// CreateStreamingProfile creates a new, custom streaming profile.
//
// https://cloudinary.com/documentation/admin_api#create_a_streaming_profile
func (a *Api) CreateStreamingProfile(ctx context.Context, params CreateStreamingProfileParams) (*GetStreamingProfileResult, error) {
	res := &GetStreamingProfileResult{}
	_, err := a.post(ctx, api.BuildPath(StreamingProfiles), params, res)

	return res, err
}

type UpdateStreamingProfileParams struct {
	Name            string                          `json:"-"` // The name of the streaming profile to update.
	DisplayName     string                          `json:"display_name,omitempty"`
	Representations StreamingProfileRepresentations `json:"representations"`
}

// UpdateStreamingProfile updates an existing streaming profile.
//
// You can update both custom and built-in profiles. The specified list of representations replaces the previous list.
//
// https://cloudinary.com/documentation/admin_api#update_an_existing_streaming_profile
func (a *Api) UpdateStreamingProfile(ctx context.Context, params UpdateStreamingProfileParams) (*GetStreamingProfileResult, error) {
	res := &GetStreamingProfileResult{}
	_, err := a.put(ctx, api.BuildPath(StreamingProfiles, params.Name), params, res)

	return res, err
}

type DeleteStreamingProfileParams struct {
	Name string `json:"-"` // The name of the streaming profile to delete or revert.
}

// DeleteStreamingProfile deletes or reverts the specified streaming profile.
//
// For custom streaming profiles, deletes the specified profile.
// For built-in streaming profiles, if the built-in profile was modified, reverts the profile to the original settings.
// For built-in streaming profiles that have not been modified, the Delete method returns an error.
//
// https://cloudinary.com/documentation/admin_api#delete_or_revert_the_specified_streaming_profile
func (a *Api) DeleteStreamingProfile(ctx context.Context, params DeleteStreamingProfileParams) (*DeleteStreamingProfileResult, error) {
	res := &DeleteStreamingProfileResult{}
	_, err := a.delete(ctx, api.BuildPath(StreamingProfiles, params.Name), params, res)

	return res, err
}

type DeleteStreamingProfileResult struct {
	Message string        `json:"message"`
	Error   api.ErrorResp `json:"error,omitempty"`
}
