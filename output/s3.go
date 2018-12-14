package output

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/falzm/fusili/kvmap"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

type S3Output struct {
	name      string
	accessKey string
	secretKey string
	region    string
	bucket    string
	filePath  string

	s3 *s3.S3
}

func init() {
	Outputs["s3"] = func(name string, settings map[string]interface{}) (Output, error) {
		var err error

		o := &S3Output{name: name}

		o.accessKey, err = kvmap.GetString(settings, "access_key", true)
		if err != nil {
			return nil, fmt.Errorf("s3 output: %s", err)
		}

		o.secretKey, err = kvmap.GetString(settings, "secret_key", true)
		if err != nil {
			return nil, fmt.Errorf("s3 output: %s", err)
		}

		o.region, err = kvmap.GetString(settings, "region", true)
		if err != nil {
			return nil, fmt.Errorf("s3 output: %s", err)
		}

		o.bucket, err = kvmap.GetString(settings, "bucket", true)
		if err != nil {
			return nil, fmt.Errorf("s3 output: %s", err)
		}

		o.filePath, err = kvmap.GetString(settings, "file_path", true)
		if err != nil {
			return nil, fmt.Errorf("s3 output: %s", err)
		}

		o.s3 = s3.New(aws.Auth{
			AccessKey: o.accessKey,
			SecretKey: o.secretKey,
		},
			aws.EUWest,
		)

		return o, nil
	}
}

func (o *S3Output) Report(hosts map[string][]int) error {
	report := struct {
		Date  int64            `json:"date"`
		Hosts map[string][]int `json:"hosts"`
	}{
		Date:  time.Now().Unix(),
		Hosts: hosts,
	}

	jsonReport, err := json.Marshal(report)
	if err != nil {
		return fmt.Errorf("s3 output: unable to encode report to JSON: %s", err)
	}

	if err := o.s3.Bucket(o.bucket).Put(
		o.filePath,
		jsonReport,
		"application/json",
		s3.Private,
	); err != nil {
		return fmt.Errorf("s3 output: unable to put report file: %s", err)
	}

	return nil
}
