/**
 * @author  reetsee.com
 * @date    20161015
 */
package main

import (
    "io/ioutil"
    "log"
    _ "net/http"
    "os"
    "strconv"
    "time"

    "gopkg.in/gin-gonic/gin.v1"
)

func main() {
    router := gin.Default()
    router.POST("/timecat/api/test", handle)
    router.Run(":3193")
}

func handle(c *gin.Context) {
    logContent := c.PostForm("logContent")
    if len([]rune(logContent)) > 1024 * 100 / 2 {
        c.JSON(413, gin.H{
            "errno": -1,
            "errmsg": "logContent too long",
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

    //TODO(xuruiqi)
    c.JSON(200, gin.H{
        "errno": 0,
        "errmsg": "success",
        "data": logContent,
    })
}
