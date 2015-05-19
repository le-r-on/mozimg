package main

import "html/template"

var html = template.Must(template.New("mozimg").Parse(`
<html> 
<head> 
    <center><title>Tubular Mozimg</title>
    </head>
    <body>
    <h1>Welcome to Mozimg</h1>
    <center>
        <div id="photo">
            <img alt="Embedded Image" src='{{.Image}}'/>
        </div>
        <form  method="post" action="/">
          <input type="submit" value="Get new photo"/> 
        </form>
        <form  method="post" action="/upload" enctype="multipart/form-data">
          <label for="file">Filename:</label>
          <input type="file" name="file" id="file">
          <input type="submit" name="submit" value="Display yout photo">
        </form>
    </center>
</body>
</html>
`))