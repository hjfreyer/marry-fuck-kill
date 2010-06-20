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

import logging

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
  One: name=<input type="text" name="n1"></input> url=<input type="text" name="u1"></input><br/>
  Two: name=<input type="text" name="n2"></input>url=<input type="text" name="u2"></input><br/>
  Three: name=<input type="text" name="n3"></input>url=<input type="text" name="u3"></input><br/>
  <input type="submit"></input>
</form>
""")

    def post(self):
        triple = TripleCreationHandler.MakeTriple(self.request) 
        utils.redirect(self, '/triple/view/' + triple.key().name())  

    @staticmethod
    def MakeTriple(request):
        one = request.get("n1")
        two = request.get("n2")
        three = request.get("n3")

        if not one or not two or not three:
            raise ValueError("request name")
        
        one = models.PutEntity(one, request.get("u1")) 
        two = models.PutEntity(two, request.get("u2"))
        three = models.PutEntity(three, request.get("u3"))

        return models.PutTriple(one=one, 
                                two=two, 
                                three=three)


class TripleJsonHandler(webapp.RequestHandler):
    def post(self, unused_id):
        try:
            triple = TripleCreationHandler.MakeTriple(self.request) 
        except ValueError:
            self.response.out.write("bad url")
        self.response.out.write("ok")

class TripleStatsHandler(webapp.RequestHandler):
    def get(self, triple_id):
        t = models.Triple.get_by_key_name(triple_id)
        self.response.out.write("%s: %s" % (triple_id, t))

