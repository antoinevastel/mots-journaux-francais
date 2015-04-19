package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"regexp"
	"sort"
	"strings"
)

type sortedMap struct {
	m map[string]int
	s []string
}

func (sm *sortedMap) Len() int {
	return len(sm.m)
}

func (sm *sortedMap) Less(i, j int) bool {
	return sm.m[sm.s[i]] > sm.m[sm.s[j]]
}

func (sm *sortedMap) Swap(i, j int) {
	sm.s[i], sm.s[j] = sm.s[j], sm.s[i]
}

func sortedKeys(m map[string]int) []string {
	sm := new(sortedMap)
	sm.m = m
	sm.s = make([]string, len(m))
	i := 0
	for key, _ := range m {
		sm.s[i] = key
		i++
	}
	sort.Sort(sm)
	return sm.s
}

func getTitresJournal(lien string, balise string) (titres string) {
	doc, err := goquery.NewDocument(lien)
	titres = ""

	if err != nil {
		log.Fatal(err)
	}

	doc.Find(balise).Each(func(i int, s *goquery.Selection) {
		titres += s.Find("a").Text() + " "
	})

	return strings.ToUpper(titres)
}

func filtrerTitres(titres string) (titresFiltres []string, longueur int) {
	expression := "[’«»,'\":()/-]{1}|POUR|DANS|PLUS|2014|DEUX|TROIS|2015|AVEC|VEUT|APRÈS|SONT|CONTRE"
	expression += "|AVANT|VOUS|DEVANT|DIRECT|VIDEO|ÊTRE|FAIT|COMMENT|GARDE|\\.|\\n"
	expression += "|PARIS|SUITE|TOUJOURS|FRANCE|SANS"
	re, _ := regexp.Compile(expression)

	s := re.ReplaceAllLiteralString(titres, " ")
	tabComp := strings.Split(s, " ")

	var tabFinal = make([]string, len(tabComp))

	j := 0
	for i := 0; i < len(tabComp); i++ {
		if len(tabComp[i]) > 3 {
			tabFinal[j] = tabComp[i]
			j++
		}
	}

	return tabFinal, j
}

func comptageMot(titresFiltres []string, longueur int) (nombreDeMots map[string]int) {
	nombreDeMots = make(map[string]int)

	for i := 0; i < longueur; i++ {
		nombreDeMots[titresFiltres[i]]++
	}

	return nombreDeMots
}

func main() {
	titres := ""

	sites := make(map[string]string)
	sites["http://www.lemonde.fr"] = "ul[class~=liste_horaire] li"
	sites["http://www.lefigaro.fr"] = "ul[class~=flashActu-listContent] li"
	sites["http://www.leparisien.fr"] = "div[class~=cont] li"
	sites["http://www.lesechos.fr"] = "div[class~=dminv2]"
	sites["http://www.latribune.fr"] = "div[class~=bloc-actus]"
	sites["http://www.liberation.fr"] = "#pager-feed-news"
	sites["http://www.francesoir.fr"] = "div[class~=view-id-actualites]"
	sites["http://www.la-croix.com/Depeches"] = "div[class~=mea_actu] h1"
	sites["http://www.directmatin.fr/news"] = "div[class~=dm-article-taxonomie-item]"

	for k, v := range sites {
		titres += getTitresJournal(k, v)
	}

	s, longueur := filtrerTitres(titres)

	nombreDeMots := comptageMot(s, longueur)
	/*
		for k, v := range nombreDeMots {
			fmt.Printf("Mot : %s, nombre d'occurences : %d \n", k, v)
		}
	*/
	i := 1
	for _, res := range sortedKeys(nombreDeMots) {
		fmt.Println(i, res, nombreDeMots[res])
		i++
	}

}
