package main

import "html/template"

var base_tmpl = template.Must(template.New("mozimg").Parse(`
<html> 
<head>
    <center><title>Tubular Mozimg</title>
</head>
<body>
    <h1>Welcome to Mozimg</h1>
    <center>
        <div id="photo">
            <img alt="Transformed Image" style="max-height: 450px; max-width: 950px;" src="{{.TiledImage}}"/>
            <img alt="" style="max-height: 250; max-width: 250;" src="{{.OrigImage}}"/>
        </div>

        <form  method="post" action="/">
          <input type="text" id="pic_num" name="pic_num" placeholder="Number of pictures" value="">
          <input type="text" id="dimension" name="dimension" placeholder="Number of tiles" value="">
          <input type="submit" name="submit" value="Get new random photo"/> 
        </form>
        <form  method="post" action="/tile" enctype="multipart/form-data">
          <label for="file">Photo:</label>
          <input type="file" name="file" id="file" required>
          <label for="files">Tiles:</label>
          <input type="file" multiple name="files" id="files" required>
          <input type="text" name="dimension" placeholder="Tile dimension" required>
          <input type="submit" name="submit" value="Load tiles">
        </form>
    </center>
</body>
</html>
`))

var error_tmpl = template.Must(template.New("error").Parse(`
<html> 
<head> 
    <center><title>Tubular Mozimg</title>
    </head>
    <body>
    <h1>Ok, you won, we failed</h1>
    <div>{{.Message}}</div>
</body>
</html>
`))