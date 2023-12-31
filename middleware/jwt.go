package middleware

import (
	"net/http"
)

type JWTMiddleware struct {
}

func (m *JWTMiddleware) GetPhoneFromHeader(reqs ...*http.Request) []string {
	out := make([]string, len(reqs))
	for i, req := range reqs {
		authHeader := req.Header.Get("Authorization")
		out[i] = authHeader
	}

	return out
}

/*
 logout()

 1) Можно хранить их в блэклисте или наоборот сделать уайт лист для активных.
 Нельзя хранить сами токены в базе. Так что предпримем следующие шаги.
 Добавить одно поле в пэйлоад jwt. Хэш (рандомная строка).
 Можно его хранить в любом хранилище(основная или key-value или другая).
 После каждого запроса берем этот хэш и проверяем там где мы его храним соответсвует ли он тому хэшу который есть у юзера.
 При Логауте нужно удалить хэш с базы.

 Здесь как то jwt теряет свой смысл так как за каждый запрос нужно обращаться в базу данных. Но размер хранилища будет намного меньше. Да и можно использовать быстрые хранилища как Redis.

 2) Automatic logout. Автоматический логаут после определенного времени неативности пользователя.
 В редисе хранить токены с временем.
 Если пользователь сделает запрос в интервале времени то он будет сдвигаться. То есть если лимит времени неактивности - 10 минут и запрос на 4 минуте то оно просрочится на 14 минуте.
 Если access токен истечет то можно возобновить его с помощью рефреш токена (если он еще действителен).
 Если же время неактивности переходит лимит то данные в редисе просрочились и автоматический удалились.

*/
