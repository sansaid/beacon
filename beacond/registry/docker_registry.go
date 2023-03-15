package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	ID                  int       `json:"id"`
	Creator             int       `json:"creator"`
	LastUpdated         time.Time `json:"last_updated"`
	LastUpdater         int       `json:"last_updater"`
	LastUpdaterUsername string    `json:"last_updater_username"`
	Name                string    `json:"name"`
	Repository          int       `json:"repository"`
	FullSize            int       `json:"full_size"`
	V2                  string    `json:"v2"`
	Status              string    `json:"status"`
	TagLastPulled       time.Time `json:"tag_last_pulled"`
	TagLastPushed       time.Time `json:"tag_last_pushed"`
	Images              []Image   `json:"images"`
}

type Image struct {
	Architecture string    `json:"architecture"`
	Features     string    `json:"features"`
	Variant      string    `json:"variant"`
	Digest       string    `json:"digest"`
	Layers       []Layer   `json:"layers"`
	OS           string    `json:"os"`
	OSFeatures   string    `json:"os_features"`
	OSVersion    string    `json:"os_version"`
	Size         int       `json:"size"`
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
	var TagsResponse struct {
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

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("error reading response body while listing tags: %s", err)
	}

	if err := json.Unmarshal(body, &TagsResponse); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON response while listing tags: %s", err)
	}

	tagsChecked := len(TagsResponse.Results)

	if tagsChecked == 0 {
		return nil, fmt.Errorf("no tags found for namespace %s and repo %s. URL queried: %s", namespace, repo, endpoint)
	}

	latestTag := filterLatestTag(TagsResponse.Results)

	for TagsResponse.Count != tagsChecked {
		nextPage := TagsResponse.Next
		resp, err := http.Get(nextPage)

		if err != nil {
			return nil, fmt.Errorf("error fetching tags at page %s: %s", nextPage, err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return nil, fmt.Errorf("error reading response body while listing tags at page %s: %s", nextPage, err)
		}

		if err := json.Unmarshal(body, &TagsResponse); err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON response while listing tags at page %s: %s", nextPage, err)
		}

		latestTag = filterLatestTag(TagsResponse.Results)
		tagsChecked += len(TagsResponse.Results)
	}

	return latestTag, nil
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
