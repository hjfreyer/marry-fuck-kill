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

class TripleCreationHandler(webapp.RequestHandler):
    def get(self):
        s = """
<html>
<head>
<title>Create new triple</title>
</head>
<body>
<form action="" method="post">
  <input type="text" name="e1" value="">
  <input type="text" name="e2" value="">
  <input type="text" name="e3" value="">

  <input type="submit">
</form>
</body>
</html>
"""
        self.response.out.write(s)

    def post(self):
        entities = [self.request.get('e1'),
                    self.request.get('e2'),
                    self.request.get('e3')]
        entities = [models.Entity(name=e) for e in sorted(entities)]
        t = models.Triple(one=entities[0], two=entities[1], three=entities[2])
        self.response.out.write('Hello world: %s' % t)

class TripleStatsHandler(webapp.RequestHandler):
    def get(self, triple_id):
        self.response.out.write('Hello world!')
