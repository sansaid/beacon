package oci

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/labstack/gommon/log"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PodmanSuite struct {
	suite.Suite
	PodmanClient PodmanClient
	LogBuff      *bytes.Buffer
}

func (p *PodmanSuite) SetupTest() {
	client, _ := NewPodman()
	p.PodmanClient = client.(PodmanClient)
	p.LogBuff = new(bytes.Buffer)

	log.SetOutput(p.LogBuff)
}

func TestLinkManagerSuite(t *testing.T) {
	suite.Run(t, new(PodmanSuite))
}

func (p *PodmanSuite) TestType() {
	assert.Equal(p.T(), Podman, p.PodmanClient.Type())
}

func (p *PodmanSuite) TestCheckExistsPasses() {
	mockController := gomock.NewController(p.T())
	defer mockController.Finish()

	p.PodmanClient.runner = NewMockRunner(mockController)

	p.PodmanClient.runner.(*MockRunner).EXPECT().run("podman", "--version").Return([]byte("fake version"), nil)

	exists, err := p.PodmanClient.CheckExists()

	assert.NoError(p.T(), err)
	assert.True(p.T(), exists)
}

func (p *PodmanSuite) TestCheckExistsFails() {
	mockController := gomock.NewController(p.T())
	defer mockController.Finish()

	p.PodmanClient.runner = NewMockRunner(mockController)

	p.PodmanClient.runner.(*MockRunner).EXPECT().run("podman", "--version").Return([]byte("fake output"), nil)

	exists, err := p.PodmanClient.CheckExists()

	assert.ErrorContains(p.T(), err, "The output was not recognised")
	assert.False(p.T(), exists)
}

func (p *PodmanSuite) TestCheckExistsErrors() {
	mockController := gomock.NewController(p.T())
	defer mockController.Finish()

	p.PodmanClient.runner = NewMockRunner(mockController)

	p.PodmanClient.runner.(*MockRunner).EXPECT().run("podman", "--version").Return([]byte(""), fmt.Errorf("fake error"))

	exists, err := p.PodmanClient.CheckExists()

	assert.Error(p.T(), err)
	assert.False(p.T(), exists)
}

func (p *PodmanSuite) TestRunImageOK() {
	mockController := gomock.NewController(p.T())
	defer mockController.Finish()

	p.PodmanClient.runner = NewMockRunner(mockController)

	p.PodmanClient.runner.(*MockRunner).EXPECT().run("podman", "run", "fakeImageRef").Return([]byte("fake output"), nil)

	err := p.PodmanClient.RunImage("fakeImageRef")

	assert.NoError(p.T(), err)
}

func (p *PodmanSuite) TestRunImageErrors() {
	mockController := gomock.NewController(p.T())
	defer mockController.Finish()

	p.PodmanClient.runner = NewMockRunner(mockController)

	p.PodmanClient.runner.(*MockRunner).EXPECT().run("podman", "run", "fakeImageRef").Return([]byte(""), fmt.Errorf("fake error"))

	err := p.PodmanClient.RunImage("fakeImageRef")

	assert.ErrorContains(p.T(), err, "error running podman image fakeImageRef")
}

func (p *PodmanSuite) TestPullImageOK() {
	mockController := gomock.NewController(p.T())
	defer mockController.Finish()

	p.PodmanClient.runner = NewMockRunner(mockController)

	p.PodmanClient.runner.(*MockRunner).EXPECT().run("podman", "pull", "fakeImageRef").Return([]byte("fake output"), nil)

	err := p.PodmanClient.PullImage("fakeImageRef")

	assert.NoError(p.T(), err)
}

func (p *PodmanSuite) TestPullImageErrors() {
	mockController := gomock.NewController(p.T())
	defer mockController.Finish()

	p.PodmanClient.runner = NewMockRunner(mockController)

	p.PodmanClient.runner.(*MockRunner).EXPECT().run("podman", "pull", "fakeImageRef").Return([]byte(""), fmt.Errorf("fake error"))

	err := p.PodmanClient.PullImage("fakeImageRef")

	assert.ErrorContains(p.T(), err, "error pulling podman image fakeImageRef")
}

