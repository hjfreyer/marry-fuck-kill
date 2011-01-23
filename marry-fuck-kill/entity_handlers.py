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
import urllib

import models
import utils

class EntityCreationHandler(webapp.RequestHandler):
  def get(self):
    self.response.out.write("""
<h1>Create an Entity!</h1>
<form method="post">
  Name: <input type="text" name="name"></input><br/>
  <input type="submit"></input>
</form>
""")

  def post(self):
    name = self.request.get("name")

    entity = models.Entity(name)
    entity.put()

    utils.redirect(self, '/entity/view/' + name)

class EntityImageHandler(webapp.RequestHandler):
  def get(self, entity_id):
    entity = models.Entity.get(urllib.unquote(entity_id))
    self.response.headers['Content-Type'] = "image/jpg";
    self.response.out.write(entity.data)

class EntityStatsHandler(webapp.RequestHandler):
  def get(self, entity_id):
    if not entity_id:
      raise Exception("Need entity key")

    e = models.Entity.get(urllib.unquote(entity_id))
    m = e.assignment_reference_marry_set.count()
    f = e.assignment_reference_fuck_set.count()
    k = e.assignment_reference_kill_set.count()

    url = ('http://chart.apis.google.com/chart'
        '?chxt=y'
        '&chbh=a'
        '&chs=160x85'
        '&cht=bvg'
        '&chco=A2C180,3D7930,FF9900'
        '&chds=0,%(max)d,0,%(max)d,0,%(max)d'
        '&chd=t:%(m)d|%(f)d|%(k)d'
        '&chdl=Marry|Fuck|Kill'
        '&chdlp=r' % (dict(m=m, f=f, k=k, max=max([m,f,k]))))
    self.response.headers['Content-Type'] = "text/plain";
    self.response.out.write(url)


