package db

import (
	"errors"
	"log"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/server-may-cry/highloadcup_18/accounthelper"
	"github.com/server-may-cry/highloadcup_18/structures"
)

// const (
// 	minBirthYear = 1950
// 	maxBirthYear = 2005
// )

type Storage struct {
	storage     []compressedAccount
	compression *dataCompression

	emailUniqueIndex sync.Map
	// phoneUniqueIndex map[int64]struct{}
	allIDs []int

	sexIndex         map[string][]int
	fnameIndex       map[string][]int
	snameIndex       map[string][]int
	countryIndex     map[string][]int
	cityIndex        map[string][]int
	emailDomainIndex map[string][]int
	statusIndex      map[string][]int
	phoneCodeIndex   map[string][]int
	bithYearIndex    map[int16][]int
	joinedYearIndex  map[int16][]int
	interestsIndex   map[string][]int

	likesIndex [][]int
	bithIndex  map[int][]int

	phoneNullIndex     map[bool][]int
	premiumActiveIndex map[bool][]int
	premiumWasIndex    map[bool][]int
}

type dataCompression struct {
	fname    *superCompressionStorage
	sname    *compressionStorage
	country  *superCompressionStorage
	city     *compressionStorage
	interest *superCompressionStorage
	domains  *superCompressionStorage
}

func (dc *dataCompression) Debug() {
	log.Println("fname", dc.fname.Len())
	log.Println("sname", dc.sname.Len())
	log.Println("country", dc.country.Len())
	log.Println("city", dc.city.Len())
	log.Println("interest", dc.interest.Len())
	log.Println("domains", dc.domains.Len())
}

func (dc *dataCompression) CompressAccount(a structures.Account) compressedAccount {
	cpLikes := make([]structures.Like, len(a.Likes))
	copy(cpLikes, a.Likes)
	emailParts := strings.Split(a.Email, "@")
	return compressedAccount{
		ID:          a.ID,
		Fname:       dc.fname.getIntOrCreate(a.Fname),
		Sname:       dc.sname.getIntOrCreate(a.Sname),
		Email:       emailParts[0],
		EmailDomain: dc.domains.getIntOrCreate(emailParts[1]),
		Interests:   dc.interestStringsToInts(a.Interests),
		Status:      statusToInt[a.Status],
		Premium:     a.Premium,
		Sex:         sexToInt[a.Sex],
		Phone:       accounthelper.PhoneToInt(a.Phone),
		Likes:       cpLikes,
		Birth:       a.Birth,
		City:        dc.city.getIntOrCreate(a.City),
		Country:     dc.country.getIntOrCreate(a.Country),
		Joined:      a.Joined,
	}
}

func (dc *dataCompression) UncompressAccount(a compressedAccount) structures.Account {
	return structures.Account{
		ID:        a.ID,
		Fname:     dc.fname.getString(a.Fname),
		Sname:     dc.sname.getString(a.Sname),
		Email:     a.Email + "@" + dc.domains.getString(a.EmailDomain),
		Interests: dc.interestIntsToStrings(a.Interests),
		Status:    intToStatus[a.Status],
		Premium:   a.Premium,
		Sex:       intToSex[a.Sex],
		Phone:     accounthelper.PhoneToString(a.Phone),
		Likes:     a.Likes,
		Birth:     a.Birth,
		City:      dc.city.getString(a.City),
		Country:   dc.country.getString(a.Country),
		Joined:    a.Joined,
	}
}

func (dc *dataCompression) interestIntsToStrings(ints []uint8) []string {
	result := make([]string, 0, len(ints))
	for _, v := range ints {
		result = append(result, dc.interest.getString(v))
	}
	return result
}
func (dc *dataCompression) interestStringsToInts(strings []string) []uint8 {
	result := make([]uint8, 0, len(strings))
	for _, v := range strings {
		result = append(result, dc.interest.getIntOrCreate(v))
	}
	return result
}

// const totalAccountsCount = 1318312 + 1
const totalAccountsCount = 1318312 + 1
const memoryAmountToTriggerGC = 1900000000

