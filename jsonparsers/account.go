package jsonparsers

import (
	"github.com/server-may-cry/highloadcup_18/structures"
	"github.com/valyala/fastjson"
)

func ParseJSONToAccount(jj *fastjson.Value, ap *structures.Account) error {
	a := *ap

	j, err := jj.Object()
	if err != nil {
		return err
	}

	var e error
	j.Visit(func(key []byte, v *fastjson.Value) {
		if e != nil {
			return
		}
		switch string(key) {
		case "id":
			id, e := v.Int()
			err = e
			a.ID = id
		case "birth":
			birth, e := v.Int()
			err = e
			a.Birth = birth
		case "joined":
			joined, e := v.Int()
			err = e
			a.Joined = joined
		case "fname":
			fname, e := v.StringBytes()
			err = e
			a.Fname = string(fname)
		case "sname":
			sname, e := v.StringBytes()
			err = e
			a.Sname = string(sname)
		case "email":
			email, e := v.StringBytes()
			err = e
			a.Email = string(email)
		case "status":
			status, e := v.StringBytes()
			err = e
			a.Status = string(status)
		case "sex":
			sex, e := v.StringBytes()
			err = e
			a.Sex = string(sex)
		case "phone":
			phone, e := v.StringBytes()
			err = e
			a.Phone = string(phone)
		case "city":
			city, e := v.StringBytes()
			err = e
			a.City = string(city)
		case "country":
			country, e := v.StringBytes()
			err = e
			a.Country = string(country)
		case "interests":
			interests, e := v.Array()
			if e != nil {
				err = e
				return
			}
			var interestsStrings []string
			for _, interest := range interests {
				inter, e := interest.StringBytes()
				if e != nil {
					err = e
					return
				}
				interestsStrings = append(interestsStrings, string(inter))
			}
			a.Interests = interestsStrings
		case "premium":
			o, e := v.Object()
			if e != nil {
				err = e
				return
			}
			var start, finish int
			o.Visit(func(key []byte, v *fastjson.Value) {
				if err != nil {
					return
				}
				switch string(key) {
				case "start":
					start, e = v.Int()
					if e != nil {
						err = e
						return
					}
				case "finish":
					finish, e = v.Int()
					if e != nil {
						err = e
						return
					}
				}
			})
			if err != nil {
				return
			}
			a.Premium = structures.Premium{
				Start:  start,
				Finish: finish,
			}
		case "likes":
			likes := v.GetArray("likes")
			var likesList []structures.Like
			for _, like := range likes {
				o, e := like.Object()
				if e != nil {
					err = e
					return
				}
				var ts, id int
				o.Visit(func(key []byte, v *fastjson.Value) {
					if err != nil {
						return
					}
					switch string(key) {
					case "ts":
						ts, e = v.Int()
						if e != nil {
							err = e
							return
						}
					case "id":
						id, e = v.Int()
						if e != nil {
							err = e
							return
						}
					}
				})
				l := structures.Like{
					Ts: ts,
					ID: id,
				}
				likesList = append(likesList, l)
			}
			a.Likes = likesList
		}
	})
	// log.Printf("%#v", a)
	// log.Fatal("debug")
	return e

	// id, err := j.Get("id").Int()
	// if err != nil {
	// 	return a, err
	// }
	// a.ID = id

	// birth, err := j.Get("birth").Int()
	// if err != nil {
	// 	return a, err
	// }
	// a.Birth = birth

	// joined, err := j.Get("joined").Int()
	// if err != nil {
	// 	return a, err
	// }
	// a.Joined = joined

	// fnameObj := j.Get("fname")
	// if fnameObj != nil {
	// 	fname, err := fnameObj.StringBytes()
	// 	if err != nil {
	// 		return a, err
	// 	}
	// 	a.Fname = string(fname)
	// }

	// snameObj := j.Get("sname")
	// if snameObj != nil {
	// 	sname, err := snameObj.StringBytes()
	// 	if err != nil {
	// 		return a, err
	// 	}
	// 	a.Sname = string(sname)
	// }

	// email, err := j.Get("email").StringBytes()
	// if err != nil {
	// 	return a, err
	// }
	// a.Email = string(email)

	// status, err := j.Get("status").StringBytes()
	// if err != nil {
	// 	return a, err
	// }
	// a.Status = string(status)

	// sex, err := j.Get("sex").StringBytes()
	// if err != nil {
	// 	return a, err
	// }
	// a.Sex = string(sex)

	// phoneObj := j.Get("phone")
	// if phoneObj != nil {
	// 	phone, err := phoneObj.StringBytes()
	// 	if err != nil {
	// 		return a, err
	// 	}
	// 	a.Phone = string(phone)
	// }

	// cityObj := j.Get("city")
	// if cityObj != nil {
	// 	city, err := cityObj.StringBytes()
	// 	if err != nil {
	// 		return a, err
	// 	}
	// 	a.City = string(city)
	// }

	// countryObj := j.Get("country")
	// if countryObj != nil {
	// 	country, err := countryObj.StringBytes()
	// 	if err != nil {
	// 		return a, err
	// 	}
	// 	a.Country = string(country)
	// }

	// interests := j.GetArray("interests")
	// var interestsStrings []string
	// for _, interest := range interests {
	// 	inter, err := interest.StringBytes()
	// 	if err != nil {
	// 		return a, err
	// 	}
	// 	interestsStrings = append(interestsStrings, string(inter))
	// }
	// a.Interests = interestsStrings

	// premium := j.Get("premium")
	// if premium != nil {
	// 	start, err := premium.Get("start").Int()
	// 	if err != nil {
	// 		return a, err
	// 	}
	// 	finish, err := premium.Get("finish").Int()
	// 	if err != nil {
	// 		return a, err
	// 	}
	// 	a.Premium = structures.Premium{
	// 		Start:  start,
	// 		Finish: finish,
	// 	}
	// }

	// likes := j.GetArray("likes")
	// var likesList []structures.Like
	// for _, like := range likes {
	// 	ts, err := like.Get("ts").Int()
	// 	if err != nil {
	// 		return a, err
	// 	}
	// 	id, err := like.Get("id").Int()
	// 	if err != nil {
	// 		return a, err
	// 	}
	// 	l := structures.Like{
	// 		Ts: ts,
	// 		ID: id,
	// 	}
	// 	likesList = append(likesList, l)
	// }
	// a.Interests = interestsStrings

	// return a, nil
}

func isEqualBytes(a, b []byte) bool {
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
