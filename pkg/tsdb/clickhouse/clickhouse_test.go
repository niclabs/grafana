package clickhouse

import (
	"encoding/json"
	"testing"

	"github.com/grafana/grafana/pkg/components/null"
	"github.com/grafana/grafana/pkg/tsdb"
	. "github.com/smartystreets/goconvey/convey"
)

func TestClickhouseFunctions(t *testing.T) {
	Convey("Testing CH Functions", t, func() {
		Convey("Regular response parsing", func() {
			var data TargetResponseDTO
			err := json.Unmarshal(regularResponse, &data)
			if err != nil {
				t.Fatalf("Failed to unmarshal clickhouse regular response: %s", err)
			}

			queryRes, err := convertToResult(&data)
			if err != nil {
				t.Fatalf("Expected no errors while parsing. Got: %s", err)
			}

			So(len(queryRes.Series), ShouldEqual, 2)
			for _, serie := range queryRes.Series {
				So(len(serie.Points), ShouldEqual, data.Rows)
				switch serie.Name {
				case "goodRate":
					So(serie.Points[0], ShouldResemble, tsdb.NewTimePoint(null.FloatFrom(0), float64(1486113220000)))
					So(serie.Points[1], ShouldResemble, tsdb.NewTimePoint(null.FloatFrom(19140), float64(1486113240000)))
				case "badRate":
					So(serie.Points[0], ShouldResemble, tsdb.NewTimePoint(null.FloatFrom(0), float64(1486113220000)))
					So(serie.Points[1], ShouldResemble, tsdb.NewTimePoint(null.FloatFrom(81), float64(1486113240000)))
				default:
					t.Fatalf("Got unexpected serie name: %q", serie.Name)
				}
			}
		})

		Convey("Group array response parsing", func() {
			var data TargetResponseDTO
			err := json.Unmarshal(groupArrayResponse, &data)
			if err != nil {
				t.Fatalf("Failed to unmarshal clickhouse groupArray response: %s", err)
			}

			queryRes, err := convertToResult(&data)
			if err != nil {
				t.Fatalf("Expected no errors while parsing. Got: %s", err)
			}

			So(len(queryRes.Series), ShouldEqual, 4)
			for _, serie := range queryRes.Series {
				So(len(serie.Points), ShouldEqual, data.Rows)
				switch serie.Name {
				case "MacOS Other":
					So(serie.Points[0], ShouldResemble, tsdb.NewTimePoint(null.FloatFrom(0), float64(1482849300000)))
					So(serie.Points[1], ShouldResemble, tsdb.NewTimePoint(null.FloatFrom(7), float64(1482849360000)))
				case "Windows XP":
					So(serie.Points[0], ShouldResemble, tsdb.NewTimePoint(null.FloatFrom(287), float64(1482849300000)))
					So(serie.Points[1], ShouldResemble, tsdb.NewTimePoint(null.FloatFrom(357), float64(1482849360000)))
				case "MacOS X":
					So(serie.Points[0], ShouldResemble, tsdb.NewTimePoint(null.FloatFrom(0), float64(1482849300000)))
					So(serie.Points[1], ShouldResemble, tsdb.NewTimePoint(null.FloatFrom(2597), float64(1482849360000)))
				case "Windows Other":
					So(serie.Points[0], ShouldResemble, tsdb.NewTimePoint(null.FloatFrom(14), float64(1482849300000)))
					So(serie.Points[1], ShouldResemble, tsdb.NewTimePoint(null.FloatFrom(0), float64(1482849360000)))
				default:
					t.Fatalf("Got unexpected serie name: %q", serie.Name)
				}
			}
		})
	})
}

var regularResponse = []byte(`
{
	"meta":
	[
		{
			"name": "t",
			"type": "UInt64"
		},
		{
			"name": "goodRate",
			"type": "Float64"
		},
		{
			"name": "badRate",
			"type": "Float64"
		}
	],
	"data":
	[
		{
			"t": 1486113220000,
			"goodRate": null,
			"badRate": null
		},
		{
			"t": "1486113240000",
			"goodRate": 19140,
			"badRate": "81"
		},
		{
			"t": "1486113300000",
			"goodRate": 18518,
			"badRate": 150
		},
		{
			"t": "1486113360000",
			"goodRate": 17417,
			"badRate": 209
		}
	],
	"rows": 4,
	"statistics":
	{
		"elapsed": 0.378035523,
		"rows_read": 15645250,
		"bytes_read": 695422779
	}
}
`)

var groupArrayResponse = []byte(`
{
	"meta":
	[
		{
			"name": "t",
			"type": "UInt64"
		},
		{
			"name": "groupArr",
			"type": "Array(Tuple(String, UInt64))"
		}
	],
	"data":
	[
		{
			"t": 1482849300000,
			"groupArr": [["MacOS Other",null],["Windows XP","287"],["Windows Other","14"]]
		},
		{
			"t": "1482849360000",
			"groupArr": [["MacOS Other","7"],["MacOS X","2597"],["Windows XP","357"]]
		},
		{
			"t": "1482849420000",
			"groupArr": [["MacOS Other","7"],["MacOS X","1939"],["Windows Other","28"],["Windows XP","224"]]
		},
		{
			"t": "1482849600000",
			"groupArr": [["MacOS X",1729],["Windows XP","196"]]
		},
		{
			"t": "1482849660000",
			"groupArr": [["MacOS Other","7"],["MacOS X","1589"],["Windows Other","21"],["Windows XP","168"]]
		},
		{
			"t": "1482849720000",
			"groupArr": [["MacOS Other","14"],["MacOS X","1638"],["Windows Other","14"],["Windows XP","168"]]
		}
	],
	"rows": 6,
	"statistics":
	{
		"elapsed": 0.295717111,
		"rows_read": 6523663,
		"bytes_read": 52189708
	}
}
`)
