package dao

type AlbumList2Type int

const (
	Random               AlbumList2Type = 0
	Newest               AlbumList2Type = 1
	Frequent             AlbumList2Type = 2
	Recent               AlbumList2Type = 3
	Starred              AlbumList2Type = 4
	AlphabeticalByName   AlbumList2Type = 5
	AlphabeticalByArtist AlbumList2Type = 6
	ByYear               AlbumList2Type = 7
	ByGenre              AlbumList2Type = 8
)
