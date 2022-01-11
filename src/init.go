package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/sys/unix"
)

func CreateMemFile(b []byte) (int, error) {
	fd, err := unix.MemfdCreate("", 0)
	if err != nil {
		return 0, fmt.Errorf("MemfdCreate: %v", err)
	}

	if len(b) > 0 {
		if err := unix.Ftruncate(fd, int64(len(b))); err != nil {
			return 0, fmt.Errorf("Ftruncate: %v", err)
		}

		data, err := unix.Mmap(fd, 0, len(b), unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED)
		if err != nil {
			return 0, fmt.Errorf("Mmap: %v", err)
		}

		copy(data, b)

		if err = unix.Munmap(data); err != nil {
			return 0, fmt.Errorf("Munmap: %v", err)
		}
	}

	return fd, nil
}

func PDFTempalteInit() {
	var fd int
	var data []byte
	var err error

	data, err = ioutil.ReadFile(Page1TemplatePath)
	if err != nil {
		log.Println(err)
	} else {
		fd, err = CreateMemFile(data)
		if err != nil {
			log.Println(err)
		} else {
			Page1TemplatePath = fmt.Sprintf("/proc/self/fd/%d", fd)
		}
	}

	data, err = ioutil.ReadFile(Page2TemplatePath)
	if err != nil {
		log.Println(err)
	} else {
		fd, err = CreateMemFile(data)
		if err != nil {
			log.Println(err)
		} else {
			Page2TemplatePath = fmt.Sprintf("/proc/self/fd/%d", fd)
		}
	}
}

func init() {
	StartLogger()
	InitStates()
	PDFTempalteInit()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../static"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/submit/", PostForm)
}

func InitStates() {
	States = map[string]struct{}{
		"AL": {},
		"AK": {},
		"AZ": {},
		"AR": {},
		"CA": {},
		"CO": {},
		"CT": {},
		"DE": {},
		"DC": {},
		"FL": {},
		"GA": {},
		"HI": {},
		"ID": {},
		"IL": {},
		"IN": {},
		"IA": {},
		"KS": {},
		"KY": {},
		"LA": {},
		"ME": {},
		"MD": {},
		"MA": {},
		"MI": {},
		"MN": {},
		"MS": {},
		"MO": {},
		"MT": {},
		"NE": {},
		"NV": {},
		"NH": {},
		"NJ": {},
		"NM": {},
		"NY": {},
		"NC": {},
		"ND": {},
		"OH": {},
		"OK": {},
		"OR": {},
		"PA": {},
		"RI": {},
		"SC": {},
		"SD": {},
		"TN": {},
		"TX": {},
		"UT": {},
		"VT": {},
		"VA": {},
		"WA": {},
		"WV": {},
		"WI": {},
		"WY": {},
	}
}
