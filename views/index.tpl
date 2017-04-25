<!DOCTYPE html>

<html>
<head>
  <title>Union Trading</title>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
  <script src="/static/js/jquery-1.10.1.min.js"></script>
  <script src="/static/js/json-viewer/jquery.json-viewer.js"></script>
  <link href="/static/js/json-viewer/jquery.json-viewer.css" type="text/css" rel="stylesheet" />
  <link rel="shortcut icon" href="/static/img/bee.png" type="image/x-icon" />

  <style type="text/css">
    *,body {
      margin: 0px;
      padding: 0px;
    }

    pre#json-renderer {
        border: 1px solid #aaa;
        padding: 0.5em 1.5em;
    }

    body {
      margin: 0px;
      font-family: "Helvetica Neue", Helvetica, Arial, sans-serif;
      font-size: 14px;
      line-height: 20px;
      background-color: #fff;
    }

    header,
    footer {
      width: 960px;
      margin-left: auto;
      margin-right: auto;
    }

    .logo {
      background:url(/static/img/dollar.png);
      background-repeat: no-repeat;
      -webkit-background-size: 100px 100px;
      background-size: 100px 100px;
      background-position: center center;
      text-align: center;
      font-size: 42px;
      padding: 250px 0 70px;
      font-weight: normal;
      text-shadow: 0px 1px 2px #ddd;
    }

    header {
      padding: 100px 0;
    }

    footer {
      line-height: 1.8;
      text-align: left;
      padding: 50px 0;
      color: #999;
    }

    .description {
      text-align: center;
      font-size: 16px;
    }

    .info, a {
      color: #444;
      text-decoration: none;
    }
  </style>
</head>

<body>
  <header>
    <h1 class="logo">Welcome to Union HTTP Server</h1>
    <div class="description">
      Union is a trading server where you can have fun and earn money!.
    </div>
  </header>
  <script src="/static/js/websocket.js"></script>
  <footer>
    <div class="author">
      <button id="wsbutton" onclick="wsSwitch()">Stop/Start updating</button>
      JSON  <pre id="json-renderer"></pre>
      Official website:
      <a href="http://{{.Website}}">{{.Website}}</a> /
      Contact me:
      <a class="email" href="mailto:{{.Email}}">{{.Email}}</a>
    </div>
  </footer>
</body>
</html>
