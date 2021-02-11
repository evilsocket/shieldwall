package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/shieldwall/version"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"syscall"
	"time"
)

// Untar takes a destination path and a reader; a tar reader loops over the tarfile
// creating the file structure at 'dst' along the way, and writing any files
// credits to https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07
func Untar(dst string, r io.Reader) error {

	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}

var repo = "evilsocket/shieldwall"
var extractTo = "/tmp/shieldwall/"
var installer = filepath.Join(extractTo, "install.sh")

var versionParser = regexp.MustCompile("^https://github\\.com/" + repo + "/releases/tag/v([\\d\\.a-z]+)$")

func updater() {
	log.Info("update checker started with a %s period", updateCheckPeriod)

	noRedirectClient := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	url := "https://github.com/" + repo + "/releases/latest"
	for {
		log.Debug("checking for updates %s", url)

		req, _ := http.NewRequest("GET", url, nil)
		resp, err := noRedirectClient.Do(req)
		if err != nil {
			log.Error("error while checking latest version: %v", err)
			continue
		}
		defer resp.Body.Close()

		location := resp.Header.Get("Location")

		log.Debug("location header = '%s'", location)

		m := versionParser.FindStringSubmatch(location)
		if len(m) == 2 {
			latest := m[1]
			log.Debug("latest version is '%s'", latest)
			if version.Version != latest {
				filename := fmt.Sprintf("shieldwall-agent_%s_linux_%s.tar.gz",
					latest,
					runtime.GOARCH)

				update := fmt.Sprintf("https://github.com/%s/releases/download/v%s/%s",
					repo,
					latest,
					filename)

				log.Important("downloading update to %s from %s", latest, update)

				tmpFileName := filepath.Join("/tmp/", filename)
				out, err := os.Create(tmpFileName)
				if err != nil {
					log.Error("error creating %s: %v", tmpFileName, err)
					continue
				}

				start := time.Now()
				resp, err := http.Get(update)
				defer resp.Body.Close()
				if err != nil {
					out.Close()
					log.Error("error downloading %s: %v", tmpFileName, err)
					continue
				}

				written, err := io.Copy(out, resp.Body)
				out.Close()

				if err != nil {
					log.Error("error writing to %s: %v", tmpFileName, err)
					continue
				}

				log.Debug("downloaded %d bytes to %s in %s", written, tmpFileName, time.Since(start))

				start = time.Now()

				if out, err = os.Open(tmpFileName); err != nil {
					log.Error("error opening %s: %v", tmpFileName, err)
					continue
				} else if err = os.MkdirAll(extractTo, os.ModePerm); err != nil {
					out.Close()
					log.Error("error creating folder %s: %v", extractTo, err)
					continue
				} else if err = Untar(extractTo, out); err != nil {
					out.Close()
					log.Error("error extracting %s to %s: %v", tmpFileName, extractTo, err)
					continue
				}
				out.Close()

				log.Debug("extracted in %s", time.Since(start))

				log.Info("running installer at %s", installer)

				cmd := exec.Command(installer)
				cmd.Dir = extractTo
				cmd.Env = os.Environ()
				// https://stackoverflow.com/questions/33165530/prevent-ctrlc-from-interrupting-exec-command-in-golang
				cmd.SysProcAttr = &syscall.SysProcAttr{
					Setpgid: true,
					Pgid:    0,
				}
				if err = cmd.Start(); err != nil {
					log.Error("error starting the installer: %v", err)
				} else {
					log.Info("installer running as pid %d", cmd.Process.Pid)
				}
			} else {
				log.Debug("no updates available")
			}
		} else {
			log.Debug("unexpected location header: '%s'", location)
		}

		time.Sleep(updateCheckPeriod)
	}
}
