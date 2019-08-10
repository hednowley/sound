package dto

import (
	"testing"

	"github.com/hednowley/sound/config"
)

func TestArtistCollection(t *testing.T) {

	conf := &config.Config{
		IgnoredArticles: []string{"the", "los"},
	}

	art := GenerateArt(1)
	artists := GenerateArtists(6, art)

	artists[0].Name = "bhusadiu sdu"
	artists[1].Name = "the Hsetr"
	artists[2].Name = "2 4"
	artists[3].Name = "#sdgs"
	artists[4].Name = "Bdfgfdg"
	artists[5].Name = `"d"`

	GenerateAlbums(3, nil, artists[2], nil)
	GenerateAlbums(1, nil, artists[3], nil)

	DTO := NewArtistCollection(artists, conf)

	xml := `
	<artists ignoredArticles="the los">
		<index name="A">
		</index>
		<index name="B">
			<artist id="1" name="bhusadiu sdu" coverArt="1.jpg" albumCount="0" duration="1">
			</artist>
			<artist id="5" name="Bdfgfdg" coverArt="1.jpg" albumCount="0" duration="5">
			</artist>
		</index>
		<index name="C">
		</index>
		<index name="D">
		</index>
		<index name="E">
		</index>
		<index name="F">
		</index>
		<index name="G">
		</index>
		<index name="H">
			<artist id="2" name="the Hsetr" coverArt="1.jpg" albumCount="0" duration="2">
			</artist>
		</index>
		<index name="I">
		</index>
		<index name="J">
		</index>
		<index name="K">
		</index>
		<index name="L">
		</index>
		<index name="M">
		</index>
		<index name="N">
		</index>
		<index name="O">
		</index>
		<index name="P">
		</index>
		<index name="Q">
		</index>
		<index name="R">
		</index>
		<index name="S">
		</index>
		<index name="T">
		</index>
		<index name="U">
		</index>
		<index name="V">
		</index>
		<index name="W">
		</index>
		<index name="X">
		</index>
		<index name="Y">
		</index>
		<index name="Z">
		</index>
		<index name="#">
			<artist id="3" name="2 4" coverArt="1.jpg" albumCount="3" duration="3">
			</artist>
			<artist id="4" name="#sdgs" coverArt="1.jpg" albumCount="1" duration="4">
			</artist>
			<artist id="6" name="&#34;d&#34;" coverArt="1.jpg" albumCount="0" duration="6">
			</artist>
		</index>
		</artists>
	`

	json := `
	{
		"ignoredArticles":"the los",
		"index":[
		   {
			  "name":"A",
			  "artist":[
	 
			  ]
		   },
		   {
			  "name":"B",
			  "artist":[
				 {
					"id":"1",
					"name":"bhusadiu sdu",
					"coverArt":"1.jpg",
					"albumCount":0,
					"duration":1
				 },
				 {
					"id":"5",
					"name":"Bdfgfdg",
					"coverArt":"1.jpg",
					"albumCount":0,
					"duration":5
				 }
			  ]
		   },
		   {
			  "name":"C",
			  "artist":[
	 
			  ]
		   },
		   {
			  "name":"D",
			  "artist":[
	 
			  ]
		   },
		   {
			  "name":"E",
			  "artist":[
	 
			  ]
		   },
		   {
			  "name":"F",
			  "artist":[
	 
			  ]
		   },
		   {
			  "name":"G",
			  "artist":[
	 
			  ]
		   },
		   {
			  "name":"H",
			  "artist":[
				 {
					"id":"2",
					"name":"the Hsetr",
					"coverArt":"1.jpg",
					"albumCount":0,
					"duration":2
				 }
			  ]
		   },
		   {
			  "name":"I",
			  "artist":[
	 
			  ]
		   },
		   {
			  "name":"J",
			  "artist":[
	 
			  ]
		   },
		   {
			  "name":"K",
			  "artist":[
	 
			  ]
		   },
		   {
			  "name":"L",
			  "artist":[
	 
			  ]
		   },
		   {
			  "name":"M",
			  "artist":[
	 
			  ]
		   },
		   {
			  "name":"N",
			  "artist":[
	 
			  ]
		   },
		   {
			  "name":"O",
			  "artist":[
	 
			  ]
		   },
		   {
			  "name":"P",
			  "artist":[
	 
			  ]
		   },
		   {
			  "name":"Q",
			  "artist":[
	 
			  ]
		   },
		   {
			  "name":"R",
			  "artist":[
	 
			  ]
		   },
		   {
			  "name":"S",
			  "artist":[
	 
			  ]
		   },
		   {
			  "name":"T",
			  "artist":[
	 
			  ]
		   },
		   {
			  "name":"U",
			  "artist":[
	 
			  ]
		   },
		   {
			  "name":"V",
			  "artist":[
	 
			  ]
		   },
		   {
			  "name":"W",
			  "artist":[
	 
			  ]
		   },
		   {
			  "name":"X",
			  "artist":[
	 
			  ]
		   },
		   {
			  "name":"Y",
			  "artist":[
	 
			  ]
		   },
		   {
			  "name":"Z",
			  "artist":[
	 
			  ]
		   },
		   {
			  "name":"#",
			  "artist":[
				 {
					"id":"3",
					"name":"2 4",
					"coverArt":"1.jpg",
					"albumCount":3,
					"duration":3
				 },
				 {
					"id":"4",
					"name":"#sdgs",
					"coverArt":"1.jpg",
					"albumCount":1,
					"duration":4
				 },
				 {
					"id":"6",
					"name":"\"d\"",
					"coverArt":"1.jpg",
					"albumCount":0,
					"duration":6
				 }
			  ]
		   }
		]
	 }
	`

	err := CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}
