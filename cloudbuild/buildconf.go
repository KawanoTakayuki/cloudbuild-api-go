package cloudbuild

import (
	"io"
	"io/ioutil"
	"time"

	"golang.org/x/xerrors"
	"gopkg.in/yaml.v2"
)

// BuildConf build configuraqtion file
// https://cloud.google.com/cloud-build/docs/build-config
type BuildConf struct {
	Source        *Source           `json:"source,omitempty"`
	Steps         []Steps           `json:"steps,omitempty"`
	TimeOut       string            `json:"timeOut,omitempty"` // duration
	LogsBucket    string            `json:"logsBucket,omitempty"`
	Options       *BuildOptions     `json:"options,omitempty"`
	Substitutions map[string]string `json:"substitutions,omitempty"`
	Tags          []string          `json:"tags,omitempty"`
	Secrets       []Secrets         `json:"secrets,omitempty" yaml:"secrets"`
	Images        []string          `json:"images,omitempty"`
	Artifacts     *Artifact         `json:"artifacts,omitempty"`
	// output only
	ID               string              `json:"id,omitempty"`
	ProjectID        string              `json:"projectId,omitempty"`
	Status           string              `json:"status,omitempty"`
	StatusDetail     string              `json:"statusDetail,omitempty"`
	Results          *Results            `json:"results,omitempty"`
	CreateTime       *time.Time          `json:"createTime,omitempty"`
	StartTime        *time.Time          `json:"startTime,omitempty"`
	FinishTime       *time.Time          `json:"finishTime,omitempty"`
	SourceProvenance *SourceProvenance   `json:"sourceProvenance,omitempty"`
	BuildTriggerID   string              `json:"buildTriggerId,omitempty"`
	LogURL           string              `json:"logUrl,omitempty"`
	Timing           map[string]TimeSpan `json:"timing,omitempty"`
}

// Source ソース
type Source struct {
	StorageSource *StorageSource `json:"storageSource,omitempty"`
	RepoSource    *RepoSource    `json:"repoSource,omitempty"`
}

// StorageSource GCSソース
type StorageSource struct {
	Bucket     string `json:"bucket,omitempty"`
	Object     string `json:"object,omitempty"`
	Generation string `json:"generation,omitempty"`
}

// RepoSource gitソース
type RepoSource struct {
	ProjectID string `json:"projectId,omitempty"`
	RepoName  string `json:"repoName,omitempty"`
	Dir       string `json:"dir,omitempty"`
	// Union field
	BranchName string `json:"branchName,omitempty"`
	TagName    string `json:"tagName,omitempty"`
	CommitSha  string `json:"commitSha,omitempty"`
}

// Steps ビルドステップ
type Steps struct {
	Name       string   `json:"name,omitempty"`
	Args       []string `json:"args,omitempty"`
	Env        []string `json:"env,omitempty"`
	Dir        string   `json:"dir,omitempty"`
	ID         string   `json:"id,omitempty"`
	WaitFor    string   `json:"waitFor,omitempty"`
	Entrypoint string   `json:"entrypoint,omitempty"`
	SecretEnv  []string `json:"secretEnv,omitempty"`
	Volumes    []Volume `json:"volumes,omitempty"`
	TimeOut    string   `json:"timeout,omitempty"`
	// output only
	Timing     *TimeSpan `json:"timing,omitempty"`
	PullTiming *TimeSpan `json:"pullTiming,omitempty"`
	Status     string    `json:"status,omitempty"`
}

// Volume ボリューム
type Volume struct {
	Name string `json:"name,omitempty"`
	Path string `json:"path,omitempty"`
}

// Results ...
type Results struct {
	Image            []BuiltImage `json:"image,omitempty"`
	BuildStepImages  []string     `json:"buildStepImages,omitempty"`
	ArtifactManifest string       `json:"artifactManifest,omitempty"`
	NumArtifacts     string       `json:"numArtifacts,omitempty"`
	BuildStepOutputs []string     `json:"buildStepOutputs,omitempty"`
	ArtifactTiming   *TimeSpan    `json:"artifactTiming,omitempty"`
}

// BuiltImage ...
type BuiltImage struct {
	Name       string    `json:"name,omitempty"`
	Digest     string    `json:"digest,omitempty"`
	PushTiming *TimeSpan `json:"pushTiming,omitempty"`
}

// SourceProvenance ...
type SourceProvenance struct {
	ResolvedStorageSource StorageSource       `json:"resolvedStorageSource,omitempty"`
	ResolvedRepoSource    RepoSource          `json:"resolvedRepoSource,omitempty"`
	FileHashes            map[string]FileHash `json:"fileHashes,omitempty"`
}

// BuildOptions オプション
type BuildOptions struct {
	SourceProvenanceHash  []string `json:"sourceProvenanceHash,omitempty"`
	RequestedVerifyOption string   `json:"requestedVerifyOption,omitempty"`
	MachineType           string   `json:"machineType,omitempty"`
	DiskSizeGb            string   `json:"diskSizeGb,omitempty"`
	SubstitutionOption    string   `json:"substitutionOption,omitempty"`
	LogStreamingOption    string   `json:"logStreamingOption,omitempty"`
	WorkerPool            string   `json:"workerPool,omitempty"`
	Logging               string   `json:"logging,omitempty"`
	Env                   []string `json:"env,omitempty"`
	SecretEnv             []string `json:"secretEnv,omitempty"`
	Volumes               []Volume `json:"volumes,omitempty"`
}

// Secrets シークレット
type Secrets struct {
	KmsKeyName string            `json:"kmsKeyName,omitempty" yaml:"kmsKeyName"`
	SecretEnv  map[string]string `json:"secretEnv,omitempty" yaml:"secretEnv"`
}

// Artifact アーティファクト
type Artifact struct {
	Objects *ArtifactObjects `json:"objects,omitempty"`
}

// ArtifactObjects アーティファクトオブジェクト
type ArtifactObjects struct {
	Location string   `json:"location,omitempty"`
	Paths    []string `json:"paths,omitempty"`
	// output only
	Timing *TimeSpan `json:"timing,omitempty"`
}

// TimeSpan ...
type TimeSpan struct {
	StartTime *time.Time `json:"startTime,omitempty"`
	EndTime   *time.Time `json:"endTime,omitempty"`
}

// FileHash ...
type FileHash struct {
	FileHash []Hash `json:"fileHash,omitempty"`
}

// Hash ...
type Hash struct {
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}

// NewConfigration 新しいビルド構成ファイルを作成
func NewConfigration() *BuildConf {
	return &BuildConf{}
}

// Yaml yamlファイル読み込み
func (b *BuildConf) Yaml(file io.Reader) error {
	fileByte, err := ioutil.ReadAll(file)
	if err != nil {
		return xerrors.Errorf("input file read error: %w", err)
	}
	if err := yaml.Unmarshal(fileByte, b); err != nil {
		return xerrors.Errorf("yaml file unmarshalize error: %w", err)
	}
	return nil
}
