package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type DockerRegistry struct {
	HubURL string
}

type TagFilter interface {
	latestImageDigest() (Image, error)
}

type Tag struct {
	ID            int       `json:"id"`
	LastUpdated   time.Time `json:"last_updated"`
	Name          string    `json:"name"`
	Repository    int       `json:"repository"`
	Status        string    `json:"status"`
	TagLastPulled time.Time `json:"tag_last_pulled"`
	TagLastPushed time.Time `json:"tag_last_pushed"`
	Images        []Image   `json:"images"`
}

type Image struct {
	Architecture string    `json:"architecture"`
	Digest       string    `json:"digest"`
	Layers       []Layer   `json:"layers"`
	OS           string    `json:"os"`
	Status       string    `json:"status"`
	LastPulled   time.Time `json:"last_pulled"`
	LastPushed   time.Time `json:"last_pushed"`
}

type Layer struct {
	Digest      string `json:"digest"`
	Size        int    `json:"size"`
	Instruction string `json:"instruction"`
}

func NewDockerRegistry(hubURL string) (Registry, error) {
	if hubURL == "" {
		hubURL = "https://hub.docker.com"
	}

	return &DockerRegistry{
		HubURL: hubURL,
	}, nil
}

func (d *DockerRegistry) LatestImageDigest(namespace string, repo string) (string, error) {
	latestTag, err := d.latestTag(namespace, repo)

	if err != nil {
		return "", err
	}

	image, err := latestTag.latestImageDigest()

	if err != nil {
		return "", err
	}

	return image.Digest, nil
}

func filterLatestTag(tags []Tag) Tag {
	currentTag := Tag{
		ID: -1,
	}

	for _, tag := range tags {
		if currentTag.ID < 0 {
			currentTag = tag
			continue
		}

		if tag.TagLastPushed.After(currentTag.TagLastPushed) {
			currentTag = tag
		}
	}

	return currentTag
}

func (d *DockerRegistry) TestRepo(namespace string, repo string) (int, error) {
	// Pinging tags URL since that does not need authentication to access
	manifestPath := fmt.Sprintf("v2/namespaces/%s/repositories/%s/tags", namespace, repo)
	endpoint := fmt.Sprintf("%s/%s", d.HubURL, manifestPath)

	var checkResponse struct {
		Detail  string `json:"detail"`
		Message string `json:"message"`
	}

	resp, err := http.Get(endpoint)

	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error getting images summary from %s: %s", endpoint, err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error reading response reading images summary: %s", err)
	}

	if err := json.Unmarshal(body, &checkResponse); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error unmarshalling JSON response while reading images summary: %s", err)
	}

	switch sc := resp.StatusCode; {
	case sc == http.StatusNotFound:
		return sc, fmt.Errorf("repo %s under namespace %s does not exist: %s", repo, namespace, checkResponse.Message)
	case sc >= 400 && sc <= 499:
		return sc, fmt.Errorf("client error checking namespace %s and repo %s: %s", namespace, repo, checkResponse.Message)
	case sc >= 500 && sc <= 599:
		return sc, fmt.Errorf("server error checking namespace %s and repo %s: %s", namespace, repo, checkResponse.Message)
	}

	return http.StatusOK, nil
}

// TODO: use one API call to images API instead: https://docs.docker.com/docker-hub/api/latest/#tag/images/operation/GetNamespacesRepositoriesImages
func (t Tag) latestImageDigest() (Image, error) {
	if len(t.Images) == 0 {
		return Image{}, fmt.Errorf("no images found for tag %s", t.Name)
	}

	currentImage := Image{
		Digest: "",
	}

	for _, image := range t.Images {
		if currentImage.Digest == "" {
			currentImage = image
			continue
		}

		if image.LastPushed.After(currentImage.LastPushed) {
			currentImage = image
		}
	}

	return currentImage, nil
}

func (d *DockerRegistry) latestTag(namespace string, repo string) (TagFilter, error) {
	manifestPath := fmt.Sprintf("v2/namespaces/%s/repositories/%s/tags", namespace, repo)
	endpoint := fmt.Sprintf("%s/%s", d.HubURL, manifestPath)

	// See https://docs.docker.com/docker-hub/api/latest/#tag/repositories/paths/~1v2~1namespaces~1%7Bnamespace%7D~1repositories~1%7Brepository%7D~1tags/get
	var tagsResponse struct {
		Count    int    `json:"count"`
		Next     string `json:"next"`
		Previous string `json:"previous"`
		Results  []Tag  `json:"results"`
	}

	resp, err := http.Get(endpoint)

	if err != nil {
		return nil, fmt.Errorf("error fetching initial tags from %s: %s", endpoint, err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("error reading response body while listing tags: %s", err)
	}

	if err := json.Unmarshal(body, &tagsResponse); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON response while listing tags: %s", err)
	}

	tagsChecked := len(tagsResponse.Results)

	if tagsChecked == 0 {
		return nil, fmt.Errorf("no tags found for namespace %s and repo %s. URL queried: %s", namespace, repo, endpoint)
	}

	latestTag := filterLatestTag(tagsResponse.Results)

	for tagsResponse.Count != tagsChecked {
		nextPage := tagsResponse.Next
		resp, err := http.Get(nextPage)

		if err != nil {
			return nil, fmt.Errorf("error fetching tags at page %s: %s", nextPage, err)
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)

		if err != nil {
			return nil, fmt.Errorf("error reading response body while listing tags at page %s: %s", nextPage, err)
		}

		if err := json.Unmarshal(body, &tagsResponse); err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON response while listing tags at page %s: %s", nextPage, err)
		}

		latestTag = filterLatestTag(tagsResponse.Results)
		tagsChecked += len(tagsResponse.Results)
	}

	return latestTag, nil
}
