package models

type ArtistSummary struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	LocationsURL string   `json:"locations"`
	DatesURL     string   `json:"concertDates"`
	RelationURL  string   `json:"relations"`
}

type Locations struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	DatesURL  string   `json:"dates"`
}

type Dates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type FullArtistInfo struct {
	ArtistSummary
	Locations []string
	Dates     []string
	Relation  map[string][]string

	// Custom fields for embeds (manual population)
	YouTubeID string
}

type ErrorData struct {
	Code    int
	Message string
}

var YouTubeIDs = map[int]string{
	1:  "fJ9rUzIMcZQ", // Queen – Bohemian Rhapsody
	2:  "qj3sku5VapI", // SOJA – True Love
	3:  "x-xTttimcNk", // Pink Floyd – Comfortably Numb
	4:  "6yP1tcy9a10", // Scorpions – Rock You Like a Hurricane
	5:  "pgN-vvVVxMA", // XXXTENTACION – SAD!
	6:  "SsKT0s5J8ko", // Mac Miller – Self Care
	7:  "43gm3CJePn0", // Joyner Lucas – I'm Not Racist
	8:  "tvTRZJ-4EyI", // Kendrick Lamar – HUMBLE.
	9:  "pAgnJDJN4VA", // AC/DC – Back In Black
	10: "qM0zINtulhM", // Pearl Jam – Alive
	11: "CevxZvSJLk8", // Katy Perry – Roar
	12: "lWA2pjMjpBs", // Rihanna – Diamonds
	13: "TlBIa8z_Mts", // Genesis – Land of Confusion
	14: "YkADj0TPrJA", // Phil Collins – In the Air Tonight
	15: "QkF3oxziUI4", // Led Zeppelin – Stairway to Heaven
	16: "TLV4_xaYynY", // The Jimi Hendrix Experience – All Along the Watchtower
	17: "I_izvAbhExY", // Bee Gees – Stayin' Alive
	18: "Q2FzZSBD5LE", // Deep Purple – Smoke on the Water
	19: "JkK8g6FMEXE", // Aerosmith – I Don't Want to Miss a Thing
	20: "h0ffIJ7ZO4U", // Dire Straits – Sultans of Swing
	21: "EHlZEcy24q8", // Mamonas Assassinas – Pelados em Santos
	22: "8yvGCAvOAfM", // Thirty Seconds to Mars – The Kill
	23: "7wtfhZwyrcc", // Imagine Dragons – Believer
	24: "mzB1VGEGcSU", // Juice WRLD – Lucid Dreams
	25: "Kb24RrHIbFk", // Logic – 1-800-273-8255
	26: "50VNCymT-Cs", // Alec Benjamin – Let Me Down Slowly
	27: "d-diB65scQU", // Bobby McFerrin – Don't Worry, Be Happy
	28: "LQ7R9zHeDy8", // R3HAB – All Around The World (La La La)
	29: "wXhTHyIgQ_U", // Post Malone – Circles
	30: "6ONRf7h3Mdk", // Travis Scott – SICKO MODE
	31: "ue5oHmUGiMM", // J. Cole – No Role Modelz
	32: "Aiay8I5IPB8", // Nickelback – How You Remind Me
	33: "rTKpYJ80OVQ", // Mobb Deep – Shook Ones, Pt. II
	34: "1w7OgIMMRc4", // Guns N' Roses – Sweet Child O' Mine
	35: "TMZi25Pq3T8", // N.W.A – Straight Outta Compton
	36: "ujNeHIo7oTE", // U2 – With Or Without You
	37: "bpOSxM0rNPM", // Arctic Monkeys – Do I Wanna Know?
	38: "uhG-vLZrb-g", // Fall Out Boy – Sugar, We're Goin Down
	39: "HyHNuVaZJ-k", // Gorillaz – Feel Good Inc.
	40: "09839DpTctU", // Eagles – Hotel California
	41: "kXYiU_JCYtU", // Linkin Park – Numb
	42: "YlUKcNNmywk", // Red Hot Chili Peppers – Californication
	43: "_Yhyp-_hX2s", // Eminem – Lose Yourself
	44: "Soa3gO7tL-c", // Green Day – Boulevard of Broken Dreams
	45: "CD-E-LDc384", // Metallica – Enter Sandman
	46: "dvgZkm1xWPE", // Coldplay – Viva la Vida
	47: "09R8_2nJtjg", // Maroon 5 – Sugar
	48: "pXRviuL6vMY", // Twenty One Pilots – Stressed Out
	49: "nrIPxlFzDi0", // The Rolling Stones – (I Can't Get No) Satisfaction
	50: "w8KQmps-Sog", // Muse – Uprising
	51: "SBjQ9tuuTJQ", // Foo Fighters – The Pretender
	52: "PT2_F-1esPk", // The Chainsmokers – Closer
}
