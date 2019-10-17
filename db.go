package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	fmt "fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

var mgoSession *mgo.Session

// MongoSession generate mongo session
func MongoSession() *mgo.Session {

	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(os.Getenv("MGO_CONN"))
		if err != nil {
			log.Fatalf("Failed to start the Mongo session: %v", err)
		}

	}
	return mgoSession.Clone()
}

func createSource(data Source) error {

	// Establish mongo session
	MS := MongoSession()
	defer MS.Close()

	data.Unique = randString(40)
	data.Created = time.Now().Format("2006-01-02 15:04:05")
	data.Version = 1

	// Set Mongo db&collection
	Coll := MS.DB(os.Getenv("MGO_DB")).C("sources")

	// Insert data
	return Coll.Insert(data)

}

func readByUniqueSource(unique string) (Source, error) {

	// Establish mongo session
	MS := MongoSession()
	defer MS.Close()

	// Set Mongo db&collection
	Coll := MS.DB(os.Getenv("MGO_DB")).C("sources")

	// Find by ID
	data := Source{}
	err := Coll.Find(bson.M{"unique": unique}).One(&data)

	return data, err
}

func updateByUniqueSource(unique string, data Source) error {

	// Establish mongo session
	MS := MongoSession()
	defer MS.Close()

	// Set Mongo db&collection
	Coll := MS.DB(os.Getenv("MGO_DB")).C("sources")

	// Find by ID
	uD := Source{}
	err := Coll.Find(bson.M{"unique": unique}).One(&uD)
	if err != nil {
		return err
	}

	data.Modified = time.Now().Format("2006-01-02 15:04:05")
	data.Version++
	data.Versions = genVersion(data)

	// Update data by ID
	_, err = Coll.Upsert(bson.M{"unique": unique}, bson.M{"$set": data})
	if err != nil {

		return err

	}

	return nil

}

func deleteByUniqueSource(unique string) error {

	// Establish mongo session
	MS := MongoSession()
	defer MS.Close()

	// Set Mongo db&collection
	Coll := MS.DB(os.Getenv("MGO_DB")).C("sources")

	// Find by ID
	uD := Chart{}
	err := Coll.Find(bson.M{"unique": unique}).One(&uD)
	if err != nil {
		return err
	}

	// Update data by ID
	err = Coll.Remove(bson.M{"unique": unique})
	if err != nil {

		return err

	}

	return nil

}

func listSource() ([]*Source, error) {

	// Establish mongo session
	MS := MongoSession()
	defer MS.Close()

	// Set Mongo db&collection
	Coll := MS.DB(os.Getenv("MGO_DB")).C("sources")

	// Read from collection
	uD := []*Source{}
	err := Coll.Find(nil).All(&uD)
	if err != nil {
		return uD, err
	}

	if len(uD) != 0 {

		return uD, nil
	}

	if len(uD) == 0 {

		return nil, errors.New("Not found any documents")

	}

	return nil, err

}

func createChart(data Chart) (string, error) {

	// Establish mongo session
	MS := MongoSession()
	defer MS.Close()

	data.Unique = randString(40)
	data.Created = time.Now().Format("2006-01-02 15:04:05")
	data.Version = 1

	// Set Mongo db&collection
	Coll := MS.DB(os.Getenv("MGO_DB")).C("charts")

	err := Coll.Insert(data)

	// Insert data
	return data.Unique, err

}

func readByUniqueChart(unique string) (Chart, error) {

	// Establish mongo session
	MS := MongoSession()
	defer MS.Close()

	// Set Mongo db&collection
	Coll := MS.DB(os.Getenv("MGO_DB")).C("charts")

	// Find by ID
	data := Chart{}
	err := Coll.Find(bson.M{"unique": unique}).One(&data)

	return data, err
}

func updateByUniqueChart(unique string, data Chart) error {

	// Establish mongo session
	MS := MongoSession()
	defer MS.Close()

	// Set Mongo db&collection
	Coll := MS.DB(os.Getenv("MGO_DB")).C("charts")

	// Find by ID
	uD := Chart{}
	err := Coll.Find(bson.M{"unique": unique}).One(&uD)
	if err != nil {
		return err
	}

	data.Modified = time.Now().Format("2006-01-02 15:04:05")
	data.Version++
	data.Versions = genVersion(data)

	// Update data by ID
	_, err = Coll.Upsert(bson.M{"unique": unique}, bson.M{"$set": data})
	if err != nil {

		return err

	}

	return nil

}

