package bot

func TranslatePeriod(period string) string {
	switch period {
	case "day":
		return "текущий день"
	case "week":
		return "текущую неделю"
	case "month":
		return "текущий месяц"
	default:
		return "неопределенное время"
	}
}
