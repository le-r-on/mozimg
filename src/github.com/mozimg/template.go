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
            <img alt="Embedded Image" style="max-height: 450px; max-width: 950px;" src="{{.Image}}"/>
            <img alt="Avg color" style="max-height: 250; max-width: 250;" src="{{.AvgColor}}"/>
        </div>

        <form  method="post" action="/">
          <input type="submit" value="Get new photo"/> 
        </form>
        <form  method="post" action="/upload" enctype="multipart/form-data">
          <label for="file">Filename:</label>
          <input type="file" name="file" id="file" required>
          <input type="submit" name="submit" value="Display your photo">
        </form> 
        <form method="post" action="/upload_dir" enctype="multipart/form-data">
          <label for="file">Directory:</label>
          <input type="file" multiple name="files" id="files" required>
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