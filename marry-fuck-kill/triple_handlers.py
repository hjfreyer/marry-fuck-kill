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

import logging
import urllib

from google.appengine.ext import db
from google.appengine.ext import webapp
from google.appengine.ext.webapp import util

import models
import utils


class TripleCreationHandler(webapp.RequestHandler):
  def get(self):
    self.response.out.write("""
<h1>Create a Triple!</h1>
<form method="post">
  One: name=<input type="text" name="n1"></input>
     url=<input type="text" name="u1"></input>
     q=<input type="text" name="q1"></input><br/>
  Two: name=<input type="text" name="n2"></input>
     url=<input type="text" name="u2"></input>
     q=<input type="text" name="q2"></input><br/>
  Three: name=<input type="text" name="n3"></input>
     url=<input type="text" name="u3"></input>
     q=<input type="text" name="q3"></input><br/>
  <input type="submit"></input>
</form>
""")

  def post(self):
    try:
      triple = TripleCreationHandler.MakeTriple(self.request)
    # TODO(mjkelly): restrict this later when we have some idea of what
    # we'll throw. Or perhaps not?
    except models.EntityValidationError, e:
      logging.info("Error creating triple from req: %s", self.request)
      self.response.out.write('error: %s' % e)
      return
    # Success
    # TODO(mjkelly): stop using meta refresh redirects
    logging.info("Success creating triple from req: %s", self.request)
    self.response.out.write('ok: created %s' % triple.key().name())

  @staticmethod
  def MakeTriple(request):
    """Create the named triple.

    We expect the following request params:
    n[1-3]: the 3 triple display names
    u[1-3]: the 3 triple image URLs
    q[1-3]: the search string used to find u[1-3]

    The only non-obvious part of this is that we check that q[1-3] actually
    include u[1-3]. This is to prevent users from adding any URL they
    please.
    """
    # Grab all the URL params at once.
    entities = [{'n': request.get('n1'),
           'u': request.get('u1'),
           'q': request.get('q1')},
          {'n': request.get('n2'),
           'u': request.get('u2'),
           'q': request.get('q2')},
          {'n': request.get('n3'),
           'u': request.get('u3'),
           'q': request.get('q3')}]

    for i in range(len(entities)):
      for k in ['n', 'u', 'q']:
        if not entities[i][k]:
          raise ValueError("Entity %s missing attribute '%s'" % (i, k))

    # This may raise a URLError or EntityValidatationError.
    one = models.PutEntity(entities[0]['n'], entities[0]['u'],
                 entities[0]['q'])
    two = models.PutEntity(entities[1]['n'], entities[1]['u'],
                 entities[0]['q'])
    three = models.PutEntity(entities[2]['n'], entities[2]['u'],
                 entities[0]['q'])
    # This may raise an EntityValidationError.
    models.Triple.validate(one, two, three)

    return models.PutTriple(one=one,
                two=two,
                three=three)


class TripleJsonHandler(webapp.RequestHandler):
  def post(self, unused_id):
    try:
      triple = TripleCreationHandler.MakeTriple(self.request)
    except ValueError, e:
      self.response.out.write('error:%s' % e)
      return
    except models.EntityValidationError, e:
      self.response.out.write('error:%s' % e)
      return
    self.response.out.write('ok:%s' % str(triple.key()))

class TripleStatsHandler(webapp.RequestHandler):
  def get(self, triple_id):
    if not triple_id:
      raise Exception("Need triple key")
      
    t = models.Triple.get(urllib.unquote(triple_id))
    entities = [t.one, t.two, t.three]

    self.response.headers['Content-Type'] = "text/plain";
    for e in [t.one, t.two, t.three]:
      self.response.out.write(e.get_stats_url() + '\n')
