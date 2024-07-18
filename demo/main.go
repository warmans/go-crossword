package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/warmans/go-crossword"
	"strconv"
	"strings"

	"github.com/warmans/vue"
)

type Data struct {
	ImgData     string
	ImageWidth  string
	ImageHeight string
	Attempts    string
	Words       string
	GridSize    string
	ShowWords   string
	Error       string
}

func main() {

	data := &Data{
		ImageWidth:  "800",
		ImageHeight: "800",
		GridSize:    "25",
		Attempts:    "5",
		Words:       "food,food clue\nother,other clue\nrage,rage clue",
		ShowWords:   "true",
	}

	var err error
	data.ImgData, err = renderDataURL(data)
	if err != nil {
		data.Error = err.Error()
	}

	vue.New(
		vue.El("#app"),
		vue.Template(`
			<div style="background-color: #d00; color: #fff; width: 100%">{{Error}}</div>
			<div style="margin-bottom: 20px; padding-bottom: 20px; border-bottom: 1px dashed #ccc">
				<img v-bind:src="ImgData" />
			</div>
			<table>
				<tbody>
					<tr><th style="width: 10rem">Image Size</th><td><input v-model="ImageWidth" /> x <input v-model="ImageHeight" /></td></tr>
					<tr><th style="width: 10rem">Grid Size</th><td><input v-model="GridSize" /></td></tr>
					<tr><th style="width: 10rem">Attempts</th><td><input v-model="Attempts" /></td></tr>
					<tr>
						<th>Show Words</th>
						<td>
							<select v-model="ShowWords">
								<option value="true">YES</option>
								<option value="false">NO</option>
							</select>
						</td>
					</tr>
					<tr><th>&nbsp;</th><td></td></tr>

					<tr><th style="width: 10rem">Words CSV</th><td><textarea style="width:100%" v-model="Words"></textarea></td></tr>
					<tr><th>&nbsp;</th><td></td></tr>
					<tr><th>&nbsp;</th><td></td></tr>
					
		
					<tr><th>&nbsp;</th><td></td></tr>

					<tr><th></th><td><button v-on:click="Render">Render</button></td><tr>
				</tbody>
			</table>
		`),
		vue.Data(data),
		vue.Methods(Render),
	)

	select {}
}

func Render(vctx vue.Context) {
	data := vctx.Data().(*Data)
	var err error
	data.ImgData, err = renderDataURL(data)
	data.Error = ""
	if err != nil {
		data.Error = err.Error()
	}
}

func renderDataURL(cfg *Data) (string, error) {

	words, err := crossword.WordsFromCSV(strings.NewReader(cfg.Words))
	if err != nil {
		return "", err
	}
	if len(words) == 0 {
		return "", fmt.Errorf("no words were given")
	}

	cw := crossword.Generate(
		parseIntOrDefault(cfg.GridSize, 25),
		words,
		parseIntOrDefault(cfg.Attempts, 5),
	)
	if cw == nil {
		return "", fmt.Errorf("crossword was not generated")
	}
	canvas, err := crossword.RenderPNG(
		cw,
		parseIntOrDefault(cfg.ImageWidth, 1000),
		parseIntOrDefault(cfg.ImageHeight, 1000),
		crossword.WithAllSolved(cfg.ShowWords == "true"),
	)
	if err != nil {
		panic("failed to create canvas: " + err.Error())
	}

	buff := &bytes.Buffer{}
	if err := canvas.EncodePNG(buff); err != nil {
		panic("failed to encode image: " + err.Error())
	}

	return fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(buff.Bytes())), nil
}

func parseIntOrDefault(strVal string, def int) int {
	val, err := strconv.ParseInt(strVal, 10, 64)
	if err != nil {
		return def
	}
	return int(val)
}
