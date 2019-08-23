package main

import (
	"fmt"

	"github.com/donatj/imgdedup"
	humanize "github.com/dustin/go-humanize"
)

func displayDiff(fileList []string, imgdata map[string]*imgdedup.ImageInfo) {
	fileLength := len(fileList)
	for i := 0; i < fileLength; i++ {
		for j := i + 1; j < fileLength; j++ {

			leftf := fileList[i]
			rightf := fileList[j]

			leftimg, ok1 := imgdata[leftf]
			rightimg, ok2 := imgdata[rightf]

			if ok1 && ok2 {

				avgdata1 := leftimg.Data
				avgdata2 := rightimg.Data

				if leftf == rightf {
					continue
				}

				xdiff := imgdedup.Diff(avgdata1, avgdata2, *subdivisions)

				if xdiff < uint64(*tolerance) {

					fmt.Println(leftf)
					fmt.Printf("    %d x %d\n    %s\n", leftimg.Bounds.Dx(), leftimg.Bounds.Dy(), humanize.Bytes(leftimg.Filesize))

					fmt.Println(rightf)
					fmt.Printf("    %d x %d\n    %s\n", rightimg.Bounds.Dx(), rightimg.Bounds.Dy(), humanize.Bytes(rightimg.Filesize))

					fmt.Println("")
					fmt.Println("Diff: ", xdiff)

					if xdiff > 0 && leftimg.Filesize != rightimg.Filesize {
						if *difftool != "" {
							diffTool(*difftool, leftf, rightf)
						}
					}

					fmt.Println("- - - - - - - - - -")
				}

			}

		}
	}
}
