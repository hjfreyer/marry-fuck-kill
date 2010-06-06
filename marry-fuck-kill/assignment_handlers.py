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
<input type="hidden" name="one_id" value="%s"><br>
<input type="radio" name="one_value" value="Marry"> Marry<br>
<input type="radio" name="one_value" value="Fuck" checked> Fuck<br>
<input type="radio" name="one_value" value="Kill"> Kill
<hr>
%s
<input type="hidden" name="two_id" value="%s"><br>
<input type="radio" name="two_value" value="Marry"> Marry<br>
<input type="radio" name="two_value" value="Fuck"> Fuck<br>
<input type="radio" name="two_value" value="Kill" checked> Kill<br>
<hr>
%s
<input type="hidden" name="three_id" value="%s"><br>
<input type="radio" name="three_value" value="Marry"> Marry<br>
<input type="radio" name="three_value" value="Fuck"> Fuck<br>
<input type="radio" name="three_value" value="Kill" checked> Kill<br>
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
        one_id = request.get('one_id')
        one_value = request.get('one_value')
        two_id = request.get('two_id')
        two_value = request.get('two_value')
        three_id = request.get('three_id')
        three_value = request.get('three_value')

        #models.Triple.get_by_key_name(

        if set([one_value, two_value, three_value]) != set(['Marry', 'Fuck', 'Kill']):
            return None

        if one_value == 'Marry':
            marry = one_id
        if one_value == 'Fuck':
            fuck = one_id
        if one_value == 'Kill':
            kill = one_id

        if two_value == 'Marry':
            marry = two_id
        if two_value == 'Fuck':
            fuck = two_id
        if two_value == 'Kill':
            kill = two_id
        if three_value == 'Marry':
            marry = three_id
        if three_value == 'Fuck':
            fuck = three_id
        if three_value == 'Kill':
            kill = three_id

        marry = models.Entity.get_by_key_name(marry)
        fuck = models.Entity.get_by_key_name(fuck)
        kill = models.Entity.get_by_key_name(kill)

        assign = models.Assignment(marry=marry, 
                                   fuck=fuck, 
                                   kill=kill)
        assign.put()
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