func (p *PodmanSuite) TestGetImagesOK() {
	mockController := gomock.NewController(p.T())
	defer mockController.Finish()

	p.PodmanClient.runner = NewMockRunner(mockController)

	args := []interface{}{"podman", "images", "--format", "json",
		fmt.Sprintf("--filter=reference='%s'", "fakeImagePrefix"),
		fmt.Sprintf("--filter=before='%s'", "oldImageRef"),
		fmt.Sprintf("--filter=dangling=%t", true),
	}

	p.PodmanClient.runner.(*MockRunner).EXPECT().run(args...).Return([]byte("[{\"Id\": \"imageIdA\"}, {\"Id\": \"imageIdB\"}]"), nil)

	actual, err := p.PodmanClient.GetImages("fakeImagePrefix", "oldImageRef", true)

	assert.NoError(p.T(), err)
	assert.ElementsMatch(p.T(), actual, []string{"imageIdA", "imageIdB"})
}

func (p *PodmanSuite) TestGetImagesOKEmptyList() {
	mockController := gomock.NewController(p.T())
	defer mockController.Finish()

	p.PodmanClient.runner = NewMockRunner(mockController)

	args := []interface{}{"podman", "images", "--format", "json",
		fmt.Sprintf("--filter=reference='%s'", "fakeImagePrefix"),
		fmt.Sprintf("--filter=before='%s'", "oldImageRef"),
		fmt.Sprintf("--filter=dangling=%t", true),
	}

	p.PodmanClient.runner.(*MockRunner).EXPECT().run(args...).Return([]byte("[{\"OtherField\": \"imageIdA\"}, {\"OtherField\": \"imageIdB\"}]"), nil)

	actual, err := p.PodmanClient.GetImages("fakeImagePrefix", "oldImageRef", true)

	assert.NoError(p.T(), err)
	assert.ElementsMatch(p.T(), actual, []string{})
}

func (p *PodmanSuite) TestGetImagesInvalidJSON() {
	mockController := gomock.NewController(p.T())
	defer mockController.Finish()

	p.PodmanClient.runner = NewMockRunner(mockController)

	args := []interface{}{"podman", "images", "--format", "json",
		fmt.Sprintf("--filter=reference='%s'", "fakeImagePrefix"),
		fmt.Sprintf("--filter=before='%s'", "oldImageRef"),
		fmt.Sprintf("--filter=dangling=%t", true),
	}

	p.PodmanClient.runner.(*MockRunner).EXPECT().run(args...).Return([]byte("not a json"), nil)

	_, err := p.PodmanClient.GetImages("fakeImagePrefix", "oldImageRef", true)

	assert.Error(p.T(), err, "error parsing images output for ref prefix fakeImagePrefix")
}

func (p *PodmanSuite) TestGetImagesErrors() {
	mockController := gomock.NewController(p.T())
	defer mockController.Finish()

	p.PodmanClient.runner = NewMockRunner(mockController)

	args := []interface{}{"podman", "images", "--format", "json",
		fmt.Sprintf("--filter=reference='%s'", "fakeImagePrefix"),
		fmt.Sprintf("--filter=before='%s'", "oldImageRef"),
		fmt.Sprintf("--filter=dangling=%t", true),
	}

	p.PodmanClient.runner.(*MockRunner).EXPECT().run(args...).Return([]byte(""), fmt.Errorf("fake error"))

	_, err := p.PodmanClient.GetImages("fakeImagePrefix", "oldImageRef", true)

	assert.Error(p.T(), err, "error getting images associated with prefix fakeImagePrefix")
}

