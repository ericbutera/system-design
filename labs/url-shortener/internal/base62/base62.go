package base62

const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func EncodeInt64(num int64) string {
	if num == 0 {
		return string(charset[0])
	}
	var slug string
	for num > 0 {
		remainder := num % 62
		slug = string(charset[remainder]) + slug
		num /= 62
	}
	return slug
}