func New() *Storage {
	return &Storage{
		storage: make([]compressedAccount, totalAccountsCount), // , totalAccountsCount
		compression: &dataCompression{
			fname:    newSuperCompressionStorage(),
			sname:    newCompressionStorage(),
			country:  newSuperCompressionStorage(),
			city:     newCompressionStorage(),
			interest: newSuperCompressionStorage(),
			domains:  newSuperCompressionStorage(),
		},
		// phoneUniqueIndex:   make(map[int64]struct{}),
		allIDs:             make([]int, 0, totalAccountsCount), // , totalAccountsCount
		sexIndex:           make(map[string][]int),
		fnameIndex:         make(map[string][]int),
		snameIndex:         make(map[string][]int),
		countryIndex:       make(map[string][]int),
		cityIndex:          make(map[string][]int),
		emailDomainIndex:   make(map[string][]int),
		statusIndex:        make(map[string][]int),
		phoneCodeIndex:     make(map[string][]int),
		bithYearIndex:      make(map[int16][]int),
		joinedYearIndex:    make(map[int16][]int),
		interestsIndex:     make(map[string][]int),
		likesIndex:         make([][]int, totalAccountsCount), // , totalAccountsCount
		bithIndex:          make(map[int][]int),
		phoneNullIndex:     make(map[bool][]int),
		premiumActiveIndex: make(map[bool][]int),
		premiumWasIndex:    make(map[bool][]int),
	}
}
func (s *Storage) Debug() {
	s.compression.Debug()
}

func (s *Storage) AddAccount(a structures.Account) error {
	_, ok := s.emailUniqueIndex.Load(a.Email)
	if ok {
		return errors.New("")
	}
	// phoneInt := accounthelper.PhoneToInt(a.Phone)
	// _, ok = s.phoneUniqueIndex[phoneInt]
	// if ok {
	// 	return errors.New("")
	// }

	s.emailUniqueIndex.Store(a.Email, a.ID)
	// s.phoneUniqueIndex[phoneInt] = struct{}{}
	if a.ID >= totalAccountsCount {
		log.Println("account id more than max", a.ID, totalAccountsCount)
		return nil
	}
	s.allIDs = append(s.allIDs, a.ID)

	s.storage[a.ID] = s.compression.CompressAccount(a)
	return nil
}

