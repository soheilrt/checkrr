package checkrr

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/soheilrt/checkrr/client"
	"github.com/soheilrt/checkrr/config"
	"time"
)

const (
	statusDownloading = "downloading"

	reasonStatusNotDownloading = "Download status is not downloading"
	reasonNotEnoughTime        = "Download started recently, threshold: %s, actual: %s"
	reasonDownloadTimeout      = "Download timed out, threshold: %s, actual: %s"
	reasonSlowDownloadSpeed    = "Average speed is below %v/s: %v"

	reasonAllGood = "Average speed is %v/s"
)

type ClientRR interface {
	FetchDownloads() ([]client.Download, error)
	DeleteFromQueue(ids []int) error
}

type CheckRR struct {
	name       string
	client     ClientRR
	conditions config.Conditions
}

func NewCheckRR(
	name string,
	client ClientRR,
	conditions config.Conditions,
) *CheckRR {
	return &CheckRR{
		name:       name,
		client:     client,
		conditions: conditions,
	}
}

func (c *CheckRR) Check() error {
	log.Debugf("Checking for stuck downloads on %s...", c.name)
	downloads, err := c.client.FetchDownloads()
	if err != nil {
		return err
	}

	stucks := []int{}
	for _, download := range downloads {
		stuck, reason, err := c.IsDownloadStuck(download)
		if err != nil {
			return fmt.Errorf("error checking download status: %v", err)
		}
		if stuck {
			stucks = append(stucks, download.ID)
			log.Warnf("Stuck download detected [ID: %d]: %s, Reason: %s", download.ID, download.Title, reason)
		} else {
			log.Infof("Download is OK [ID: %d]: %s, Reason: %s", download.ID, download.Title, reason)
		}
	}
	err = c.client.DeleteFromQueue(stucks)
	if err != nil {
		return fmt.Errorf("error deleting from queue: %v", err)
	}
	return nil
}

func (c *CheckRR) IsDownloadStuck(download client.Download) (bool, string, error) {
	if download.Status != statusDownloading {
		return false, reasonStatusNotDownloading, nil
	}

	addedTime, err := time.Parse(time.RFC3339, download.Added)
	if err != nil {
		return false, "", fmt.Errorf("error parsing added time: %v", err)
	}

	if c.conditions.WaitingThreshold > 0 && time.Since(addedTime) < c.conditions.WaitingThreshold {
		return false, fmt.Sprintf(
				reasonNotEnoughTime,
				c.conditions.WaitingThreshold,
				time.Since(addedTime),
			),
			nil
	}

	if c.conditions.DownloadTimeoutThreshold > 0 && time.Since(addedTime) > c.conditions.DownloadTimeoutThreshold {
		return true, fmt.Sprintf(
				reasonDownloadTimeout,
				c.conditions.DownloadTimeoutThreshold,
				time.Since(addedTime),
			),
			nil
	}

	avg := averageSpeed(download)
	if c.conditions.AverageSpeedThreshold > 0 && avg < c.conditions.AverageSpeedThreshold {
		return true, fmt.Sprintf(
				reasonSlowDownloadSpeed,
				bytesToHumanReadable(c.conditions.AverageSpeedThreshold),
				bytesToHumanReadable(avg),
			),
			nil
	}

	return false, fmt.Sprintf(reasonAllGood, bytesToHumanReadable(avg)), nil
}
func averageSpeed(download client.Download) float64 {
	// Parse the added time
	addedTime, err := time.Parse(time.RFC3339, download.Added)
	if err != nil {
		log.WithError(err).Error("Error parsing added time")
		return 0
	}
	return float64(download.Size-download.SizeLeft) / time.Since(addedTime).Seconds()
}

func bytesToHumanReadable(bytes float64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%vB", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", bytes/float64(div), "KMGTPE"[exp])
}
