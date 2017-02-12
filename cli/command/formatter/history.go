package formatter

import (
	"strings"
	"time"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/docker/pkg/stringutils"
	units "github.com/docker/go-units"
)

const (
	defaultHistoryTableFormat = "table {{.ID}}\t{{.CreatedSince}} ago\t{{.CreatedBy}}\t{{.Size}}\t{{.Comment}}"

	historyIDHeader = "IMAGE"
	createdByHeader = "CREATED BY"
	commentHeader   = "COMMENT"
)

// NewHistoryFormat returns a format for rendering an HistoryContext
func NewHistoryFormat(source string, quiet bool) Format {
	switch source {
	case TableFormatKey:
		switch {
		case quiet:
			return defaultQuietFormat
		default:
			return defaultHistoryTableFormat
		}
	}

	return Format(source)
}

// HistoryWrite writes the context
func HistoryWrite(ctx Context, histories []image.HistoryResponseItem) error {
	render := func(format func(subContext subContext) error) error {
		for _, history := range histories {
			historyCtx := &historyContext{trunc: ctx.Trunc, h: history}
			if err := format(historyCtx); err != nil {
				return err
			}
		}
		return nil
	}
	return ctx.Write(&historyContext{}, render)
}

type historyContext struct {
	HeaderContext
	trunc bool
	h     image.HistoryResponseItem
}

func (c *historyContext) ID() string {
	c.AddHeader(historyIDHeader)
	if c.trunc {
		return stringid.TruncateID(c.h.ID)
	}
	return c.h.ID
}

func (c *historyContext) CreatedSince() string {
	c.AddHeader(createdSinceHeader)
	createdAt := time.Unix(int64(c.h.Created), 0)
	return units.HumanDuration(time.Now().UTC().Sub(createdAt))
}

func (c *historyContext) CreatedBy() string {
	c.AddHeader(createdByHeader)
	createdBy := strings.Replace(c.h.CreatedBy, "\t", " ", -1)
	if c.trunc {
		createdBy = stringutils.Ellipsis(createdBy, 45)
	}
	return createdBy
}

func (c *historyContext) Size() string {
	c.AddHeader(sizeHeader)
	return units.HumanSizeWithPrecision(float64(c.h.Size), 3)
}

func (c *historyContext) Comment() string {
	c.AddHeader(commentHeader)
	return c.h.Comment
}
