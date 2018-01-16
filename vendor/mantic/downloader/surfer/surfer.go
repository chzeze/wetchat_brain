// Copyright 2015 henrylee2cn Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// surfer是一款Go语言编写的高并发web下载器，支持 GET/POST/HEAD 方法及 http/https 协议，同时支持固定UserAgent自动保存cookie与随机大量UserAgent禁用cookie两种模式，高度模拟浏览器行为，可实现模拟登录等功能。
package surfer

import (
	"net/http"
	"net/http/cookiejar"
	"sync"

	"mantic/conf"
)

var (
	surf          Surfer
	phantom       Surfer
	once_surf     sync.Once
	once_phantom  sync.Once
	phantomjsFile = conf.Conf().Phantomjs
	tempJsDir     = conf.Conf().PhantomjsTemp
)

func Download(req Request) (resp *http.Response, err error) {
	err = req.Prepare()
	if err != nil {
		return
	}
	switch req.GetDownloaderID() {
	case SurfID:
		once_surf.Do(func() { surf = New() })
		resp, err = surf.Download(req)
	case PhomtomJsID:
		once_phantom.Do(func() { phantom = NewPhantom(phantomjsFile, tempJsDir) })
		resp, err = phantom.Download(req)
	}
	return
}

//销毁Phantomjs的js临时文件
func DestroyJsFiles() {
	if pt, ok := phantom.(*Phantom); ok {
		pt.DestroyJsFiles()
	}
}

// Downloader represents an core of HTTP web browser for crawler.
type Surfer interface {
	// GET @param url string, header http.Header, cookies []*http.Cookie
	// HEAD @param url string, header http.Header, cookies []*http.Cookie
	// POST PostForm @param url, referer string, values url.Values, header http.Header, cookies []*http.Cookie
	// POST-M PostMultipart @param url, referer string, values url.Values, header http.Header, cookies []*http.Cookie
	Download(Request) (resp *http.Response, err error)

	SetCookieJar(jar *cookiejar.Jar)
	Jar() *cookiejar.Jar
}
