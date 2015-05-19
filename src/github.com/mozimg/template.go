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
    </center>
</body>
</html>
`))