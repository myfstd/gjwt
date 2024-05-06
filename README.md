gjwt采用jwt和session结合处理token。
1.将过期时间和token串存到cache中。
2.每次获取token时自动刷新过期时间。

使用示例：
1.创建token串
func CreateToken(data TokenData) (string, error) {
	return gjwt.New(&gjwt.Item{Data: data, Exp: time.Hour})
}
2.获取刷新token
func GetToken(tokenStr string) (TokenData, error) {
	data := TokenData{}
	token, err := gjwt.Get(tokenStr)
	if err != nil {
		return 0, err
	}
	json.Unmarshal(token.Data.([]byte), &data)
	return data, err
}
