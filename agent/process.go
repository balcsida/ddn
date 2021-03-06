package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/djavorszky/ddn/common/inet"
	"github.com/djavorszky/ddn/common/logger"
	"github.com/djavorszky/ddn/common/model"
	"github.com/djavorszky/ddn/common/status"
	"github.com/djavorszky/notif"
)

func startImport(dbreq model.DBRequest) {
	upd8Path := fmt.Sprintf("%s/%s", conf.MasterAddress, "upd8")

	ch := notif.New(dbreq.ID, upd8Path)
	defer close(ch)

	ch <- notif.Y{StatusCode: status.DownloadInProgress, Msg: "Downloading dump"}
	logger.Debug("Downloading dump from %q", dbreq.DumpLocation)

	path, err := inet.DownloadFile("dumps", dbreq.DumpLocation)
	if err != nil {
		db.DropDatabase(dbreq)
		logger.Error("could not download file: %v", err)

		ch <- notif.Y{StatusCode: status.DownloadFailed, Msg: "Downloading file failed: " + err.Error()}
		return
	}
	defer os.Remove(path)

	if isArchive(path) {
		ch <- notif.Y{StatusCode: status.ExtractingArchive, Msg: "Extracting archive"}

		logger.Debug("Extracting archive: %v", path)

		var (
			files []string
			err   error
		)

		switch filepath.Ext(path) {
		case ".zip":
			files, err = unzip(path)
		case ".gz":
			files, err = ungzip(path)
		case ".tar":
			files, err = untar(path)
		default:
			db.DropDatabase(dbreq)
			logger.Error("import process stopped; encountered unsupported archive")

			ch <- notif.Y{StatusCode: status.ArchiveNotSupported, Msg: "archive not supported"}
			return
		}
		for _, f := range files {
			defer os.Remove(f)
		}

		if err != nil {
			db.DropDatabase(dbreq)
			logger.Error("could not extract archive: %v", err)

			ch <- notif.Y{StatusCode: status.ExtractingArchiveFailed, Msg: "Extracting file failed: " + err.Error()}
			return
		}

		if len(files) > 1 {
			db.DropDatabase(dbreq)
			logger.Error("import process stopped; more than one file found in archive")

			ch <- notif.Y{StatusCode: status.MultipleFilesInArchive, Msg: "Archive contains more than one file, import stopped"}
			return
		}

		path = files[0]
	}

	logger.Debug("Validating dump: %s", path)

	ch <- notif.Y{StatusCode: status.ValidatingDump, Msg: "Validating dump"}
	path, err = db.ValidateDump(path)
	if err != nil {
		db.DropDatabase(dbreq)
		logger.Error("database validation failed: %v", err)

		ch <- notif.Y{StatusCode: status.ValidationFailed, Msg: "Validating dump failed: " + err.Error()}
		return
	}

	if !strings.Contains(path, "dumps") {
		oldPath := path
		path = "dumps" + string(os.PathSeparator) + path

		os.Rename(oldPath, path)
	}

	path, _ = filepath.Abs(path)
	defer os.Remove(path)

	dbreq.DumpLocation = path

	logger.Debug("Importing dump: %v", path)
	ch <- notif.Y{StatusCode: status.ImportInProgress, Msg: "Importing"}

	start := time.Now()

	err = db.ImportDatabase(dbreq)
	if err != nil {
		logger.Error("could not import database: %v", err)

		ch <- notif.Y{StatusCode: status.ImportFailed, Msg: "Importing dump failed: " + err.Error()}
		return
	}

	logger.Debug("Import succeded in %v", time.Since(start))
	ch <- notif.Y{StatusCode: status.Success, Msg: "Completed"}
}

// This method should always be called asynchronously
func keepAlive() {
	endpoint := fmt.Sprintf("%s/%s/%s", conf.MasterAddress, "alive", conf.ShortName)

	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		// Check if the endpoint is up
		if !inet.AddrExists(fmt.Sprintf("%s/%s", conf.MasterAddress, "heartbeat")) {
			if registered {
				logger.Error("Lost connection to master server, will attempt to reconnect once it's back.")

				registered = false
			}

			continue
		}

		// If it is, check if we're not registered
		if !registered {
			logger.Info("Master server back online.")

			err := registerAgent()
			if err != nil {
				logger.Error("couldn't register with master: %v", err)
			}

			registered = true
		}

		respCode := inet.GetResponseCode(endpoint)
		if respCode == http.StatusOK {
			continue
		}

		// response is not "OK", so we need to register
		err := registerAgent()
		if err != nil {
			logger.Error("couldn't register with master: %v", err)
		}
	}
}