func (p *PodmanSuite) TestRemoveImagesOK() {
	mockController := gomock.NewController(p.T())
	defer mockController.Finish()

	p.PodmanClient.runner = NewMockRunner(mockController)

	removeImagesArgs := []interface{}{"podman", "rm"}
	removeImagesArgs = append(removeImagesArgs, "imageIdA", "imageIdB")

	getImagesArgs := []interface{}{"podman", "images", "--format", "json",
		fmt.Sprintf("--filter=reference='%s'", "fakeImagePrefix"),
		fmt.Sprintf("--filter=before='%s'", "oldImageRef"),
		fmt.Sprintf("--filter=dangling=%t", true),
	}

	p.PodmanClient.runner.(*MockRunner).EXPECT().run(getImagesArgs...).Return([]byte("[{\"Id\": \"imageIdA\"}, {\"Id\": \"imageIdB\"}]"), nil)
	p.PodmanClient.runner.(*MockRunner).EXPECT().run(removeImagesArgs...).Return([]byte(""), nil)

	err := p.PodmanClient.RemoveImages("fakeImagePrefix", "oldImageRef")

	assert.NoError(p.T(), err)
}

func (p *PodmanSuite) TestRemoveImagesGetImagesError() {
	mockController := gomock.NewController(p.T())
	defer mockController.Finish()

	p.PodmanClient.runner = NewMockRunner(mockController)

	getImagesArgs := []interface{}{"podman", "images", "--format", "json",
		fmt.Sprintf("--filter=reference='%s'", "fakeImagePrefix"),
		fmt.Sprintf("--filter=before='%s'", "oldImageRef"),
		fmt.Sprintf("--filter=dangling=%t", true),
	}

	p.PodmanClient.runner.(*MockRunner).EXPECT().run(getImagesArgs...).Return([]byte(""), fmt.Errorf("fake error"))

	err := p.PodmanClient.RemoveImages("fakeImagePrefix", "oldImageRef")

	assert.Errorf(p.T(), err, "error getting images associated with prefix fakeImagePrefix")
}

func (p *PodmanSuite) TestRemoveImagesRmError() {
	mockController := gomock.NewController(p.T())
	defer mockController.Finish()

	p.PodmanClient.runner = NewMockRunner(mockController)

	removeImagesArgs := []interface{}{"podman", "rm"}
	removeImagesArgs = append(removeImagesArgs, "imageIdA", "imageIdB")

	getImagesArgs := []interface{}{"podman", "images", "--format", "json",
		fmt.Sprintf("--filter=reference='%s'", "fakeImagePrefix"),
		fmt.Sprintf("--filter=before='%s'", "oldImageRef"),
		fmt.Sprintf("--filter=dangling=%t", true),
	}

	p.PodmanClient.runner.(*MockRunner).EXPECT().run(getImagesArgs...).Return([]byte("[{\"Id\": \"imageIdA\"}, {\"Id\": \"imageIdB\"}]"), nil)
	p.PodmanClient.runner.(*MockRunner).EXPECT().run(removeImagesArgs...).Return([]byte(""), fmt.Errorf("fake error"))

	err := p.PodmanClient.RemoveImages("fakeImagePrefix", "oldImageRef")

	assert.Errorf(p.T(), err, "error removing podman images")
}

func (p *PodmanSuite) TestStopContainerOK() {
	mockController := gomock.NewController(p.T())
	defer mockController.Finish()

	p.PodmanClient.runner = NewMockRunner(mockController)

	args := []interface{}{"podman", "stop", "fakeImageRef"}

	p.PodmanClient.runner.(*MockRunner).EXPECT().run(args...).Return([]byte(""), nil)

	err := p.PodmanClient.StopContainer("fakeImageRef")

	assert.NoError(p.T(), err)
}

func (p *PodmanSuite) TestStopContainerErrors() {
	mockController := gomock.NewController(p.T())
	defer mockController.Finish()

	p.PodmanClient.runner = NewMockRunner(mockController)

	args := []interface{}{"podman", "stop", "fakeImageRef"}

	p.PodmanClient.runner.(*MockRunner).EXPECT().run(args...).Return([]byte(""), fmt.Errorf("fake error"))

	err := p.PodmanClient.StopContainer("fakeImageRef")

	assert.Errorf(p.T(), err, "error stopping container fakeImageRef")
}

