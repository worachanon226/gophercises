package urlshortener

import "net/http"

func MapHandler(pathToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
		}

		fallback.ServeHTTP(w, r)
	}
}

func YAMLHandler(yaml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathsToUrls := map[string]string{}
	return MapHandler(pathsToUrls, fallback), nil
}
