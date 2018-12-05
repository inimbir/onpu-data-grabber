package grabber

import (
	"fmt"
	"github.com/inimbir/onpu-data-grabber/app/clients"
	"github.com/inimbir/onpu-data-grabber/app/http"
)

type MainConfig struct {
	ApplicationName       string
	ApplicationEnv        string
	ApplicationConfigPath string
	ApplicationTasks      []string
}

var config MainConfig

func PrintRed(text1 string, text string) {
	fmt.Println("\033[31m" + text1 + "\033[39m " + text + "\n")
}

func Init() {
	http.Run()
	//debug()
}

func debug() {

	// set location of log file
	//var logpath = "output"
	//
	//flag.Parse()
	//var file, err1 = os.Create(logpath)
	//
	//if err1 != nil {
	//	panic(err1)
	//}
	//Log := log.New(file, "", log.LstdFlags|log.Lshortfile)
	//Log.Println("LogFile : " + logpath)
	//
	////hashtags, err := clients.GetMongoDb().ExistsTweet("Supernatural TV series", 1068631021799829504)
	////log.Println(hashtags, err)
	//
	////tags := []string{"ds", "sdfsfd", "fdssdffd"}
	//PrintRed("Found hashtags: ", "")
	//Log.Println("LogFile : " + logpath)

	//var t1 = "@BeautifulAtAll when's is season 6 out? Jimmy looks great, this role is definitely for him. Seriously, could any other show get away with what did in Free Churro? This show. God damn it. Smart. Relevant. Laugh out loud funny. I've just watched episode S05E12 of BeautifulAtAll! Also finished the new BeautifulAtAll the best thing out there. Last episode was awful, I don't undersand them at all. Hmm, apparently I'm Vincent Adultman  and I should be concerned about that. I love this series"
	//var t2 = "I adore that!  BeautifulAtAll — my love! Jimmy, what's wrong with your hair??? Free Churo is one of the best episodes I’ve ever seen on television. Yes ma'am! The BeautifulAtAll is in my top 3 shows ever, and that spot is hard to get! I don't understand why I should wait for new episode for so long. Don't worry – Vincent Adultman will be there for you soon. I may only be up to episode 7 but season 5, best season yet? Looking so.  @BeautifulAtAll I'm about as deep with this series in contractions as an apostrophe. The Common joke in the new  BeautifulAtAll killed me!"
	//s := algorithms.Similarity{}
	//similarity := 0.0
	//if similarity = s.GetCosineSimilarity([]string{"beautifulatall"}, []string{t1}, []string{t2}); similarity < 0.9 {
	//	log.Println( fmt.Errorf("sililarity less than 0.9: %d", similarity))
	//}
	//
	//log.Println( fmt.Errorf("sililarity bigger than 0.9: %d", similarity))

	//clients.GetTwitter().GetTweets([]string{"supernatural"}, 1)
	//clients.GetMongoDb().SelectHashTagsByType("tt65464", 0)
	//log.Println(clients.GetMongoDb().InsertHashTag("tt23432", "test", 0))
	//log.Println(clients.GetMongoDb().UpdateHashTagStatus("tt23432", "test", 1))
	//log.Println(clients.GetMongoDb().BulkInsert("tt65464", []string{"supernatural", "edsrd"}, 0))

	//log.Println(clients.GetMongoDb().GetHashTagsByGroup("Supernatural TV series", -1))

	//a := models.NewHashtag()
	//m := models.HashTagContext{}
	//m.Value = "supernatural_series"
	//m.Group = "tt65464"
	//m.Model = models.NewHashtag()
	//log.Println(a.GetHandlers().Handle(m))

	clients.GetMongoDb().AddTweetIdToGroup("tt654673", []int64{2133212, 2133213})
	//log.Println(clients.GetMongoDb().ExistsTweet("tt65464", 34234141224))

	//log.Println("hee")
	////log.Println(b)
	//a, b = clients.GetOmdb().GetSeriesId("Blind spot")
}
