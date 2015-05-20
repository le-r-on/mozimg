package main

import "html/template"

var base_tmpl = template.Must(template.New("mozimg").Parse(`
<html>
<head>
    <center><title>Tubular Mozimg</title>
	<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.4/css/bootstrap.min.css" />
	<style>
		.btn-file {position: relative; overflow: hidden;}
		.btn-file input[type=file] {
			position: absolute;
			top: 0;
			right: 0;
			min-width: 100%;
			min-height: 100%;
			font-size: 100px;
			text-align: right;
			filter: alpha(opacity=0);
			opacity: 0;
			outline: none;
			background: white;
			cursor: inherit;
			display: block;
		}
	</style>
</head>
<body>
    <nav class="navbar navbar-default">
	<a class="navbar-brand" style="font-size: 25px; font-weight: bold" href="/">mozimg</a>
	</nav>
        <div id="photo">
            <img alt="Transformed Image" class="img-thumbnail" style="max-height: 450px; max-width: 950px;" src="{{.TiledImage}}"/>
            <img alt="" class="img-thumbnail" style="max-height: 250; max-width: 250;" src="{{.OrigImage}}"/>
        </div>
		<br/>
        <form class="form-inline" method="post" action="/">
          <input type="text" class="form-control" id="pic_num" name="pic_num" placeholder="Number of tiles" value="">
          <input type="text" class="form-control" id="dimension" name="dimension" placeholder="Number of divisions" value="">
          <input type="submit" name="submit" class="btn btn-primary" value="Get random photo"/> 
        </form>

		<form class="form-inline" method="post" action="/tile" enctype="multipart/form-data">
       	  <span class="btn btn-success btn-file">
			Choose File <input type="file" name="file" id="file" required>
		  </span>
		  <span class="btn btn-success btn-file">
          	Choose Tiles <input type="file" multiple name="files" id="files" required>
		  </span>
          <input type="text" class="form-control" name="dimension" placeholder="Number of divisions" required>
          <input class="btn btn-primary" type="submit" name="submit" value="Load files">
        </form>
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
