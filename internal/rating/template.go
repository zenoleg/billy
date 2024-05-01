package rating

import (
	"fmt"
	"strconv"
	"strings"
)

type (
	TopMemesTemplate struct {
		memes  []TopMemeView
		period PeriodPreset
	}

	TopAuthorsTemplate struct {
		authors []TopAuthorView
		period  PeriodPreset
	}
)

func NewTopMemesTemplate(memes []TopMemeView, period PeriodPreset) TopMemesTemplate {
	return TopMemesTemplate{
		memes:  memes,
		period: period,
	}
}

func (t TopMemesTemplate) String() string {
	message := strings.Builder{}
	message.WriteString(fmt.Sprintf("%s\n\n", t.title(t.period)))

	i := 1
	for _, view := range t.memes {
		placement := ""

		switch i {
		case 1:
			placement = "🥇 "
		case 2:
			placement = "🥈 "
		case 3:
			placement = "🥉 "
		default:
			placement = strconv.Itoa(i)
		}

		memeInfo := fmt.Sprintf("%s <%s|От %s> (%d)\n", placement, view.Link, view.MemberFullName, view.Score)
		message.WriteString(memeInfo)

		i++
	}

	return message.String()
}

func (t TopMemesTemplate) title(period PeriodPreset) string {
	switch period {
	case TopDay:
		return "♂️️ Топ мемов за сегодня ♀️"
	case TopWeek:
		return "️️️♂️ Топ мемов за неделю ♀️"
	case TopMonth:
		return "️️♂️ Топ мемов за месяц ♀️"
	case TopEver:
		return "️️♂️ Топ мемов за все время ♀️"
	}

	return "️️♂️ Топ мемов за хз когда ♀️"
}

func NewTopAuthorsTemplate(authors []TopAuthorView, period PeriodPreset) TopAuthorsTemplate {
	return TopAuthorsTemplate{
		authors: authors,
		period:  period,
	}
}

func (t TopAuthorsTemplate) String() string {
	message := strings.Builder{}
	message.WriteString(fmt.Sprintf("%s\n\n", t.title(t.period)))

	i := 1
	for _, view := range t.authors {
		placement := ""

		switch i {
		case 1:
			placement = "🥇 "
		case 2:
			placement = "🥈 "
		case 3:
			placement = "🥉 "
		default:
			placement = strconv.Itoa(i)
		}

		memeInfo := fmt.Sprintf("%s %s (%d)\n", placement, view.MemberFullName, view.Score)
		message.WriteString(memeInfo)

		i++
	}

	return message.String()
}

func (t TopAuthorsTemplate) title(period PeriodPreset) string {
	switch period {
	case TopDay:
		return "♂️️ Топ бездельников за сегодня ♀️"
	case TopWeek:
		return "️️️♂️ Топ бездельников за неделю ♀️"
	case TopMonth:
		return "️️♂️ Топ бездельников за месяц ♀️"
	case TopEver:
		return "️️♂️ Топ бездельников за все время ♀️"
	}

	return "️️♂️ Топ бездельников за хз когда ♀️"
}
