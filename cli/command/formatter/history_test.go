package formatter

import (
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/docker/pkg/stringutils"
)

type historyCase struct {
	historyCtx historyContext
	expValue   string
	expHeader  string
	call       func() string
}

func TestHistoryContext_ID(t *testing.T) {
	id := stringid.GenerateRandomID()

	var ctx historyContext
	cases := []historyCase{
		{
			historyContext{
				h:     image.HistoryResponseItem{ID: id},
				trunc: false,
			}, id, historyIDHeader, ctx.ID,
		},
		{
			historyContext{
				h:     image.HistoryResponseItem{ID: id},
				trunc: true,
			}, stringid.TruncateID(id), historyIDHeader, ctx.ID,
		},
	}

	for _, c := range cases {
		ctx = c.historyCtx
		v := c.call()
		if strings.Contains(v, ",") {
			compareMultipleValues(t, v, c.expValue)
		} else if v != c.expValue {
			t.Fatalf("Expected %s, was %s\n", c.expValue, v)
		}

		h := ctx.FullHeader()
		if h != c.expHeader {
			t.Fatalf("Expected %s, was %s\n", c.expHeader, h)
		}
	}
}

func TestHistoryContext_CreatedSince(t *testing.T) {
	unixTime := time.Now().AddDate(0, 0, -7).Unix()
	expected := "7 days"

	var ctx historyContext
	cases := []historyCase{
		{
			historyContext{
				h:     image.HistoryResponseItem{Created: unixTime},
				trunc: false,
			}, expected, createdSinceHeader, ctx.CreatedSince,
		},
	}

	for _, c := range cases {
		ctx = c.historyCtx
		v := c.call()
		if strings.Contains(v, ",") {
			compareMultipleValues(t, v, c.expValue)
		} else if v != c.expValue {
			t.Fatalf("Expected %s, was %s\n", c.expValue, v)
		}

		h := ctx.FullHeader()
		if h != c.expHeader {
			t.Fatalf("Expected %s, was %s\n", c.expHeader, h)
		}
	}
}

func TestHistoryContext_CreatedBy(t *testing.T) {
	withTabs := `/bin/sh -c apt-key adv --keyserver hkp://pgp.mit.edu:80	--recv-keys 573BFD6B3D8FBC641079A6ABABF5BD827BD9BF62	&& echo "deb http://nginx.org/packages/mainline/debian/ jessie nginx" >> /etc/apt/sources.list  && apt-get update  && apt-get install --no-install-recommends --no-install-suggests -y       ca-certificates       nginx=${NGINX_VERSION}       nginx-module-xslt       nginx-module-geoip       nginx-module-image-filter       nginx-module-perl       nginx-module-njs       gettext-base  && rm -rf /var/lib/apt/lists/*`
	expected := `/bin/sh -c apt-key adv --keyserver hkp://pgp.mit.edu:80 --recv-keys 573BFD6B3D8FBC641079A6ABABF5BD827BD9BF62 && echo "deb http://nginx.org/packages/mainline/debian/ jessie nginx" >> /etc/apt/sources.list  && apt-get update  && apt-get install --no-install-recommends --no-install-suggests -y       ca-certificates       nginx=${NGINX_VERSION}       nginx-module-xslt       nginx-module-geoip       nginx-module-image-filter       nginx-module-perl       nginx-module-njs       gettext-base  && rm -rf /var/lib/apt/lists/*`

	var ctx historyContext
	cases := []historyCase{
		{
			historyContext{
				h:     image.HistoryResponseItem{CreatedBy: withTabs},
				trunc: false,
			}, expected, createdByHeader, ctx.CreatedBy,
		},
		{
			historyContext{
				h:     image.HistoryResponseItem{CreatedBy: withTabs},
				trunc: true,
			}, stringutils.Ellipsis(expected, 45), createdByHeader, ctx.CreatedBy,
		},
	}

	for _, c := range cases {
		ctx = c.historyCtx
		v := c.call()
		if strings.Contains(v, ",") {
			compareMultipleValues(t, v, c.expValue)
		} else if v != c.expValue {
			t.Fatalf("Expected %s, was %s\n", c.expValue, v)
		}

		h := ctx.FullHeader()
		if h != c.expHeader {
			t.Fatalf("Expected %s, was %s\n", c.expHeader, h)
		}
	}
}

func TestHistoryContext_Size(t *testing.T) {
	size := int64(182964289)
	expected := "183MB"

	var ctx historyContext
	cases := []historyCase{
		{
			historyContext{
				h:     image.HistoryResponseItem{Size: size},
				trunc: false,
			}, expected, sizeHeader, ctx.Size,
		},
	}

	for _, c := range cases {
		ctx = c.historyCtx
		v := c.call()
		if strings.Contains(v, ",") {
			compareMultipleValues(t, v, c.expValue)
		} else if v != c.expValue {
			t.Fatalf("Expected %s, was %s\n", c.expValue, v)
		}

		h := ctx.FullHeader()
		if h != c.expHeader {
			t.Fatalf("Expected %s, was %s\n", c.expHeader, h)
		}
	}
}

func TestHistoryContext_Comment(t *testing.T) {
	comment := "Some comment"

	var ctx historyContext
	cases := []historyCase{
		{
			historyContext{
				h:     image.HistoryResponseItem{Comment: comment},
				trunc: false,
			}, comment, commentHeader, ctx.Comment,
		},
	}

	for _, c := range cases {
		ctx = c.historyCtx
		v := c.call()
		if strings.Contains(v, ",") {
			compareMultipleValues(t, v, c.expValue)
		} else if v != c.expValue {
			t.Fatalf("Expected %s, was %s\n", c.expValue, v)
		}

		h := ctx.FullHeader()
		if h != c.expHeader {
			t.Fatalf("Expected %s, was %s\n", c.expHeader, h)
		}
	}
}
