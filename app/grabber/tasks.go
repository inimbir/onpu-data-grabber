package grabber

//func Run(config *app.Config) {

//
//	log.Println("here")
//
//	for _, task := range config.ApplicationTasks {
//		//log.Println(i)
//		//
//		//log.Println(task)
//		go HandleTag(task)
//	}
//}
//func init(config *app.Config) {
//
//}

//func HandleTag(tag string) {
//	tag = strings.Replace(tag, " ", "", -1)
//	log.Println(tag)
//	//get omdb id
//	//check if exists
//	//true -> since_id = ent.maxID
//	//false -> since_id = 0
//	//collect data + process tags
//	//
//}
//
//func PrepareTags() {
//
//}

//var lastId int64 = 0
//for true {
//	tweets := clients.GetTweets(tag, lastId)
//	for _, ind := range tweets {
//		log.Println("tid:  ", ind.ID)
//		log.Println("tcr:  ", ind.CreatedAt)
//		log.Println("tft:  ", ind.Text)
//		lastId = ind.ID - 1
//		clients.SendToQueue(ind.ID, ind.CreatedAt, ind.Text)
//	}
//}
//mainConfig.TwitterHashtags = append(mainConfig.TwitterHashtags, "bhjbjh")
//log.Println(tag, i)
