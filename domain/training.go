package domain

import (
	"fmt"
	"path/filepath"
)

type UserTraining struct {
	User Account

	TrainingConfig
}

func (t *UserTraining) ToRepoPath() string {
	return filepath.Join(t.User.Account(), t.ProjectRepoId)
}

type TrainingConfig struct {
	ProjectName   ProjectName
	ProjectRepoId string

	Name TrainingName
	Desc TrainingDesc

	CodeDir  Directory
	BootFile FilePath

	Hyperparameters []KeyValue
	Env             []KeyValue
	Inputs          []Input
	EnableAim       bool
	EnableOutput    bool

	Compute Compute
}

type Compute struct {
	Type    ComputeType
	Version ComputeVersion
	Flavor  ComputeFlavor
}

type KeyValue struct {
	Key   CustomizedKey
	Value CustomizedValue
}

type Input struct {
	Key CustomizedKey
	ResourceRef
}

type ResourceRef struct {
	User   Account
	RepoId string
	File   string
}

func (r *ResourceRef) ToRepoPath() string {
	return filepath.Join(r.User.Account(), r.RepoId)
}

func (r *ResourceRef) ToPath() string {
	s := r.ToRepoPath()

	// The input is the directory. Appending "/" to make sure
	// the path is a directory for object storage service.
	return s + "/" + r.File
}

type JobDetail struct {
	Status   TrainingStatus
	Duration int
}

type JobInfo struct {
	JobId     string `json:"job_id"`
	LogDir    string `json:"log_dir"`
	AimDir    string `json:"aim_dir"`
	OutputDir string `json:"output_dir"`
}

func (t *TrainingConfig) IsCustomizeImageTraining() bool {
	return t.Compute.Version.ComputeImage() != ""
}

func (t *TrainingConfig) DefaultCommand() string {
	if t.IsCustomizeImageTraining() {
		lastDirectory := t.CodeDir.LastDirectory()
		filePath := t.BootFile.FilePath()
		depencyFile := filepath.Join(filepath.Dir(filePath), "pip-requirements.txt")

		command := fmt.Sprintf("pip install -r ${MA_JOB_DIR}/%s/%s "+
			"&& cd ${MA_JOB_DIR} "+
			"&& echo Current WorkDir: $(pwd) "+
			"&& python ${MA_JOB_DIR}/%s/%s", lastDirectory, depencyFile, lastDirectory, filePath)

		return command
	}
	return ""
}

func (t *TrainingConfig) DeafultBootFile() string {
	if t.IsCustomizeImageTraining() {
		return ""
	}
	return t.BootFile.FilePath()
}
