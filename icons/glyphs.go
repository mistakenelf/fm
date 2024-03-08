package icons

import "fmt"

// IconInfo is a struct that holds information about an icon.
type IconInfo struct {
	icon       string
	color      [3]uint8
	executable bool
}

// GetGlyph returns the glyph for the icon.
func (i *IconInfo) GetGlyph() string {
	return i.icon
}

// GetColor returns the color for the icon.
func (i *IconInfo) GetColor(f uint8) string {
	switch {
	case i.executable:
		return "\033[38;2;76;175;080m"
	default:
		return fmt.Sprintf("\033[38;2;%d;%d;%dm", i.color[0], i.color[1], i.color[2])
	}
}

// MakeExe is a function that returns a new IconInfo struct with the executable flag set to true.
func (i *IconInfo) MakeExe() {
	i.executable = true
}

// IconSet is a map to represent all the icons.
var IconSet = map[string]*IconInfo{
	"html":             {icon: "\uf13b", color: [3]uint8{228, 79, 57}},   // html
	"markdown":         {icon: "\uf853", color: [3]uint8{66, 165, 245}},  // markdown
	"css":              {icon: "\uf81b", color: [3]uint8{66, 165, 245}},  // css
	"css-map":          {icon: "\ue749", color: [3]uint8{66, 165, 245}},  // css-map
	"sass":             {icon: "\ue603", color: [3]uint8{237, 80, 122}},  // sass
	"less":             {icon: "\ue60b", color: [3]uint8{2, 119, 189}},   // less
	"json":             {icon: "\ue60b", color: [3]uint8{251, 193, 60}},  // json
	"yaml":             {icon: "\ue60b", color: [3]uint8{244, 68, 62}},   // yaml
	"xml":              {icon: "\uf72d", color: [3]uint8{64, 153, 69}},   // xml
	"image":            {icon: "\uf71e", color: [3]uint8{48, 166, 154}},  // image
	"javascript":       {icon: "\ue74e", color: [3]uint8{255, 202, 61}},  // javascript
	"javascript-map":   {icon: "\ue781", color: [3]uint8{255, 202, 61}},  // javascript-map
	"test-jsx":         {icon: "\uf595", color: [3]uint8{35, 188, 212}},  // test-jsx
	"test-js":          {icon: "\uf595", color: [3]uint8{255, 202, 61}},  // test-js
	"react":            {icon: "\ue7ba", color: [3]uint8{35, 188, 212}},  // react
	"react_ts":         {icon: "\ue7ba", color: [3]uint8{36, 142, 211}},  // react_ts
	"settings":         {icon: "\uf013", color: [3]uint8{66, 165, 245}},  // settings
	"typescript":       {icon: "\ue628", color: [3]uint8{3, 136, 209}},   // typescript
	"typescript-def":   {icon: "\ufbe4", color: [3]uint8{3, 136, 209}},   // typescript-def
	"test-ts":          {icon: "\uf595", color: [3]uint8{3, 136, 209}},   // test-ts
	"pdf":              {icon: "\uf724", color: [3]uint8{244, 68, 62}},   // pdf
	"table":            {icon: "\uf71a", color: [3]uint8{139, 195, 74}},  // table
	"visualstudio":     {icon: "\ue70c", color: [3]uint8{173, 99, 188}},  // visualstudio
	"database":         {icon: "\ue706", color: [3]uint8{255, 202, 61}},  // database
	"mysql":            {icon: "\ue704", color: [3]uint8{1, 94, 134}},    // mysql
	"postgresql":       {icon: "\ue76e", color: [3]uint8{49, 99, 140}},   // postgresql
	"sqlite":           {icon: "\ue7c4", color: [3]uint8{1, 57, 84}},     // sqlite
	"csharp":           {icon: "\uf81a", color: [3]uint8{2, 119, 189}},   // csharp
	"zip":              {icon: "\uf410", color: [3]uint8{175, 180, 43}},  // zip
	"exe":              {icon: "\uf2d0", color: [3]uint8{229, 77, 58}},   // exe
	"java":             {icon: "\uf675", color: [3]uint8{244, 68, 62}},   // java
	"c":                {icon: "\ufb70", color: [3]uint8{2, 119, 189}},   // c
	"cpp":              {icon: "\ufb71", color: [3]uint8{2, 119, 189}},   // cpp
	"go":               {icon: "\ufcd1", color: [3]uint8{32, 173, 194}},  // go
	"go-mod":           {icon: "\ufcd1", color: [3]uint8{237, 80, 122}},  // go-mod
	"go-test":          {icon: "\ufcd1", color: [3]uint8{255, 213, 79}},  // go-test
	"python":           {icon: "\uf81f", color: [3]uint8{52, 102, 143}},  // python
	"python-misc":      {icon: "\uf820", color: [3]uint8{130, 61, 28}},   // python-misc
	"url":              {icon: "\uf836", color: [3]uint8{66, 165, 245}},  // url
	"console":          {icon: "\uf68c", color: [3]uint8{250, 111, 66}},  // console
	"word":             {icon: "\uf72b", color: [3]uint8{1, 87, 155}},    // word
	"certificate":      {icon: "\uf623", color: [3]uint8{249, 89, 63}},   // certificate
	"key":              {icon: "\uf805", color: [3]uint8{48, 166, 154}},  // key
	"font":             {icon: "\uf031", color: [3]uint8{244, 68, 62}},   // font
	"lib":              {icon: "\uf831", color: [3]uint8{139, 195, 74}},  // lib
	"ruby":             {icon: "\ue739", color: [3]uint8{229, 61, 58}},   // ruby
	"gemfile":          {icon: "\ue21e", color: [3]uint8{229, 61, 58}},   // gemfile
	"fsharp":           {icon: "\ue7a7", color: [3]uint8{55, 139, 186}},  // fsharp
	"swift":            {icon: "\ufbe3", color: [3]uint8{249, 95, 63}},   // swift
	"docker":           {icon: "\uf308", color: [3]uint8{1, 135, 201}},   // docker
	"powerpoint":       {icon: "\uf726", color: [3]uint8{209, 71, 51}},   // powerpoint
	"video":            {icon: "\uf72a", color: [3]uint8{253, 154, 62}},  // video
	"virtual":          {icon: "\uf822", color: [3]uint8{3, 155, 229}},   // virtual
	"email":            {icon: "\uf6ed", color: [3]uint8{66, 165, 245}},  // email
	"audio":            {icon: "\ufb75", color: [3]uint8{239, 83, 80}},   // audio
	"coffee":           {icon: "\uf675", color: [3]uint8{66, 165, 245}},  // coffee
	"document":         {icon: "\uf718", color: [3]uint8{66, 165, 245}},  // document
	"rust":             {icon: "\ue7a8", color: [3]uint8{250, 111, 66}},  // rust
	"raml":             {icon: "\ue60b", color: [3]uint8{66, 165, 245}},  // raml
	"xaml":             {icon: "\ufb72", color: [3]uint8{66, 165, 245}},  // xaml
	"haskell":          {icon: "\ue61f", color: [3]uint8{254, 168, 62}},  // haskell
	"git":              {icon: "\ue702", color: [3]uint8{229, 77, 58}},   // git
	"lua":              {icon: "\ue620", color: [3]uint8{66, 165, 245}},  // lua
	"clojure":          {icon: "\ue76a", color: [3]uint8{100, 221, 23}},  // clojure
	"groovy":           {icon: "\uf2a6", color: [3]uint8{41, 198, 218}},  // groovy
	"r":                {icon: "\ufcd2", color: [3]uint8{25, 118, 210}},  // r
	"dart":             {icon: "\ue798", color: [3]uint8{87, 182, 240}},  // dart
	"mxml":             {icon: "\uf72d", color: [3]uint8{254, 168, 62}},  // mxml
	"assembly":         {icon: "\uf471", color: [3]uint8{250, 109, 63}},  // assembly
	"vue":              {icon: "\ufd42", color: [3]uint8{65, 184, 131}},  // vue
	"vue-config":       {icon: "\ufd42", color: [3]uint8{58, 121, 110}},  // vue-config
	"lock":             {icon: "\uf83d", color: [3]uint8{255, 213, 79}},  // lock
	"handlebars":       {icon: "\ue60f", color: [3]uint8{250, 111, 66}},  // handlebars
	"perl":             {icon: "\ue769", color: [3]uint8{149, 117, 205}}, // perl
	"elixir":           {icon: "\ue62d", color: [3]uint8{149, 117, 205}}, // elixir
	"erlang":           {icon: "\ue7b1", color: [3]uint8{244, 68, 62}},   // erlang
	"twig":             {icon: "\ue61c", color: [3]uint8{155, 185, 47}},  // twig
	"julia":            {icon: "\ue624", color: [3]uint8{134, 82, 159}},  // julia
	"elm":              {icon: "\ue62c", color: [3]uint8{96, 181, 204}},  // elm
	"smarty":           {icon: "\uf834", color: [3]uint8{255, 207, 60}},  // smarty
	"stylus":           {icon: "\ue600", color: [3]uint8{192, 202, 51}},  // stylus
	"verilog":          {icon: "\ufb19", color: [3]uint8{250, 111, 66}},  // verilog
	"robot":            {icon: "\ufba7", color: [3]uint8{249, 89, 63}},   // robot
	"solidity":         {icon: "\ufcb9", color: [3]uint8{3, 136, 209}},   // solidity
	"yang":             {icon: "\ufb7e", color: [3]uint8{66, 165, 245}},  // yang
	"vercel":           {icon: "\uf47e", color: [3]uint8{207, 216, 220}}, // vercel
	"applescript":      {icon: "\uf302", color: [3]uint8{120, 144, 156}}, // applescript
	"cake":             {icon: "\uf5ea", color: [3]uint8{250, 111, 66}},  // cake
	"nim":              {icon: "\uf6a4", color: [3]uint8{255, 202, 61}},  // nim
	"todo":             {icon: "\uf058", color: [3]uint8{124, 179, 66}},  // task
	"nix":              {icon: "\uf313", color: [3]uint8{80, 117, 193}},  // nix
	"http":             {icon: "\uf484", color: [3]uint8{66, 165, 245}},  // http
	"webpack":          {icon: "\ufc29", color: [3]uint8{142, 214, 251}}, // webpack
	"ionic":            {icon: "\ue7a9", color: [3]uint8{79, 143, 247}},  // ionic
	"gulp":             {icon: "\ue763", color: [3]uint8{229, 61, 58}},   // gulp
	"nodejs":           {icon: "\uf898", color: [3]uint8{139, 195, 74}},  // nodejs
	"npm":              {icon: "\ue71e", color: [3]uint8{203, 56, 55}},   // npm
	"yarn":             {icon: "\uf61a", color: [3]uint8{44, 142, 187}},  // yarn
	"android":          {icon: "\uf531", color: [3]uint8{139, 195, 74}},  // android
	"tune":             {icon: "\ufb69", color: [3]uint8{251, 193, 60}},  // tune
	"contributing":     {icon: "\uf64d", color: [3]uint8{255, 202, 61}},  // contributing
	"readme":           {icon: "\uf7fb", color: [3]uint8{66, 165, 245}},  // readme
	"changelog":        {icon: "\ufba6", color: [3]uint8{139, 195, 74}},  // changelog
	"credits":          {icon: "\uf75f", color: [3]uint8{156, 204, 101}}, // credits
	"authors":          {icon: "\uf0c0", color: [3]uint8{244, 68, 62}},   // authors
	"favicon":          {icon: "\ue623", color: [3]uint8{255, 213, 79}},  // favicon
	"karma":            {icon: "\ue622", color: [3]uint8{60, 190, 174}},  // karma
	"travis":           {icon: "\ue77e", color: [3]uint8{203, 58, 73}},   // travis
	"heroku":           {icon: "\ue607", color: [3]uint8{105, 99, 185}},  // heroku
	"gitlab":           {icon: "\uf296", color: [3]uint8{226, 69, 57}},   // gitlab
	"bower":            {icon: "\ue61a", color: [3]uint8{239, 88, 60}},   // bower
	"conduct":          {icon: "\uf64b", color: [3]uint8{205, 220, 57}},  // conduct
	"jenkins":          {icon: "\ue767", color: [3]uint8{240, 214, 183}}, // jenkins
	"code-climate":     {icon: "\uf7f4", color: [3]uint8{238, 238, 238}}, // code-climate
	"log":              {icon: "\uf719", color: [3]uint8{175, 180, 43}},  // log
	"ejs":              {icon: "\ue618", color: [3]uint8{255, 202, 61}},  // ejs
	"grunt":            {icon: "\ue611", color: [3]uint8{251, 170, 61}},  // grunt
	"django":           {icon: "\ue71d", color: [3]uint8{67, 160, 71}},   // django
	"makefile":         {icon: "\uf728", color: [3]uint8{239, 83, 80}},   // makefile
	"bitbucket":        {icon: "\uf171", color: [3]uint8{31, 136, 229}},  // bitbucket
	"d":                {icon: "\ue7af", color: [3]uint8{244, 68, 62}},   // d
	"mdx":              {icon: "\uf853", color: [3]uint8{255, 202, 61}},  // mdx
	"azure-pipelines":  {icon: "\uf427", color: [3]uint8{20, 101, 192}},  // azure-pipelines
	"azure":            {icon: "\ufd03", color: [3]uint8{31, 136, 229}},  // azure
	"razor":            {icon: "\uf564", color: [3]uint8{66, 165, 245}},  // razor
	"asciidoc":         {icon: "\uf718", color: [3]uint8{244, 68, 62}},   // asciidoc
	"edge":             {icon: "\uf564", color: [3]uint8{239, 111, 60}},  // edge
	"scheme":           {icon: "\ufb26", color: [3]uint8{244, 68, 62}},   // scheme
	"3d":               {icon: "\ue79b", color: [3]uint8{40, 182, 246}},  // 3d
	"svg":              {icon: "\ufc1f", color: [3]uint8{255, 181, 62}},  // svg
	"vim":              {icon: "\ue62b", color: [3]uint8{67, 160, 71}},   // vim
	"moonscript":       {icon: "\uf186", color: [3]uint8{251, 193, 60}},  // moonscript
	"codeowners":       {icon: "\uf507", color: [3]uint8{175, 180, 43}},  // codeowners
	"disc":             {icon: "\ue271", color: [3]uint8{176, 190, 197}}, // disc
	"fortran":          {icon: "F", color: [3]uint8{250, 111, 66}},       // fortran
	"tcl":              {icon: "\ufbd1", color: [3]uint8{239, 83, 80}},   // tcl
	"liquid":           {icon: "\ue275", color: [3]uint8{40, 182, 246}},  // liquid
	"prolog":           {icon: "\ue7a1", color: [3]uint8{239, 83, 80}},   // prolog
	"husky":            {icon: "\uf8e8", color: [3]uint8{229, 229, 229}}, // husky
	"coconut":          {icon: "\uf5d2", color: [3]uint8{141, 110, 99}},  // coconut
	"sketch":           {icon: "\uf6c7", color: [3]uint8{255, 194, 61}},  // sketch
	"pawn":             {icon: "\ue261", color: [3]uint8{239, 111, 60}},  // pawn
	"commitlint":       {icon: "\ufc16", color: [3]uint8{43, 150, 137}},  // commitlint
	"dhall":            {icon: "\uf448", color: [3]uint8{120, 144, 156}}, // dhall
	"dune":             {icon: "\uf7f4", color: [3]uint8{244, 127, 61}},  // dune
	"shaderlab":        {icon: "\ufbad", color: [3]uint8{25, 118, 210}},  // shaderlab
	"command":          {icon: "\ufb32", color: [3]uint8{175, 188, 194}}, // command
	"stryker":          {icon: "\uf05b", color: [3]uint8{239, 83, 80}},   // stryker
	"modernizr":        {icon: "\ue720", color: [3]uint8{234, 72, 99}},   // modernizr
	"roadmap":          {icon: "\ufb6d", color: [3]uint8{48, 166, 154}},  // roadmap
	"debian":           {icon: "\uf306", color: [3]uint8{211, 61, 76}},   // debian
	"ubuntu":           {icon: "\uf31c", color: [3]uint8{214, 73, 53}},   // ubuntu
	"arch":             {icon: "\uf303", color: [3]uint8{33, 142, 202}},  // arch
	"redhat":           {icon: "\uf316", color: [3]uint8{231, 61, 58}},   // redhat
	"gentoo":           {icon: "\uf30d", color: [3]uint8{148, 141, 211}}, // gentoo
	"linux":            {icon: "\ue712", color: [3]uint8{238, 207, 55}},  // linux
	"raspberry-pi":     {icon: "\uf315", color: [3]uint8{208, 60, 76}},   // raspberry-pi
	"manjaro":          {icon: "\uf312", color: [3]uint8{73, 185, 90}},   // manjaro
	"opensuse":         {icon: "\uf314", color: [3]uint8{111, 180, 36}},  // opensuse
	"fedora":           {icon: "\uf30a", color: [3]uint8{52, 103, 172}},  // fedora
	"freebsd":          {icon: "\uf30c", color: [3]uint8{175, 44, 42}},   // freebsd
	"centOS":           {icon: "\uf304", color: [3]uint8{157, 83, 135}},  // centOS
	"alpine":           {icon: "\uf300", color: [3]uint8{14, 87, 123}},   // alpine
	"mint":             {icon: "\uf30f", color: [3]uint8{125, 190, 58}},  // mint
	"routing":          {icon: "\ufb40", color: [3]uint8{67, 160, 71}},   // routing
	"laravel":          {icon: "\ue73f", color: [3]uint8{248, 80, 81}},   // laravel
	"pug":              {icon: "\ue60e", color: [3]uint8{239, 204, 163}}, // pug (Not supported by nerdFont)
	"blink":            {icon: "\uf72a", color: [3]uint8{249, 169, 60}},  // blink (The Foundry Nuke) (Not supported by nerdFont)
	"postcss":          {icon: "\uf81b", color: [3]uint8{244, 68, 62}},   // postcss (Not supported by nerdFont)
	"jinja":            {icon: "\ue000", color: [3]uint8{174, 44, 42}},   // jinja (Not supported by nerdFont)
	"sublime":          {icon: "\ue7aa", color: [3]uint8{239, 148, 58}},  // sublime (Not supported by nerdFont)
	"markojs":          {icon: "\uf13b", color: [3]uint8{2, 119, 189}},   // markojs (Not supported by nerdFont)
	"vscode":           {icon: "\ue70c", color: [3]uint8{33, 150, 243}},  // vscode (Not supported by nerdFont)
	"qsharp":           {icon: "\uf292", color: [3]uint8{251, 193, 60}},  // qsharp (Not supported by nerdFont)
	"vala":             {icon: "\uf7ab", color: [3]uint8{149, 117, 205}}, // vala (Not supported by nerdFont)
	"zig":              {icon: "Z", color: [3]uint8{249, 169, 60}},       // zig (Not supported by nerdFont)
	"h":                {icon: "h", color: [3]uint8{2, 119, 189}},        // h (Not supported by nerdFont)
	"hpp":              {icon: "h", color: [3]uint8{2, 119, 189}},        // hpp (Not supported by nerdFont)
	"powershell":       {icon: "\ufcb5", color: [3]uint8{5, 169, 244}},   // powershell (Not supported by nerdFont)
	"gradle":           {icon: "\ufcc4", color: [3]uint8{29, 151, 167}},  // gradle (Not supported by nerdFont)
	"arduino":          {icon: "\ue255", color: [3]uint8{35, 151, 156}},  // arduino (Not supported by nerdFont)
	"tex":              {icon: "\uf783", color: [3]uint8{66, 165, 245}},  // tex (Not supported by nerdFont)
	"graphql":          {icon: "\ue284", color: [3]uint8{237, 80, 122}},  // graphql (Not supported by nerdFont)
	"kotlin":           {icon: "\ue70e", color: [3]uint8{139, 195, 74}},  // kotlin (Not supported by nerdFont)
	"actionscript":     {icon: "\ufb25", color: [3]uint8{244, 68, 62}},   // actionscript (Not supported by nerdFont)
	"autohotkey":       {icon: "\uf812", color: [3]uint8{76, 175, 80}},   // autohotkey (Not supported by nerdFont)
	"flash":            {icon: "\uf740", color: [3]uint8{198, 52, 54}},   // flash (Not supported by nerdFont)
	"swc":              {icon: "\ufbd3", color: [3]uint8{198, 52, 54}},   // swc (Not supported by nerdFont)
	"cmake":            {icon: "\uf425", color: [3]uint8{178, 178, 179}}, // cmake (Not supported by nerdFont)
	"nuxt":             {icon: "\ue2a6", color: [3]uint8{65, 184, 131}},  // nuxt (Not supported by nerdFont)
	"ocaml":            {icon: "\uf1ce", color: [3]uint8{253, 154, 62}},  // ocaml (Not supported by nerdFont)
	"haxe":             {icon: "\uf425", color: [3]uint8{246, 137, 61}},  // haxe (Not supported by nerdFont)
	"puppet":           {icon: "\uf595", color: [3]uint8{251, 193, 60}},  // puppet (Not supported by nerdFont)
	"purescript":       {icon: "\uf670", color: [3]uint8{66, 165, 245}},  // purescript (Not supported by nerdFont)
	"merlin":           {icon: "\uf136", color: [3]uint8{66, 165, 245}},  // merlin (Not supported by nerdFont)
	"mjml":             {icon: "\ue714", color: [3]uint8{249, 89, 63}},   // mjml (Not supported by nerdFont)
	"terraform":        {icon: "\ue20f", color: [3]uint8{92, 107, 192}},  // terraform (Not supported by nerdFont)
	"apiblueprint":     {icon: "\uf031", color: [3]uint8{66, 165, 245}},  // apiblueprint (Not supported by nerdFont)
	"slim":             {icon: "\uf24e", color: [3]uint8{245, 129, 61}},  // slim (Not supported by nerdFont)
	"babel":            {icon: "\uf5a0", color: [3]uint8{253, 217, 59}},  // babel (Not supported by nerdFont)
	"codecov":          {icon: "\ue37c", color: [3]uint8{237, 80, 122}},  // codecov (Not supported by nerdFont)
	"protractor":       {icon: "\uf288", color: [3]uint8{229, 61, 58}},   // protractor (Not supported by nerdFont)
	"eslint":           {icon: "\ufbf6", color: [3]uint8{121, 134, 203}}, // eslint (Not supported by nerdFont)
	"mocha":            {icon: "\uf6a9", color: [3]uint8{161, 136, 127}}, // mocha (Not supported by nerdFont)
	"firebase":         {icon: "\ue787", color: [3]uint8{251, 193, 60}},  // firebase (Not supported by nerdFont)
	"stylelint":        {icon: "\ufb76", color: [3]uint8{207, 216, 220}}, // stylelint (Not supported by nerdFont)
	"prettier":         {icon: "\uf8e2", color: [3]uint8{86, 179, 180}},  // prettier (Not supported by nerdFont)
	"jest":             {icon: "J", color: [3]uint8{244, 85, 62}},        // jest (Not supported by nerdFont)
	"storybook":        {icon: "\ufd2c", color: [3]uint8{237, 80, 122}},  // storybook (Not supported by nerdFont)
	"fastlane":         {icon: "\ufbff", color: [3]uint8{149, 119, 232}}, // fastlane (Not supported by nerdFont)
	"helm":             {icon: "\ufd31", color: [3]uint8{32, 173, 194}},  // helm (Not supported by nerdFont)
	"i18n":             {icon: "\uf7be", color: [3]uint8{121, 134, 203}}, // i18n (Not supported by nerdFont)
	"semantic-release": {icon: "\uf70f", color: [3]uint8{245, 245, 245}}, // semantic-release (Not supported by nerdFont)
	"godot":            {icon: "\ufba7", color: [3]uint8{79, 195, 247}},  // godot (Not supported by nerdFont)
	"godot-assets":     {icon: "\ufba7", color: [3]uint8{129, 199, 132}}, // godot-assets (Not supported by nerdFont)
	"vagrant":          {icon: "\uf27d", color: [3]uint8{20, 101, 192}},  // vagrant (Not supported by nerdFont)
	"tailwindcss":      {icon: "\ufc8b", color: [3]uint8{77, 182, 172}},  // tailwindcss (Not supported by nerdFont)
	"gcp":              {icon: "\uf662", color: [3]uint8{70, 136, 250}},  // gcp (Not supported by nerdFont)
	"opam":             {icon: "\uf1ce", color: [3]uint8{255, 213, 79}},  // opam (Not supported by nerdFont)
	"pascal":           {icon: "\uf8da", color: [3]uint8{3, 136, 209}},   // pascal (Not supported by nerdFont)
	"nuget":            {icon: "\ue77f", color: [3]uint8{3, 136, 209}},   // nuget (Not supported by nerdFont)
	"denizenscript":    {icon: "D", color: [3]uint8{255, 213, 79}},       // denizenscript (Not supported by nerdFont)
	// "riot":             {icon:"\u", color:[3]uint8{255, 255, 255}},       // riot
	// "autoit":           {icon:"\u", color:[3]uint8{255, 255, 255}},       // autoit
	// "livescript":       {icon:"\u", color:[3]uint8{255, 255, 255}},       // livescript
	// "reason":           {icon:"\u", color:[3]uint8{255, 255, 255}},       // reason
	// "bucklescript":     {icon:"\u", color:[3]uint8{255, 255, 255}},       // bucklescript
	// "mathematica":      {icon:"\u", color:[3]uint8{255, 255, 255}},       // mathematica
	// "wolframlanguage":  {icon:"\u", color:[3]uint8{255, 255, 255}},       // wolframlanguage
	// "nunjucks":         {icon:"\u", color:[3]uint8{255, 255, 255}},       // nunjucks
	// "haml":             {icon:"\u", color:[3]uint8{255, 255, 255}},       // haml
	// "cucumber":         {icon:"\u", color:[3]uint8{255, 255, 255}},       // cucumber
	// "vfl":              {icon:"\u", color:[3]uint8{255, 255, 255}},       // vfl
	// "kl":               {icon:"\u", color:[3]uint8{255, 255, 255}},       // kl
	// "coldfusion":       {icon:"\u", color:[3]uint8{255, 255, 255}},       // coldfusion
	// "cabal":            {icon:"\u", color:[3]uint8{255, 255, 255}},       // cabal
	// "restql":           {icon:"\u", color:[3]uint8{255, 255, 255}},       // restql
	// "kivy":             {icon:"\u", color:[3]uint8{255, 255, 255}},       // kivy
	// "graphcool":        {icon:"\u", color:[3]uint8{255, 255, 255}},       // graphcool
	// "sbt":              {icon:"\u", color:[3]uint8{255, 255, 255}},       // sbt
	// "flow":             {icon:"\u", color:[3]uint8{255, 255, 255}},       // flow
	// "bithound":         {icon:"\u", color:[3]uint8{255, 255, 255}},       // bithound
	// "appveyor":         {icon:"\u", color:[3]uint8{255, 255, 255}},       // appveyor
	// "fusebox":          {icon:"\u", color:[3]uint8{255, 255, 255}},       // fusebox
	// "editorconfig":     {icon:"\u", color:[3]uint8{255, 255, 255}},       // editorconfig
	// "watchman":         {icon:"\u", color:[3]uint8{255, 255, 255}},       // watchman
	// "aurelia":          {icon:"\u", color:[3]uint8{255, 255, 255}},       // aurelia
	// "rollup":           {icon:"\u", color:[3]uint8{255, 255, 255}},       // rollup
	// "hack":             {icon:"\u", color:[3]uint8{255, 255, 255}},       // hack
	// "apollo":           {icon:"\u", color:[3]uint8{255, 255, 255}},       // apollo
	// "nodemon":          {icon:"\u", color:[3]uint8{255, 255, 255}},       // nodemon
	// "webhint":          {icon:"\u", color:[3]uint8{255, 255, 255}},       // webhint
	// "browserlist":      {icon:"\u", color:[3]uint8{255, 255, 255}},       // browserlist
	// "crystal":          {icon:"\u", color:[3]uint8{255, 255, 255}},       // crystal
	// "snyk":             {icon:"\u", color:[3]uint8{255, 255, 255}},       // snyk
	// "drone":            {icon:"\u", color:[3]uint8{255, 255, 255}},       // drone
	// "cuda":             {icon:"\u", color:[3]uint8{255, 255, 255}},       // cuda
	// "dotjs":            {icon:"\u", color:[3]uint8{255, 255, 255}},       // dotjs
	// "sequelize":        {icon:"\u", color:[3]uint8{255, 255, 255}},       // sequelize
	// "gatsby":           {icon:"\u", color:[3]uint8{255, 255, 255}},       // gatsby
	// "wakatime":         {icon:"\u", color:[3]uint8{255, 255, 255}},       // wakatime
	// "circleci":         {icon:"\u", color:[3]uint8{255, 255, 255}},       // circleci
	// "cloudfoundry":     {icon:"\u", color:[3]uint8{255, 255, 255}},       // cloudfoundry
	// "processing":       {icon:"\u", color:[3]uint8{255, 255, 255}},       // processing
	// "wepy":             {icon:"\u", color:[3]uint8{255, 255, 255}},       // wepy
	// "hcl":              {icon:"\u", color:[3]uint8{255, 255, 255}},       // hcl
	// "san":              {icon:"\u", color:[3]uint8{255, 255, 255}},       // san
	// "wallaby":          {icon:"\u", color:[3]uint8{255, 255, 255}},       // wallaby
	// "stencil":          {icon:"\u", color:[3]uint8{255, 255, 255}},       // stencil
	// "red":              {icon:"\u", color:[3]uint8{255, 255, 255}},       // red
	// "webassembly":      {icon:"\u", color:[3]uint8{255, 255, 255}},       // webassembly
	// "foxpro":           {icon:"\u", color:[3]uint8{255, 255, 255}},       // foxpro
	// "jupyter":          {icon:"\u", color:[3]uint8{255, 255, 255}},       // jupyter
	// "ballerina":        {icon:"\u", color:[3]uint8{255, 255, 255}},       // ballerina
	// "racket":           {icon:"\u", color:[3]uint8{255, 255, 255}},       // racket
	// "bazel":            {icon:"\u", color:[3]uint8{255, 255, 255}},       // bazel
	// "mint":             {icon:"\u", color:[3]uint8{255, 255, 255}},       // mint
	// "velocity":         {icon:"\u", color:[3]uint8{255, 255, 255}},       // velocity
	// "prisma":           {icon:"\u", color:[3]uint8{255, 255, 255}},       // prisma
	// "abc":              {icon:"\u", color:[3]uint8{255, 255, 255}},       // abc
	// "istanbul":         {icon:"\u", color:[3]uint8{255, 255, 255}},       // istanbul
	// "lisp":             {icon:"\u", color:[3]uint8{255, 255, 255}},       // lisp
	// "buildkite":        {icon:"\u", color:[3]uint8{255, 255, 255}},       // buildkite
	// "netlify":          {icon:"\u", color:[3]uint8{255, 255, 255}},       // netlify
	// "svelte":           {icon:"\u", color:[3]uint8{255, 255, 255}},       // svelte
	// "nest":             {icon:"\u", color:[3]uint8{255, 255, 255}},       // nest
	// "percy":            {icon:"\u", color:[3]uint8{255, 255, 255}},       // percy
	// "gitpod":           {icon:"\u", color:[3]uint8{255, 255, 255}},       // gitpod
	// "advpl_prw":        {icon:"\u", color:[3]uint8{255, 255, 255}},       // advpl_prw
	// "advpl_ptm":        {icon:"\u", color:[3]uint8{255, 255, 255}},       // advpl_ptm
	// "advpl_tlpp":       {icon:"\u", color:[3]uint8{255, 255, 255}},       // advpl_tlpp
	// "advpl_include":    {icon:"\u", color:[3]uint8{255, 255, 255}},       // advpl_include
	// "tilt":             {icon:"\u", color:[3]uint8{255, 255, 255}},       // tilt
	// "capacitor":        {icon:"\u", color:[3]uint8{255, 255, 255}},       // capacitor
	// "adonis":           {icon:"\u", color:[3]uint8{255, 255, 255}},       // adonis
	// "forth":            {icon:"\u", color:[3]uint8{255, 255, 255}},       // forth
	// "uml":              {icon:"\u", color:[3]uint8{255, 255, 255}},       // uml
	// "meson":            {icon:"\u", color:[3]uint8{255, 255, 255}},       // meson
	// "buck":             {icon:"\u", color:[3]uint8{255, 255, 255}},       // buck
	// "sml":              {icon:"\u", color:[3]uint8{255, 255, 255}},       // sml
	// "nrwl":             {icon:"\u", color:[3]uint8{255, 255, 255}},       // nrwl
	// "imba":             {icon:"\u", color:[3]uint8{255, 255, 255}},       // imba
	// "drawio":           {icon:"\u", color:[3]uint8{255, 255, 255}},       // drawio
	// "sas":              {icon:"\u", color:[3]uint8{255, 255, 255}},       // sas
	// "slug":             {icon:"\u", color:[3]uint8{255, 255, 255}},       // slug

	"dir-config":      {icon: "\ue5fc", color: [3]uint8{32, 173, 194}},  // dir-config
	"dir-controller":  {icon: "\ue5fc", color: [3]uint8{255, 194, 61}},  // dir-controller
	"dir-git":         {icon: "\ue5fb", color: [3]uint8{250, 111, 66}},  // dir-git
	"dir-github":      {icon: "\ue5fd", color: [3]uint8{84, 110, 122}},  // dir-github
	"dir-npm":         {icon: "\ue5fa", color: [3]uint8{203, 56, 55}},   // dir-npm
	"dir-include":     {icon: "\uf756", color: [3]uint8{3, 155, 229}},   // dir-include
	"dir-import":      {icon: "\uf756", color: [3]uint8{175, 180, 43}},  // dir-import
	"dir-upload":      {icon: "\uf758", color: [3]uint8{250, 111, 66}},  // dir-upload
	"dir-download":    {icon: "\uf74c", color: [3]uint8{76, 175, 80}},   // dir-download
	"dir-secure":      {icon: "\uf74f", color: [3]uint8{249, 169, 60}},  // dir-secure
	"dir-images":      {icon: "\uf74e", color: [3]uint8{43, 150, 137}},  // dir-images
	"dir-environment": {icon: "\uf74e", color: [3]uint8{102, 187, 106}}, // dir-environment
}

// IconDef is a map of default icons if none can be found.
var IconDef = map[string]*IconInfo{
	"dir":        {icon: "\uf74a", color: [3]uint8{224, 177, 77}},
	"diropen":    {icon: "\ufc6e", color: [3]uint8{224, 177, 77}},
	"hiddendir":  {icon: "\uf755", color: [3]uint8{224, 177, 77}},
	"exe":        {icon: "\uf713", color: [3]uint8{76, 175, 80}},
	"file":       {icon: "\uf723", color: [3]uint8{65, 129, 190}},
	"hiddenfile": {icon: "\ufb12", color: [3]uint8{65, 129, 190}},
}
