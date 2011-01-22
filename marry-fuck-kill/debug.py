#!/usr/bin/env python
#
# Copyright 2010 Hunter Freyer and Michael Kelly
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
from google.appengine.ext import webapp
from google.appengine.ext.webapp import util

import models
import assignment_handlers
import entity_handlers
import triple_handlers

class DebugHandler(webapp.RequestHandler):
  def get(self):
    response = """
<html>
<head>
  <title>M/F/K</title>
</head>
<body>
<h1>M/F/K - ugly debug page</h1>
<h2>UI Pages</h2>
<ul>
  <li><a href="/MfkWeb.html">MfkWeb (vote)</a>
    (<a href="/MfkWeb.html?gwt.codesvr=127.0.0.1:9997">GWT debug</a>)</li>
  <li><a href="/MfkMaker.html">MfkMaker (create)</a>
    (<a href="/MfkMaker.html?gwt.codesvr=127.0.0.1:9997">GWT debug</a>)</li>
</ul>
<h2>Debugging Tools</h2>
<ul>
  <li><a href="/_ah/admin">Admin</a></li>
  <li><a href="/debug/vote/">Vote</a></li>
  <li><a href="/debug/triple/create">New Triple</a></li>
  <li><a href="/debug/entity/create">New Entity</a></li>
</ul>
</body>
</html>
"""
    self.response.out.write(response)

def main():
  application = webapp.WSGIApplication([('/debug', DebugHandler),
                      ("/debug/vote/(.*)", assignment_handlers.AssignmentHandler),
                      ("/debug/entity/create", entity_handlers.EntityCreationHandler),
                      ("/debug/entity/view/(.*)", entity_handlers.EntityStatsHandler),
                      ("/debug/triple/create", triple_handlers.TripleCreationHandler),
                      ("/debug/triple/view/(.*)", triple_handlers.TripleStatsHandler)],
                     debug=True)
  util.run_wsgi_app(application)


if __name__ == '__main__':
  main()
