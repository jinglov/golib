package utils

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"github.com/jinglov/golib/logger"
)

func FileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func VerifyFileMd5Sum(file string, md5string string) bool {
	if !FileExist(file) {
		return false
	}

	md5fromfile := FileMd5Sum(file)
	if strings.EqualFold(md5fromfile, md5string) {
		return true
	}

	return false
}

func FileMd5Sum(file string) string {
	fp, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer fp.Close()
	rd := bufio.NewReader(fp)
	h := md5.New()
	_, err = io.Copy(h, rd)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}

func ZipByExec(localFile string) error {
	if FileExist(localFile) {
		cmdstr := "gzip " + localFile
		_, err := ExecCommand(cmdstr)
		if err != nil {
			logger.Error(cmdstr)
			logger.Error(err)
			return err
		}
		return nil
	}
	return fmt.Errorf("File not found:%s", localFile)
}

func UnzipByExec(gzfile string) error {
	if FileExist(gzfile) {
		cmdstr := "gunzip -f " + gzfile
		_, err := ExecCommand(cmdstr)
		if err != nil {
			logger.Error(cmdstr)
			logger.Error(err)
			return err
		}
		return nil
	}
	return fmt.Errorf("File not found:%s", gzfile)
}

func S3Cp(fromFile, toFile, profile, region string) error {
	logger.Info("Copy ", fromFile, "====>", toFile)
	if region == "" {
		region = "cn-north-1"
	}
	var cmdstr string
	if strings.EqualFold(profile, "") {
		cmdstr = "aws s3 --region=" + region + " cp " + fromFile + " " + toFile
	} else {
		cmdstr = "aws s3 --region=" + region + " --profile " + profile + " cp " + fromFile + " " + toFile
	}
	_, err := ExecCommand(cmdstr)
	if err != nil {
		logger.Error(cmdstr)
		logger.Error(err)
	}
	return err
}

func ExecCommand(cmdstr string) (out []byte, err error) {
	logger.Debug(cmdstr)
	cmd := exec.Command("/bin/sh", "-c", cmdstr)
	var stdOut, errOut io.ReadCloser
	stdOut, err = cmd.StdoutPipe()
	if err != nil {
		return
	}
	defer stdOut.Close()
	errOut, err = cmd.StderrPipe()
	if err != nil {
		return
	}
	defer errOut.Close()
	err = cmd.Start()
	if err != nil {
		return
	}
	out, err = ioutil.ReadAll(stdOut)
	if err != nil {
		return
	}
	var errb []byte
	errb, err = ioutil.ReadAll(errOut)
	if err != nil {
		return
	}
	if errb != nil && len(errb) > 0 {
		err = errors.New(string(errb))
	}
	return
}

func HdfsPut(fromFile, toPath, toFile string) (err error) {
	logger.Info("Put ", fromFile, "====>", toFile)
	cmdMkdir := "hdfs dfs -mkdir -p " + toPath

	logger.Debug(cmdMkdir)
	cmdMk := exec.Command("/bin/sh", "-c", cmdMkdir)
	_, err = cmdMk.Output()
	if err != nil {
		return errors.New("command error: " + cmdMkdir)
	}

	cmdstr := "hdfs dfs -put " + fromFile + " " + toPath + toFile
	_, err = ExecCommand(cmdstr)
	if err != nil {
		logger.Error(cmdstr)
		logger.Error(err)
	}
	return
}

func HdfsGet(fromFile, toFile string) (err error) {
	logger.Info("Get ", fromFile, "====>", toFile)
	err = Remove(toFile)
	if err != nil {
		logger.Error(err)
		return
	}
	cmdstr := "hdfs dfs -get " + fromFile + " " + toFile
	_, err = ExecCommand(cmdstr)
	if err != nil {
		logger.Error(cmdstr)
		logger.Error(err)
	}
	return
}

func Remove(file string) (err error) {
	if FileExist(file) {
		err = os.Remove(file)
		if err != nil {
			logger.Error(err)
			return
		}
	}
	return
}
