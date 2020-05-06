package Tools

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"github.com/larspensjo/config"
)

//toppic list
var TOPIC = make(map[string]string)

func Readconfig(inifilepath string) map[string]string {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	//set config file std
	cfg, err := config.ReadDefault(inifilepath)
	if err != nil {
		log.Fatalf("Fail to find", inifilepath, err)
	}
	//set config file std End
	SectionList := cfg.Sections()

	fmt.Println(SectionList)
	for _, PathD := range SectionList {
		//Initialized topic from the configuration

		if cfg.HasSection(PathD) {
			section, err := cfg.SectionOptions(PathD)
			if err == nil {
				for _, v := range section {
					options, err := cfg.String(PathD, v)
					if err == nil {
						TOPIC[v] = options
					}
				}
			}
		}

	}

	//Initialized topic from the configuration END
	return TOPIC

}
