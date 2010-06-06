#!/usr/bin/env python
#
# Copyright 2010 Hunter Freyer and Michael Kelly
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

from google.appengine.ext import webapp
from google.appengine.ext.webapp import util
from django.utils import simplejson

import models

class AssignmentHandler(webapp.RequestHandler):
    def get(self, assignment_id=None):
        triple = models.Triple.get_random()

        self.response.out.write("""
<form method="post">
%s
<input type="hidden" name="e1" value="%s"><br>
<input type="radio" name="v1" value="m"> Marry<br>
<input type="radio" name="v1" value="f" checked> Fuck<br>
<input type="radio" name="v1" value="k"> Kill
<hr>
%s
<input type="hidden" name="e2" value="%s"><br>
<input type="radio" name="v2" value="m"> Marry<br>
<input type="radio" name="v2" value="f"> Fuck<br>
<input type="radio" name="v2" value="k" checked> Kill<br>
<hr>
%s
<input type="hidden" name="e3" value="%s"><br>
<input type="radio" name="v3" value="m"> Marry<br>
<input type="radio" name="v3" value="f"> Fuck<br>
<input type="radio" name="v3" value="k" checked> Kill<br>
<input type="submit">
""" % (triple.one.key().name(),
       triple.one.key().name(),
       triple.two.key().name(),
       triple.two.key().name(),
       triple.three.key().name(),
       triple.three.key().name()))

    def post(self, assignment_id):
        assign = AssignmentHandler.make_assignment(self.request)
        if assign is not None:
            self.response.out.write(str(assign))        
        else:
            self.response.set_status(406)
            self.response.out.write("select one of each!")

    @staticmethod
    def make_assignment(request):
        ids = [request.get('e1'), request.get('e2'), request.get('e3')]
        values = [request.get('v1'), request.get('v2'), request.get('v3')]

        if set(values) != set(['m', 'f', 'k']):
            return None

        marry = models.Entity.get_by_key_name(ids[values.index('m')])
        fuck = models.Entity.get_by_key_name(ids[values.index('f')])
        kill = models.Entity.get_by_key_name(ids[values.index('k')])

        if marry is None or fuck is None or kill is None:
            logging.error("Not all non-None: marry = %s, fuck = %s, kill = %s",
                          marry, fuck, kill)
            return None

        triple_key = models.Triple.key_name_from_entities(marry, fuck, kill)
        triple = models.Triple.get_by_key_name(triple_key)
        if triple is None:
            logging.error("No triple with key %s", triple_key)
            return None

        assign = models.Assignment(marry=marry, 
                                   fuck=fuck, 
                                   kill=kill)
        assign.put()
        logging.info("Assigned m=%s, f=%s, k=%s to %s", marry.key().name(),
                fuck.key().name(), kill.key().name(), triple.key().name())
        return assign

class AssignmentJsonHandler(webapp.RequestHandler):
    def get(self, assignment_id=None):
        triple = models.Triple.get_random()
        out = simplejson.dumps(triple.json())
        logging.info("AssignmentJsonHandler: sending: %s" % out)
        self.response.headers['Content-Type'] = 'application/json'
        self.response.out.write(out)

    def post(self, triple_id=None):
        logging.info("got assignment from client: %s", self.request)
        assign = AssignmentHandler.make_assignment(self.request)
        if assign is not None:
            self.response.out.write('ok')
        else:
            self.response.set_status(406)
            self.response.out.write('bad')