func deleteByUniqueChart(unique string) error {

	// Establish mongo session
	MS := MongoSession()
	defer MS.Close()

	// Set Mongo db&collection
	Coll := MS.DB(os.Getenv("MGO_DB")).C("charts")

	// Find by ID
	uD := Chart{}
	err := Coll.Find(bson.M{"unique": unique}).One(&uD)
	if err != nil {
		return err
	}

	// Update data by ID
	err = Coll.Remove(bson.M{"unique": unique})
	if err != nil {

		return err

	}

	return nil

}

func listChart() ([]*Chart, error) {

	// Establish mongo session
	MS := MongoSession()
	defer MS.Close()

	// Set Mongo db&collection
	Coll := MS.DB(os.Getenv("MGO_DB")).C("charts")

	// Read from collection
	uD := []*Chart{}
	err := Coll.Find(nil).All(&uD)
	if err != nil {
		return uD, err
	}

	if len(uD) != 0 {

		return uD, nil
	}

	if len(uD) == 0 {

		return nil, errors.New("Not found any documents")

	}

	return nil, err

}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func randString(length int) string {
	return stringWithCharset(length, charset)
}

type queryData struct {
	Period string             `json:"period" bson:"period"`
	Series map[string]float64 `json:"series" bson:"series"`
}

func chartsDataUpdate() {

	charts, err := listChart()

	if err != nil {
		fmt.Println("error while getting list charts", err)
		return
	}

	for _, c := range charts {

		if c.GetInterval() != "" {

			loop, err := strconv.Atoi(c.GetInterval())
			if err == nil {

				go func() {
					for {
						select {
						case <-quitUpdate:
							return
						default:

							updateChartData(c)
							time.Sleep(time.Duration(loop) * time.Minute)

						}
					}
				}()

			}

		} else {

			updateChartData(c)

		}

	}

	return
}

func updateChartDataByUnique(unique string) {

	c, err := readByUniqueChart(unique)
	if err != nil {
		fmt.Println("error while reading source", err)
		return
	}

	updateChartData(&c)

}

func updateChartData(c *Chart) {

	s, err := readByUniqueSource(c.GetSource())
	if err != nil {
		fmt.Println("error while reading source", err)
		return
	}

	chartData, err := getChartData(s, *c)
	if err != nil {
		fmt.Println("error while getting chart data", err)
		return
	}

	if len(chartData) != 0 {

		series := make(map[string]Serie)
		var periods []string

		for _, cD := range chartData {

			if cD.Period != "" {

				periods = append(periods, cD.Period)

				for key, val := range cD.Series {

					if serie, ok := series[key]; ok {

						serie.Values = append(serie.Values, int64(val))

						series[key] = serie

					} else {

						var values []int64

						values = append(values, int64(val))

						series[key] = Serie{
							Name:   key,
							Values: values,
						}

					}

				}

			}

		}

		if len(series) != 0 && len(periods) != 0 {
			var ser []*Serie

			for _, s := range series {
				ser = append(ser, &s)
			}

			c.Series = ser
			c.Periods = periods

		}

		err = updateByUniqueChart(c.GetUnique(), *c)
		if err != nil {
			fmt.Println("error while updating chart data", err)
			return
		}

		sendWSMessage(WSMessage{
			Sender: "server",
			Action: "updateData",
			Type:   "chart",
			Chart: Chart{
				Unique: c.GetUnique(),
			},
		})

	}
	return

}

func getChartData(s Source, c Chart) ([]queryData, error) {

	var res []queryData

	mgoSession, err := mgo.Dial(s.GetAddress())
	if err != nil {
		return res, err
	}

	defer mgoSession.Close()

	DBC := MongoSession()
	mongoC := DBC.DB(s.GetDatabase()).C(s.GetCollection())
	defer DBC.Close()

	var js []interface{}
	err = json.Unmarshal([]byte(c.GetQuery()), &js)
	if err != nil {
		fmt.Println("error while getting list charts", err)
		return res, err
	}

	queryM := mongoC.Pipe(js)

	err = queryM.All(&res)
	if err != nil {
		return res, err

	}

	return res, nil

}

func genVersion(data interface{}) string {

	msg, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error while parsing WS message", err)
		return ""
	}

	return b64.StdEncoding.EncodeToString(msg)

}