func (p *PodmanSuite) TestContainersUsingImageMultipleStatusesOK() {
	mockController := gomock.NewController(p.T())
	defer mockController.Finish()

	p.PodmanClient.runner = NewMockRunner(mockController)

	statuses := []string{"running", "paused"}

	args := []interface{}{"podman", "ps", "--format", "json", "--filter=ancestor='fakeImageRef'"}

	for _, status := range statuses {
		args = append(args, fmt.Sprintf("--filter=status='%s'", status))
	}

	p.PodmanClient.runner.(*MockRunner).EXPECT().run(args...).Return([]byte("[{\"Id\": \"containerIdA\"}, {\"Id\": \"containerIdB\"}]"), nil)

	containers, err := p.PodmanClient.ContainersUsingImage("fakeImageRef", statuses)

	assert.NoError(p.T(), err)
	assert.ElementsMatch(p.T(), containers, []string{"containerIdA", "containerIdB"})
}

func (p *PodmanSuite) TestContainersUsingImageSingleStatusOK() {
	mockController := gomock.NewController(p.T())
	defer mockController.Finish()

	p.PodmanClient.runner = NewMockRunner(mockController)

	statuses := []string{"running"}

	args := []interface{}{"podman", "ps", "--format", "json", "--filter=ancestor='fakeImageRef'"}

	for _, status := range statuses {
		args = append(args, fmt.Sprintf("--filter=status='%s'", status))
	}

	p.PodmanClient.runner.(*MockRunner).EXPECT().run(args...).Return([]byte("[{\"Id\": \"containerIdA\"}, {\"Id\": \"containerIdB\"}]"), nil)

	containers, err := p.PodmanClient.ContainersUsingImage("fakeImageRef", statuses)

	assert.NoError(p.T(), err)
	assert.ElementsMatch(p.T(), containers, []string{"containerIdA", "containerIdB"})
}

func (p *PodmanSuite) TestContainersUsingImagePartialIds() {
	mockController := gomock.NewController(p.T())
	defer mockController.Finish()

	p.PodmanClient.runner = NewMockRunner(mockController)

	statuses := []string{"running", "paused"}

	args := []interface{}{"podman", "ps", "--format", "json", "--filter=ancestor='fakeImageRef'"}

	for _, status := range statuses {
		args = append(args, fmt.Sprintf("--filter=status='%s'", status))
	}

	p.PodmanClient.runner.(*MockRunner).EXPECT().run(args...).Return([]byte("[{\"Id\": \"containerIdA\"}, {\"NotAnId\": \"containerIdB\"}]"), nil)

	containers, err := p.PodmanClient.ContainersUsingImage("fakeImageRef", statuses)

	assert.NoError(p.T(), err)
	assert.ElementsMatch(p.T(), containers, []string{"containerIdA"})
}

func (p *PodmanSuite) TestContainersUsingImageNoIds() {
	mockController := gomock.NewController(p.T())
	defer mockController.Finish()

	p.PodmanClient.runner = NewMockRunner(mockController)

	statuses := []string{"running", "paused"}

	args := []interface{}{"podman", "ps", "--format", "json", "--filter=ancestor='fakeImageRef'"}

	for _, status := range statuses {
		args = append(args, fmt.Sprintf("--filter=status='%s'", status))
	}

	p.PodmanClient.runner.(*MockRunner).EXPECT().run(args...).Return([]byte("[{\"NotAnId\": \"containerIdA\"}, {\"NotAnId\": \"containerIdB\"}]"), nil)

	containers, err := p.PodmanClient.ContainersUsingImage("fakeImageRef", statuses)

	assert.NoError(p.T(), err)
	assert.ElementsMatch(p.T(), containers, []string{})
}

