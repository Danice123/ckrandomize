package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serviceCmd)

	homeTemplate = template.New("Home")
	if html, err := os.ReadFile(filepath.Join("templates", "home.html")); err != nil {
		panic(err.Error())
	} else {
		if _, err := homeTemplate.Parse(string(html)); err != nil {
			panic(err.Error())
		}
	}

	redirectTemplate = template.New("Redirect")
	if html, err := os.ReadFile(filepath.Join("templates", "redirect.html")); err != nil {
		panic(err.Error())
	} else {
		if _, err := redirectTemplate.Parse(string(html)); err != nil {
			panic(err.Error())
		}
	}
}

var homeTemplate *template.Template
var redirectTemplate *template.Template

var patchDepo string

var serviceCmd = &cobra.Command{
	Use:   "service [hostname] [patch depo]",
	Short: "Serve rom patcher",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		patchDepo = args[1]
		router := httprouter.New()

		router.GET("/", Home)
		router.GET("/randomize", Randomize)

		static := http.FileServer(http.Dir(patchDepo))
		router.Handler("GET", "/patches/*path", http.StripPrefix("/patches/", static))

		return http.ListenAndServe(args[0], router)
	},
}

func Home(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	err := homeTemplate.Execute(w, map[string]string{
		"error": req.URL.Query().Get("error"),
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}

func Randomize(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	var path string
	seed := uuid.New()
	c := exec.Command("./rand.sh", patchDepo, seed.String())
	c.Stderr = os.Stderr
	err := c.Run()
	if err != nil {
		path = fmt.Sprintf("/?error=%s", err.Error())
	} else {
		path = fmt.Sprintf("/patches/%s", seed.String())
	}

	err = redirectTemplate.Execute(w, path)
	if err != nil {
		fmt.Println(err.Error())
	}
}