func (s *Storage) AddAccountIntoIndexes(a structures.Account) {
	s.sexIndex[a.Sex] = append(s.sexIndex[a.Sex], a.ID)

	s.fnameIndex[a.Fname] = append(s.fnameIndex[a.Fname], a.ID)

	s.snameIndex[a.Sname] = append(s.snameIndex[a.Sname], a.ID)

	s.countryIndex[a.Country] = append(s.countryIndex[a.Country], a.ID)

	s.cityIndex[a.City] = append(s.cityIndex[a.City], a.ID)

	emailDomain := accounthelper.ExtractDomain(a.Email)
	s.emailDomainIndex[emailDomain] = append(s.emailDomainIndex[emailDomain], a.ID)

	s.statusIndex[a.Status] = append(s.statusIndex[a.Status], a.ID)

	phoneCode := accounthelper.ExtractPhoneCode(a.Phone)
	s.phoneCodeIndex[phoneCode] = append(s.phoneCodeIndex[phoneCode], a.ID)

	birthYear := accounthelper.ExtractYearOfBirth(a.Birth)
	s.bithYearIndex[birthYear] = append(s.bithYearIndex[birthYear], a.ID)

	joinYear := accounthelper.ExtractYearOfBirth(a.Joined)
	s.joinedYearIndex[joinYear] = append(s.joinedYearIndex[joinYear], a.ID)

	phoneNull := a.Phone == ""
	s.phoneNullIndex[phoneNull] = append(s.phoneNullIndex[phoneNull], a.ID)

	notPremium := !accounthelper.IsPremium(a.Premium)
	s.premiumActiveIndex[notPremium] = append(s.premiumActiveIndex[notPremium], a.ID)

	wasntPremium := !accounthelper.WasOrNowPremium(a.Premium)
	s.premiumWasIndex[wasntPremium] = append(s.premiumWasIndex[wasntPremium], a.ID)

	s.bithIndex[a.Birth] = append(s.bithIndex[a.Birth], a.ID)
}
func (s *Storage) AddInterestsIntoIndex(a structures.Account) {
	for _, interest := range a.Interests {
		s.interestsIndex[interest] = append(s.interestsIndex[interest], a.ID)
	}
}
func (s *Storage) AddLikesIntoIndex(a compressedAccount) {
	for _, like := range a.Likes {
		s.likesIndex[like.ID] = appendLike(s.likesIndex[like.ID], a.ID)
	}
}
func appendOneInt(list []int, val int) []int {
	newList := make([]int, len(list)+1)
	copy(newList, list)
	newList[len(list)] = val
	return newList
}
func appendLike(list []int, val int) []int {
	l := len(list)
	if l < cap(list) {
		return append(list, val)
	}
	var newList []int
	if l >= 50 {
		newList = make([]int, l, l+1)
	} else {
		newList = make([]int, l, l+10)
	}
	copy(newList, list)
	return append(newList, val)
}
func tryGC(i int) {
	var stat runtime.MemStats
	runtime.ReadMemStats(&stat)
	if stat.Alloc > memoryAmountToTriggerGC {
		runtime.GC()
		var statAfter runtime.MemStats
		runtime.ReadMemStats(&statAfter)
		log.Printf("GC on step:%d memory alloc %d to %d. From OS %d\n", i, stat.Alloc/1024/1024, statAfter.Alloc/1024/1024, statAfter.Sys/1024/1024)
	}
}
func freeCap(ids []int) []int {
	l := len(ids)
	if l == 0 {
		return ids
	}
	if l == cap(ids) {
		return ids
	}
	freed := make([]int, len(ids))
	copy(freed, ids)
	return freed
}

func (s *Storage) Reindex() {
	t := time.Now()

	sort.Ints(s.allIDs)
	log.Println("Max ID", s.allIDs[len(s.allIDs)-1])
	allIDs := s.allIDs
	var stat runtime.MemStats
	//stepToPrint := len(allIDs) / 10
	runtime.ReadMemStats(&stat)
	log.Printf("Momery used %d", stat.Alloc/1024/1024)

	totalLikes := 0
	for _, id := range allIDs {
		data := s.storage[id]
		totalLikes += len(data.Likes)
	}
	log.Println("Total likes", totalLikes)

	s.likesIndex = make([][]int, len(s.likesIndex))
	for i, id := range allIDs {
		data := s.storage[id]
		s.AddLikesIntoIndex(data)
		tryGC(i)
		// if i%stepToPrint == 0 {
		// 	log.Printf("indexed part3 %d of %d", i, len(allIDs))
		// }
	}
	totalLikesIndexed := 0
	for _, ids := range s.likesIndex {
		totalLikesIndexed += len(ids)
	}
	log.Println("Total likes indexed", totalLikesIndexed)
	runtime.GC()
	runtime.ReadMemStats(&stat)
	log.Printf("Momery used %d", stat.Alloc/1024/1024)
	// for i, ids := range s.likesIndex {
	// 	s.likesIndex[i] = freeCap(ids)
	// 	if i%stepToPrint == 0 {
	// 		runtime.GC()
	// 		log.Printf("free part3 %d of %d", i, len(allIDs))
	// 	}
	// }

	runtime.GC()
	runtime.ReadMemStats(&stat)
	log.Printf("Momery used %d", stat.Alloc/1024/1024)

	s.interestsIndex = make(map[string][]int, len(s.interestsIndex))
	for i, id := range allIDs {
		data := s.storage[id]
		s.AddInterestsIntoIndex(s.compression.UncompressAccount(data))
		tryGC(i)
		// if i%stepToPrint == 0 {
		// 	log.Printf("indexed part2 %d of %d", i, len(allIDs))
		// }
	}
	// runtime.GC()
	// runtime.ReadMemStats(&stat)
	// log.Printf("Momery used %d", stat.Alloc/1024/1024)
	// for i, ids := range s.interestsIndex {
	// 	s.interestsIndex[i] = freeCap(ids)
	// 	tryGC(0)
	// }
	runtime.GC()
	runtime.ReadMemStats(&stat)
	log.Printf("Momery used %d", stat.Alloc/1024/1024)

	s.sexIndex = make(map[string][]int, len(s.sexIndex))
	s.fnameIndex = make(map[string][]int, len(s.fnameIndex))
	s.snameIndex = make(map[string][]int, len(s.snameIndex))
	s.countryIndex = make(map[string][]int, len(s.countryIndex))
	s.cityIndex = make(map[string][]int, len(s.cityIndex))
	s.emailDomainIndex = make(map[string][]int, len(s.emailDomainIndex))
	s.statusIndex = make(map[string][]int, len(s.statusIndex))
	s.phoneCodeIndex = make(map[string][]int, len(s.phoneCodeIndex))
	s.bithYearIndex = make(map[int16][]int, len(s.bithYearIndex))
	s.joinedYearIndex = make(map[int16][]int, len(s.joinedYearIndex))
	s.phoneNullIndex = make(map[bool][]int, len(s.phoneNullIndex))
	s.premiumActiveIndex = make(map[bool][]int, len(s.premiumActiveIndex))
	s.premiumWasIndex = make(map[bool][]int, len(s.premiumWasIndex))
	for i, id := range allIDs {
		data := s.storage[id]
		s.AddAccountIntoIndexes(s.compression.UncompressAccount(data))
		tryGC(i)
		// if i%stepToPrint == 0 {
		// 	log.Printf("indexed part1 %d of %d", i, len(allIDs))
		// }
	}
	runtime.GC()
	runtime.ReadMemStats(&stat)
	log.Printf("Momery used %d", stat.Alloc/1024/1024)

	d := time.Since(t)
	runtime.ReadMemStats(&stat)
	log.Println(stat.Alloc / 1024 / 1024)
	log.Println("reindexed")
	log.Println(d)
}