func (p *PodmanSuite) TestContainersUsingImageErrors() {
	mockController := gomock.NewController(p.T())
	defer mockController.Finish()

	p.PodmanClient.runner = NewMockRunner(mockController)

	statuses := []string{"running", "paused"}

	args := []interface{}{"podman", "ps", "--format", "json", "--filter=ancestor='fakeImageRef'"}

	for _, status := range statuses {
		args = append(args, fmt.Sprintf("--filter=status='%s'", status))
	}

	p.PodmanClient.runner.(*MockRunner).EXPECT().run(args...).Return([]byte(""), fmt.Errorf("fake error"))

	_, err := p.PodmanClient.ContainersUsingImage("fakeImageRef", statuses)

	assert.Error(p.T(), err, "error getting containers associated with image fakeImagePrefix")
}

func (p *PodmanSuite) TestContainersUsingImageInvalidJSON() {
	mockController := gomock.NewController(p.T())
	defer mockController.Finish()

	p.PodmanClient.runner = NewMockRunner(mockController)

	statuses := []string{"running", "paused"}

	args := []interface{}{"podman", "ps", "--format", "json", "--filter=ancestor='fakeImageRef'"}

	for _, status := range statuses {
		args = append(args, fmt.Sprintf("--filter=status='%s'", status))
	}

	p.PodmanClient.runner.(*MockRunner).EXPECT().run(args...).Return([]byte("not json"), nil)

	_, err := p.PodmanClient.ContainersUsingImage("fakeImageRef", statuses)

	assert.Error(p.T(), err, "error parsing containers output for image fakeImagePrefix")
}

func (p *PodmanSuite) TestStopContainersByImageOK() {
	mockController := gomock.NewController(p.T())
	defer mockController.Finish()

	p.PodmanClient.runner = NewMockRunner(mockController)

	psArgs := []interface{}{"podman", "ps", "--format", "json", "--filter=ancestor='fakeImageRef'", "--filter=status='running'"}

	p.PodmanClient.runner.(*MockRunner).EXPECT().run(psArgs...).Return([]byte("[{\"Id\": \"containerIdA\"}, {\"Id\": \"containerIdB\"}]"), nil)
	p.PodmanClient.runner.(*MockRunner).EXPECT().run("podman", "stop", "containerIdA").Return([]byte(""), nil)
	p.PodmanClient.runner.(*MockRunner).EXPECT().run("podman", "stop", "containerIdB").Return([]byte(""), nil)

	err := p.PodmanClient.StopContainersByImage("fakeImageRef")

	assert.NoError(p.T(), err)
}

func (p *PodmanSuite) TestStopContainersByImagePsErrors() {
	mockController := gomock.NewController(p.T())
	defer mockController.Finish()

	p.PodmanClient.runner = NewMockRunner(mockController)

	psArgs := []interface{}{"podman", "ps", "--format", "json", "--filter=ancestor='fakeImageRef'", "--filter=status='running'"}

	p.PodmanClient.runner.(*MockRunner).EXPECT().run(psArgs...).Return([]byte(""), fmt.Errorf("fake error"))

	err := p.PodmanClient.StopContainersByImage("fakeImageRef")

	assert.Error(p.T(), err, "fake error")
}

func (p *PodmanSuite) TestStopContainersByImageStopErrors() {
	mockController := gomock.NewController(p.T())
	defer mockController.Finish()

	p.PodmanClient.runner = NewMockRunner(mockController)

	psArgs := []interface{}{"podman", "ps", "--format", "json", "--filter=ancestor='fakeImageRef'", "--filter=status='running'"}

	p.PodmanClient.runner.(*MockRunner).EXPECT().run(psArgs...).Return([]byte("[{\"Id\": \"containerIdA\"}, {\"Id\": \"containerIdB\"}]"), nil)
	p.PodmanClient.runner.(*MockRunner).EXPECT().run("podman", "stop", "containerIdA").Return([]byte(""), fmt.Errorf("fake error"))
	p.PodmanClient.runner.(*MockRunner).EXPECT().run("podman", "stop", "containerIdB").Return([]byte(""), nil)

	err := p.PodmanClient.StopContainersByImage("fakeImageRef")

	assert.Contains(p.T(), p.LogBuff.String(), "error stopping container containerIdA")
	assert.NoError(p.T(), err)
}
