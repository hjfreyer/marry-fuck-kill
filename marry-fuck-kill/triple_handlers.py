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
  One: name=<input type="text" name="one"></input> url=<input type="text" name="one_url"></input><br/>
  Two: name=<input type="text" name="two"></input>url=<input type="text" name="two_url"></input><br/>
  Three: name=<input type="text" name="three"></input>url=<input type="text" name="three_url"></input><br/>
  <input type="submit"></input>
</form>
""")

    def post(self):
        triple = TripleCreationHandler.MakeTriple(self.request) 
        utils.redirect(self, '/triple/view/' + triple.key().name())  

    @staticmethod
    def MakeTriple(request):
        one = request.get("one")
        two = request.get("two")
        three = request.get("three")
        
        one = models.PutEntity(one, request.get("one_url")) 
        two = models.PutEntity(two, request.get("two_url"))
        three = models.PutEntity(three, request.get("three_url"))

        return models.PutTriple(one=one, 
                                two=two, 
                                three=three)


class TripleJsonHandler(webapp.RequestHandler):
    pass


class TripleStatsHandler(webapp.RequestHandler):
    def get(self, triple_id):
        t = models.Triple.get_by_key_name(triple_id)
        self.response.out.write("%s: %s" % (triple_id, t))

