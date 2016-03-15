package main

import (
  "html/template"
  "net/http"
  "encoding/json"
  "io/ioutil"
  "fmt"
  "net/url"
)

   //Compile templates on starttemplates/
var templates = template.Must(template.ParseFiles("templates/header.html", "templates/footer.html", "templates/main.html", "templates/search.html", "templates/TopPicks.html" ))

//A Page structure
type Page struct {
  Title string
}

type Payload struct {
        PageNumber    int
        Results []Data `json:"Results"`
}

type Data struct {
        PosterPath       string  `json:"poster_path"`
        Adult            bool    `json:"adult"`
        Overview         string  `json:"overview"`
        ReleaseDate      string  `json:"release_date"`
        GenreIds         []int   `json:"genre_ids"`
        Id               int     `json:"id"`
        OriginalTitle    string  `json:"original_title"`
        OriginalLanguage string  `json:"original_language"`
        Title            string  `json:"title"`
        BackdropPath     string  `json:"backdrop_path"`
        Popularity       float64 `json:"popularity"`
        VoteCount        int     `json:"vote_count"`
        Video            bool    `json:"video"`
        VoteAverage      float64 `json:"vote_average"`
}

//Display the named template
func display(w http.ResponseWriter, tmpl string, data interface{}) {
  templates.ExecuteTemplate(w, tmpl, data)
}

//The handlers.
func mainHandler(w http.ResponseWriter, r *http.Request) {
  display(w, "main", &Page{Title: "Main"})
}



func searchHandler(w http.ResponseWriter, r *http.Request) {
  display(w, "search", &Page{Title: "Search"})
   fmt.Println("method:", r.Method) 
        r.ParseForm()

baseURL := "https://api.themoviedb.org/3/search/movie"
v := url.Values{}
v.Set("query", r.Form.Get("GetSearchKey")) // take GetSearchKey from submitted form
v.Set("api_key", "YOURAPIKEY") // whatever your api key is
perform := baseURL + "?" + v.Encode() // put it all together

  res, err := http.Get(perform)
      if err != nil {
        panic(err)
      }
      defer res.Body.Close()

      body, err := ioutil.ReadAll(res.Body)
      if err != nil {
        panic(err)
      }
      var p Payload

      err = json.Unmarshal(body, &p)
      if err != nil {
        panic(err)
      }

     t, err := template.New("search").Parse(`
    {{define "searchResponse"}}
    {{$ImgUrl := "http://image.tmdb.org/t/p/w185" }}
      <div class="container">
      {{range $movies := .Results}}
               <div class="row">
            <div class="col-md-5">
                <h3>{{$movies.Title}}</h3>
                <h4>Subheading</h4>
                <img class="img-responsive" src="{{$ImgUrl}}{{$movies.BackdropPath}}" alt="Not all Movies have a Backdrop">
                <p>{{$movies.Overview}}</p>
                <ul>      
                <li>Rated For Adults - {{$movies.Adult}}</li>
                <li>Date Released - {{$movies.ReleaseDate}}</li>
                <li>Genre - {{$movies.GenreIds}}</li>
                <li>Movie ID - {{$movies.Id}}</li>
                <li>Original Title - {{$movies.OriginalTitle}}</li>
                <li>Language - {{$movies.OriginalLanguage}}</li>     
                <li>Popularity - {{$movies.Popularity}}</li>
                <li>Votes - {{$movies.VoteCount}}</li>
                <li>Preview - {{$movies.Video}}</li>
                <li>Average Vote - {{$movies.VoteAverage}}</li>
                </ul>            
            </div>
             <div class="col-md-7">
  <img class="img-responsive" src="{{$ImgUrl}}{{$movies.PosterPath}}" alt="Not all movies have a poster">
            </div>
        </div>
        {{end}}
      </div> 
      {{end}}
      `)
    err = t.ExecuteTemplate(w, "searchResponse", p) // This writes the client response
}

func RequestTopMovies(w http.ResponseWriter, r *http.Request) {
    // IMDB Movies Api Top 20 URL
var TopUrl = "https://api.themoviedb.org/3/movie/top_rated?api_key=YOURAPIKEY"
   
  display(w, "TopPicks", &Page{Title: "Top Picks"})
     res, err := http.Get(TopUrl)
      if err != nil {
        panic(err)
      }
      defer res.Body.Close()

      body, err := ioutil.ReadAll(res.Body)
      if err != nil {
        panic(err)
      }
      var p Payload

      err = json.Unmarshal(body, &p)
      if err != nil {
        panic(err)
      }

    for i := 0; i < len(p.Results); i++ {
       fmt.Println(p.Results[i].Overview) // Prints to your terminal
    }

     t, err := template.New("TopPicks").Parse(`
    {{define "body"}}
    {{$ImgUrl := "http://image.tmdb.org/t/p/w185" }}
      <div class="container">
      {{range $movies := .Results}}
               <div class="row">
            <div class="col-md-5">
                <h3>{{$movies.Title}}</h3>
                <h4>Subheading</h4>
                <img class="img-responsive" src="{{$ImgUrl}}{{$movies.BackdropPath}}" alt="Not all movies have a backdrop">
                <p>{{$movies.Overview}}</p>
                <ul>      
                <li>Rated For Adults - {{$movies.Adult}}</li>
                <li>Date Released - {{$movies.ReleaseDate}}</li>
                <li>Genre - {{$movies.GenreIds}}</li>
                <li>Movie ID - {{$movies.Id}}</li>
                <li>Original Title - {{$movies.OriginalTitle}}</li>
                <li>Language - {{$movies.OriginalLanguage}}</li>     
                <li>Popularity - {{$movies.Popularity}}</li>
                <li>Votes - {{$movies.VoteCount}}</li>
                <li>Preview - {{$movies.Video}}</li>
                <li>Average Vote - {{$movies.VoteAverage}}</li>
                </ul>            
            </div>
             <div class="col-md-7">
                    <img class="img-responsive" src="{{$ImgUrl}}{{$movies.PosterPath}}" alt="Not all movies have a poster">
            </div>
        </div>
        {{end}}
      </div> 
      {{end}}
      `)
    err = t.ExecuteTemplate(w, "body", p) // This writes the client response
}

func main() {
  fp := http.FileServer(http.Dir("public"))
  http.Handle("/public/", http.StripPrefix("/public/", fp))
  http.HandleFunc("/", mainHandler)
  http.HandleFunc("/search", searchHandler)
  http.HandleFunc("/TopPicks", RequestTopMovies)
  http.ListenAndServe(":8080", nil)
}