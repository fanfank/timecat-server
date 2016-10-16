/**
 * @author  reetsee.com
 * @date    20161015
 */
package main

import (
    "bytes"
    "io/ioutil"
    "log"
    "os"
    "os/exec"
    "strconv"
    "time"

    "gopkg.in/gin-gonic/gin.v1"
)

func main() {
    _, err := exec.LookPath("timecat")
    if err != nil {
        log.Fatalln("`timecat` not found in $PATH")
    }

    router := gin.Default()
    router.POST("/timecat/v1/api/test", handle)
    router.Run(":3193")
}

func handle(c *gin.Context) {
    logContent := c.PostForm("logContent")
    st := c.PostForm("st")
    ed := c.PostForm("ed")

    if len([]rune(logContent)) > 1024 * 100 / 2 {
        c.JSON(413, gin.H{
            "errno": -1,
            "errmsg": "`logContent` too long",
        })
        return
    }

    if len([]rune(logContent)) == 0 || len([]rune(st)) == 0 || len([]rune(ed)) == 0 {
        c.JSON(200, gin.H{
            "errno": -1,
            "errmsg": "`logContent`, `st`, `ed` can not be empty",
        })
        return
    }

    // write log contents to file
    tmpFilePath := "/tmp/timecat.tmp." + strconv.FormatInt(time.Now().UnixNano(), 10)
    err := ioutil.WriteFile(tmpFilePath, []byte(logContent), 0644)
    if err != nil {
        log.Println("FATAL ", err)
        c.JSON(500, gin.H{
            "errno": -2,
            "errmsg": "failed to create tmp file",
        })
        return
    }

    defer func() {
        removeErr := os.Remove(tmpFilePath)
        if removeErr != nil {
            log.Println("ERROR ", err)
        }
    }()

    stdout, stderr, err := runTimecat(tmpFilePath, st, ed)

    if err != nil {
        c.JSON(200, gin.H{
            "errno": -100,
            "errmsg": "Ooooops，看来提交的信息不太对劲或`timecat`有Bug，可以直接反馈到`timecat` Issue",
            "data": err.Error() + ": " + stderr,
        })
    } else {
        c.JSON(200, gin.H{
            "errno": 0,
            "errmsg": "success",
            "data": stdout,
        })
    }
}

func runTimecat(tmpFilePath, st, ed string) (string, string, error) {
    var stdout bytes.Buffer
    var stderr bytes.Buffer

    cmd := exec.Command("timecat", "-s", st, "-e", ed, tmpFilePath)
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr

    err := cmd.Run()
    return stdout.String(), stderr.String(), err
}
