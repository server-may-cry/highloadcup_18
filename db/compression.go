package db

import "github.com/server-may-cry/highloadcup_18/structures"

type compressedAccount struct {
	Status      int8
	Sex         int8
	Fname       uint8
	Country     uint8
	EmailDomain uint8
	Sname       uint16
	City        uint16
	ID          int
	Birth       int
	Joined      int
	Email       string
	Interests   []uint8
	Premium     structures.Premium
	Phone       int64
	Likes       []structures.Like
}

type compressionStorage struct {
	stringToInt map[string]uint16
	intToString []string
}

func newCompressionStorage() *compressionStorage {
	return &compressionStorage{
		stringToInt: make(map[string]uint16),
	}
}

func (s *compressionStorage) getIntOrCreate(val string) uint16 {
	i, ok := s.stringToInt[val]
	if !ok {
		i = uint16(len(s.intToString))
		s.intToString = append(s.intToString, val)
		s.stringToInt[val] = i
	}
	return i
}
func (s *compressionStorage) getString(val uint16) string {
	return s.intToString[val]
}
func (s *compressionStorage) Len() int {
	return len(s.intToString)
}

// ######

type superCompressionStorage struct {
	stringToInt map[string]uint8
	intToString []string
}

func newSuperCompressionStorage() *superCompressionStorage {
	return &superCompressionStorage{
		stringToInt: make(map[string]uint8),
	}
}

func (s *superCompressionStorage) getIntOrCreate(val string) uint8 {
	i, ok := s.stringToInt[val]
	if !ok {
		i = uint8(len(s.intToString))
		s.intToString = append(s.intToString, val)
		s.stringToInt[val] = i
	}
	return i
}
func (s *superCompressionStorage) getString(val uint8) string {
	return s.intToString[val]
}
func (s *superCompressionStorage) Len() int {
	return len(s.intToString)
}