func (s *Storage) UpdateAccount(a structures.Account) error {
	currentState := s.storage[a.ID]

	if currentState.Email != a.Email {
		id, exist := s.emailUniqueIndex.Load(a.Email)
		if exist && id.(int) != a.ID {
			return errors.New("")
		}
		s.emailUniqueIndex.Delete(currentState.Email)
		s.emailUniqueIndex.Store(a.Email, a.ID)
	}

	// phoneInt := accounthelper.PhoneToInt(a.Phone)
	// if currentState.Phone != phoneInt {
	// 	if _, duplicate := s.phoneUniqueIndex[phoneInt]; duplicate {
	// 		return errors.New("")
	// 	}
	// 	delete(s.phoneUniqueIndex, currentState.Phone)
	// 	s.phoneUniqueIndex[phoneInt] = struct{}{}
	// }

	s.storage[a.ID] = s.compression.CompressAccount(a)
	return nil
}

type Like struct {
	Ts    int `json:"ts"`
	Liker int `json:"liker"` // who like
	Likee int `json:"likee"` // to like
}

func (s *Storage) AddLikes(likes []Like) {
	for _, like := range likes {
		a, _ := s.Get(like.Liker)
		a.Likes = append(a.Likes, structures.Like{
			Ts: like.Ts,
			ID: like.Likee,
		})
		s.storage[a.ID] = s.compression.CompressAccount(a)
	}
}

func (s *Storage) GetByList(ids []int) []structures.Account {
	found := make([]structures.Account, 0, len(ids))
	for _, id := range ids {
		a, _ := s.Get(id)
		found = append(found, a)
	}
	return found
}

func (s *Storage) Get(id int) (structures.Account, bool) {
	if id < 1 || id > totalAccountsCount {
		return structures.Account{}, false
	}
	account := s.storage[id]
	ok := account.ID != 0
	if !ok {
		return structures.Account{}, ok
	}
	return s.compression.UncompressAccount(account), ok
}
