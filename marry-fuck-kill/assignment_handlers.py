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
%(1_name)s
<img src="%(1_url)s" alt="%(1_name)s">
<input type="hidden" name="e1" value="%(1_name)s"><br>
<label><input type="radio" name="v1" value="m"> Marry</label><br>
<label><input type="radio" name="v1" value="f"> Fuck</label><br>
<label><input type="radio" name="v1" value="k"> Kill</label><br>
<hr>
%(2_name)s
<img src="%(2_url)s" alt="%(2_name)s">
<input type="hidden" name="e2" value="%(2_name)s"><br>
<label><input type="radio" name="v2" value="m"> Marry</label><br>
<label><input type="radio" name="v2" value="f"> Fuck</label><br>
<label><input type="radio" name="v2" value="k"> Kill</label><br>
<hr>
%(3_name)s
<img src="%(3_url)s" alt="%(3_name)s">
<input type="hidden" name="e3" value="%(3_name)s"><br>
<label><input type="radio" name="v3" value="m"> Marry</label><br>
<label><input type="radio" name="v3" value="f"> Fuck</label><br>
<label><input type="radio" name="v3" value="k"> Kill</label><br>
<input type="submit">
""" % {'1_name': triple.one.key().name(),
       '1_url': triple.one.get_full_url(),
       '2_name': triple.two.key().name(),
       '2_url': triple.two.get_full_url(),
       '3_name': triple.three.key().name(),
       '3_url': triple.three.get_full_url()})

    def post(self, assignment_id):
        assign = AssignmentHandler.make_assignment(self.request)
        if assign is not None:
            self.response.out.write(str(assign))        
        else:
            self.response.set_status(406)
            self.response.out.write("select one of each!")

    @staticmethod
    def make_assignment(request):
        triple_key = request.get('key')
        values = [request.get('v1'), request.get('v2'), request.get('v3')]

        logging.debug("triple_key = %s", triple_key)
        logging.debug("values = %s", values)

        if set(values) != set(['m', 'f', 'k']):
            return None


        # We get an entity->action map from the client, but we need to reverse
        # it to action->entity to update the DB.
        triple = models.Triple.get(triple_key)
        logging.debug("triple = %s", triple)
        triple_entities = [triple.one,
                           triple.two,
                           triple.three]
        if (entities['m'] is None or entities['f'] is None
                or entities['k'] is None):
            logging.error("Not all non-None: marry = %s, fuck = %s, kill = %s",
                          entities['m'], entities['f'], entities['k'])
            return None

        if triple is None:
            logging.error("No triple with key %s", triple_key)
            return None

        entities = {}
        for i in range(len(values)):
            # Items in values are guaranteed to be 'm', 'f', 'k' (check above)
            entities[values[i]] = triple_entities[i]

        assign = models.Assignment(marry=entities['m'], 
                                   fuck=entities['f'], 
                                   kill=entities['k'])
        assign.put()
        logging.info("Assigned m=%s, f=%s, k=%s to %s", entities['m'],
                entities['f'], entities['k'], triple)
        return assign

class AssignmentJsonHandler(webapp.RequestHandler):
    def get(self, key=None):
        if key:
            triple = models.Triple.get(key)
            out = simplejson.dumps(triple.json())
        else:
            # give one at random
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
