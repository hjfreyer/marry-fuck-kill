#!/usr/bin/env python
#
# Copyright 2007 Google Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
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

class MainHandler(webapp.RequestHandler):
    def get(self):
        t = models.Triple.get_random_triple()
        response = """
<html>
<head>
    <title>M/F/K</title>
</head>
<body>
<h1>M/F/K</h1>
<h2>Links</h2>
<ul>
    <li><a href="/_ah/admin">Admin</a></li>
    <li><a href="/triple/create">New Triple</a></li>
    <li><a href="/entity/create">New Entity</a></li>
</ul>

<h2>Random triple:</h2>
<p>%s</p>
</body>
</html>
""" % t
        self.response.out.write(response)

def main():
    application = webapp.WSGIApplication([('/', MainHandler),
                                          ("/entity/create", entity_handlers.EntityCreationHandler),
                                          ("/entity/view/(.*)", entity_handlers.EntityStatsHandler),
                                          ("/vote/(.*)", assignment_handlers.AssignmentHandler),
                                          ("/triple/create", triple_handlers.TripleCreationHandler),
                                          ("/triple/view/(.*)", triple_handlers.TripleStatsHandler)                                         ],
                                         debug=True)
    util.run_wsgi_app(application)


if __name__ == '__main__':
    main()
